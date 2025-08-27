package clientsystems

import (
	"math"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/animations"
	"github.com/TheBitDrifter/concrete_echos/shared/combatdata"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/scenes"
)

type MobAnimationSystem struct{}

func (s MobAnimationSystem) Run(lc coldbrew.LocalClient, scene coldbrew.Scene) error {
	query := warehouse.Factory.NewQuery().And(
		client.Components.SpriteBundle,
		components.CharacterKeyComponent,
		components.MobTag,
		warehouse.Factory.NewQuery().Not(components.PlayerTag),
	)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		bundle := client.Components.SpriteBundle.GetFromCursor(cursor)
		key := components.CharacterKeyComponent.GetFromCursor(cursor)

		dodging := components.DodgeComponent.CheckCursor(cursor)

		atk, attacking := combat.Components.Attack.GetFromCursorSafe(cursor)
		defeated := combat.Components.Defeat.CheckCursor(cursor)

		if defeated {
			bundle.Blueprints[0].TryAnimation(animations.Registry.Characters[*key][animations.Defeat])
			continue
		}

		if dodging {
			bundle.Blueprints[0].TryAnimation(animations.Registry.Characters[*key][animations.Dodge])
			continue
		}

		if attacking {
			if components.MobLazySkullyComponent.CheckCursor(cursor) {
				s.handleLazySkully(atk, bundle)
				continue
			}
			_, hasPrimaryAtk := combatdata.PrimarySeqs[*key]
			if hasPrimaryAtk && (atk.ID == combatdata.PrimarySeqs[*key].First().ID) {
				bundle.Blueprints[0].TryAnimation(animations.Registry.Characters[*key][animations.PrimaryAttack])
			}
			_, hasSecondaryAtk := combatdata.SecondarySeqs[*key]
			if hasSecondaryAtk && (atk.ID == combatdata.SecondarySeqs[*key].First().ID) {
				bundle.Blueprints[0].TryAnimation(animations.Registry.Characters[*key][animations.SecondaryAttack])
			}
			continue
		}

		throwerMob, isThrower := components.MobSimpleThrowerComponent.GetFromCursorSafe(cursor)
		throwing := isThrower && scene.CurrentTick() < throwerMob.LastStartedThrow+throwerMob.ThrowDuration

		if throwing {
			bundle.Blueprints[0].TryAnimation(animations.Registry.Characters[*key][animations.PrimaryRanged])
			continue
		}

		shooterMob, isThrower := components.MobSimpleHoriShooterComponent.GetFromCursorSafe(cursor)
		shooting := isThrower && scene.CurrentTick() < shooterMob.LastStartedShot+shooterMob.ShootDuration

		swapStunned := components.SwapVulnerableComponent.CheckCursor(cursor)

		if shooting && !swapStunned {
			bundle.Blueprints[0].TryAnimation(animations.Registry.Characters[*key][animations.PrimaryRanged])
			continue
		}

		hurt := combat.Components.Hurt.CheckCursor(cursor)
		_, hasHurtAnim := animations.Registry.Characters[*key][animations.Hurt]
		if hurt && hasHurtAnim {
			bundle.Blueprints[0].TryAnimation(animations.Registry.Characters[*key][animations.Hurt])
			continue
		}

		dyn, dynOK := motion.Components.Dynamics.GetFromCursorSafe(cursor)
		mxm, mxmOK := components.MobXMobCollisionComponent.GetFromCursorSafe(cursor)
		_, hasRunAnim := animations.Registry.Characters[*key][animations.Run]

		if mxmOK && dynOK && hasRunAnim && mxm.ShouldStop && false && !swapStunned {
			bundle.Blueprints[0].TryAnimation(animations.Registry.Characters[*key][animations.Run])
			continue
		}

		if (dynOK && hasRunAnim && math.Abs(dyn.Vel.X) > 20) && !swapStunned {
			bundle.Blueprints[0].TryAnimation(animations.Registry.Characters[*key][animations.Run])
			continue
		}

		if inConvo, inConvoOK := components.InConversationComponent.GetFromCursorSafe(cursor); inConvoOK {
			convoStartAnim, ok := animations.Registry.Characters[*key][animations.ConvoStart]
			if ok {

				startUpDurationTicks := convoStartAnim.FrameCount * convoStartAnim.Speed

				if scene.CurrentTick() <= inConvo.StartedTick+startUpDurationTicks {
					bundle.Blueprints[0].TryAnimation(convoStartAnim)
					continue
				} else {
					if inConvoAnim, ok := animations.Registry.Characters[*key][animations.InConvo]; ok {
						bundle.Blueprints[0].TryAnimation(inConvoAnim)
						continue
					}
				}
			}

			// For if they don't have a startup, but do have a inConvo
			if inConvoAnim, ok := animations.Registry.Characters[*key][animations.InConvo]; ok {
				bundle.Blueprints[0].TryAnimation(inConvoAnim)
				continue
			}

		}

		if _, isLR := components.MobLazySkullyComponent.GetFromCursorSafe(cursor); isLR {
			agro := components.FriendlyAgroComponent.GetFromCursor(cursor)
			if agro.StartTick+scenes.LazyRoninWakeUpTicks > scene.CurrentTick() && agro.StartTick != 0 {
				bundle.Blueprints[0].TryAnimation(animations.LazySkullRoninAnims.Animations[8])
				continue
			}

			if agro.StartTick != 0 {
				bundle.Blueprints[0].TryAnimation(animations.LazySkullRoninAnims.Animations[4])
				continue
			}
		}
		bundle.Blueprints[0].TryAnimation(animations.Registry.Characters[*key][animations.Idle])
	}
	return nil
}

func (s MobAnimationSystem) handleLazySkully(a *combat.Attack, bundle *client.SpriteBundle) {
	if a.ID == combatdata.LazySkullyAttacks[combatdata.LazySkullyAttackKeys.PrimaryCombo].ID {
		bundle.Blueprints[0].TryAnimation(animations.LazySkullRoninAnims.Animations[1])
	}
	if a.ID == combatdata.LazySkullyAttacks[combatdata.LazySkullyAttackKeys.SecondaryCombo].ID {
		bundle.Blueprints[0].TryAnimation(animations.LazySkullRoninAnims.Animations[2])
	}
	if a.ID == combatdata.LazySkullyAttacks[combatdata.LazySkullyAttackKeys.TertCombo].ID {
		bundle.Blueprints[0].TryAnimation(animations.LazySkullRoninAnims.Animations[3])
	}
}
