package scenes

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/animations"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
	"github.com/TheBitDrifter/concrete_echos/shared/combatdata"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/persistence"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

var DEFAULT_MOB_SOUND_BUNDLE = client.NewSoundBundle().
	AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Default][sounds.Run]).
	AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Default][sounds.Jump]).
	AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Default][sounds.Land]).
	AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Default][sounds.PrimaryAttack]).
	AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Default][sounds.SecondaryAttack]).
	AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Default][sounds.Hurt]).
	AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Default][sounds.Dodge]).
	AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Default][sounds.Defeat])

func NewMudHandMob(x, y float64, sto warehouse.Storage, isLeft bool, optionalMinX, optionalMaxX float64) (warehouse.Entity, error) {
	const (
		MUD_HAND_MOB_SPRITE_SHEET_PATH = "images/mobs/mud_hand_sheet.png"
		MUD_HAND_CASH_DROP             = 20
	)

	DEFAULT_MUD_HAND_SPR_BUNDLE := client.NewSpriteBundle().
		AddSprite(MUD_HAND_MOB_SPRITE_SHEET_PATH, true).
		WithAnimations(
			animations.Registry.Characters[characterkeys.MudHand][animations.Idle],
			animations.Registry.Characters[characterkeys.MudHand][animations.Hurt],
			animations.Registry.Characters[characterkeys.MudHand][animations.Defeat],
		).
		SetActiveAnimation(animations.Registry.Characters[characterkeys.MudHand][animations.Idle]).
		WithOffset(vector.Two{X: -31, Y: -34}).
		WithPriority(15).
		AddSprite(INTERACTION_MARKER_SHEET_PATH, true).
		WithAnimations(animations.InteractionMarkerAnim).
		WithOffset(vector.Two{X: -20, Y: -65}).
		WithCustomRenderer()

	comps := []warehouse.Component{
		components.PacerComponent,
		components.MobBoundsComponent,
	}
	comps = append(comps, DefaultMobComposition...)

	arche, err := sto.NewOrExistingArchetype(comps...)
	if err != nil {
		return nil, err
	}

	dir := spatial.NewDirectionRight()
	if isLeft {
		dir.SetLeft()
	}
	entities, err := arche.GenerateAndReturnEntity(1,
		spatial.NewPosition(x, y),
		spatial.NewRectangle(18, 58),
		motion.NewDynamics(10),
		dir,
		combatdata.HurtBoxes[characterkeys.MudHand],
		DEFAULT_MUD_HAND_SPR_BUNDLE,
		DEFAULT_MOB_SOUND_BUNDLE,
		combat.Health{Value: 20},
		characterkeys.MudHand,
		components.Pacer{
			IsLeft: isLeft,
			Speed:  80,
		},
		components.MobBounds{
			MinX: optionalMinX,
			MaxX: optionalMaxX,
		},
		components.Drop{
			MoneyDrop: MUD_HAND_CASH_DROP,
		},
	)
	if err != nil {
		return nil, err
	}
	return entities[0], nil
}

// -------------------------------------------------------------------------------------------------------------------
const (
	SKULL_RONIN_SPRITE_SHEET_PATH = "images/mobs/skull_ronin_sheet.png"
	SKULL_RONIN_STARTING_HEALTH   = 30

	SKULL_RONIN_MONEY_DROP = 50
)

var DEFAULT_SKULL_RONIN_SPR_BUNDLE = client.NewSpriteBundle().
	AddSprite(SKULL_RONIN_SPRITE_SHEET_PATH, true).
	WithAnimations(
		animations.Registry.Characters[characterkeys.SkullRonin][animations.Idle],
		animations.Registry.Characters[characterkeys.SkullRonin][animations.Run],
		animations.Registry.Characters[characterkeys.SkullRonin][animations.Defeat],
		animations.Registry.Characters[characterkeys.SkullRonin][animations.Hurt],
		animations.Registry.Characters[characterkeys.SkullRonin][animations.PrimaryAttack],
	).
	SetActiveAnimation(animations.Registry.Characters[characterkeys.SkullRonin][animations.Idle]).
	WithOffset(vector.Two{X: -150, Y: -100}).
	WithPriority(15).
	AddSprite(INTERACTION_MARKER_SHEET_PATH, true).
	WithAnimations(animations.InteractionMarkerAnim).
	WithOffset(vector.Two{X: -20, Y: -65}).
	WithCustomRenderer()

