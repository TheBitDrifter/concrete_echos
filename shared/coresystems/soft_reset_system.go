package coresystems

import (
	"log"
	"math"

	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

const SOFT_RESET_PLAYER_DURATION_TICKS = 180

type SoftResetSystem struct{}

func (s SoftResetSystem) Run(scene blueprint.Scene, dt float64) error {
	query := warehouse.Factory.NewQuery().And(
		components.SoftResetComponent,
	)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {

		dyn := motion.Components.Dynamics.GetFromCursor(cursor)
		dyn.Vel.X = 0
		dyn.Vel.Y = 0

		sr := components.SoftResetComponent.GetFromCursor(cursor)
		if scene.CurrentTick() >= sr.StartedTick+SOFT_RESET_PLAYER_DURATION_TICKS {

			playerEn, err := cursor.CurrentEntity()
			if err != nil {
				return err
			}
			playerEn.EnqueueRemoveComponent(components.SoftResetComponent)

			pos := spatial.Components.Position.GetFromCursor(cursor)

			softResetCheckpointQuery := warehouse.Factory.NewQuery().And(
				components.SoftResetCheckpointComponent,
			)
			softResetCheckpointCursor := scene.NewCursor(softResetCheckpointQuery)

			var closestCheckpointPos *spatial.Position
			minDistSq := math.MaxFloat64

			for range softResetCheckpointCursor.Next() {
				targetPos := spatial.Components.Position.GetFromCursor(softResetCheckpointCursor)
				currentDistSq := pos.Two.Sub(targetPos.Two).MagSquared()
				softResetCheckPointState := components.SoftResetCheckpointComponent.GetFromCursor(softResetCheckpointCursor)

				if currentDistSq < minDistSq && softResetCheckPointState.Activated {
					minDistSq = currentDistSq
					closestCheckpointPos = targetPos
				}
			}

			if closestCheckpointPos != nil {
				pos.X = closestCheckpointPos.X
				pos.Y = closestCheckpointPos.Y
				health := combat.Components.Health.GetFromCursor(cursor)
				health.Value -= 10
				playerEn.EnqueueAddComponentWithValue(combat.Components.Hurt, combat.Hurt{StartTick: scene.CurrentTick()})

				playerEn.EnqueueAddComponentWithValue(combat.Components.Invincible, combat.Invincible{StartTick: scene.CurrentTick()})
			} else {
				log.Println("Soft Reset Warning: No checkpoints found!")
			}
		}
	}
	return nil
}

type SoftResetBoundsSystem struct{}

func (s SoftResetBoundsSystem) Run(scene blueprint.Scene, dt float64) error {
	query := warehouse.Factory.NewQuery().And(components.PlayerTag)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {

		pos := spatial.Components.Position.GetFromCursor(cursor)
		cookedY := pos.Y < -500 || pos.Y > float64(scene.Height()+500)
		cookedX := pos.X < -500 || pos.X > float64(scene.Width()+500)

		playerEn, err := cursor.CurrentEntity()
		if err != nil {
			return err
		}

		if cookedX || cookedY && !components.SoftResetComponent.CheckCursor(cursor) && !combat.Components.Defeat.CheckCursor(cursor) {
			playerEn.EnqueueAddComponentWithValue(components.SoftResetComponent, components.SoftReset{StartedTick: scene.CurrentTick() - 150})
		}

	}
	return nil
}

type SoftResetCheckpointActivationSystem struct{}

func (s SoftResetCheckpointActivationSystem) Run(scene blueprint.Scene, dt float64) error {
	queryForSRC := warehouse.Factory.NewQuery().And(
		components.SoftResetCheckpointComponent,
	)
	outerCursor := scene.NewCursor(queryForSRC)

	queryForPlayer := warehouse.Factory.NewQuery().And(
		components.PlayerTag,
	)
	innerCursor := scene.NewCursor(queryForPlayer)
	for range outerCursor.Next() {
		srCheckpoint := components.SoftResetCheckpointComponent.GetFromCursor(outerCursor)
		if srCheckpoint.Activated {
			continue
		}

		for range innerCursor.Next() {
			playerShape := spatial.Components.Shape.GetFromCursor(innerCursor)
			checkpointShape := spatial.Components.Shape.GetFromCursor(outerCursor)

			playerPos := spatial.Components.Position.GetFromCursor(innerCursor)
			checkpointPos := spatial.Components.Position.GetFromCursor(outerCursor)

			if ok, _ := spatial.Detector.Check(
				*playerShape, *checkpointShape, playerPos.Two, checkpointPos.Two,
			); ok {
				srCheckpoint.Activated = true
			}
		}
	}

	return nil
}
