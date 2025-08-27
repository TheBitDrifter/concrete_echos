package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/callbacks"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

type TrapDoorSystem struct{}

func (s TrapDoorSystem) Run(scene blueprint.Scene, dt float64) error {
	query := warehouse.Factory.NewQuery().And(components.TrapDoorComponent)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		td := components.TrapDoorComponent.GetFromCursor(cursor)
		if td.Open {
			continue
		}
		if callbacks.CheckTrapDoor(scene.Storage(), float64(scene.Width()), float64(scene.Height()), *td) {
			td.Open = true
			td.LastChangedTick = scene.CurrentTick()
		}
	}
	return nil
}
