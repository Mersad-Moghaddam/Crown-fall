package engine

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	matchdomain "crownfall/backend/internal/game/domain/match"
	platformrandom "crownfall/backend/internal/platform/random"
)

const (
	CommandJoinRoom        = "JOIN_ROOM"
	CommandSetReady        = "SET_READY"
	CommandStartMatch      = "START_MATCH"
	CommandAcknowledgeRole = "ACKNOWLEDGE_ROLE"
)

var (
	ErrDuplicateCommandConflict = errors.New("command id reused with different content")
	ErrStaleRevision            = errors.New("expected revision does not match")
	ErrCommandNotAllowed        = errors.New("command not allowed in phase")
	ErrInvalidPlayerCount       = errors.New("match requires 6 to 10 players")
	ErrPlayersNotReady          = errors.New("all players must be ready")
	ErrHostOnly                 = errors.New("only the host may start the match")
	ErrNotMember                = errors.New("player is not a match member")
	ErrClientSequence           = errors.New("client sequence must increase monotonically")
)

type Command struct {
	CommandID        string         `json:"commandId"`
	MatchID          string         `json:"matchId"`
	PlayerID         string         `json:"playerId"`
	ExpectedRevision uint64         `json:"expectedRevision"`
	CommandType      string         `json:"commandType"`
	Payload          map[string]any `json:"payload"`
	ClientTimestamp  time.Time      `json:"clientTimestamp"`
	ClientSequence   uint64         `json:"clientSequence"`
}

type Event struct {
	Type     string            `json:"type"`
	Phase    matchdomain.Phase `json:"phase"`
	Sequence uint64            `json:"sequence"`
	Revision uint64            `json:"revision"`
	Payload  map[string]any    `json:"payload"`
}

type Result struct {
	State         matchdomain.State
	PublicEvents  []Event
	PrivateEvents map[string][]Event
	DomainEvents  []Event
	Duplicate     bool
}

type Engine struct{}

func (Engine) Handle(ctx context.Context, state matchdomain.State, command Command) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	if err := validateIdentity(state, command); err != nil {
		return Result{}, err
	}
	fingerprint, err := fingerprint(command)
	if err != nil {
		return Result{}, err
	}
	if accepted, exists := state.Accepted[command.CommandID]; exists {
		if accepted.Fingerprint != fingerprint {
			return Result{}, ErrDuplicateCommandConflict
		}
		return Result{State: state, Duplicate: true, PrivateEvents: map[string][]Event{}}, nil
	}
	if command.ExpectedRevision != state.Revision {
		return Result{}, ErrStaleRevision
	}
	if player, exists := state.Players[command.PlayerID]; exists && command.ClientSequence <= player.LastClientSequence {
		return Result{}, ErrClientSequence
	}
	state = cloneState(state)

	var public []Event
	private := map[string][]Event{}
	switch command.CommandType {
	case CommandJoinRoom:
		if state.Phase != matchdomain.PhaseLobby {
			return Result{}, ErrCommandNotAllowed
		}
		if _, exists := state.Players[command.PlayerID]; exists {
			return Result{}, errors.New("player already joined")
		}
		if len(state.Players) >= matchdomain.MaximumPlayers {
			return Result{}, ErrInvalidPlayerCount
		}
		state.Players[command.PlayerID] = matchdomain.Player{ID: command.PlayerID, Connected: true, LastClientSequence: command.ClientSequence}
		state.PlayerOrder = append(state.PlayerOrder, command.PlayerID)
		if state.HostPlayerID == "" {
			state.HostPlayerID = command.PlayerID
		}
		public = append(public, event(&state, "room.playerJoined", map[string]any{"playerId": command.PlayerID, "host": state.HostPlayerID == command.PlayerID}))
	case CommandSetReady:
		if state.Phase != matchdomain.PhaseLobby {
			return Result{}, ErrCommandNotAllowed
		}
		player, exists := state.Players[command.PlayerID]
		if !exists {
			return Result{}, ErrNotMember
		}
		ready, ok := command.Payload["ready"].(bool)
		if !ok {
			return Result{}, errors.New("ready payload must be boolean")
		}
		player.Ready, player.LastClientSequence = ready, command.ClientSequence
		state.Players[command.PlayerID] = player
		public = append(public, event(&state, "room.playerReadyChanged", map[string]any{"playerId": command.PlayerID, "ready": ready}))
	case CommandStartMatch:
		if state.Phase != matchdomain.PhaseLobby {
			return Result{}, ErrCommandNotAllowed
		}
		if command.PlayerID != state.HostPlayerID {
			return Result{}, ErrHostOnly
		}
		if len(state.Players) < matchdomain.MinimumPlayers || len(state.Players) > matchdomain.MaximumPlayers {
			return Result{}, ErrInvalidPlayerCount
		}
		for _, player := range state.Players {
			if !player.Ready {
				return Result{}, ErrPlayersNotReady
			}
		}
		assignRoles(&state)
		state.Phase = matchdomain.PhaseRoleDeal
		public = append(public, event(&state, "match.roleDealStarted", map[string]any{"seedCommitment": state.SeedCommitment, "playerCount": len(state.Players)}))
		for _, playerID := range state.PlayerOrder {
			player := state.Players[playerID]
			private[playerID] = []Event{event(&state, "role.assigned", map[string]any{"roleId": player.RoleID, "faction": player.Faction, "objective": player.Objective})}
		}
	case CommandAcknowledgeRole:
		if state.Phase != matchdomain.PhaseRoleDeal {
			return Result{}, ErrCommandNotAllowed
		}
		player, exists := state.Players[command.PlayerID]
		if !exists {
			return Result{}, ErrNotMember
		}
		if player.RoleAcknowledged {
			return Result{}, errors.New("role already acknowledged")
		}
		player.RoleAcknowledged, player.LastClientSequence = true, command.ClientSequence
		state.Players[command.PlayerID] = player
		public = append(public, event(&state, "role.acknowledged", map[string]any{"playerId": command.PlayerID}))
		if allAcknowledged(state) {
			state.Phase = matchdomain.PhaseChapterStart
			public = append(public, event(&state, "match.chapterStarted", map[string]any{"chapter": 1}))
		}
	default:
		return Result{}, ErrCommandNotAllowed
	}

	state.Revision++
	for index := range public {
		public[index].Revision = state.Revision
	}
	for playerID := range private {
		for index := range private[playerID] {
			private[playerID][index].Revision = state.Revision
		}
	}
	state.Accepted[command.CommandID] = matchdomain.AcceptedCommand{Fingerprint: fingerprint, Revision: state.Revision, EventSequence: state.EventSequence}
	domain := append([]Event(nil), public...)
	for _, playerID := range state.PlayerOrder {
		domain = append(domain, private[playerID]...)
	}
	sort.Slice(domain, func(i, j int) bool { return domain[i].Sequence < domain[j].Sequence })
	return Result{State: state, PublicEvents: public, PrivateEvents: private, DomainEvents: domain}, nil
}

