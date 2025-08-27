package coresystems

import (
	"sort"

	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/actions"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

type TeleportSwapSystem struct{}

func (TeleportSwapSystem) Run(scene blueprint.Scene, dt float64) error {
	const maxDistance = 350.0
	const maxDistanceSq = maxDistance * maxDistance
	const cd = 240

	playerQuery := warehouse.Factory.NewQuery().And(
		components.TeleportSwapComponent, components.PlayerTag,
		components.WarpSwapUnlockedTag,
		warehouse.Factory.NewQuery().Not(
			combat.Components.Defeat,
			INPUT_DISABLED_COMPONENTS,
		),
	)
	playerCursor := scene.NewCursor(playerQuery)

	mobQuery := warehouse.Factory.NewQuery().And(
		warehouse.Factory.NewQuery().Or(
			components.MobTag,
			components.WarpTotemComponent,
		),
		warehouse.Factory.NewQuery().Not(
			combat.Components.Defeat,
			components.IgnoreSwapTag,
		),
	)
	mobCursor := scene.NewCursor(mobQuery)

	for range playerCursor.Next() {
		actionBuffer := input.Components.ActionBuffer.GetFromCursor(playerCursor)
		_, actionShiftLeft := actionBuffer.ConsumeAction(actions.ShiftTeleTargetLeft)
		_, actionShiftRight := actionBuffer.ConsumeAction(actions.ShiftTeleTargetRight)
		_, actionShiftNear := actionBuffer.ConsumeAction(actions.ShiftTeleTargetNear)
		_, actionSwap := actionBuffer.ConsumeAction(actions.TeleSwap)

		playerPos := spatial.Components.Position.GetFromCursor(playerCursor)
		teleState := components.TeleportSwapComponent.GetFromCursor(playerCursor)

		if scene.CurrentTick() < teleState.StartTick+cd {
			return nil
		}

		type mobWithDist struct {
			entity warehouse.Entity
			distSq float64
		}

		var nearbyMobs []mobWithDist

		for range mobCursor.Next() {
			mobPos := spatial.Components.Position.GetFromCursor(mobCursor)
			distVec := playerPos.Sub(mobPos.Two)
			distSq := distVec.MagSquared()

			if distSq <= maxDistanceSq {
				mobEntity, err := mobCursor.CurrentEntity()
				if err != nil {
					return err
				}
				nearbyMobs = append(nearbyMobs, mobWithDist{entity: mobEntity, distSq: distSq})
			}
		}

		if len(nearbyMobs) == 0 {
			teleState.HasTarget = false
			return nil
		}

		sort.Slice(nearbyMobs, func(i, j int) bool {
			return nearbyMobs[i].distSq < nearbyMobs[j].distSq
		})

		if !teleState.HasTarget {
			closestMobEntity := nearbyMobs[0].entity
			teleState.ActiveTarget = closestMobEntity
			teleState.HasTarget = true
		} else {
			if teleState.ActiveTarget == nil {
				continue
			}
			isDefeated := combat.Components.Defeat.Check(teleState.ActiveTarget.Table())

			targetPos := spatial.Components.Position.GetFromEntity(teleState.ActiveTarget)
			distVec := playerPos.Sub(targetPos.Two)
			distSq := distVec.MagSquared()

			if isDefeated || distSq > maxDistanceSq {
				teleState.HasTarget = false
			}
		}

		if teleState.HasTarget {
			currentIndex := -1
			for i, mob := range nearbyMobs {
				if mob.entity.ID() == teleState.ActiveTarget.ID() {
					currentIndex = i
					break
				}
			}

			if actionShiftLeft {
				if currentIndex != -1 {
					newIndex := (currentIndex + 1) % len(nearbyMobs)
					teleState.ActiveTarget = nearbyMobs[newIndex].entity
				}
			}

			if actionShiftRight {
				if currentIndex != -1 {
					newIndex := (currentIndex - 1 + len(nearbyMobs)) % len(nearbyMobs)
					teleState.ActiveTarget = nearbyMobs[newIndex].entity
				}
			}

			if actionShiftNear {
				teleState.ActiveTarget = nearbyMobs[0].entity
			}

			if actionSwap && teleState.HasTarget {

				target := teleState.ActiveTarget
				target, _ = scene.Storage().Entity(int(target.ID()))

				if combat.Components.Defeat.Check(target.Table()) {
					continue
				}

				if components.SwapVulnerableComponent.Check(target.Table()) {
					sv := components.SwapVulnerableComponent.GetFromEntity(target)
					sv.StartTick = scene.CurrentTick()
				} else {
					target.EnqueueAddComponentWithValue(components.SwapVulnerableComponent, components.SwapVulnerable{StartTick: scene.CurrentTick()})
				}

				isTotem := components.WarpTotemComponent.Check(target.Table())

				targetPos := spatial.Components.Position.GetFromEntity(target)
				playerOldX := playerPos.X
				playerOldY := playerPos.Y
				playerPos.X = targetPos.X
				playerPos.Y = targetPos.Y

				mOK := components.MobBoundsComponent.Check(target.Table())
				if mOK {
					mobBounds := components.MobBoundsComponent.GetFromEntity(target)

					if !mobBounds.NoResetOnSwap {
						mobBounds.MinX = 0
						mobBounds.MaxX = 0
					}
				}

				if !isTotem {
					targetPos.X = playerOldX
					targetPos.Y = playerOldY
				}

				teleState.StartTick = scene.CurrentTick()

				if iv, ok := combat.Components.Invincible.GetFromCursorSafe(playerCursor); ok {
					iv.StartTick = scene.CurrentTick()
				} else {
					playerEn, _ := playerCursor.CurrentEntity()
					playerEn.EnqueueAddComponentWithValue(combat.Components.Invincible, combat.Invincible{StartTick: scene.CurrentTick()})

				}

			}
		}
	}
	return nil
}
