package coresystems

import (
	"math"
	"math/rand"

	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/combatdata"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/scenes"
)

const (
	StateChasing = iota
	StateFleeing
)

const (
	attackRangeSquared         = 14000.0
	fleeRangeSquared           = 80.0 * 80.0
	mobSpeed                   = 80.0
	stateChangeCooldownInTicks = 20
	fleeEndRangeSquared        = (80.0 + 20.0) * (80.0 + 20.0)
	noAttackChance             = 0.25
)

type MobLazySkullySystem struct{}

func (s MobLazySkullySystem) Run(scene blueprint.Scene, dt float64) error {
	mobQuery := warehouse.Factory.NewQuery().And(
		components.MobLazySkullyComponent,
		warehouse.Factory.NewQuery().Not(
			combat.Components.Defeat,
			components.SwapVulnerableComponent,
		),
	)
	playerQuery := warehouse.Factory.NewQuery().And(components.PlayerTag)

	mobCursor := scene.NewCursor(mobQuery)
	playerCursor := scene.NewCursor(playerQuery)

	for range playerCursor.Next() {
		playerPos := spatial.Components.Position.GetFromCursor(playerCursor)
		playerOnGround := components.OnGroundComponent.CheckCursor(playerCursor)
		execCount := components.PlayerExecutionCountComponent.GetFromCursor(playerCursor).Count

		for range mobCursor.Next() {
			s.processMob(scene, mobCursor, playerPos, scene.CurrentTick(), playerOnGround, execCount)
			mobPos := spatial.Components.Position.GetFromCursor(mobCursor)
			bounds := components.MobBoundsComponent.GetFromCursor(mobCursor)
			if mobPos.X < bounds.MinX {
				mobPos.X = bounds.MinX
			}
			if mobPos.X > bounds.MaxX {
				mobPos.X = bounds.MaxX
			}
		}
	}

	return nil
}

func (s MobLazySkullySystem) processMob(
	scene blueprint.Scene,
	mobCursor *warehouse.Cursor,
	playerPos *spatial.Position,
	currentTick int,
	playerOnGround bool,
	playerExecCount int,
) {
	mobEn, err := mobCursor.CurrentEntity()
	if err != nil {
		return
	}

	mobPos := spatial.Components.Position.GetFromCursor(mobCursor)
	mobDyn := motion.Components.Dynamics.GetFromCursor(mobCursor)
	lazySkully := components.MobLazySkullyComponent.GetFromCursor(mobCursor)
	currAttack, isAttacking := combat.Components.Attack.GetFromCursorSafe(mobCursor)
	mobHealth := combat.Components.Health.GetFromCursor(mobCursor)
	agro := components.FriendlyAgroComponent.GetFromCursor(mobCursor)
	mobBounds := components.MobBoundsComponent.GetFromCursor(mobCursor)

	mobIsFullHealth := mobHealth.Value == scenes.LazyRoninMaxHealth
	playerIsEvil := (playerExecCount >= scenes.PLAYER_EVIL_THRESHOLD)

	if agro.StartTick == 0 && mobIsFullHealth && !playerIsEvil {
		return
	}
	inRange := playerPos.X >= mobBounds.MinX+30
	if !inRange {
		return
	}

	if agro.StartTick == 0 && (!mobIsFullHealth || playerIsEvil) {
		agro.StartTick = currentTick
		en, err := mobCursor.CurrentEntity()
		if err != nil {
			return
		}
		en.EnqueueRemoveComponent(components.DialogueManualStepperComponent)
		en.EnqueueRemoveComponent(components.DialogueActivationComponent)
		return
	}

	// Wait for delay
	if agro.StartTick+scenes.LazyRoninWakeUpTicks > currentTick {
		return
	}

	// Reset health as awoken
	if agro.StartTick+scenes.LazyRoninWakeUpTicks == currentTick {
		mobHealth.Value = scenes.LazyRoninMaxHealth
	}

	if isAttacking {
		elapsedTicks := currentTick - currAttack.StartTick
		idx := elapsedTicks / currAttack.Speed

		if currAttack.ID == combatdata.LazySkullyAttacks[combatdata.LazySkullyAttackKeys.PrimaryCombo].ID {
			if idx >= 10 && idx <= 12 && playerOnGround {
				s.updateMobDirection(mobCursor, mobPos, playerPos)
			}
		}

		return
	}

	s.updateMobDirection(mobCursor, mobPos, playerPos)

	distanceSquared := playerPos.Sub(mobPos.Two).MagSquared()

	s.updateMobMovement(mobDyn, lazySkully, mobPos, playerPos, distanceSquared, currentTick)

	isInAttackRange := distanceSquared < attackRangeSquared
	if isInAttackRange {
		s.initiateAttack(scene, mobEn, lazySkully, scene.CurrentTick(), mobPos, playerPos)
	}
}

