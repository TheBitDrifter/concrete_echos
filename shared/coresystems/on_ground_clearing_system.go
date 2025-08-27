package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

type OnGroundClearingSystem struct{}

func (OnGroundClearingSystem) Run(scene blueprint.Scene, dt float64) error {
	onGroundQuery := warehouse.Factory.NewQuery().And(components.OnGroundComponent)
	onGroundCursor := scene.NewCursor(onGroundQuery)

	for range onGroundCursor.Next() {
		onGround := components.OnGroundComponent.GetFromCursor(onGroundCursor)
		if scene.CurrentTick() > onGround.LastTouch {
			groundedEntity, _ := onGroundCursor.CurrentEntity()
			err := groundedEntity.EnqueueRemoveComponent(components.OnGroundComponent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