func NewSkullRoninMob(x, y float64, sto warehouse.Storage, left bool, optionalMinX, optionalMaxX, customVision float64) (warehouse.Entity, error) {
	comps := []warehouse.Component{
		components.MobSimpleAttackerComponent,
		components.MobBoundsComponent,
		components.MobFollowerComponent,
		components.MobXMobCollisionComponent,
	}

	comps = append(comps, DefaultMobComposition...)
	arche, err := sto.NewOrExistingArchetype(comps...)
	if err != nil {
		return nil, err
	}

	dir := spatial.NewDirectionRight()
	if left {
		dir = spatial.NewDirectionLeft()
	}

	if customVision == 0 {
		customVision = 300
	}
	entities, err := arche.GenerateAndReturnEntity(1,
		spatial.NewPosition(x, y),
		spatial.NewRectangle(18, 58),
		combatdata.HurtBoxes[characterkeys.SkullRonin],

		motion.NewDynamics(10),
		dir,
		DEFAULT_SKULL_RONIN_SPR_BUNDLE,
		DEFAULT_MOB_SOUND_BUNDLE,
		combat.Health{Value: SKULL_RONIN_STARTING_HEALTH},
		characterkeys.SkullRonin,
		components.MobSimpleAttacker{
			AttackVisionRadius: 65,
			Speed:              70,
			Delay:              120,
		},
		components.MobFollower{
			VisionRadius: customVision,
			StopRadius:   35,
			Speed:        38,
		},
		components.MobBounds{
			MinX: optionalMinX,
			MaxX: optionalMaxX,
		},
		components.Drop{
			MoneyDrop: SKULL_RONIN_MONEY_DROP,
		},
	)
	if err != nil {
		return nil, err
	}
	return entities[0], nil
}

// -------------------------------------------------------------------------------------------------------------------
const (
	CAN_THROWER_SPRITE_SHEET_PATH = "images/mobs/can_thrower_sheet.png"
	CAN_THROWER_STARTING_HEALTH   = 20
	CAN_THROWER_MONEY_DROP        = 35
)

var DEFAULT_CAN_THROWER_SPR_BUNDLE = client.NewSpriteBundle().
	AddSprite(CAN_THROWER_SPRITE_SHEET_PATH, true).
	WithAnimations(
		animations.Registry.Characters[characterkeys.CanThrower][animations.Idle],
		animations.Registry.Characters[characterkeys.CanThrower][animations.Defeat],
		animations.Registry.Characters[characterkeys.CanThrower][animations.Hurt],
		animations.Registry.Characters[characterkeys.CanThrower][animations.PrimaryAttack],
		animations.Registry.Characters[characterkeys.CanThrower][animations.PrimaryRanged],
	).
	SetActiveAnimation(animations.Registry.Characters[characterkeys.CanThrower][animations.Idle]).
	WithOffset(vector.Two{X: -54, Y: -54}).
	WithPriority(15).
	AddSprite(INTERACTION_MARKER_SHEET_PATH, true).
	WithAnimations(animations.InteractionMarkerAnim).
	WithOffset(vector.Two{X: -20, Y: -65}).
	WithCustomRenderer()

func NewCanThrower(x, y float64, sto warehouse.Storage, left bool, customVision, delay float64) (warehouse.Entity, error) {
	comps := []warehouse.Component{
		components.MobSimpleThrowerComponent,
	}

	comps = append(comps, DefaultMobComposition...)
	arche, err := sto.NewOrExistingArchetype(comps...)
	if err != nil {
		return nil, err
	}

	dir := spatial.NewDirectionRight()
	if left {
		dir = spatial.NewDirectionLeft()
	}

	if customVision == 0 {
		customVision = 600
	}
	if delay == 0 {
		delay = 240
	}

	entities, err := arche.GenerateAndReturnEntity(1,
		spatial.NewPosition(x, y),
		spatial.NewRectangle(20, 50),
		motion.NewDynamics(10),
		dir,
		combatdata.HurtBoxes[characterkeys.CanThrower],
		DEFAULT_CAN_THROWER_SPR_BUNDLE,
		DEFAULT_MOB_SOUND_BUNDLE,
		combat.Health{Value: CAN_THROWER_STARTING_HEALTH},
		characterkeys.CanThrower,
		components.MobSimpleThrower{
			VisionRadius:  customVision,
			Delay:         int(delay),
			ThrowRelease:  50,
			ThrowDuration: 60,
			SpawnOffset: vector.Two{
				X: 10,
				Y: -20,
			},
		},
		components.Drop{
			MoneyDrop: CAN_THROWER_MONEY_DROP,
		},
	)
	if err != nil {
		return nil, err
	}
	return entities[0], nil
}

