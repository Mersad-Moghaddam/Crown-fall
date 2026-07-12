package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	matchdomain "crownfall/backend/internal/game/domain/match"
	platformrandom "crownfall/backend/internal/platform/random"
)

func newState(seed string) matchdomain.State {
	return matchdomain.New("match-1", []byte(seed), platformrandom.Commitment([]byte(seed)))
}

func command(id, player, commandType string, revision, sequence uint64, payload map[string]any) Command {
	return Command{CommandID: id, MatchID: "match-1", PlayerID: player, ExpectedRevision: revision, CommandType: commandType, Payload: payload, ClientTimestamp: time.Unix(1, 0).UTC(), ClientSequence: sequence}
}

func bootstrapLobby(t *testing.T, seed string, count int) matchdomain.State {
	t.Helper()
	state := newState(seed)
	for index := 1; index <= count; index++ {
		playerID := fmt.Sprintf("player-%d", index)
		result, err := (Engine{}).Handle(context.Background(), state, command("join-"+playerID, playerID, CommandJoinRoom, state.Revision, 1, nil))
		if err != nil {
			t.Fatal(err)
		}
		state = result.State
		result, err = (Engine{}).Handle(context.Background(), state, command("ready-"+playerID, playerID, CommandSetReady, state.Revision, 2, map[string]any{"ready": true}))
		if err != nil {
			t.Fatal(err)
		}
		state = result.State
	}
	return state
}

func TestBootstrapTransitionsToChapterStartExactlyOnce(t *testing.T) {
	state := bootstrapLobby(t, "fixed-seed", 6)
	started, err := (Engine{}).Handle(context.Background(), state, command("start", "player-1", CommandStartMatch, state.Revision, 3, nil))
	if err != nil {
		t.Fatal(err)
	}
	if started.State.Phase != matchdomain.PhaseRoleDeal || started.State.Revision != state.Revision+1 {
		t.Fatalf("unexpected role deal: %+v", started.State)
	}
	if started.PublicEvents[0].Payload["seedCommitment"] == "" || len(started.PrivateEvents) != 6 {
		t.Fatal("role deal did not publish commitment and recipient events")
	}
	state = started.State
	for index := 1; index <= 6; index++ {
		playerID := fmt.Sprintf("player-%d", index)
		result, err := (Engine{}).Handle(context.Background(), state, command("ack-"+playerID, playerID, CommandAcknowledgeRole, state.Revision, 3, nil))
		if err != nil {
			t.Fatal(err)
		}
		state = result.State
		if index < 6 && state.Phase != matchdomain.PhaseRoleDeal {
			t.Fatal("chapter started before all acknowledgements")
		}
	}
	if state.Phase != matchdomain.PhaseChapterStart {
		t.Fatal("final acknowledgement did not start chapter")
	}
	revision := state.Revision
	if _, err := (Engine{}).Handle(context.Background(), state, command("another-ack", "player-1", CommandAcknowledgeRole, revision, 4, nil)); !errors.Is(err, ErrCommandNotAllowed) {
		t.Fatalf("expected phase rejection, got %v", err)
	}
	if state.Revision != revision {
		t.Fatal("rejection mutated revision")
	}
}

func TestStartPreconditions(t *testing.T) {
	tests := []struct {
		name   string
		state  func(*testing.T) matchdomain.State
		player string
		want   error
	}{
		{"invalid count", func(t *testing.T) matchdomain.State { return bootstrapLobby(t, "seed", 5) }, "player-1", ErrInvalidPlayerCount},
		{"non host", func(t *testing.T) matchdomain.State { return bootstrapLobby(t, "seed", 6) }, "player-2", ErrHostOnly},
		{"unready", func(t *testing.T) matchdomain.State {
			state := bootstrapLobby(t, "seed", 6)
			player := state.Players["player-6"]
			player.Ready = false
			state.Players[player.ID] = player
			return state
		}, "player-1", ErrPlayersNotReady},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.state(t)
			revision := state.Revision
			_, err := (Engine{}).Handle(context.Background(), state, command("start", test.player, CommandStartMatch, revision, 3, nil))
			if !errors.Is(err, test.want) {
				t.Fatalf("want %v, got %v", test.want, err)
			}
			if state.Revision != revision {
				t.Fatal("rejection mutated state")
			}
		})
	}
}

