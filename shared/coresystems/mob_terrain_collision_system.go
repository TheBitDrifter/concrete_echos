package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/scenes"
)

type MobTerrainCollisionSystem struct{}

func (s MobTerrainCollisionSystem) Run(scene blueprint.Scene, dt float64) error {
	blockTerrainQuery := warehouse.Factory.NewQuery().Or(components.BlockTerrainTag, components.PlatformTag)
	blockTerrainCursor := scene.NewCursor(blockTerrainQuery)

	enemyQuery := warehouse.Factory.NewQuery().And(scenes.DefaultMobComposition, warehouse.Factory.NewQuery().Not(input.Components.ActionBuffer))
	enemyCursor := scene.NewCursor(enemyQuery)

	for range blockTerrainCursor.Next() {
		for range enemyCursor.Next() {
			// Delegate to helper
			err := s.resolve(scene, blockTerrainCursor, enemyCursor)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (MobTerrainCollisionSystem) resolve(scene blueprint.Scene, blockCursor, enemyCursor *warehouse.Cursor) error {
	enemyPos := spatial.Components.Position.GetFromCursor(enemyCursor)
	enemyShape := spatial.Components.Shape.GetFromCursor(enemyCursor)
	enemyDyn := motion.Components.Dynamics.GetFromCursor(enemyCursor)

	blockPosition := spatial.Components.Position.GetFromCursor(blockCursor)
	blockShape := spatial.Components.Shape.GetFromCursor(blockCursor)
	blockDynamics := motion.Components.Dynamics.GetFromCursor(blockCursor)

	// Check for a collision
	ignore := components.IgnoreTerrainCollisionsMob.CheckCursor(enemyCursor)
	if ok, collisionResult := spatial.Detector.Check(
		*enemyShape, *blockShape, enemyPos.Two, blockPosition.Two,
	); ok && !ignore {
		motion.Resolver.Resolve(
			&enemyPos.Two,
			&blockPosition.Two,
			enemyDyn,
			blockDynamics,
			collisionResult,
		)
	}
	return nil
}
