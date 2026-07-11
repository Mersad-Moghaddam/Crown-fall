package match

type Phase string

const (
	PhaseLobby        Phase = "LOBBY"
	PhaseRoleDeal     Phase = "ROLE_DEAL"
	PhaseChapterStart Phase = "CHAPTER_START"
)

type Player struct {
	ID        string
	Ready     bool
	RoleID    string
	Objective string
}

type State struct {
	ID             string
	Phase          Phase
	Revision       uint64
	EventSequence  uint64
	SeedCommitment string
	Players        map[string]Player
	Accepted       map[string]uint64
}

func New(id, commitment string) State {
	return State{
		ID: id, Phase: PhaseLobby, SeedCommitment: commitment,
		Players: make(map[string]Player), Accepted: make(map[string]uint64),
	}
}
