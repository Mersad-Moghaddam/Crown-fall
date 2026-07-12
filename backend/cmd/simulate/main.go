package main

import (
	"context"
	"fmt"
	"time"

	matchdomain "crownfall/backend/internal/game/domain/match"
	"crownfall/backend/internal/game/engine"
	platformrandom "crownfall/backend/internal/platform/random"
)

func main() {
	seed := []byte("crownfall-deterministic-simulation")
	state := matchdomain.New("simulation", seed, platformrandom.Commitment(seed))
	processor := engine.Engine{}
	sequence := map[string]uint64{}
	handle := func(id, playerID, commandType string, payload map[string]any) {
		sequence[playerID]++
		result, err := processor.Handle(context.Background(), state, engine.Command{CommandID: id, MatchID: state.ID, PlayerID: playerID, ExpectedRevision: state.Revision, CommandType: commandType, Payload: payload, ClientTimestamp: time.Unix(1, 0).UTC(), ClientSequence: sequence[playerID]})
		if err != nil {
			panic(err)
		}
		state = result.State
	}
	for index := 1; index <= 6; index++ {
		id := fmt.Sprintf("player-%d", index)
		handle("join-"+id, id, engine.CommandJoinRoom, nil)
		handle("ready-"+id, id, engine.CommandSetReady, map[string]any{"ready": true})
	}
	handle("start", "player-1", engine.CommandStartMatch, nil)
	for index := 1; index <= 6; index++ {
		id := fmt.Sprintf("player-%d", index)
		handle("ack-"+id, id, engine.CommandAcknowledgeRole, nil)
	}
	fmt.Printf("phase=%s revision=%d events=%d commitment=%s\n", state.Phase, state.Revision, state.EventSequence, state.SeedCommitment)
}
