package engine

import (
	"context"
	"errors"
	"fmt"
	"time"

	matchdomain "crownfall/backend/internal/game/domain/match"
)

const CommandStartMatch = "START_MATCH"

var (
	ErrDuplicateCommand  = errors.New("command already accepted")
	ErrStaleRevision     = errors.New("expected revision does not match")
	ErrCommandNotAllowed = errors.New("command not allowed in phase")
)

type Command struct {
	CommandID        string         `json:"command_id"`
	MatchID          string         `json:"match_id"`
	PlayerID         string         `json:"player_id"`
	ExpectedRevision uint64         `json:"expected_revision"`
	CommandType      string         `json:"command_type"`
	Payload          map[string]any `json:"payload"`
	ClientTimestamp  time.Time      `json:"client_timestamp"`
	ClientSequence   uint64         `json:"client_sequence"`
}

type Event struct {
	Type     string            `json:"type"`
	Phase    matchdomain.Phase `json:"phase"`
	Sequence uint64            `json:"sequence"`
	Payload  map[string]any    `json:"payload"`
}

type Result struct {
	State         matchdomain.State
	PublicEvents  []Event
	PrivateEvents map[string][]Event
	DomainEvents  []Event
}

type Engine struct{}

func (Engine) Handle(_ context.Context, state matchdomain.State, command Command) (Result, error) {
	if command.MatchID != state.ID || command.CommandID == "" || command.PlayerID == "" {
		return Result{}, fmt.Errorf("invalid command identity")
	}
	if _, exists := state.Accepted[command.CommandID]; exists {
		return Result{}, ErrDuplicateCommand
	}
	if command.ExpectedRevision != state.Revision {
		return Result{}, ErrStaleRevision
	}
	if state.Phase != matchdomain.PhaseLobby || command.CommandType != CommandStartMatch {
		return Result{}, ErrCommandNotAllowed
	}
	if len(state.Players) < 1 {
		return Result{}, errors.New("cannot start an empty match")
	}

	state.Revision++
	state.EventSequence++
	state.Phase = matchdomain.PhaseRoleDeal
	state.Accepted[command.CommandID] = state.Revision
	event := Event{Type: "match.roleDealStarted", Phase: state.Phase, Sequence: state.EventSequence, Payload: map[string]any{"revision": state.Revision}}
	private := make(map[string][]Event, len(state.Players))
	for id, player := range state.Players {
		private[id] = []Event{{Type: "role.assigned", Phase: state.Phase, Sequence: state.EventSequence, Payload: map[string]any{"role_id": player.RoleID, "objective": player.Objective}}}
	}
	return Result{State: state, PublicEvents: []Event{event}, PrivateEvents: private, DomainEvents: []Event{event}}, nil
}
