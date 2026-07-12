package rooms

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	matchdomain "crownfall/backend/internal/game/domain/match"
	"crownfall/backend/internal/game/engine"
)

type blockingProcessor struct {
	entered chan struct{}
	release chan struct{}
	panic   bool
}

func (processor *blockingProcessor) Handle(ctx context.Context, state matchdomain.State, command engine.Command) (engine.Result, error) {
	if processor.panic {
		panic("test panic")
	}
	processor.entered <- struct{}{}
	select {
	case <-processor.release:
		state.Revision++
		return engine.Result{State: state}, nil
	case <-ctx.Done():
		return engine.Result{}, ctx.Err()
	}
}

func TestMailboxOverflowShutdownAndPanicContainment(t *testing.T) {
	state := matchdomain.New("match", []byte("seed"), "commitment")
	processor := &blockingProcessor{entered: make(chan struct{}, 1), release: make(chan struct{})}
	room := NewRoomWithProcessor(state, 1, processor)
	firstDone := make(chan struct{})
	go func() { _, _ = room.Handle(context.Background(), engine.Command{}); close(firstDone) }()
	<-processor.entered
	secondDone := make(chan struct{})
	go func() { _, _ = room.Handle(context.Background(), engine.Command{}); close(secondDone) }()
	time.Sleep(10 * time.Millisecond)
	if _, err := room.Handle(context.Background(), engine.Command{}); !errors.Is(err, ErrMailboxFull) {
		t.Fatalf("expected overflow, got %v", err)
	}
	close(processor.release)
	<-firstDone
	<-secondDone
	room.Close()
	if _, err := room.Handle(context.Background(), engine.Command{}); !errors.Is(err, ErrRoomClosed) {
		t.Fatalf("expected closed, got %v", err)
	}

	panicking := NewRoomWithProcessor(state, 1, &blockingProcessor{panic: true})
	if _, err := panicking.Handle(context.Background(), engine.Command{}); !errors.Is(err, ErrActorPanic) {
		t.Fatalf("expected recovered panic, got %v", err)
	}
	panicking.Close()
}

func TestConcurrentCommandsAreSerialized(t *testing.T) {
	state := matchdomain.New("match", []byte("seed"), "commitment")
	processor := &blockingProcessor{entered: make(chan struct{}, 10), release: make(chan struct{}, 10)}
	room := NewRoomWithProcessor(state, 10, processor)
	defer room.Close()
	var wait sync.WaitGroup
	for index := 0; index < 5; index++ {
		wait.Add(1)
		go func() { defer wait.Done(); _, _ = room.Handle(context.Background(), engine.Command{}) }()
	}
	for index := 0; index < 5; index++ {
		<-processor.entered
		processor.release <- struct{}{}
	}
	wait.Wait()
	current, err := room.State(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if current.Revision != 5 {
		t.Fatalf("want 5 serialized revisions, got %d", current.Revision)
	}
}
