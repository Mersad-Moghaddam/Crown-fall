package match

type Phase string

const (
	PhaseLobby        Phase = "LOBBY"
	PhaseRoleDeal     Phase = "ROLE_DEAL"
	PhaseChapterStart Phase = "CHAPTER_START"
)

const (
	MinimumPlayers = 6
	MaximumPlayers = 10
)

type Player struct {
	ID                 string
	Ready              bool
	Connected          bool
	RoleID             string
	Faction            string
	Objective          string
	RoleAcknowledged   bool
	LastClientSequence uint64
}

type AcceptedCommand struct {
	Fingerprint   string
	Revision      uint64
	EventSequence uint64
}

// State is internal engine state. Transport packages must use explicit projections.
type State struct {
	ID             string
	HostPlayerID   string
	Phase          Phase
	Revision       uint64
	EventSequence  uint64
	Seed           []byte
	SeedCommitment string
	Players        map[string]Player
	PlayerOrder    []string
	Accepted       map[string]AcceptedCommand
}

func New(id string, seed []byte, commitment string) State {
	return State{
		ID: id, Phase: PhaseLobby, Seed: append([]byte(nil), seed...), SeedCommitment: commitment,
		Players: make(map[string]Player), Accepted: make(map[string]AcceptedCommand),
	}
}
