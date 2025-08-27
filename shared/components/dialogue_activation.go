package components

import "github.com/TheBitDrifter/bappa/blueprint/dialogue"

type DialogueActivation struct {
	SlidesID        dialogue.SlidesEnum
	ConvoCallbackID dialogue.CallbackEnum

	Range       float64
	MinRange    float64
	MustBeRight bool
	MustBeLeft  bool
}
