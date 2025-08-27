package components

import (
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
)

var (
	OnGroundComponent            = warehouse.FactoryNewComponent[OnGround]()
	OnWallComponent              = warehouse.FactoryNewComponent[OnWall]()
	PlayerSceneTransferComponent = warehouse.FactoryNewComponent[PlayerSceneTransfer]()
	JumpStateComponent           = warehouse.FactoryNewComponent[JumpState]()
	PlayerSpawnComponent         = warehouse.FactoryNewComponent[PlayerSpawn]()
	PacerComponent               = warehouse.FactoryNewComponent[Pacer]()
	CharacterKeyComponent        = warehouse.FactoryNewComponent[characterkeys.CharEnum]()
	DodgeComponent               = warehouse.FactoryNewComponent[Dodge]()
	MovementCooldownComponent    = warehouse.FactoryNewComponent[MovementCooldowns]()
	AttackCooldownComponent      = warehouse.FactoryNewComponent[AttackCooldowns]()
	ContactComponent             = warehouse.FactoryNewComponent[Contact]()

	DialogueAutoStepperComponent   = warehouse.FactoryNewComponent[DialogueAutoStepper]()
	DialogueManualStepperComponent = warehouse.FactoryNewComponent[DialogueManualStepper]()

	CannotDashThroughComponent = warehouse.FactoryNewComponent[CannotDashThrough]()
	MobDemonAntComponent       = warehouse.FactoryNewComponent[MobDemonAnt]()
	MobSimpleAttackerComponent = warehouse.FactoryNewComponent[MobSimpleAttacker]()
	MobSimpleThrowerComponent  = warehouse.FactoryNewComponent[MobSimpleThrower]()

	MobSimpleHoriShooterComponent = warehouse.FactoryNewComponent[MobSimpleHoriShooter]()
	MobBoundsComponent            = warehouse.FactoryNewComponent[MobBounds]()
	MobFollowerComponent          = warehouse.FactoryNewComponent[MobFollower]()
	VanishingPlatformComponent    = warehouse.FactoryNewComponent[VanishingPlatform]()
	MobXMobCollisionComponent     = warehouse.FactoryNewComponent[MobXMobCollision]()
	TeleportSwapComponent         = warehouse.FactoryNewComponent[TeleportSwap]()
	SwapVulnerableComponent       = warehouse.FactoryNewComponent[SwapVulnerable]()

	EntityRefComponent = warehouse.FactoryNewComponent[EntityRef]()

	EntityReferencesComponent = warehouse.FactoryNewComponent[EntityReferences]()

	ObstacleComponent = warehouse.FactoryNewComponent[Obstacle]()

	SoftResetComponent           = warehouse.FactoryNewComponent[SoftReset]()
	SoftResetCheckpointComponent = warehouse.FactoryNewComponent[SoftResetCheckpoint]()
	TrapDoorComponent            = warehouse.FactoryNewComponent[TrapDoor]()

	WalletComponent = warehouse.FactoryNewComponent[Wallet]()
	DropComponent   = warehouse.FactoryNewComponent[Drop]()

	DialogueActivationComponent = warehouse.FactoryNewComponent[DialogueActivation]()

	InConversationComponent = warehouse.FactoryNewComponent[InConversation]()
	SceneTitleComponent     = warehouse.FactoryNewComponent[SceneTitle]()
	IsSavingComponent       = warehouse.FactoryNewComponent[IsSaving]()
	SaveActivationComponent = warehouse.FactoryNewComponent[SaveActivation]()

	PersistenceComponent = warehouse.FactoryNewComponent[Persistence]()

	WarpTotemComponent = warehouse.FactoryNewComponent[WarpTotem]()

	SimpleNotificationComponent = warehouse.FactoryNewComponent[SimpleNotification]()

	MobExecuteComponent           = warehouse.FactoryNewComponent[MobExecute]()
	PlayerIsExecutingComponent    = warehouse.FactoryNewComponent[PlayerIsExecuting]()
	PlayerExecutionCountComponent = warehouse.FactoryNewComponent[PlayerExecutionCount]()
	MobLazySkullyComponent        = warehouse.FactoryNewComponent[MobLazySkully]()

	FriendlyAgroComponent  = warehouse.FactoryNewComponent[FriendlyAgro]()
	MusicPlaylistComponent = warehouse.FactoryNewComponent[MusicPlaylist]()

	BossDefeatedComponent = warehouse.FactoryNewComponent[BossDefeated]()

	FastTravelActivationComponent = warehouse.FactoryNewComponent[FastTravelActivation]()

	LastCombatComponent = warehouse.FactoryNewComponent[LastCombat]()

	PlayerCamStateComponent = warehouse.FactoryNewComponent[PlayerCamState]()

	KeyRebindingStateComponent = warehouse.FactoryNewComponent[KeyRebindingState]()
)