const (
	CAN_STRAIGHT_THROWER_SPRITE_SHEET_PATH = "images/mobs/can_straight_thrower_sheet.png"
	CAN_STRAIGHT_THROWER_STARTING_HEALTH   = 20
)

var DEFAULT_CAN_STRAIGHT_THROWER_SPR_BUNDLE = client.NewSpriteBundle().
	AddSprite(CAN_STRAIGHT_THROWER_SPRITE_SHEET_PATH, true).
	WithAnimations(
		animations.Registry.Characters[characterkeys.CanStraightThrower][animations.Idle],
		animations.Registry.Characters[characterkeys.CanStraightThrower][animations.Defeat],
		animations.Registry.Characters[characterkeys.CanStraightThrower][animations.Hurt],
		animations.Registry.Characters[characterkeys.CanStraightThrower][animations.PrimaryAttack],
		animations.Registry.Characters[characterkeys.CanStraightThrower][animations.PrimaryRanged],
	).
	SetActiveAnimation(animations.Registry.Characters[characterkeys.CanStraightThrower][animations.Idle]).
	WithOffset(vector.Two{X: -54, Y: -54}).
	WithPriority(15).
	AddSprite(INTERACTION_MARKER_SHEET_PATH, true).
	WithAnimations(animations.InteractionMarkerAnim).
	WithOffset(vector.Two{X: -20, Y: -65}).
	WithCustomRenderer()

func NewCanStraightThrower(x, y float64, sto warehouse.Storage, left bool, customVision float64, customDelay float64) (warehouse.Entity, error) {
	comps := []warehouse.Component{
		components.MobSimpleHoriShooterComponent,
	}

	comps = append(comps, DefaultMobComposition...)
	arche, err := sto.NewOrExistingArchetype(comps...)
	if err != nil {
		return nil, err
	}

	dir := spatial.NewDirectionRight()
	if left {
		dir = spatial.NewDirectionLeft()
	}

	if customVision == 0 {
		customVision = 3000
	}
	if customDelay == 0 {
		customDelay = 240
	}
	entities, err := arche.GenerateAndReturnEntity(1,
		spatial.NewPosition(x, y),
		spatial.NewRectangle(20, 50),
		motion.NewDynamics(10),
		dir,
		combatdata.HurtBoxes[characterkeys.CanStraightThrower],
		DEFAULT_CAN_STRAIGHT_THROWER_SPR_BUNDLE,
		DEFAULT_MOB_SOUND_BUNDLE,
		combat.Health{Value: CAN_STRAIGHT_THROWER_STARTING_HEALTH},
		characterkeys.CanThrower,
		components.MobSimpleHoriShooter{
			VisionRadius:  customVision,
			Delay:         int(customDelay),
			ShootRelease:  50,
			ShootDuration: 60,
			SpawnOffset: vector.Two{
				X: 10,
				Y: -20,
			},
		},
		components.Drop{
			MoneyDrop: CAN_THROWER_MONEY_DROP,
		},
	)
	if err != nil {
		return nil, err
	}
	return entities[0], nil
}

// -------------------------------------------------------------------------------------------------------------------

