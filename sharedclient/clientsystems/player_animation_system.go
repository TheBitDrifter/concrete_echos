package clientsystems

import (
	"math"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/animations"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
	"github.com/TheBitDrifter/concrete_echos/shared/combatdata"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

type PlayerAnimationSystem struct{}

func (PlayerAnimationSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	// Handle Non Combat Anim
	query1 := warehouse.Factory.NewQuery().And(
		components.PlayerTag,
		warehouse.Factory.NewQuery().Not(combat.Components.Attack),
	)
	cursor := scene.NewCursor(query1)

	for range cursor.Next() {
		bundle := client.Components.SpriteBundle.GetFromCursor(cursor)
		spriteBlueprint := &bundle.Blueprints[0]

		if components.SoftResetComponent.CheckCursor(cursor) || combat.Components.Defeat.CheckCursor(cursor) {
			spriteBlueprint.TryAnimation(animations.Registry.Characters[characterkeys.BoxHead][animations.Defeat])
			continue
		}

		teleState := components.TeleportSwapComponent.GetFromCursor(cursor)
		if scene.CurrentTick() < teleState.StartTick+40 && teleState.StartTick != 0 {
			spriteBlueprint.TryAnimation(animations.Registry.Characters[characterkeys.BoxHead][animations.Teleport])
			continue
		}

		if combat.Components.Hurt.CheckCursor(cursor) {
			spriteBlueprint.TryAnimation(animations.Registry.Characters[characterkeys.BoxHead][animations.Hurt])
			continue
		}

		if components.PlayerIsExecutingComponent.CheckCursor(cursor) {
			spriteBlueprint.TryAnimation(animations.Registry.Characters[characterkeys.BoxHead][animations.Execute])
			continue
		}

		if components.DodgeComponent.CheckCursor(cursor) {
			spriteBlueprint.TryAnimation(animations.Registry.Characters[characterkeys.BoxHead][animations.Dodge])
			continue
		}

		dyn := motion.Components.Dynamics.GetFromCursor(cursor)
		onGround, grounded := components.OnGroundComponent.GetFromCursorSafe(cursor)
		if grounded {
			grounded = scene.CurrentTick()-onGround.LastTouch <= 3
		}

		if components.InConversationComponent.CheckCursor(cursor) {
			spriteBlueprint.TryAnimation(animations.Registry.Characters[characterkeys.BoxHead][animations.InConvo])
			continue
		}

		if components.IsSavingComponent.CheckCursor(cursor) {
			spriteBlueprint.TryAnimation(animations.Registry.Characters[characterkeys.BoxHead][animations.IsSaving])
			continue
		}

		// Player is moving horizontal and grounded (running)
		if math.Abs(dyn.Vel.X) > 20 && grounded {
			spriteBlueprint.TryAnimation(animations.Registry.Characters[characterkeys.BoxHead][animations.Run])

			// Player is moving down and not grounded (falling)
		} else if dyn.Vel.Y > 0 && !grounded {
			spriteBlueprint.TryAnimation(animations.Registry.Characters[characterkeys.BoxHead][animations.Fall])

			// Player is moving up and not grounded (jumping)
		} else if dyn.Vel.Y <= 0 && !grounded {
			spriteBlueprint.TryAnimation(animations.Registry.Characters[characterkeys.BoxHead][animations.Jump])

			// Default: player is idle
		} else {
			spriteBlueprint.TryAnimation(animations.Registry.Characters[characterkeys.BoxHead][animations.Idle])
		}
	}

	// Handle Combat Anim
	query2 := warehouse.Factory.NewQuery().And(combat.Components.Attack, client.Components.SpriteBundle, components.PlayerTag)
	cursor = scene.NewCursor(query2)

	for range cursor.Next() {
		bundle := client.Components.SpriteBundle.GetFromCursor(cursor)

		atk := combat.Components.Attack.GetFromCursor(cursor)
		spriteBlueprint := &bundle.Blueprints[0]
		spriteBlueprint.TryAnimation(combatdata.BoxHeadAttackAnimationMapping[atk.ID])
	}

	return nil
}
