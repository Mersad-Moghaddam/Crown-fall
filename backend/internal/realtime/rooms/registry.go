package rooms

import (
	"context"
	"errors"
	"fmt"
	"sync"

	matchdomain "crownfall/backend/internal/game/domain/match"
	"crownfall/backend/internal/game/engine"
)

var (
	ErrMailboxFull = errors.New("room mailbox is full")
	ErrRoomClosed  = errors.New("room actor is closed")
	ErrActorPanic  = errors.New("room actor recovered a panic")
)

type Processor interface {
	Handle(context.Context, matchdomain.State, engine.Command) (engine.Result, error)
}

type request struct {
	context  context.Context
	command  engine.Command
	result   chan response
	snapshot bool
}
type response struct {
	result engine.Result
	state  matchdomain.State
	err    error
}

type Room struct {
	mailbox   chan request
	done      chan struct{}
	closeOnce sync.Once
	mu        sync.RWMutex
	closed    bool
	processor Processor
}

func NewRoom(state matchdomain.State, capacity int) *Room {
	return NewRoomWithProcessor(state, capacity, engine.Engine{})
}

func NewRoomWithProcessor(state matchdomain.State, capacity int, processor Processor) *Room {
	if capacity < 1 {
		capacity = 1
	}
	room := &Room{mailbox: make(chan request, capacity), done: make(chan struct{}), processor: processor}
	go room.run(state)
	return room
}

func (room *Room) run(state matchdomain.State) {
	defer close(room.done)
	for request := range room.mailbox {
		if request.snapshot {
			request.result <- response{state: state}
			continue
		}
		func() {
			defer func() {
				if recovered := recover(); recovered != nil {
					request.result <- response{err: fmt.Errorf("%w: %v", ErrActorPanic, recovered)}
				}
			}()
			result, err := room.processor.Handle(request.context, state, request.command)
			if err == nil {
				state = result.State
			}
			request.result <- response{result: result, err: err}
		}()
	}
}

func (room *Room) Handle(ctx context.Context, command engine.Command) (engine.Result, error) {
	response, err := room.submit(ctx, request{context: ctx, command: command, result: make(chan response, 1)})
	return response.result, err
}

func (room *Room) State(ctx context.Context) (matchdomain.State, error) {
	response, err := room.submit(ctx, request{context: ctx, snapshot: true, result: make(chan response, 1)})
	return response.state, err
}

func (room *Room) submit(ctx context.Context, request request) (response, error) {
	room.mu.RLock()
	if room.closed {
		room.mu.RUnlock()
		return response{}, ErrRoomClosed
	}
	select {
	case room.mailbox <- request:
		room.mu.RUnlock()
	case <-ctx.Done():
		room.mu.RUnlock()
		return response{}, ctx.Err()
	default:
		room.mu.RUnlock()
		return response{}, ErrMailboxFull
	}
	select {
	case result := <-request.result:
		return result, result.err
	case <-ctx.Done():
		return response{}, ctx.Err()
	}
}

func (room *Room) Close() {
	room.closeOnce.Do(func() { room.mu.Lock(); room.closed = true; close(room.mailbox); room.mu.Unlock(); <-room.done })
}

type Registry struct {
	mu    sync.RWMutex
	rooms map[string]*Room
}

func NewRegistry() *Registry { return &Registry{rooms: make(map[string]*Room)} }
func (registry *Registry) Add(id string, room *Room) bool {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	if _, exists := registry.rooms[id]; exists {
		return false
	}
	registry.rooms[id] = room
	return true
}
func (registry *Registry) Get(id string) (*Room, bool) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	room, ok := registry.rooms[id]
	return room, ok
}
func (registry *Registry) Remove(id string) bool {
	registry.mu.Lock()
	room, ok := registry.rooms[id]
	if ok {
		delete(registry.rooms, id)
	}
	registry.mu.Unlock()
	if ok {
		room.Close()
	}
	return ok
}
func (registry *Registry) Close() {
	registry.mu.Lock()
	current := registry.rooms
	registry.rooms = make(map[string]*Room)
	registry.mu.Unlock()
	for _, room := range current {
		room.Close()
	}
}
