package clientsystems

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

type TrapDoorAnimationStateSystem struct{}

func (sys TrapDoorAnimationStateSystem) Run(lc coldbrew.LocalClient, scene coldbrew.Scene) error {
	query := warehouse.Factory.NewQuery().And(components.TrapDoorComponent)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		td := components.TrapDoorComponent.GetFromCursor(cursor)
		sprBundle := client.Components.SpriteBundle.GetFromCursor(cursor)

		if td.Open {
			sprBundle.Blueprints[0].TryAnimationFromIndex(1)
			continue
		}
		if !td.Open {
			sprBundle.Blueprints[0].TryAnimationFromIndex(0)
			continue
		}
	}
	return nil
}
