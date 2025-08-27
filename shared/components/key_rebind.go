package components

import "github.com/TheBitDrifter/bappa/blueprint/input"

type KeyRebindingState struct {
	Active        bool
	CurrentAction input.Action
	ActionsToBind []input.Action
}
