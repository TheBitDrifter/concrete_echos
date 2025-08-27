package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/tteokbokki/tteo_coresystems"
)

var DefaultCoreSystems = []blueprint.CoreSystem{
	MobCollisionMarkerSystem{},
	GravitySystem{},                      // Apply gravity forces
	FrictionSystem{},                     // Apply Friction forces
	PlayerMovementSystem{},               // Apply player input forces
	tteo_coresystems.IntegrationSystem{}, // Update velocities and positions
	tteo_coresystems.TransformSystem{},   // Update collision shapes
	PlayerBlockCollisionSystem{},         // Handle collisions
	NewPlayerPlatformCollisionSystem(),   // Handle collisions — func returns ptr because system is not pure (has state)
	OnGroundClearingSystem{},             // Clear onGround
	ClearDroppingStateSystem{},           // Clear ignorePlatform
	PlayerAttackSystem{},                 // Player Attacks (Triggering)
	MobTerrainCollisionSystem{},
	MobPacerSystem{},
	MobSimpleThrowerSystem{},
	MobSimpleHorizontalShooterSystem{},
	AttackHurtBoxCollisionSystem{},
	ProjectileCollisionSystem{},
	PlayerTrapDoorCollisionSystem{},
	IframeRemoverSystem{},
	HurtRemoverSystem{},
	ContactDamageSystem{},
	DefeatSystem{},
	MobBoundsPlatformDetectionSystem{},
	MobSimpleAttackerSystem{},
	MobDemonAntMovementSystem{},
	MobFollowerSystem{},
	MobCollisionHandlerSystem{},
	TeleportSwapSystem{},
	SwapVurnRemoverSystem{},
	SwapVurnStunSystem{},
	PlayerObstacleCollisionSystem{},
	ObstacleMovementSystem{},
	SoftResetSystem{},
	SoftResetBoundsSystem{},
	SoftResetCheckpointActivationSystem{},
	TrapDoorSystem{},
	DropSystem{},
	MobKBSystem{},
	PlayerAerialInterruptSystem{},
	MobExecuteSystem{},
	MobLazySkullySystem{},
}
