package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

type MobBoundsPlatformDetectionSystem struct{}

func (s MobBoundsPlatformDetectionSystem) Run(scene blueprint.Scene, dt float64) error {
	platformQuery := warehouse.Factory.NewQuery().Or(
		components.BlockTerrainTag,
		components.PlatformTag,
	)
	platformCursor := scene.NewCursor(platformQuery)

	mobQuery := warehouse.Factory.NewQuery().And(
		components.MobBoundsComponent,
		spatial.Components.Position,
		spatial.Components.Shape,
	)
	mobCursor := scene.NewCursor(mobQuery)

	for range mobCursor.Next() {

		bounds, hasBounds := components.MobBoundsComponent.GetFromCursorSafe(mobCursor)
		if bounds.MaxX != 0 || bounds.MinX != 0 {
			continue
		}
	MobLoop:
		for range platformCursor.Next() {
			mobPos := spatial.Components.Position.GetFromCursor(mobCursor)
			mobShape := spatial.Components.Shape.GetFromCursor(mobCursor)

			platformPos := spatial.Components.Position.GetFromCursor(platformCursor)
			platformShape := spatial.Components.Shape.GetFromCursor(platformCursor)

			checkPos := mobPos.Two.Add(vector.Two{Y: 1})

			if ok, _ := spatial.Detector.Check(*mobShape, *platformShape, checkPos, platformPos.Two); ok {
				if !hasBounds {
					continue MobLoop
				}

				platformWidth := platformShape.WorldAAB.Width
				leftX := platformPos.X - platformWidth/2
				rightX := platformPos.X + platformWidth/2

				bounds.MinX = leftX
				bounds.MaxX = rightX
				break MobLoop
			}
		}
	}
	return nil
}
