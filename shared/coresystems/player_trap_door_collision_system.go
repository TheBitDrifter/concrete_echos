package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

type PlayerTrapDoorCollisionSystem struct{}

func (s PlayerTrapDoorCollisionSystem) Run(scene blueprint.Scene, dt float64) error {
	tdQuery := warehouse.Factory.NewQuery().And(components.TrapDoorComponent)
	tdCursor := scene.NewCursor(tdQuery)
	playerQuery := warehouse.Factory.NewQuery().And(
		components.PlayerTag,
		warehouse.Factory.NewQuery().Not(components.SoftResetComponent),
	)
	playerCursor := scene.NewCursor(playerQuery)

	for range tdCursor.Next() {
		td := components.TrapDoorComponent.GetFromCursor(tdCursor)
		if td.Open {
			continue
		}
		for range playerCursor.Next() {
			err := s.resolve(scene, tdCursor, playerCursor)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (PlayerTrapDoorCollisionSystem) resolve(scene blueprint.Scene, tdCursor, playerCursor *warehouse.Cursor) error {
	playerPosition := spatial.Components.Position.GetFromCursor(playerCursor)
	playerShape := spatial.Components.Shape.GetFromCursor(playerCursor)
	playerDynamics := motion.Components.Dynamics.GetFromCursor(playerCursor)

	tdPos := spatial.Components.Position.GetFromCursor(tdCursor)
	tdShape := spatial.Components.Shape.GetFromCursor(tdCursor)
	tdDyn := motion.Components.Dynamics.GetFromCursor(tdCursor)

	if ok, collisionResult := spatial.Detector.Check(
		*playerShape, *tdShape, playerPosition.Two, tdPos.Two,
	); ok {
		motion.Resolver.Resolve(
			&playerPosition.Two,
			&tdPos.Two,
			playerDynamics,
			tdDyn,
			collisionResult,
		)
	}
	return nil
}
