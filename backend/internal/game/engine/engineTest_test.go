package engine

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	matchdomain "crownfall/backend/internal/game/domain/match"
)

func state() matchdomain.State {
	value := matchdomain.New("match-1", "commitment")
	value.Players["player-1"] = matchdomain.Player{ID: "player-1", Ready: true, RoleID: "usurper", Objective: "hidden"}
	return value
}

func command(revision uint64) Command {
	return Command{CommandID: "command-1", MatchID: "match-1", PlayerID: "player-1", ExpectedRevision: revision, CommandType: CommandStartMatch}
}

func TestStartMatchTransition(t *testing.T) {
	result, err := (Engine{}).Handle(context.Background(), state(), command(0))
	if err != nil {
		t.Fatal(err)
	}
	if result.State.Phase != matchdomain.PhaseRoleDeal || result.State.Revision != 1 || result.State.EventSequence != 1 {
		t.Fatalf("unexpected state: %+v", result.State)
	}
	if result.PrivateEvents["player-1"][0].Payload["role_id"] != "usurper" {
		t.Fatal("missing private role")
	}
	public := ProjectPublic(result.State)
	if len(public.Players) != 1 {
		t.Fatal("missing public player")
	}
}

func TestRejectsStaleAndDuplicateCommands(t *testing.T) {
	if _, err := (Engine{}).Handle(context.Background(), state(), command(1)); !errors.Is(err, ErrStaleRevision) {
		t.Fatalf("expected stale revision, got %v", err)
	}
	accepted, _ := (Engine{}).Handle(context.Background(), state(), command(0))
	if _, err := (Engine{}).Handle(context.Background(), accepted.State, command(1)); !errors.Is(err, ErrDuplicateCommand) {
		t.Fatalf("expected duplicate, got %v", err)
	}
}

func TestPublicProjectionDoesNotExposeRole(t *testing.T) {
	view := ProjectPublic(state())
	if view.Players[0].ID != "player-1" {
		t.Fatal("unexpected player")
	}
	payload, err := json.Marshal(view)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(payload), "usurper") || strings.Contains(string(payload), "hidden") {
		t.Fatalf("public projection leaked private state: %s", payload)
	}
}
