package models

type AssetState interface {
	StateKind() StateKind
	SoldierNumber() int
	CurrentAssetName() string
}

type StateKind string

const (
	IdleState       StateKind = "IDLE"       // can infinitely loop through the same frame or change to another state
	MovingState     StateKind = "MOVING"     // requires character movement
	TransitionState StateKind = "TRANSITION" //soldier must change
)

// Concrete implementations of AssetState

type IdleStateImpl struct {
	stateKind        StateKind
	soldierNumber    int
	currentAssetName string
}

func (i *IdleStateImpl) StateKind() StateKind {
	return i.stateKind
}

func (i *IdleStateImpl) SoldierNumber() int {
	return i.soldierNumber
}

func (i *IdleStateImpl) CurrentAssetName() string {
	return i.currentAssetName
}

type MovingStateImpl struct {
	stateKind        StateKind
	soldierNumber    int
	currentAssetName string
}

func (m *MovingStateImpl) StateKind() StateKind {
	return m.stateKind
}

func (m *MovingStateImpl) SoldierNumber() int {
	return m.soldierNumber
}

func (m *MovingStateImpl) CurrentAssetName() string {
	return m.currentAssetName
}

type TransitionStateImpl struct {
	stateKind        StateKind
	soldierNumber    int
	currentAssetName string
}

func (t *TransitionStateImpl) StateKind() StateKind {
	return t.stateKind
}

func (t *TransitionStateImpl) SoldierNumber() int {
	return t.soldierNumber
}

func (t *TransitionStateImpl) CurrentAssetName() string {
	return t.currentAssetName
}

// NewStateFromKind creates the appropriate state implementation based on state kind string
func NewStateFromKind(stateKindStr string, soldierNumber int, currentAssetName string) AssetState {
	switch StateKind(stateKindStr) {
	case IdleState:
		return &IdleStateImpl{
			stateKind:        IdleState,
			soldierNumber:    soldierNumber,
			currentAssetName: currentAssetName,
		}
	case MovingState:
		return &MovingStateImpl{
			stateKind:        MovingState,
			soldierNumber:    soldierNumber,
			currentAssetName: currentAssetName,
		}
	case TransitionState:
		return &TransitionStateImpl{
			stateKind:        TransitionState,
			soldierNumber:    soldierNumber,
			currentAssetName: currentAssetName,
		}
	default:
		// Default to Idle if unknown
		return &IdleStateImpl{
			stateKind:        IdleState,
			soldierNumber:    soldierNumber,
			currentAssetName: currentAssetName,
		}
	}
}
