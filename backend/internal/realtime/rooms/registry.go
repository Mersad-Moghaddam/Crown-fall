package rooms

import (
	"context"
	"errors"
	"sync"

	matchdomain "crownfall/backend/internal/game/domain/match"
	"crownfall/backend/internal/game/engine"
)

var ErrMailboxFull = errors.New("room mailbox is full")

type request struct {
	context context.Context
	command engine.Command
	result  chan response
}

type response struct {
	result engine.Result
	err    error
}

type Room struct {
	mailbox chan request
	done    chan struct{}
}

func NewRoom(state matchdomain.State, capacity int) *Room {
	room := &Room{mailbox: make(chan request, capacity), done: make(chan struct{})}
	go room.run(state)
	return room
}

func (room *Room) run(state matchdomain.State) {
	defer close(room.done)
	processor := engine.Engine{}
	for request := range room.mailbox {
		result, err := processor.Handle(request.context, state, request.command)
		if err == nil {
			state = result.State
		}
		request.result <- response{result: result, err: err}
	}
}

func (room *Room) Handle(ctx context.Context, command engine.Command) (engine.Result, error) {
	result := make(chan response, 1)
	select {
	case room.mailbox <- request{context: ctx, command: command, result: result}:
	default:
		return engine.Result{}, ErrMailboxFull
	}
	select {
	case response := <-result:
		return response.result, response.err
	case <-ctx.Done():
		return engine.Result{}, ctx.Err()
	}
}

func (room *Room) Close() { close(room.mailbox); <-room.done }

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

func (registry *Registry) Close() {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	for id, room := range registry.rooms {
		room.Close()
		delete(registry.rooms, id)
	}
}