func NewDemonAnt(x, y float64, sto warehouse.Storage, left bool, optionalMinX, optionalMaxX, customVision float64) (warehouse.Entity, error) {
	const (
		DEMON_ANT_SPRITE_SHEET_PATH = "images/mobs/demon_ant_sheet.png"
		DEMON_ANT_STARTING_HEALTH   = 40
		DEMON_ANT_MONEY_DROP        = 100
	)

	DEMON_ANT_SPR_BUNDLE := client.NewSpriteBundle().
		AddSprite(DEMON_ANT_SPRITE_SHEET_PATH, true).
		WithAnimations(
			animations.Registry.Characters[characterkeys.DemonAnt][animations.Idle],
			animations.Registry.Characters[characterkeys.DemonAnt][animations.Defeat],
			animations.Registry.Characters[characterkeys.DemonAnt][animations.Hurt],
			animations.Registry.Characters[characterkeys.DemonAnt][animations.PrimaryAttack],
			animations.Registry.Characters[characterkeys.DemonAnt][animations.SecondaryAttack],
			animations.Registry.Characters[characterkeys.DemonAnt][animations.Dodge],
			animations.Registry.Characters[characterkeys.DemonAnt][animations.Run],
		).
		WithOffset(vector.Two{X: -200, Y: -98}).
		WithPriority(15).
		AddSprite(INTERACTION_MARKER_SHEET_PATH, true).
		WithAnimations(animations.InteractionMarkerAnim).
		WithOffset(vector.Two{X: -20, Y: -65}).
		WithCustomRenderer()

	DEMON_ANT_SND_BUND := client.NewSoundBundle().
		AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Default][sounds.Run]).
		AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Default][sounds.Jump]).
		AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Default][sounds.Land]).
		AddSoundFromConfig(sounds.Registry.Characters[characterkeys.DemonAnt][sounds.PrimaryAttack]).
		AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Default][sounds.SecondaryAttack]).
		AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Default][sounds.Hurt]).
		AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Default][sounds.Dodge]).
		AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Default][sounds.Defeat])
	comps := []warehouse.Component{
		components.MobDemonAntComponent,
		components.MobBoundsComponent,
		components.MobFollowerComponent,
		components.MobXMobCollisionComponent,
	}

	comps = append(comps, DefaultMobComposition...)
	arche, err := sto.NewOrExistingArchetype(comps...)
	if err != nil {
		return nil, err
	}

	dir := spatial.NewDirectionRight()
	if left {
		dir = spatial.NewDirectionLeft()
	}
	if customVision == 0 {
		customVision = 200
	}

	entities, err := arche.GenerateAndReturnEntity(1,
		spatial.NewPosition(x, y),
		spatial.NewRectangle(30, 63),
		motion.NewDynamics(10),
		dir,
		combatdata.HurtBoxes[characterkeys.DemonAnt],
		DEMON_ANT_SPR_BUNDLE,
		DEMON_ANT_SND_BUND,
		combat.Health{Value: DEMON_ANT_STARTING_HEALTH},
		characterkeys.DemonAnt,
		components.MobDemonAnt{
			SlashAttackDelay: 30,

			AttackVisionJabRadius: 100,
			JabAttackDelay:        90,

			DodgeVisionRadius: 65,
			DodgeSpeed:        400,
			DodgeDuration:     20,
			DodgeDelay:        45,
		},
		components.MobFollower{
			VisionRadius: customVision,
			StopRadius:   35,
			Speed:        110,
		},
		components.MobBounds{
			MinX: optionalMinX,
			MaxX: optionalMaxX,
		},
		components.Drop{
			MoneyDrop: DEMON_ANT_MONEY_DROP,
		},
	)
	if err != nil {
		return nil, err
	}
	return entities[0], nil
}

// ---------------------------

