package clientsystems

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

type VanishingPlatformAnimationSystem struct{}

func (sys VanishingPlatformAnimationSystem) Run(lc coldbrew.LocalClient, scene coldbrew.Scene) error {
	query := warehouse.Factory.NewQuery().And(components.VanishingPlatformComponent, client.Components.SpriteBundle)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {

		vp := components.VanishingPlatformComponent.GetFromCursor(cursor)
		sprBundle := client.Components.SpriteBundle.GetFromCursor(cursor)

		if vp.TimerStarted && !vp.Vanished {
			sprBundle.Blueprints[0].TryAnimationFromIndex(1)
			continue
		}

		if vp.TimerStarted && vp.Vanished {
			buffer := 30 // TODO: Make this a component field?

			vanishedDuration := scene.CurrentTick() - vp.TimerStartedTick
			if vp.RespawnDelay-vanishedDuration <= buffer {

				sprBundle.Blueprints[0].TryAnimationFromIndex(2)
				continue
			}
		}

	}
	return nil
}
