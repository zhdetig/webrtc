package webrtc

// ICERole describes the role ice.Agent is playing in selecting the
// preferred the candidate pair.
type ICERole int

const (
	// ICERoleControlling indicates that the ICE agent that is responsible
	// for selecting the final choice of candidate pairs and signaling them
	// through STUN and an updated offer, if needed. In any session, one agent
	// is always controlling. The other is the controlled agent.
	ICERoleControlling ICERole = iota + 1

	// ICERoleControlled indicates that an ICE agent that waits for the
	// controlling agent to select the final choice of candidate pairs.
	ICERoleControlled

	// ICERoleUnknown is a placeholder for error processing
	ICERoleUnknown
)

// This is done this way because of a linter.
const (
	iceRoleControllingStr = "controlling"
	iceRoleControlledStr  = "controlled"
)

func newICERole(raw string) ICERole {
	switch raw {
	case iceRoleControllingStr:
		return ICERoleControlling
	case iceRoleControlledStr:
		return ICERoleControlled
	default:
		return ICERoleUnknown
	}
}

func (t ICERole) String() string {
	switch t {
	case ICERoleControlling:
		return iceRoleControllingStr
	case ICERoleControlled:
		return iceRoleControlledStr
	default:
		return ErrUnknownType.Error()
	}
}
