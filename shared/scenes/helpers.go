package scenes

import (
	"fmt"
	"strconv"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/dialogue"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/animations"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/dialoguedata"
	"github.com/TheBitDrifter/concrete_echos/shared/persistence"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

func NewPlayerSpawn(x, y float64, sto warehouse.Storage) (warehouse.Entity, error) {
	spawnArchetype, err := sto.NewOrExistingArchetype(
		components.PlayerSpawnComponent,
	)
	entities, err := spawnArchetype.GenerateAndReturnEntity(1,
		components.PlayerSpawn{X: x, Y: y},
	)
	if err != nil {
		return nil, err
	}
	return entities[0], nil
}

func NewPlatformRotated(sto warehouse.Storage, x, y, rotation float64) error {
	platformArche, err := sto.NewOrExistingArchetype(PlatformComposition...)
	if err != nil {
		return err
	}
	return platformArche.Generate(1,
		spatial.NewPosition(x, y),
		spatial.Rotation(rotation),
		spatial.NewTriangularPlatform(144, 16),
		client.NewSpriteBundle().
			AddSprite("images/terrain/platform.png", true).
			WithOffset(vector.Two{X: -72, Y: -8}),
	)
}

func NewRamp(sto warehouse.Storage, x, y float64) error {
	// Add a sprite
	composition := []warehouse.Component{
		client.Components.SpriteBundle,
	}

	composition = append(composition, BlockTerrainComposition...)
	rampArche, err := sto.NewOrExistingArchetype(composition...)
	if err != nil {
		return err
	}

	return rampArche.Generate(1,
		spatial.NewPosition(x, y),
		spatial.NewDoubleRamp(250, 46, 0.2),
		client.NewSpriteBundle().
			AddSprite("images/terrain/ramp.png", true).
			WithOffset(vector.Two{X: -125, Y: -22}).WithPriority(10),
	)
}

func NewCollisionPlayerTransfer(
	sto warehouse.Storage, x, y, w, h, playerTargetX, playerTargetY float64, target string,
) error {
	collisionPlayerTransferArche, err := sto.NewOrExistingArchetype(
		CollisionPlayerTransferComposition...,
	)
	if err != nil {
		return err
	}
	return collisionPlayerTransferArche.Generate(1,
		spatial.NewPosition(x, y),
		spatial.NewRectangle(w, h),
		components.PlayerSceneTransfer{
			Dest: target,
			X:    playerTargetX,
			Y:    playerTargetY,
		},
	)
}

func NewAmbientWindNoise(sto warehouse.Storage) error {
	windArche, err := sto.NewOrExistingArchetype(AmbientNoiseComposition...)
	if err != nil {
		return err
	}
	return windArche.Generate(1, client.NewSoundBundle().AddSoundFromPath("sounds/wind.wav"))
}

func newVanishingPlatform(sto warehouse.Storage, x, y, rotation, width, height float64, size string) error {
	compo := []warehouse.Component{}
	compo = append(compo, PlatformComposition...)
	compo = append(compo, components.VanishingPlatformComponent)

	vanishingArche, err := sto.NewOrExistingArchetype(compo...)
	if err != nil {
		return err
	}

	spriteSheetPath := fmt.Sprintf("images/terrain/vanishing_platform_%s_sheet.png", size)
	spriteOffset := vector.Two{X: -(width / 2), Y: -(height / 2)}

	sprBundle := client.NewSpriteBundle().
		AddSprite(spriteSheetPath, true).
		WithOffset(spriteOffset).
		WithAnimations(
			client.AnimationData{
				RowIndex:    0,
				Name:        "inwait",
				FrameCount:  1,
				FrameWidth:  int(width),
				FrameHeight: int(height),
				Speed:       5,
				Freeze:      true,
			},
			client.AnimationData{
				RowIndex:    0,
				Name:        "despawn",
				FrameCount:  10,
				FrameWidth:  int(width),
				FrameHeight: int(height),
				Speed:       10,
				Freeze:      true,
			},
			client.AnimationData{
				RowIndex:    1,
				Name:        "respawn",
				FrameCount:  9,
				FrameWidth:  int(width),
				FrameHeight: int(height),
				Speed:       10,
				Freeze:      true,
			},
		)

	return vanishingArche.Generate(1,
		spatial.NewPosition(x, y),
		spatial.Rotation(rotation),
		spatial.NewTriangularPlatform(width, height),
		sprBundle,
		components.VanishingPlatform{
			LiveDuration: 100,
			RespawnDelay: 180,
		},
	)
}

func NewVanishingPlatformSmall(sto warehouse.Storage, x, y, rotation float64) error {
	return newVanishingPlatform(sto, x, y, rotation, 32, 16, "small")
}

func NewVanishingPlatformMed(sto warehouse.Storage, x, y, rotation float64) error {
	return newVanishingPlatform(sto, x, y, rotation, 80, 16, "med")
}

func NewVanishingPlatformLg(sto warehouse.Storage, x, y, rotation float64) error {
	return newVanishingPlatform(sto, x, y, rotation, 144, 16, "lg")
}

func NewSoftResetCheckpoint(sto warehouse.Storage, x, y, w, h float64, automaticActive bool) error {
	compo := []warehouse.Component{
		components.SoftResetCheckpointComponent,
		spatial.Components.Position,
		spatial.Components.Shape,
	}
	arche, err := sto.NewOrExistingArchetype(compo...)
	if err != nil {
		return err
	}
	return arche.Generate(1,
		spatial.NewPosition(x, y),
		spatial.NewRectangle(h, w),
		components.SoftResetCheckpoint{
			Activated: automaticActive,
		},
	)
}

func NewTrapDoor(sto warehouse.Storage, x, y float64, isLeft bool, callbackKey components.TrapDoorEnum, persistID persistence.PersistenceID) error {
	width := 16.0
	height := 80.0

	spriteSheetPath := "images/terrain/trap_door_sheet.png"
	spriteOffset := vector.Two{X: -(width / 2), Y: -(height / 2)}

	sprBundle := client.NewSpriteBundle().
		AddSprite(spriteSheetPath, true).
		WithOffset(spriteOffset).
		WithAnimations(
			client.AnimationData{
				RowIndex:    1,
				Name:        "close",
				FrameCount:  6,
				FrameWidth:  int(width),
				FrameHeight: int(height),
				Speed:       5,
				Freeze:      true,
			},
			client.AnimationData{
				RowIndex:    0,
				Name:        "open",
				FrameCount:  6,
				FrameWidth:  int(width),
				FrameHeight: int(height),
				Speed:       5,
				Freeze:      true,
			},
		)
	compo := TrapDoorComposition
	arche, err := sto.NewOrExistingArchetype(compo...)
	dir := spatial.NewDirectionRight()
	if isLeft {
		dir = spatial.NewDirectionLeft()
	}
	if err != nil {
		return err
	}
	return arche.Generate(1,
		spatial.NewPosition(x, y),
		spatial.NewRectangle(16, 80),
		sprBundle,
		dir,
		components.TrapDoor{IsOpenCallbackID: callbackKey},
		client.NewSoundBundle().AddSoundFromConfig(sounds.TrapDoorOpenSound),

		components.Persistence{
			EntityType: persistence.ENITY_TYPE_TRAP_DOORS,
			PersistID:  persistID,
		},
	)
}

func NewDialogueEntity(sto warehouse.Storage, slidesEnum dialogue.SlidesEnum) (warehouse.Entity, error) {
	archetype, err := sto.NewOrExistingArchetype(DialogueComposition...)

	convo := dialogue.Conversation{}
	convo.SlidesID = slidesEnum
	slides := dialogue.SlidesRegistry[dialogue.SlidesEnum(convo.SlidesID)]

	bundle := client.NewSpriteBundle()
	bundle = bundle.AddSprite("images/dialogue_box_sheet.png", true).
		WithCustomRenderer().
		WithAnimations(animations.DialogueOpenAnimation, animations.DialogueCloseAnimation)

	seenPortraits := map[dialogue.PortraitEnum]struct{}{}

	for _, slide := range slides {
		path := "images/portraits/" + strconv.Itoa(int(slide.PortraitID)) + ".png"

		if _, ok := seenPortraits[slide.PortraitID]; !ok {
			seenPortraits[slide.PortraitID] = struct{}{}

			bundle = bundle.AddSpriteAtIndex(int(slide.PortraitID), path, true).
				WithCustomRenderer()

		}
	}

	if err != nil {
		return nil, err
	}
	entities, error := archetype.GenerateAndReturnEntity(
		1,
		convo,
		bundle,
		spatial.NewDirectionRight(),
		spatial.NewPosition(dialoguedata.DEFAULT_BOX_POS.AsFloats()),
	)
	if err != nil {
		return nil, error
	}
	return entities[0], err
}

func NewDialogueEntityAutoStepper(sto warehouse.Storage, slidesEnum dialogue.SlidesEnum) (warehouse.Entity, error) {
	comp := []warehouse.Component{components.DialogueAutoStepperComponent}
	comp = append(comp, DialogueComposition...)
	archetype, err := sto.NewOrExistingArchetype(comp...)

	convo := dialogue.Conversation{}
	convo.SlidesID = slidesEnum
	slides := dialogue.SlidesRegistry[convo.SlidesID]

	bundle := client.NewSpriteBundle()
	bundle = bundle.AddSprite("images/dialogue_box_sheet.png", true).
		WithCustomRenderer().
		WithAnimations(animations.DialogueOpenAnimation, animations.DialogueCloseAnimation).
		AddSprite("images/dialogue_next_sheet.png", true).
		WithAnimations(animations.DialogueNext)

	i := 2
	convo.PortraitIDForSpriteBundleBlueprintIndex = make(map[dialogue.PortraitEnum]int)

	for _, slide := range slides {
		path := "images/portraits/" + strconv.Itoa(int(slide.PortraitID)) + ".png"

		if _, ok := convo.PortraitIDForSpriteBundleBlueprintIndex[slide.PortraitID]; !ok {
			convo.PortraitIDForSpriteBundleBlueprintIndex[slide.PortraitID] = i
			i++
			bundle = bundle.AddSprite(path, true).
				WithCustomRenderer()
		}
	}
	if err != nil {
		return nil, err
	}

	soundBundle := client.NewSoundBundle().
		AddSoundFromConfig(sounds.DialoguePing)
	entities, error := archetype.GenerateAndReturnEntity(
		1,
		convo,
		bundle,
		soundBundle,
		spatial.NewDirectionRight(),
		spatial.NewPosition(dialoguedata.DEFAULT_BOX_POS.AsFloats()),
		components.DialogueAutoStepper{},
	)
	if err != nil {
		return nil, error
	}
	return entities[0], err
}

func NewDialogueEntityManualStepper(sto warehouse.Storage, slidesEnum dialogue.SlidesEnum, callbackEnum dialogue.CallbackEnum, delay int, entities ...warehouse.Entity) (warehouse.Entity, error) {
	comp := []warehouse.Component{components.DialogueManualStepperComponent, components.EntityReferencesComponent}
	comp = append(comp, DialogueComposition...)
	archetype, err := sto.NewOrExistingArchetype(comp...)

	convo := dialogue.Conversation{}
	convo.SlidesID = slidesEnum
	convo.CallbackID = callbackEnum

	bundle := client.NewSpriteBundle()
	bundle = bundle.AddSprite("images/dialogue_box_sheet.png", true).
		WithCustomRenderer().
		WithAnimations(animations.DialogueOpenAnimation, animations.DialogueCloseAnimation).
		AddSprite("images/dialogue_next_sheet.png", true).
		WithAnimations(animations.DialogueNext)

	i := 2
	convo.PortraitIDForSpriteBundleBlueprintIndex = make(map[dialogue.PortraitEnum]int)

	slides := dialogue.SlidesRegistry[convo.SlidesID]
	for _, slide := range slides {
		path := "images/portraits/" + strconv.Itoa(int(slide.PortraitID)) + ".png"

		if _, ok := convo.PortraitIDForSpriteBundleBlueprintIndex[slide.PortraitID]; !ok {
			convo.PortraitIDForSpriteBundleBlueprintIndex[slide.PortraitID] = i
			i++
			bundle = bundle.AddSprite(path, true).
				WithCustomRenderer()

		}
	}

	if err != nil {
		return nil, err
	}

	references := components.EntityReferences{}

	for i, en := range entities {
		references.Active[i] = true
		references.Refs[i] = components.EntityRef{ID: int(en.ID()), Recycled: en.Recycled()}
	}

	soundBundle := client.NewSoundBundle().
		AddSoundFromConfig(sounds.DialoguePing)
	entities, error := archetype.GenerateAndReturnEntity(
		1,
		convo,
		bundle,
		soundBundle,
		spatial.NewDirectionRight(),
		spatial.NewPosition(dialoguedata.DEFAULT_BOX_POS.AsFloats()),
		components.DialogueManualStepper{MinDelayInTicks: delay},
		references,
	)
	if err != nil {
		return nil, error
	}
	return entities[0], err
}

func NewSceneTitle(sto warehouse.Storage, title string) (warehouse.Entity, error) {
	archetype, err := sto.NewOrExistingArchetype(components.SceneTitleComponent)

	entities, error := archetype.GenerateAndReturnEntity(
		1,
		components.SceneTitle{Value: title},
	)
	if err != nil {
		return nil, error
	}
	return entities[0], err
}

func NewSaveBench(sto warehouse.Storage, x, y float64, isLeft bool, optionalID persistence.OptionalSaveID) error {
	const SPR_PATH = "images/terrain/save_bench.png"

	SPR_BUNDLE := client.NewSpriteBundle().
		AddSprite(SPR_PATH, true).
		WithOffset(vector.Two{X: -24, Y: -24}).
		WithPriority(15).
		AddSprite("images/interaction_marker_sheet.png", true).
		WithAnimations(animations.InteractionMarkerAnim).
		WithOffset(vector.Two{X: -18, Y: -42}).
		WithCustomRenderer()
	directionLR := spatial.NewDirectionRight()
	if isLeft {
		directionLR = spatial.NewDirectionLeft()
	}

	compo := []warehouse.Component{
		client.Components.SpriteBundle,
		client.Components.SoundBundle,
		spatial.Components.Position,
		spatial.Components.Direction,
		spatial.Components.Shape,
		components.SaveActivationComponent,
	}

	archetype, err := sto.NewOrExistingArchetype(compo...)
	if err != nil {
		return err
	}

	err = archetype.Generate(1,
		SPR_BUNDLE,
		spatial.NewPosition(x, y),
		directionLR,
		spatial.NewRectangle(48, 42),
		components.SaveActivation{Range: 20, OptionalID: optionalID},
	)
	if err != nil {
		return err
	}

	return nil
}

func NewWarpTotem(sto warehouse.Storage, x, y float64, isLeft bool) error {
	const SPR_PATH = "images/terrain/save_bench.png"

	SPR_BUNDLE := client.NewSpriteBundle().
		AddSprite("images/terrain/totem.png", true).
		WithOffset(vector.Two{X: -16, Y: -32}).
		WithPriority(15)

	directionLR := spatial.NewDirectionRight()
	if isLeft {
		directionLR = spatial.NewDirectionLeft()
	}

	compo := []warehouse.Component{
		client.Components.SpriteBundle,
		spatial.Components.Position,
		spatial.Components.Direction,
		spatial.Components.Shape,
		components.WarpTotemComponent,
	}

	archetype, err := sto.NewOrExistingArchetype(compo...)
	if err != nil {
		return err
	}

	err = archetype.Generate(1,
		SPR_BUNDLE,
		spatial.NewPosition(x, y),
		directionLR,
		spatial.NewRectangle(32, 64),
	)
	if err != nil {
		return err
	}

	return nil
}

func NewSimpleNotification(sto warehouse.Storage, title, body string, titlePaddingX, titlePaddingY float64) error {
	const TITLE_MAX_W = 300
	const BODY_MAX_W = 360

	SPR_BUNDLE := client.NewSpriteBundle().
		AddSprite("images/unlock_abil_sheet.png", true).
		WithAnimations(client.AnimationData{
			FrameCount:  6,
			FrameWidth:  400,
			FrameHeight: 187,
			Speed:       5,
			Freeze:      true,
		}).
		WithStatic(true).
		WithPriority(100).
		WithOffset(vector.Two{X: 120, Y: 70})

	compo := []warehouse.Component{
		client.Components.SpriteBundle,
		spatial.Components.Position,
		spatial.Components.Direction,
		spatial.Components.Shape,
		components.SimpleNotificationComponent,
	}

	archetype, err := sto.NewOrExistingArchetype(compo...)
	if err != nil {
		return err
	}

	err = archetype.Generate(1,
		SPR_BUNDLE,
		components.SimpleNotification{Title: title, Body: body, PaddingX: titlePaddingX, PaddingY: titlePaddingY, TitleMaxWidth: TITLE_MAX_W, BodyMaxWidth: BODY_MAX_W},
	)
	if err != nil {
		return err
	}

	return nil
}

func AddPlaylist(sto warehouse.Storage, startIndex int, songs sounds.SoundCollection) error {
	compo := []warehouse.Component{
		client.Components.SoundBundle,
		components.MusicPlaylistComponent,
	}

	bundle := client.NewSoundBundle()

	var start client.SoundConfig

	for i, song := range songs.Sounds {
		if i == startIndex {
			start = song
		}
		bundle = bundle.AddSoundFromConfig(song)
	}

	archetype, err := sto.NewOrExistingArchetype(compo...)
	if err != nil {
		return err
	}

	err = archetype.Generate(1, bundle, components.MusicPlaylist{
		ActiveSong:       start,
		SongCollectionID: songs.ID,
		Collection:       songs,
		IsFading:         true,
	})
	if err != nil {
		return err
	}

	return nil
}

func NewBossDefeat(sto warehouse.Storage, ct int) error {
	SPR_BUNDLE := client.NewSpriteBundle().
		AddSprite("images/boss_defeat_sheet.png", true).
		WithAnimations(client.AnimationData{
			FrameCount:  6,
			FrameWidth:  360,
			FrameHeight: 48,
			Speed:       5,
			Freeze:      true,
		}).
		WithStatic(true).
		WithPriority(100).
		WithOffset(vector.Two{X: 130, Y: 100})

	compo := []warehouse.Component{
		client.Components.SpriteBundle,
		spatial.Components.Position,
		spatial.Components.Direction,
		spatial.Components.Shape,
		components.BossDefeatedComponent,
	}

	archetype, err := sto.NewOrExistingArchetype(compo...)
	if err != nil {
		return err
	}

	err = archetype.Generate(1,
		SPR_BUNDLE,
		components.BossDefeated{StartTick: ct},
	)
	if err != nil {
		return err
	}

	return nil
}

func NewTravelTotem(sto warehouse.Storage, x, y float64, isLeft bool, name string) error {
	SPR_BUNDLE := client.NewSpriteBundle().
		AddSprite("images/terrain/ft_totem.png", true).
		WithOffset(vector.Two{X: -25, Y: -64}).
		WithPriority(15).
		AddSprite(INTERACTION_MARKER_SHEET_PATH, true).
		WithAnimations(animations.InteractionMarkerAnim).
		WithOffset(vector.Two{X: -20, Y: -90}).
		WithCustomRenderer()

	directionLR := spatial.NewDirectionRight()
	if isLeft {
		directionLR = spatial.NewDirectionLeft()
	}

	compo := []warehouse.Component{
		client.Components.SpriteBundle,
		spatial.Components.Position,
		spatial.Components.Direction,
		spatial.Components.Shape,
		components.FastTravelActivationComponent,
	}

	archetype, err := sto.NewOrExistingArchetype(compo...)
	if err != nil {
		return err
	}

	err = archetype.Generate(1,
		SPR_BUNDLE,
		spatial.NewPosition(x, y),
		directionLR,
		spatial.NewRectangle(50, 120),
		components.FastTravelActivation{Range: 50, Name: name},
	)
	if err != nil {
		return err
	}

	return nil
}
