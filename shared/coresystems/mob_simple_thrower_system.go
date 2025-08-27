package coresystems

import (
	"math/rand"

	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/scenes"
)

type MobSimpleThrowerSystem struct{}

type thrownProjectileBuilder struct {
	pos      spatial.Position
	shape    spatial.Shape
	velocity vector.Two
	mass     float64
	mobID    int
	mobRe    int
}

func (s MobSimpleThrowerSystem) Run(scene blueprint.Scene, dt float64) error {
	mobQuery := warehouse.Factory.NewQuery().And(
		components.MobTag,
		components.MobSimpleThrowerComponent,
		warehouse.Factory.NewQuery().Not(
			components.PlayerTag,
			combat.Components.Defeat,
		),
	)
	mobCursor := scene.NewCursor(mobQuery)

	playerQuery := warehouse.Factory.NewQuery().And(components.PlayerTag)
	playerCursor := scene.NewCursor(playerQuery)

	projectilesToFire := []thrownProjectileBuilder{}

	for range mobCursor.Next() {
		mbEn, _ := mobCursor.CurrentEntity()
		mobID := mbEn.ID()
		mobRe := mbEn.Recycled()

		simpleThrower := components.MobSimpleThrowerComponent.GetFromCursor(mobCursor)
		if components.SwapVulnerableComponent.CheckCursor(mobCursor) {
			simpleThrower.IsWindingUp = false
			simpleThrower.LastStartedThrow = 0
			continue
		}

		for range playerCursor.Next() {
			playerPos := spatial.Components.Position.GetFromCursor(playerCursor)
			mobPos := spatial.Components.Position.GetFromCursor(mobCursor)
			mobAttackVisionSq := simpleThrower.VisionRadius * simpleThrower.VisionRadius
			distanceSquared := playerPos.Sub(mobPos.Two).MagSquared()

			inRangeToShoot := distanceSquared <= mobAttackVisionSq

			isReadyToStartThrow := scene.CurrentTick() > simpleThrower.LastFiredThrow+simpleThrower.Delay

			if inRangeToShoot && isReadyToStartThrow && !simpleThrower.IsWindingUp {
				simpleThrower.LastStartedThrow = scene.CurrentTick()
				simpleThrower.IsWindingUp = true
			}

			throwing := scene.CurrentTick() < simpleThrower.LastStartedThrow+simpleThrower.ThrowDuration
			if !throwing {
				if playerPos.X < mobPos.X {
					spatial.Components.Direction.GetFromCursor(mobCursor).SetLeft()
				} else {
					spatial.Components.Direction.GetFromCursor(mobCursor).SetRight()
				}
			}

			if !simpleThrower.IsWindingUp {
				continue
			}

			isReadyToReleaseThrow := scene.CurrentTick() > simpleThrower.LastStartedThrow+simpleThrower.ThrowRelease

			if isReadyToReleaseThrow {
				simpleThrower.LastFiredThrow = scene.CurrentTick()
				simpleThrower.IsWindingUp = false

				const (
					minTime               = 0.5
					maxTime               = 2.0
					maxDistanceForScaling = 1500.0
					timeGuess             = 1.5
				)

				playerVel := motion.Components.Dynamics.GetFromCursor(playerCursor).Vel
				mobPosAdj := vector.Two{
					X: mobPos.X + simpleThrower.SpawnOffset.X*spatial.Components.Direction.GetFromCursor(mobCursor).AsFloat(),
					Y: mobPos.Y + simpleThrower.SpawnOffset.Y,
				}

				// Calculate direct distance to determine lead scaling
				directDistance := playerPos.Two.Sub(mobPosAdj).Mag()
				distanceRatio := directDistance / maxDistanceForScaling
				if distanceRatio > 1.0 {
					distanceRatio = 1.0
				}

				// Scale lead factor with distance. Closer targets have less lead.
				// At max distance (ratio=1.0), lead is ~0.9-1.1.
				// At min distance (ratio=0.0), lead is ~0.1-0.3.
				leadFactorX := (0.1 + distanceRatio*0.8) + rand.Float64()*0.2

				// Reduce Y-axis lead to prevent over-throwing on jumps
				leadFactorY := leadFactorX * 0.4

				initialInterceptPos := vector.Two{
					X: playerPos.X + (playerVel.X * leadFactorX * timeGuess),
					Y: playerPos.Y + (playerVel.Y * leadFactorY * timeGuess),
				}

				distanceToTarget := initialInterceptPos.Sub(mobPosAdj).Mag()
				timeRatioForImpact := distanceToTarget / maxDistanceForScaling
				if timeRatioForImpact > 1.0 {
					timeRatioForImpact = 1.0
				}
				timeToImpact := minTime + (maxTime-minTime)*timeRatioForImpact

				finalInterceptPos := vector.Two{
					X: playerPos.X + (playerVel.X * leadFactorX * timeToImpact),
					Y: playerPos.Y + (playerVel.Y * leadFactorY * timeToImpact),
				}
				toInterceptVec := finalInterceptPos.Sub(mobPosAdj)

				actualAcceleration := DEFAULT_GRAVITY * PIXELS_PER_METER
				horizontalVelocity := toInterceptVec.X / timeToImpact
				initialVerticalVelocity := (toInterceptVec.Y / timeToImpact) - (0.5 * actualAcceleration * timeToImpact)
				launchVelocity := vector.Two{X: horizontalVelocity, Y: initialVerticalVelocity}

				builder := thrownProjectileBuilder{
					pos:      spatial.NewPosition(mobPosAdj.X, mobPosAdj.Y),
					shape:    spatial.NewRectangle(30, 30),
					velocity: launchVelocity,
					mass:     10,
					mobID:    int(mobID),
					mobRe:    mobRe,
				}
				projectilesToFire = append(projectilesToFire, builder)
			}
		}
	}

	for _, projectileB := range projectilesToFire {
		archetype, err := scene.Storage().NewOrExistingArchetype(scenes.ProjectileComposition...)
		if err != nil {
			return err
		}
		dynamics := motion.NewDynamics(projectileB.mass)
		dynamics.Vel = projectileB.velocity
		anim := client.AnimationData{
			Name:        "idle",
			RowIndex:    0,
			Speed:       5,
			FrameCount:  6,
			FrameWidth:  72,
			FrameHeight: 69,
			PositionOffset: vector.Two{
				X: -36,
				Y: -30,
			},
		}
		bundle := client.NewSpriteBundle().AddSprite("images/projectiles/fireball.png", true).WithAnimations(anim).WithPriority(80)
		direction := spatial.NewDirectionRight()
		if projectileB.velocity.X < 0 {
			direction = spatial.NewDirectionLeft()
		}
		err = archetype.Generate(
			1,
			projectileB.pos,
			projectileB.shape,
			components.ProjectileTag,
			dynamics,
			bundle,
			direction,
			components.EntityRef{
				ID:       projectileB.mobID,
				Recycled: projectileB.mobRe,
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}
