package components

import "github.com/TheBitDrifter/bappa/warehouse"

// Tags help us identify/categorize archetypes/entities when their
// composition alone isn't enough.
//
// For example its hard to tell the
// difference between a block and platform since they both have
// dynamics, shapes, sprites, etc
type blockTag struct{}

type platformTag struct{}

type musicTag struct{}

type playerTag struct{}

type mobTag struct{}

type cutsceneTag struct{}

type ambiNoiseTag struct{}

type projectileTag struct{}

type ignoreFrictionDampTag struct{}

type isDroppingThroughPlatformTag struct{}

type stillCollidingWithPlatformTag struct{}

type noGrav struct{}

type ignoreSwap struct{}

type ignoreContactDamage struct{}

type chestTag struct{}

type warpSwapUnlocked struct{}

type wallJumpUnlocked struct{}

type ignoreTerrainCollisionsMob struct{}

type isBoss struct{}

type isFastTravelButton struct{}

type ignoreExecute struct{}

var (
	BlockTerrainTag = warehouse.FactoryNewComponent[blockTag]()
	PlatformTag     = warehouse.FactoryNewComponent[platformTag]()
	MusicTag        = warehouse.FactoryNewComponent[musicTag]()

	AmbientNoiseTag = warehouse.FactoryNewComponent[ambiNoiseTag]()

	PlayerTag = warehouse.FactoryNewComponent[playerTag]()
	MobTag    = warehouse.FactoryNewComponent[mobTag]()

	CutsceneTag = warehouse.FactoryNewComponent[cutsceneTag]()

	ProjectileTag = warehouse.FactoryNewComponent[projectileTag]()

	IgnoreDefaultFrictionDampTag = warehouse.FactoryNewComponent[ignoreFrictionDampTag]()

	IsDroppingThroughPlatformTag = warehouse.FactoryNewComponent[isDroppingThroughPlatformTag]()

	StillCollidingWithPlatformTag = warehouse.FactoryNewComponent[stillCollidingWithPlatformTag]()

	NoGravityTag = warehouse.FactoryNewComponent[noGrav]()

	IgnoreSwapTag = warehouse.FactoryNewComponent[ignoreSwap]()

	IgnoreContactDamageTag = warehouse.FactoryNewComponent[ignoreContactDamage]()

	IgnoreTerrainCollisionsMob = warehouse.FactoryNewComponent[ignoreTerrainCollisionsMob]()

	ChestTag = warehouse.FactoryNewComponent[chestTag]()

	WarpSwapUnlockedTag   = warehouse.FactoryNewComponent[warpSwapUnlocked]()
	WallJumpUnlockedTag   = warehouse.FactoryNewComponent[wallJumpUnlocked]()
	IsBossTag             = warehouse.FactoryNewComponent[isBoss]()
	IsFastTravelButtonTag = warehouse.FactoryNewComponent[isFastTravelButton]()

	IgnoreExecuteTag = warehouse.FactoryNewComponent[ignoreExecute]()
)