func cloneState(state matchdomain.State) matchdomain.State {
	state.Seed = append([]byte(nil), state.Seed...)
	state.PlayerOrder = append([]string(nil), state.PlayerOrder...)
	players := make(map[string]matchdomain.Player, len(state.Players))
	for id, player := range state.Players {
		players[id] = player
	}
	state.Players = players
	accepted := make(map[string]matchdomain.AcceptedCommand, len(state.Accepted))
	for id, command := range state.Accepted {
		accepted[id] = command
	}
	state.Accepted = accepted
	return state
}

func validateIdentity(state matchdomain.State, command Command) error {
	if command.CommandID == "" || command.MatchID != state.ID || command.PlayerID == "" || command.ClientSequence == 0 || command.ClientTimestamp.IsZero() {
		return errors.New("invalid command envelope")
	}
	return nil
}

func fingerprint(command Command) (string, error) {
	payload, err := json.Marshal(struct {
		MatchID, PlayerID, Type string
		Payload                 map[string]any
		ClientSequence          uint64
	}{command.MatchID, command.PlayerID, command.CommandType, command.Payload, command.ClientSequence})
	if err != nil {
		return "", fmt.Errorf("fingerprint command: %w", err)
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}

func event(state *matchdomain.State, eventType string, payload map[string]any) Event {
	state.EventSequence++
	return Event{Type: eventType, Phase: state.Phase, Sequence: state.EventSequence, Payload: payload}
}

func assignRoles(state *matchdomain.State) {
	type ranked struct {
		id    string
		score string
	}
	rankedPlayers := make([]ranked, 0, len(state.PlayerOrder))
	stream := platformrandom.Derive(state.Seed, "roleAssignment")
	for _, id := range state.PlayerOrder {
		rankedPlayers = append(rankedPlayers, ranked{id: id, score: hex.EncodeToString(platformrandom.Derive(stream, id))})
	}
	sort.Slice(rankedPlayers, func(i, j int) bool { return rankedPlayers[i].score < rankedPlayers[j].score })
	shadowCount := 2
	if len(rankedPlayers) >= 9 {
		shadowCount = 3
	}
	for index, entry := range rankedPlayers {
		player := state.Players[entry.id]
		switch {
		case index == 0:
			player.RoleID, player.Faction, player.Objective = "usurper", "SHADOW", "Become elected Pathfinder after Dread reaches two."
		case index < shadowCount:
			player.RoleID, player.Faction, player.Objective = "nightblade", "SHADOW", "Advance the Shadow ending without exposing the Usurper."
		default:
			player.RoleID, player.Faction, player.Objective = "oathboundKnight", "CROWN", "Advance Hope and protect the realm."
		}
		state.Players[entry.id] = player
	}
}

func allAcknowledged(state matchdomain.State) bool {
	if len(state.Players) == 0 {
		return false
	}
	for _, player := range state.Players {
		if !player.RoleAcknowledged {
			return false
		}
	}
	return true
}
