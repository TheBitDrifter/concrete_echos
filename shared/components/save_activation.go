package components

import "github.com/TheBitDrifter/concrete_echos/shared/persistence"

type SaveActivation struct {
	Range      float64
	OptionalID persistence.OptionalSaveID
}