func TestIdempotencyRevisionAndSequence(t *testing.T) {
	state := newState("seed")
	join := command("join", "player-1", CommandJoinRoom, 0, 1, nil)
	accepted, err := (Engine{}).Handle(context.Background(), state, join)
	if err != nil {
		t.Fatal(err)
	}
	duplicate, err := (Engine{}).Handle(context.Background(), accepted.State, join)
	if err != nil || !duplicate.Duplicate || duplicate.State.Revision != 1 {
		t.Fatalf("bad duplicate result: %+v %v", duplicate, err)
	}
	conflict := join
	conflict.CommandType = CommandSetReady
	if _, err := (Engine{}).Handle(context.Background(), accepted.State, conflict); !errors.Is(err, ErrDuplicateCommandConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
	if _, err := (Engine{}).Handle(context.Background(), accepted.State, command("stale", "player-1", CommandSetReady, 0, 2, map[string]any{"ready": true})); !errors.Is(err, ErrStaleRevision) {
		t.Fatalf("expected stale, got %v", err)
	}
	if _, err := (Engine{}).Handle(context.Background(), accepted.State, command("sequence", "player-1", CommandSetReady, 1, 1, map[string]any{"ready": true})); !errors.Is(err, ErrClientSequence) {
		t.Fatalf("expected sequence rejection, got %v", err)
	}
}

func TestRoleAssignmentDeterminismAndSecrecy(t *testing.T) {
	assign := func(seed string) matchdomain.State {
		state := bootstrapLobby(t, seed, 6)
		result, err := (Engine{}).Handle(context.Background(), state, command("start", "player-1", CommandStartMatch, state.Revision, 3, nil))
		if err != nil {
			t.Fatal(err)
		}
		return result.State
	}
	first, second := assign("same-seed"), assign("same-seed")
	for id := range first.Players {
		if first.Players[id].RoleID != second.Players[id].RoleID {
			t.Fatal("same seed produced different assignment")
		}
	}
	publicJSON, _ := json.Marshal(ProjectPublic(first))
	if strings.Contains(string(publicJSON), "usurper") || strings.Contains(string(publicJSON), "SHADOW") || strings.Contains(string(publicJSON), "roleId") || strings.Contains(string(publicJSON), "faction") || strings.Contains(string(publicJSON), "objective") {
		t.Fatalf("public leak: %s", publicJSON)
	}
	private, ok := ProjectPrivate(first, "player-1")
	if !ok || private.RoleID == "" {
		t.Fatal("recipient role missing")
	}
	privateJSON, _ := json.Marshal(private)
	for id := range first.Players {
		if id != "player-1" && strings.Contains(string(privateJSON), id) {
			t.Fatalf("private view leaked %s", id)
		}
	}
	resync, ok := ProjectResync(first, "player-1")
	if !ok || resync.Private.PlayerID != "player-1" || resync.Private.RoleAcknowledged {
		t.Fatal("invalid reconnect projection")
	}
	if first.Revision != resync.Public.Revision || first.EventSequence != resync.Private.EventSequence {
		t.Fatal("resync versions differ")
	}
}

func TestDifferentSeedsCanChangeAssignments(t *testing.T) {
	one, two := bootstrapLobby(t, "seed-one", 6), bootstrapLobby(t, "seed-two", 6)
	oneResult, _ := (Engine{}).Handle(context.Background(), one, command("start", "player-1", CommandStartMatch, one.Revision, 3, nil))
	twoResult, _ := (Engine{}).Handle(context.Background(), two, command("start", "player-1", CommandStartMatch, two.Revision, 3, nil))
	oneRoles, twoRoles := map[string]string{}, map[string]string{}
	for id, player := range oneResult.State.Players {
		oneRoles[id] = player.RoleID
	}
	for id, player := range twoResult.State.Players {
		twoRoles[id] = player.RoleID
	}
	if reflect.DeepEqual(oneRoles, twoRoles) {
		t.Fatal("different seeds unexpectedly produced identical assignment")
	}
}

func BenchmarkStartMatch(b *testing.B) {
	for index := 0; index < b.N; index++ {
		seed := fmt.Sprint(index)
		state := newState(seed)
		for playerIndex := 1; playerIndex <= 6; playerIndex++ {
			id := fmt.Sprintf("player-%d", playerIndex)
			state.Players[id] = matchdomain.Player{ID: id, Ready: true, Connected: true, LastClientSequence: 2}
			state.PlayerOrder = append(state.PlayerOrder, id)
		}
		state.HostPlayerID, state.Revision = "player-1", 12
		_, _ = (Engine{}).Handle(context.Background(), state, command("start", "player-1", CommandStartMatch, state.Revision, 3, nil))
	}
}
