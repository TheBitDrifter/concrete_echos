package coresystems

import (
	"math"

	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/actions"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

type MobExecuteSystem struct{}

func (sys MobExecuteSystem) Run(scene blueprint.Scene, dt float64) error {
	defeatedMobsQuery := warehouse.Factory.NewQuery().And(
		combat.Components.Defeat,
		warehouse.Factory.NewQuery().Not(components.MobExecuteComponent),
	)
	defeatedMobsCursor := scene.NewCursor(defeatedMobsQuery)

	for range defeatedMobsCursor.Next() {

		mEn, err := defeatedMobsCursor.CurrentEntity()
		if err != nil {
			return err
		}
		mEn.EnqueueAddComponent(components.MobExecuteComponent)

	}

	playerQuery := warehouse.Factory.NewQuery().And(components.PlayerTag)
	playerCursor := scene.NewCursor(playerQuery)

	execMobsQuery := warehouse.Factory.NewQuery().And(
		components.MobExecuteComponent,
		warehouse.Factory.NewQuery().Not(components.IgnoreExecuteTag),
	)
	execMobsCursor := scene.NewCursor(execMobsQuery)
	for range execMobsCursor.Next() {
		mPos := spatial.Components.Position.GetFromCursor(execMobsCursor)
		mShape := spatial.Components.Shape.GetFromCursor(execMobsCursor)
		mBottom := mPos.Y + mShape.LocalAAB.Height/2
		for range playerCursor.Next() {
			playerPos := spatial.Components.Position.GetFromCursor(playerCursor)
			playerShape := spatial.Components.Shape.GetFromCursor(playerCursor)
			playerBottom := playerPos.Y + playerShape.LocalAAB.Height/2
			playerDir := spatial.Components.Direction.GetFromCursor(playerCursor)
			playerFacingRightWay := (playerDir.IsLeft() && playerPos.X > mPos.X) || (playerDir.IsRight() && playerPos.X < mPos.X)

			vertDistOK := math.Abs(playerBottom-mBottom) < 50
			currentDistSq := playerPos.Sub(mPos.Two).MagSquared()
			inRange := currentDistSq <= 45*45 && currentDistSq > 15*15 && (components.LastCombatComponent.GetFromCursor(playerCursor).StartTick+120 <= scene.CurrentTick())

			buffer := input.Components.ActionBuffer.GetFromCursor(playerCursor)

			if buffer.HasAction(actions.Interact) && inRange && vertDistOK && playerFacingRightWay {
				mExec := components.MobExecuteComponent.GetFromCursor(execMobsCursor)
				if mExec.StartTick == 0 {
					mExec.StartTick = scene.CurrentTick()
					playerEN, err := playerCursor.CurrentEntity()
					if err != nil {
						return err
					}

					playerEN.EnqueueAddComponentWithValue(components.PlayerIsExecutingComponent, components.PlayerIsExecuting{StartTick: scene.CurrentTick()})
				}
			}

		}
	}

	for range playerCursor.Next() {
		playerExec, ok := components.PlayerIsExecutingComponent.GetFromCursorSafe(playerCursor)
		if !ok {
			continue
		}
		if scene.CurrentTick() >= playerExec.StartTick+13*5 {
			playerEn, err := playerCursor.CurrentEntity()
			if err != nil {
				return err
			}
			health := combat.Components.Health.GetFromCursor(playerCursor)
			health.Value += 10
			if health.Value > 100 {
				health.Value = 100
			}
			counter := components.PlayerExecutionCountComponent.GetFromCursor(playerCursor)
			counter.Count += 1

			playerEn.EnqueueRemoveComponent(components.PlayerIsExecutingComponent)
		}
	}

	for range execMobsCursor.Next() {
		mobExec, ok := components.MobExecuteComponent.GetFromCursorSafe(execMobsCursor)
		if !ok {
			continue
		}
		if scene.CurrentTick() >= mobExec.StartTick+8*5 && mobExec.StartTick != 0 {
			mobEn, err := execMobsCursor.CurrentEntity()
			if err != nil {
				return err
			}
			scene.Storage().EnqueueDestroyEntities(mobEn)
		}
	}

	return nil
}
