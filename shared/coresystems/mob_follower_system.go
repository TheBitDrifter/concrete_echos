package coresystems

import (
	"math"

	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

type MobFollowerSystem struct{}

func (s MobFollowerSystem) Run(scene blueprint.Scene, dt float64) error {
	mobQuery := warehouse.Factory.NewQuery().And(
		components.MobFollowerComponent,
		components.MobBoundsComponent,
		warehouse.Factory.NewQuery().Not(combat.Components.Defeat),
	)
	playerQuery := warehouse.Factory.NewQuery().And(components.PlayerTag)

	mobCursor := scene.NewCursor(mobQuery)
	playerCursor := scene.NewCursor(playerQuery)

	for range playerCursor.Next() {
		playerPos := spatial.Components.Position.GetFromCursor(playerCursor)

		for range mobCursor.Next() {

			mobPos := spatial.Components.Position.GetFromCursor(mobCursor)
			// TODO: MAKE 150 and 5 not magic values but rather part of follower comp
			if math.Abs(playerPos.Y-mobPos.Y) > 90 {
				continue
			}
			if math.Abs(playerPos.X-mobPos.X) < 5 {
				continue
			}
			mobHurt := combat.Components.Hurt.CheckCursor(mobCursor)
			mobAttacking := combat.Components.Attack.CheckCursor(mobCursor)
			mobDodging := components.DodgeComponent.CheckCursor(mobCursor)

			if mobAttacking || mobDodging || mobHurt {
				continue
			}

			follower := components.MobFollowerComponent.GetFromCursor(mobCursor)
			bounds := components.MobBoundsComponent.GetFromCursor(mobCursor)
			mobDir := spatial.Components.Direction.GetFromCursor(mobCursor)
			mobDyn := motion.Components.Dynamics.GetFromCursor(mobCursor)

			mobDyn.Vel.X = 0

			distanceSquared := playerPos.Sub(mobPos.Two).MagSquared()
			inBounds := playerPos.X > bounds.MinX && playerPos.X < bounds.MaxX
			// TODO: MAKE INBOUNDS OPTIONAl
			// Make skelly use mob follower and stop short!
			isInVisionRange := distanceSquared <= (follower.VisionRadius*follower.VisionRadius) || inBounds && false

			if !isInVisionRange {
				continue
			}

			isTooClose := distanceSquared <= (follower.StopRadius * follower.StopRadius)
			if isTooClose {
				continue
			}

			if playerPos.X > mobPos.X {
				mobDir.SetRight()
			} else {
				mobDir.SetLeft()
			}

			wantsToMoveRight := mobDir.IsRight()
			atRightBoundary := mobPos.X >= bounds.MaxX
			if wantsToMoveRight && atRightBoundary {
				continue
			}

			wantsToMoveLeft := mobDir.IsLeft()
			atLeftBoundary := mobPos.X <= bounds.MinX
			if wantsToMoveLeft && atLeftBoundary {
				continue
			}

			mobDyn.Vel.X = follower.Speed * mobDir.AsFloat()
		}
	}
	return nil
}
