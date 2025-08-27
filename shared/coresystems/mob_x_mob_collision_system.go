package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

// 1. Marker system using proximity to decide who stops.
type MobCollisionMarkerSystem struct{}

func (s MobCollisionMarkerSystem) Run(scene blueprint.Scene, dt float64) error {
	type mobData struct {
		pos          *spatial.Position
		shape        *spatial.Shape
		collision    *components.MobXMobCollision
		characterKey *characterkeys.CharEnum
	}

	playerQuery := warehouse.Factory.NewQuery().And(components.PlayerTag, spatial.Components.Position)
	playerCursor := scene.NewCursor(playerQuery)

	var targetPos spatial.Position

	for range playerCursor.Next() {
		pos, hasTarget := spatial.Components.Position.GetFromCursorSafe(playerCursor)
		if !hasTarget {
			return nil // No target, no priority logic to apply
		}
		targetPos = *pos
	}

	query := warehouse.Factory.NewQuery().And(
		components.MobXMobCollisionComponent,
		components.CharacterKeyComponent,
		spatial.Components.Shape,
		spatial.Components.Position,
		warehouse.Factory.NewQuery().Not(combat.Components.Defeat, components.DodgeComponent),
	)

	// Clear state from the previous frame.
	clearCursor := scene.NewCursor(query)
	for range clearCursor.Next() {
		state := components.MobXMobCollisionComponent.GetFromCursor(clearCursor)
		state.ShouldStop = false
	}

	var mobs []mobData
	collectCursor := scene.NewCursor(query)
	for range collectCursor.Next() {
		mobs = append(mobs, mobData{
			pos:          spatial.Components.Position.GetFromCursor(collectCursor),
			shape:        spatial.Components.Shape.GetFromCursor(collectCursor),
			collision:    components.MobXMobCollisionComponent.GetFromCursor(collectCursor),
			characterKey: components.CharacterKeyComponent.GetFromCursor(collectCursor),
		})
	}

	for i := 0; i < len(mobs); i++ {
		for j := i + 1; j < len(mobs); j++ {
			mobA := mobs[i]
			mobB := mobs[j]

			if *mobA.characterKey != *mobB.characterKey {
				continue
			}

			colliding, _ := spatial.Detector.Check(*mobA.shape, *mobB.shape, mobA.pos, mobB.pos)

			if colliding {
				distA := DistSqr(&mobA.pos.Two, &targetPos.Two)
				distB := DistSqr(&mobB.pos.Two, &targetPos.Two)

				if distA > distB {
					mobA.collision.ShouldStop = true
				} else {
					mobB.collision.ShouldStop = true
				}
			}
		}
	}

	return nil
}

// 2. Handler system is now very simple.
type MobCollisionHandlerSystem struct{}

func (s MobCollisionHandlerSystem) Run(scene blueprint.Scene, dt float64) error {
	query := warehouse.Factory.NewQuery().And(
		components.MobXMobCollisionComponent,
		motion.Components.Dynamics,
		warehouse.Factory.NewQuery().Not(combat.Components.Defeat, components.DodgeComponent),
	)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		mxm := components.MobXMobCollisionComponent.GetFromCursor(cursor)
		if mxm.ShouldStop {
			motion.Components.Dynamics.GetFromCursor(cursor).Vel.X = 0
			// You may also want to stop vertical movement
			// motion.Components.Dynamics.GetFromCursor(cursor).Vel.Y = 0
		}
	}

	return nil
}

func DistSqr(a, b *vector.Two) float64 {
	dx := b.X - a.X
	dy := b.Y - a.Y
	return dx*dx + dy*dy
}
