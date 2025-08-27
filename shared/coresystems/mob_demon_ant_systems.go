package coresystems

import (
	"math/rand"

	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/combatdata"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

type MobDemonAntMovementSystem struct{}

func (s MobDemonAntMovementSystem) Run(scene blueprint.Scene, dt float64) error {
	mobQuery := warehouse.Factory.NewQuery().And(
		components.MobDemonAntComponent,
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

		for range mobCursor.Next() {
			demonAnt := components.MobDemonAntComponent.GetFromCursor(mobCursor)
			mobPos := spatial.Components.Position.GetFromCursor(mobCursor)
			mobDyn := motion.Components.Dynamics.GetFromCursor(mobCursor)
			mobDir := spatial.Components.Direction.GetFromCursor(mobCursor)

			dodgeActiveUntil := demonAnt.LastDodged + demonAnt.DodgeDuration

			isDodging := scene.CurrentTick() < dodgeActiveUntil
			isAttacking := combat.Components.Attack.CheckCursor(mobCursor)
			mobBounds := components.MobBoundsComponent.GetFromCursor(mobCursor)

			if mobPos.X > playerPos.X && !isAttacking && !isDodging {
				mobDir.SetLeft()
			}
			if mobPos.X <= playerPos.X && !isAttacking && !isDodging {
				mobDir.SetRight()
			}

			if isDodging {
				mobDyn.Vel.X = demonAnt.DodgeSpeed * demonAnt.DodgeDirection.AsFloat()

				if mobPos.X < mobBounds.MinX && mobDyn.Vel.X < 0 {
					mobDyn.Vel.X = 0
				}
				if mobPos.X > mobBounds.MaxX && mobDyn.Vel.X > 0 {
					mobDyn.Vel.X = 0
				}

				if demonAnt.DodgeDirection.AsFloat() < 0 {
					mobDir.SetRight()
				} else {
					mobDir.SetLeft()
				}
				continue
			}

			mobDyn.Vel.X = 0

			en, err := mobCursor.CurrentEntity()
			if err != nil {
				return err
			}
			if components.DodgeComponent.Check(en.Table()) {
				en.EnqueueRemoveComponent(components.DodgeComponent)
			}

			distanceSquared := playerPos.Sub(mobPos.Two).MagSquared()

			isInRangeDodge := distanceSquared <= (demonAnt.DodgeVisionRadius * demonAnt.DodgeVisionRadius)
			isInRangeJabAtk := distanceSquared <= (demonAnt.AttackVisionJabRadius*demonAnt.AttackVisionJabRadius) && !isInRangeDodge

			mobHealth := combat.Components.Health.GetFromCursor(mobCursor)
			mobIsLow := mobHealth.Value <= 10

			playerHighGround := mobPos.Y-playerPos.Y > 50

			if isInRangeJabAtk && scene.CurrentTick()-demonAnt.LastJabAttack > demonAnt.JabAttackDelay && !mobIsLow {
				if playerHighGround {
					mobKey := components.CharacterKeyComponent.GetFromCursor(mobCursor)
					attack := combatdata.SecondarySeqs[*mobKey].First()
					attack.StartTick = scene.CurrentTick()
					en.EnqueueAddComponentWithValue(combat.Components.Attack, attack)
					demonAnt.LastSlashAttack = scene.CurrentTick()
					demonAnt.LastAttack = scene.CurrentTick()
				} else {
					mobKey := components.CharacterKeyComponent.GetFromCursor(mobCursor)
					attack := combatdata.PrimarySeqs[*mobKey].First()
					attack.StartTick = scene.CurrentTick()
					en.EnqueueAddComponentWithValue(combat.Components.Attack, attack)
					demonAnt.LastJabAttack = scene.CurrentTick()
					demonAnt.LastAttack = scene.CurrentTick()
				}

				continue
			}

			canDodge := scene.CurrentTick() >= (demonAnt.LastDodged+demonAnt.DodgeDuration)+demonAnt.DodgeDelay
			canSlash := scene.CurrentTick()-demonAnt.LastSlashAttack > demonAnt.SlashAttackDelay

			if isInRangeDodge && !isAttacking {

				shouldDodge := rand.Intn(2) == 0

				if shouldDodge && canDodge {
					demonAnt.LastDodged = scene.CurrentTick()
					en.EnqueueAddComponentWithValue(components.DodgeComponent, components.Dodge{StartTick: scene.CurrentTick()})

					dodgeAwayFromPlayer := mobPos.X > playerPos.X
					mult := 1.0
					if !dodgeAwayFromPlayer {
						mult = -1.0
					}

					const edgeBuffer = 32.0
					atLeftEdge := mobPos.X <= mobBounds.MinX+edgeBuffer
					atRightEdge := mobPos.X >= mobBounds.MaxX-edgeBuffer

					if atLeftEdge && mult < 0 {
						mult = 1.0
					} else if atRightEdge && mult > 0 {
						mult = -1.0
					}

					if mult == -1.0 {
						demonAnt.DodgeDirection = spatial.NewDirectionLeft()
					} else {
						demonAnt.DodgeDirection = spatial.NewDirectionRight()
					}
					continue

				} else if !shouldDodge && canSlash {
					mobKey := components.CharacterKeyComponent.GetFromCursor(mobCursor)

					attack := combatdata.SecondarySeqs[*mobKey].First()
					attack.StartTick = scene.CurrentTick()
					en.EnqueueAddComponentWithValue(combat.Components.Attack, attack)
					demonAnt.LastSlashAttack = scene.CurrentTick()
					demonAnt.LastAttack = scene.CurrentTick()
					continue
				}
			}
		}
	}
	return nil
}
