package main

import (
	"context"
	"fmt"
	"time"

	matchdomain "crownfall/backend/internal/game/domain/match"
	"crownfall/backend/internal/game/engine"
)

func main() {
	state := matchdomain.New("simulation", "fixed-commitment")
	state.Players["player-1"] = matchdomain.Player{ID: "player-1", Ready: true, RoleID: "fixture-role"}
	result, err := (engine.Engine{}).Handle(context.Background(), state, engine.Command{CommandID: "simulation-1", MatchID: state.ID, PlayerID: "player-1", CommandType: engine.CommandStartMatch, ClientTimestamp: time.Unix(0, 0).UTC(), ClientSequence: 1})
	if err != nil {
		panic(err)
	}
	fmt.Printf("phase=%s revision=%d events=%d\n", result.State.Phase, result.State.Revision, len(result.DomainEvents))
}
