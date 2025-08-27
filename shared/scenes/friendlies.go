package scenes

import (
	"math"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/dialogue"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/animations"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

const INTERACTION_MARKER_SHEET_PATH = "images/interaction_marker_sheet.png"

func NewDrifterFriendly(x, y float64, sto warehouse.Storage, isLeft bool, slidesKey int) (warehouse.Entity, error) {
	const DRIFTER_SPR_PATH = "images/mobs/drifter_sheet.png"

	DRIFTER_SPR_BUNDLE := client.NewSpriteBundle().
		AddSprite(DRIFTER_SPR_PATH, true).
		WithAnimations(
			animations.Registry.Characters[characterkeys.Drifter][animations.Idle],
			animations.Registry.Characters[characterkeys.Drifter][animations.InConvo],
			animations.Registry.Characters[characterkeys.Drifter][animations.ConvoStart],
		).
		WithOffset(vector.Two{X: -32, Y: -48}).
		WithPriority(15).
		AddSprite(INTERACTION_MARKER_SHEET_PATH, true).
		WithAnimations(animations.InteractionMarkerAnim).
		WithOffset(vector.Two{X: -20, Y: -65}).
		WithCustomRenderer()

	comps := []warehouse.Component{
		components.IgnoreContactDamageTag,
		components.IgnoreSwapTag,
		components.DialogueActivationComponent,
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
		spatial.NewRectangle(30, 58),
		motion.NewDynamics(10),
		dir,
		DRIFTER_SPR_BUNDLE,
		combat.Health{Value: 900},
		characterkeys.Drifter,
		components.DialogueActivation{SlidesID: dialogue.SlidesEnum(slidesKey), Range: 60, MinRange: 30},
		components.MobBounds{MinX: 0, MaxX: math.MaxFloat64},
		components.MobFollower{Speed: 0, VisionRadius: 100},
	)
	if err != nil {
		return nil, err
	}
	return entities[0], nil
}

const (
	LazyRoninMaxHealth   = 100
	LazyRoninWakeUpTicks = 120
)

func NewLazySkullRoninFriendly(x, y float64, sto warehouse.Storage, slidesKey dialogue.SlidesEnum, convoCallback dialogue.CallbackEnum) (warehouse.Entity, error) {
	const SPR_PATH = "images/mobs/lazy_skull_ronin_sheet.png"

	SPR_BUNDLE := client.NewSpriteBundle().
		AddSprite(SPR_PATH, true).
		WithAnimations(
			animations.Registry.Characters[characterkeys.LazySkullRonin][animations.Idle],
			animations.Registry.Characters[characterkeys.LazySkullRonin][animations.InConvo],
			animations.LazySkullRoninAnims.Animations[1],
			animations.LazySkullRoninAnims.Animations[2],
			animations.LazySkullRoninAnims.Animations[3],
			animations.LazySkullRoninAnims.Animations[4],
			animations.LazySkullRoninAnims.Animations[5],
			animations.LazySkullRoninAnims.Animations[6],
			animations.LazySkullRoninAnims.Animations[7],
			animations.LazySkullRoninAnims.Animations[8],
		).
		WithOffset(vector.Two{X: -300, Y: -275}).
		WithPriority(15).
		AddSprite(INTERACTION_MARKER_SHEET_PATH, true).
		WithAnimations(animations.InteractionMarkerAnim).
		WithOffset(vector.Two{X: -20, Y: -65}).
		WithCustomRenderer()

	comps := []warehouse.Component{
		components.IgnoreSwapTag,
		components.DialogueActivationComponent,
		components.MobBoundsComponent,
		components.MobFollowerComponent,
		components.MobLazySkullyComponent,
		components.FriendlyAgroComponent,
		components.IsBossTag,
	}
	comps = append(comps, DefaultMobComposition...)

	arche, err := sto.NewOrExistingArchetype(comps...)
	if err != nil {
		return nil, err
	}

	dir := spatial.NewDirectionRight()

	soundBundle := client.NewSoundBundle().
		AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Default][sounds.Run]).
		AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Default][sounds.Jump]).
		AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Default][sounds.Land]).
		AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Default][sounds.Hurt]).
		AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Default][sounds.Dodge]).
		AddSoundFromConfig(sounds.Registry.Characters[characterkeys.Default][sounds.Defeat]).
		AddSoundFromConfig(sounds.LazyRoninSounds.SlashSound).
		AddSoundFromConfig(sounds.LazyRoninSounds.ExplosionSound).
		AddSoundFromConfig(sounds.LazyRoninSounds.SmallExplosionSound).
		AddSoundFromConfig(sounds.LazyRoninSounds.AgroSound)

	entities, err := arche.GenerateAndReturnEntity(1,
		spatial.NewPosition(x, y),
		spatial.NewRectangle(18, 58),
		motion.NewDynamics(10),
		dir,
		soundBundle,
		SPR_BUNDLE,
		combat.Health{Value: LazyRoninMaxHealth},
		characterkeys.LazySkullRonin,
		combat.HurtBoxes{
			combat.NewHurtBox(20, 55, 0, 0),
		},
		components.DialogueActivation{SlidesID: slidesKey, ConvoCallbackID: convoCallback, Range: 60, MinRange: 30, MustBeRight: true},
		components.MobBounds{MinX: 1750, MaxX: 2750},
	)
	if err != nil {
		return nil, err
	}
	return entities[0], nil
}