func (s MobLazySkullySystem) updateMobMovement(mobDyn *motion.Dynamics, lazySkully *components.MobLazySkully, mobPos *spatial.Position, playerPos *spatial.Position, distanceSquared float64, currentTick int) {
	if currentTick < lazySkully.LastStateChangeTick+stateChangeCooldownInTicks {
		return
	}

	directionToPlayer := 1.0
	if playerPos.X < mobPos.X {
		directionToPlayer = -1.0
	}

	switch lazySkully.MovementState {
	case StateChasing:
		if distanceSquared <= fleeRangeSquared {
			lazySkully.MovementState = StateFleeing
			lazySkully.LastStateChangeTick = currentTick
			mobDyn.Vel.X = directionToPlayer * mobSpeed
		} else {
			mobDyn.Vel.X = directionToPlayer * mobSpeed
		}

	case StateFleeing:
		if distanceSquared > fleeEndRangeSquared {
			lazySkully.MovementState = StateChasing
			lazySkully.LastStateChangeTick = currentTick
			mobDyn.Vel.X = directionToPlayer * mobSpeed
		} else {
			mobDyn.Vel.X = -directionToPlayer * mobSpeed
		}
	}
}

func (s MobLazySkullySystem) updateMobDirection(mobCursor *warehouse.Cursor, mobPos, playerPos *spatial.Position) {
	mobDirection := spatial.Components.Direction.GetFromCursor(mobCursor)
	if playerPos.X > mobPos.X {
		mobDirection.SetRight()
	} else {
		mobDirection.SetLeft()
	}
}

// initiateAttack determines if and what attack the mob should perform.
func (s MobLazySkullySystem) initiateAttack(scene blueprint.Scene, mobEn warehouse.Entity, lazySkully *components.MobLazySkully, currentTick int, mobPos, playerPos *spatial.Position) {
	// Cooldown after a successful attack.
	const attackCooldown = 120
	// Cooldown after choosing "no attack".
	const hesitationCooldown = 60

	if currentTick < lazySkully.NextAttackTick {
		return
	}

	primAttack := combatdata.LazySkullyAttacks[combatdata.LazySkullyAttackKeys.PrimaryCombo]
	secAttack := combatdata.LazySkullyAttacks[combatdata.LazySkullyAttackKeys.SecondaryCombo]
	terAttack := combatdata.LazySkullyAttacks[combatdata.LazySkullyAttackKeys.TertCombo]

	isPlayerJumpingOver := playerPos.Y < mobPos.Y-30 && math.Abs(playerPos.X-mobPos.X) < 50
	if isPlayerJumpingOver {
		s.executeAttack(mobEn, lazySkully, secAttack, currentTick, attackCooldown)
		return
	}

	roll := rand.Float64()

	if roll < 0.40 {
		// 40% chance: Perform Combo 1
		s.executeAttack(mobEn, lazySkully, primAttack, currentTick, attackCooldown)
		// 30% chance: Combo 2
	} else if roll < 0.70 {
		s.executeAttack(mobEn, lazySkully, secAttack, currentTick, attackCooldown)
		// 15% chance: combo 3
	} else if roll < 0.85 {
		s.executeAttack(mobEn, lazySkully, terAttack, currentTick, attackCooldown)
		// 15% chance: hesitate
	} else {
		lazySkully.NextAttackTick = currentTick + hesitationCooldown
	}
}

// executeAttack helper remains the same.
func (s MobLazySkullySystem) executeAttack(mobEn warehouse.Entity, lazySkully *components.MobLazySkully, attack combat.Attack, currentTick int, cooldown int) {
	attack.StartTick = currentTick
	mobEn.EnqueueAddComponentWithValue(combat.Components.Attack, attack)
	lazySkully.NextAttackTick = currentTick + cooldown
}
