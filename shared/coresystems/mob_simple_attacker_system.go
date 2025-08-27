package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/combatdata"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

type MobSimpleAttackerSystem struct{}

func (s MobSimpleAttackerSystem) Run(scene blueprint.Scene, dt float64) error {
	mobQuery := warehouse.Factory.NewQuery().And(
		components.MobTag,
		components.MobSimpleAttackerComponent,
		warehouse.Factory.NewQuery().Not(
			components.PlayerTag,
			combat.Components.Defeat,
			components.SwapVulnerableComponent,
		),
	)
	mobCursor := scene.NewCursor(mobQuery)
	playerQuery := warehouse.Factory.NewQuery().And(components.PlayerTag)
	playerCursor := scene.NewCursor(playerQuery)

	for range mobCursor.Next() {
		for range playerCursor.Next() {

			playerPos := spatial.Components.Position.GetFromCursor(playerCursor)
			mobPos := spatial.Components.Position.GetFromCursor(mobCursor)
			simpleAttacker := components.MobSimpleAttackerComponent.GetFromCursor(mobCursor)

			mobAttackVision := simpleAttacker.AttackVisionRadius
			mobAttackVisionSq := mobAttackVision * mobAttackVision

			distanceSquared := playerPos.Sub(mobPos.Two).MagSquared()

			atkInRange := distanceSquared <= mobAttackVisionSq

			isAttacking := combat.Components.Attack.CheckCursor(mobCursor)

			mxm, shouldWaitForClearance := components.MobXMobCollisionComponent.GetFromCursorSafe(mobCursor)
			shouldWaitForClearance = shouldWaitForClearance && mxm.ShouldStop && false

			isReadyToAttack := (scene.CurrentTick()-simpleAttacker.LastStarted) > simpleAttacker.Delay && !isAttacking && !shouldWaitForClearance

			if atkInRange && isReadyToAttack {
				simpleAttacker.LastStarted = scene.CurrentTick()

				mobEn, err := mobCursor.CurrentEntity()
				if err != nil {
					return err
				}

				mobKey := components.CharacterKeyComponent.GetFromCursor(mobCursor)
				attack := combatdata.PrimarySeqs[*mobKey].First()
				attack.StartTick = scene.CurrentTick()
				mobEn.EnqueueAddComponentWithValue(combat.Components.Attack, attack)

			}

		}
	}

	return nil
}
