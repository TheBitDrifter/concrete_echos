package components

import (
	"github.com/TheBitDrifter/bappa/warehouse"
)

type TeleportSwap struct {
	StartTick    int
	Duration     int
	ActiveTarget warehouse.Entity
	HasTarget    bool
}

func (js TeleportSwap) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}
