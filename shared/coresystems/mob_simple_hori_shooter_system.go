package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/scenes"
)

// MobSimpleHorizontalShooterSystem handles the logic for mobs that shoot projectiles horizontally.
type MobSimpleHorizontalShooterSystem struct{}

// shotProjectileBuilder is a helper struct to hold projectile data before creation.
type shotProjectileBuilder struct {
	pos      spatial.Position
	shape    spatial.Shape
	velocity vector.Two
	mass     float64
	mobID    int
	mobRe    int
}

// Run executes the system's logic for the current frame.
func (s MobSimpleHorizontalShooterSystem) Run(scene blueprint.Scene, dt float64) error {
	// Query for all mobs that have the simple shooter component and are not defeated.
	mobQuery := warehouse.Factory.NewQuery().And(
		components.MobTag,
		components.MobSimpleHoriShooterComponent,
		warehouse.Factory.NewQuery().Not(
			components.PlayerTag,
			combat.Components.Defeat,
		),
	)
	mobCursor := scene.NewCursor(mobQuery)

	// Query for the player entity.
	playerQuery := warehouse.Factory.NewQuery().And(components.PlayerTag)
	playerCursor := scene.NewCursor(playerQuery)

	// A slice to hold all projectiles that need to be created this frame.
	projectilesToFire := []shotProjectileBuilder{}

	// Iterate over all shooter mobs.
	for range mobCursor.Next() {
		mbEn, _ := mobCursor.CurrentEntity()
		mobID := mbEn.ID()
		mobRe := mbEn.Recycled()

		simpleShooter := components.MobSimpleHoriShooterComponent.GetFromCursor(mobCursor)
		if components.SwapVulnerableComponent.CheckCursor(mobCursor) {
			simpleShooter.IsWindingUp = false
			simpleShooter.LastStartedShot = 0
			continue
		}

		for range playerCursor.Next() {
			playerPos := spatial.Components.Position.GetFromCursor(playerCursor)
			mobPos := spatial.Components.Position.GetFromCursor(mobCursor)

			mobAttackVisionSq := simpleShooter.VisionRadius * simpleShooter.VisionRadius
			distanceSquared := playerPos.Sub(mobPos.Two).MagSquared()
			inRangeToShoot := distanceSquared <= mobAttackVisionSq

			isReadyToStartShot := scene.CurrentTick() > simpleShooter.LastFiredShot+simpleShooter.Delay

			if inRangeToShoot && isReadyToStartShot && !simpleShooter.IsWindingUp {
				simpleShooter.LastStartedShot = scene.CurrentTick()
				simpleShooter.IsWindingUp = true
			}

			if playerPos.X < mobPos.X {
				spatial.Components.Direction.GetFromCursor(mobCursor).SetLeft()
			} else {
				spatial.Components.Direction.GetFromCursor(mobCursor).SetRight()
			}

			// If the mob is not in the wind-up phase, skip to the next mob.
			if !simpleShooter.IsWindingUp {
				continue
			}

			// Check if it's time to release the projectile during the wind-up animation.
			isReadyToReleaseShot := scene.CurrentTick() > simpleShooter.LastStartedShot+simpleShooter.ShootRelease

			if isReadyToReleaseShot {
				// Mark the time the shot was fired and end the wind-up state.
				simpleShooter.LastFiredShot = scene.CurrentTick()
				simpleShooter.IsWindingUp = false

				// --- Simplified Projectile Logic ---
				// This section replaces the complex trajectory calculation from the thrower.

				// Define the constant speed for the projectile.
				const PROJECTILE_SPEED = 300.0

				// Determine the projectile's starting position using the mob's position and an offset.
				mobPosAdj := vector.Two{
					X: mobPos.X + simpleShooter.SpawnOffset.X*spatial.Components.Direction.GetFromCursor(mobCursor).AsFloat(),
					Y: mobPos.Y + simpleShooter.SpawnOffset.Y,
				}

				// Set the velocity based on the player's direction.
				// The projectile has no vertical movement (Y velocity is 0).
				launchVelocity := vector.Two{X: 0, Y: 0}
				if playerPos.X < mobPos.X {
					launchVelocity.X = -PROJECTILE_SPEED // Fire left
				} else {
					launchVelocity.X = PROJECTILE_SPEED // Fire right
				}

				// Create the projectile data and add it to the list to be spawned.
				builder := shotProjectileBuilder{
					pos:      spatial.NewPosition(mobPosAdj.X, mobPosAdj.Y),
					shape:    spatial.NewRectangle(30, 30), // Example shape
					velocity: launchVelocity,
					mass:     1, // Mass doesn't affect gravity-less projectiles
					mobID:    int(mobID),
					mobRe:    mobRe,
				}
				projectilesToFire = append(projectilesToFire, builder)
			}
		}
	}

	// This loop creates the actual projectile entities in the scene.
	for _, projectileB := range projectilesToFire {
		comps := append(scenes.ProjectileComposition, components.NoGravityTag)
		archetype, err := scene.Storage().NewOrExistingArchetype(comps...)
		if err != nil {
			return err
		}
		// Create a dynamics component and set its velocity. No gravity is applied.
		dynamics := motion.NewDynamics(projectileB.mass)
		dynamics.Vel = projectileB.velocity

		// Define animation data for the projectile's sprite.
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
		bundle := client.NewSpriteBundle().AddSprite("images/projectiles/trash_projectile.png", true).WithAnimations(anim).WithPriority(80)

		// Set the projectile's direction based on its velocity.
		direction := spatial.NewDirectionRight()
		if projectileB.velocity.X < 0 {
			direction = spatial.NewDirectionLeft()
		}

		// Generate the entity with all its components.
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
