package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

// MobHorizontalPacerSystem now only handles the movement logic, using pre-calculated bounds.
type MobPacerSystem struct{}

func (s MobPacerSystem) Run(scene blueprint.Scene, dt float64) error {
	mobQuery := warehouse.Factory.NewQuery().And(
		components.PacerComponent,
		components.MobBoundsComponent,
		components.MobTag,
		warehouse.Factory.NewQuery().Not(
			combat.Components.Defeat,
		),
	)
	mobCursor := scene.NewCursor(mobQuery)

	for range mobCursor.Next() {
		pacer := components.PacerComponent.GetFromCursor(mobCursor)
		bounds := components.MobBoundsComponent.GetFromCursor(mobCursor)
		mobPos := spatial.Components.Position.GetFromCursor(mobCursor)
		mobShape := spatial.Components.Shape.GetFromCursor(mobCursor)
		mobDyn := motion.Components.Dynamics.GetFromCursor(mobCursor)

		speed := pacer.Speed
		if combat.Components.Hurt.CheckCursor(mobCursor) {
			speed /= 4
		}

		if contactData, ok := components.ContactComponent.GetFromCursorSafe(mobCursor); ok {
			if contactData.LastHit == scene.CurrentTick()-1 && pacer.SwapDirOnHit {
				pacer.IsLeft = !pacer.IsLeft
			}
		}

		if pacer.IsVert {
			if bounds.MinY == 0 && bounds.MaxY == 0 {
				continue
			}

			mobTopY := mobPos.Y - mobShape.WorldAAB.Height/2
			mobBottomY := mobPos.Y + mobShape.WorldAAB.Height/2

			if mobTopY <= bounds.MinY {
				pacer.IsLeft = false
			}
			if mobBottomY >= bounds.MaxY {
				pacer.IsLeft = true
			}

			mobDyn.Vel.X = 0
			if pacer.IsLeft {
				mobDyn.Vel.Y = -speed
			} else {
				mobDyn.Vel.Y = speed
			}

		} else {
			if bounds.MinX == 0 && bounds.MaxX == 0 {
				continue
			}

			mobDirection := spatial.Components.Direction.GetFromCursor(mobCursor)
			mobLeftX := mobPos.X - mobShape.WorldAAB.Width/2
			mobRightX := mobPos.X + mobShape.WorldAAB.Width/2

			if mobLeftX <= bounds.MinX {
				pacer.IsLeft = false
			}
			if mobRightX >= bounds.MaxX {
				pacer.IsLeft = true
			}

			mobDyn.Vel.Y = 0
			if pacer.IsLeft {
				mobDyn.Vel.X = -speed
				mobDirection.SetLeft()
			} else {
				mobDyn.Vel.X = speed
				mobDirection.SetRight()
			}
		}
	}
	return nil
}
