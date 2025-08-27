package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

// IgnorePlatformClearingSystem clears out individual platform ignores that have expired

// ClearDroppingStateSystem ends the player's "drop through" state once they
// are no longer colliding with any platforms.
type ClearDroppingStateSystem struct{}

func (ClearDroppingStateSystem) Run(scene blueprint.Scene, dt float64) error {
	query := warehouse.Factory.NewQuery().And(components.IsDroppingThroughPlatformTag)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		playerEntity, _ := cursor.CurrentEntity()

		if _, stillColliding := components.StillCollidingWithPlatformTag.GetFromCursorSafe(cursor); stillColliding {
			playerEntity.EnqueueRemoveComponent(components.StillCollidingWithPlatformTag)
		} else {
			playerEntity.EnqueueRemoveComponent(components.IsDroppingThroughPlatformTag)
		}
	}
	return nil
}