func NewChest(x, y float64, sto warehouse.Storage, left bool, moneyDrop float64, persistID persistence.PersistenceID, callbackID int, useAltSprite bool) (warehouse.Entity, error) {
	const (
		CHEST_STARTING_HEALTH = 30
	)

	CHEST_SPRITE_SHEET_PATH := "images/mobs/chest_sheet.png"
	if useAltSprite {
		CHEST_SPRITE_SHEET_PATH = "images/mobs/chest_sheet_alt.png"
	}

	CHEST_SPR_BUNDLE := client.NewSpriteBundle().
		AddSprite(CHEST_SPRITE_SHEET_PATH, true).
		WithAnimations(
			animations.Registry.Characters[characterkeys.Chest][animations.Idle],
			animations.Registry.Characters[characterkeys.Chest][animations.Defeat],
			animations.Registry.Characters[characterkeys.Chest][animations.Hurt],
		).
		WithOffset(vector.Two{X: -100, Y: -100}).
		WithPriority(15)
	soundBund := client.NewSoundBundle().
		AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Default][sounds.Hurt]).
		AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Chest][sounds.Defeat])
	comps := []warehouse.Component{
		components.IgnoreContactDamageTag,
		components.IgnoreSwapTag,
		components.ChestTag,
		components.PersistenceComponent,
		components.IgnoreExecuteTag,
	}

	comps = append(comps, DefaultMobComposition...)
	arche, err := sto.NewOrExistingArchetype(comps...)
	if err != nil {
		return nil, err
	}

	dir := spatial.NewDirectionRight()
	if left {
		dir = spatial.NewDirectionLeft()
	}

	entities, err := arche.GenerateAndReturnEntity(1,
		spatial.NewPosition(x, y),
		spatial.NewRectangle(50, 30),
		motion.NewDynamics(10),
		dir,
		combatdata.HurtBoxes[characterkeys.Chest],
		CHEST_SPR_BUNDLE,
		soundBund,
		combat.Health{Value: CHEST_STARTING_HEALTH},
		characterkeys.Chest,
		components.Drop{
			MoneyDrop:              moneyDrop,
			HealthDrop:             50,
			CustomSpawnCallbackKey: callbackID,
		},
		components.Persistence{
			EntityType: persistence.ENITY_TYPE_CHESTS,
			PersistID:  persistID,
		},
	)
	if err != nil {
		return nil, err
	}
	return entities[0], nil
}

func NewBatFlyer(x, y float64, sto warehouse.Storage, isLeft bool, distX, distY float64) (warehouse.Entity, error) {
	const (
		SPRITE_SHEET_PATH = "images/mobs/bat_sheet.png"
		CASH_DROP         = 20
	)

	SPR_BUNDLE := client.NewSpriteBundle().
		AddSprite(SPRITE_SHEET_PATH, true).
		WithAnimations(
			animations.Registry.Characters[characterkeys.Bat][animations.Idle],
			animations.Registry.Characters[characterkeys.Bat][animations.Hurt],
			animations.Registry.Characters[characterkeys.Bat][animations.Defeat],
		).
		SetActiveAnimation(animations.Registry.Characters[characterkeys.Bat][animations.Idle]).
		WithOffset(vector.Two{X: -32, Y: -32}).
		WithPriority(15).
		AddSprite(INTERACTION_MARKER_SHEET_PATH, true).
		WithAnimations(animations.InteractionMarkerAnim).
		WithOffset(vector.Two{X: -20, Y: -65}).
		WithCustomRenderer()

	comps := []warehouse.Component{
		components.PacerComponent,
		components.MobBoundsComponent,
		components.NoGravityTag,
		components.IgnoreTerrainCollisionsMob,
	}
	comps = append(comps, DefaultMobComposition...)

	arche, err := sto.NewOrExistingArchetype(comps...)
	if err != nil {
		return nil, err
	}

	dir := spatial.NewDirectionRight()
	if isLeft {
		dir.SetLeft()
	}

	var minX, maxX, minY, maxY float64

	if distX != 0 {
		minX = x
		maxX = minX + distX
	}

	if distY != 0 {
		minY = y
		maxY = minY + distY
	}
	isVert := false
	if distY != 0 {
		isVert = true
	}

	boxDimen := 20.0
	entities, err := arche.GenerateAndReturnEntity(1,
		spatial.NewPosition(x, y),
		spatial.NewRectangle(boxDimen, boxDimen),
		motion.NewDynamics(10),
		dir,
		combat.HurtBoxes{
			combat.NewHurtBox(boxDimen, boxDimen, 0, 0),
		},
		SPR_BUNDLE,
		DEFAULT_MOB_SOUND_BUNDLE,
		combat.Health{Value: 10},
		characterkeys.Bat,
		components.Pacer{
			IsLeft:       isLeft,
			Speed:        160,
			IsVert:       isVert,
			SwapDirOnHit: false,
		},
		components.MobBounds{
			MinX:          minX,
			MaxX:          maxX,
			MinY:          minY,
			MaxY:          maxY,
			NoResetOnSwap: true,
		},
		components.Drop{
			MoneyDrop: CASH_DROP,
		},
	)
	if err != nil {
		return nil, err
	}
	return entities[0], nil
}
