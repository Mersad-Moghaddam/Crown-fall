package engine

import (
	"sort"

	matchdomain "crownfall/backend/internal/game/domain/match"
)

type PublicPlayerView struct {
	ID    string `json:"id"`
	Ready bool   `json:"ready"`
}

type PublicMatchView struct {
	ID             string             `json:"id"`
	Phase          matchdomain.Phase  `json:"phase"`
	Revision       uint64             `json:"revision"`
	SeedCommitment string             `json:"seed_commitment"`
	Players        []PublicPlayerView `json:"players"`
}

type PrivatePlayerView struct {
	PlayerID  string `json:"player_id"`
	RoleID    string `json:"role_id"`
	Objective string `json:"objective"`
}

type SpectatorView struct{ PublicMatchView }

func ProjectPublic(state matchdomain.State) PublicMatchView {
	view := PublicMatchView{ID: state.ID, Phase: state.Phase, Revision: state.Revision, SeedCommitment: state.SeedCommitment}
	ids := make([]string, 0, len(state.Players))
	for id := range state.Players {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		player := state.Players[id]
		view.Players = append(view.Players, PublicPlayerView{ID: player.ID, Ready: player.Ready})
	}
	return view
}

func ProjectPrivate(state matchdomain.State, playerID string) (PrivatePlayerView, bool) {
	player, ok := state.Players[playerID]
	return PrivatePlayerView{PlayerID: player.ID, RoleID: player.RoleID, Objective: player.Objective}, ok
}
