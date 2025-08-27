package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

type PlayerObstacleCollisionSystem struct{}

func (s PlayerObstacleCollisionSystem) Run(scene blueprint.Scene, dt float64) error {
	obsQuery := warehouse.Factory.NewQuery().And(components.ObstacleComponent)
	obsCursor := scene.NewCursor(obsQuery)
	playerQuery := warehouse.Factory.NewQuery().And(
		components.PlayerTag,
		warehouse.Factory.NewQuery().Not(components.SoftResetComponent),
	)
	playerCursor := scene.NewCursor(playerQuery)

	for range obsCursor.Next() {
		for range playerCursor.Next() {
			err := s.resolve(scene, obsCursor, playerCursor)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (PlayerObstacleCollisionSystem) resolve(scene blueprint.Scene, obsCursor, playerCursor *warehouse.Cursor) error {
	playerPosition := spatial.Components.Position.GetFromCursor(playerCursor)
	playerShape := spatial.Components.Shape.GetFromCursor(playerCursor)
	playerDynamics := motion.Components.Dynamics.GetFromCursor(playerCursor)

	obsPos := spatial.Components.Position.GetFromCursor(obsCursor)
	obsShape := spatial.Components.Shape.GetFromCursor(obsCursor)
	obsDyn := motion.Components.Dynamics.GetFromCursor(obsCursor)

	if ok, collisionResult := spatial.Detector.Check(
		*playerShape, *obsShape, playerPosition.Two, obsPos.Two,
	); ok {
		motion.Resolver.Resolve(
			&playerPosition.Two,
			&obsPos.Two,
			playerDynamics,
			obsDyn,
			collisionResult,
		)

		lilPushy := collisionResult.Normal.Scale(-130)
		playerDynamics.Vel = playerDynamics.Vel.Add(lilPushy)

		if components.DodgeComponent.CheckCursor(playerCursor) {
			playerEn, _ := playerCursor.CurrentEntity()
			playerEn.EnqueueRemoveComponent(components.DodgeComponent)

		}

		if !components.SoftResetComponent.CheckCursor(playerCursor) {
			playerEn, _ := playerCursor.CurrentEntity()
			playerEn.EnqueueAddComponentWithValue(components.SoftResetComponent, components.SoftReset{StartedTick: scene.CurrentTick()})
		}

	}
	return nil
}
