package engine

import (
	"sort"

	matchdomain "crownfall/backend/internal/game/domain/match"
)

type PublicPlayerView struct {
	ID               string `json:"id"`
	Ready            bool   `json:"ready"`
	Connected        bool   `json:"connected"`
	RoleAcknowledged bool   `json:"roleAcknowledged"`
}

type PublicMatchView struct {
	ID             string             `json:"id"`
	HostPlayerID   string             `json:"hostPlayerId"`
	Phase          matchdomain.Phase  `json:"phase"`
	Revision       uint64             `json:"revision"`
	EventSequence  uint64             `json:"eventSequence"`
	SeedCommitment string             `json:"seedCommitment,omitempty"`
	Players        []PublicPlayerView `json:"players"`
}

type PrivatePlayerView struct {
	PlayerID         string            `json:"playerId"`
	Phase            matchdomain.Phase `json:"phase"`
	Revision         uint64            `json:"revision"`
	EventSequence    uint64            `json:"eventSequence"`
	RoleID           string            `json:"roleId,omitempty"`
	Faction          string            `json:"faction,omitempty"`
	Objective        string            `json:"objective,omitempty"`
	RoleAcknowledged bool              `json:"roleAcknowledged"`
}

type SpectatorView struct{ PublicMatchView }

type ResyncView struct {
	Public  PublicMatchView   `json:"public"`
	Private PrivatePlayerView `json:"private"`
}

func ProjectPublic(state matchdomain.State) PublicMatchView {
	view := PublicMatchView{ID: state.ID, HostPlayerID: state.HostPlayerID, Phase: state.Phase, Revision: state.Revision, EventSequence: state.EventSequence}
	if state.Phase != matchdomain.PhaseLobby {
		view.SeedCommitment = state.SeedCommitment
	}
	ids := append([]string(nil), state.PlayerOrder...)
	if len(ids) == 0 {
		for id := range state.Players {
			ids = append(ids, id)
		}
		sort.Strings(ids)
	}
	for _, id := range ids {
		player := state.Players[id]
		view.Players = append(view.Players, PublicPlayerView{ID: player.ID, Ready: player.Ready, Connected: player.Connected, RoleAcknowledged: player.RoleAcknowledged})
	}
	return view
}

func ProjectPrivate(state matchdomain.State, playerID string) (PrivatePlayerView, bool) {
	player, ok := state.Players[playerID]
	if !ok {
		return PrivatePlayerView{}, false
	}
	return PrivatePlayerView{PlayerID: player.ID, Phase: state.Phase, Revision: state.Revision, EventSequence: state.EventSequence, RoleID: player.RoleID, Faction: player.Faction, Objective: player.Objective, RoleAcknowledged: player.RoleAcknowledged}, true
}

func ProjectResync(state matchdomain.State, playerID string) (ResyncView, bool) {
	private, ok := ProjectPrivate(state, playerID)
	return ResyncView{Public: ProjectPublic(state), Private: private}, ok
}
