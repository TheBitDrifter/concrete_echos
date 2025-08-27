package scenes

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/animations"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
	"github.com/TheBitDrifter/concrete_echos/shared/combatdata"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

const (
	PLAYER_SPRITE_SHEET_PATH = "images/characters/box_man_sheet.png"
	PLAYER_MAIN_HUD_PATH     = "images/characters/box_man_main_hud.png"
	PLAYER_HEART_HUD_PATH    = "images/characters/box_man_heart.png"

	PLAYER_DEFEAT_SCREEN_PATH = "images/defeat_screen_sheet.png"
	PLAYER_STARTING_HEALTH    = 100
)

var DEFAULT_PLAYER_SPR_BUNDLE = client.NewSpriteBundle().
	AddSprite(PLAYER_SPRITE_SHEET_PATH, true).
	WithAnimations(
		animations.Registry.Characters[characterkeys.BoxHead][animations.Idle],
		animations.Registry.Characters[characterkeys.BoxHead][animations.Run],
		animations.Registry.Characters[characterkeys.BoxHead][animations.Fall],
		animations.Registry.Characters[characterkeys.BoxHead][animations.Jump],
		animations.Registry.Characters[characterkeys.BoxHead][animations.PrimaryAttack],
		animations.Registry.Characters[characterkeys.BoxHead][animations.Hurt],
		animations.Registry.Characters[characterkeys.BoxHead][animations.Dodge],
		animations.Registry.Characters[characterkeys.BoxHead][animations.Aerial],
		animations.Registry.Characters[characterkeys.BoxHead][animations.AerialDownSmash],
		animations.Registry.Characters[characterkeys.BoxHead][animations.AerialDownSmashLanding],
		animations.Registry.Characters[characterkeys.BoxHead][animations.Teleport],
		animations.Registry.Characters[characterkeys.BoxHead][animations.TeleportMarker],
		animations.Registry.Characters[characterkeys.BoxHead][animations.TeleportEffect],
		animations.Registry.Characters[characterkeys.BoxHead][animations.Defeat],
		animations.Registry.Characters[characterkeys.BoxHead][animations.InConvo],
		animations.Registry.Characters[characterkeys.BoxHead][animations.IsSaving],
		animations.Registry.Characters[characterkeys.BoxHead][animations.UpAttack],
		animations.Registry.Characters[characterkeys.BoxHead][animations.UpAerial],
		animations.Registry.Characters[characterkeys.BoxHead][animations.Execute],
	).
	WithOffset(vector.Two{X: -72, Y: -101}).
	WithPriority(20).
	WithCustomRenderer().
	AddSprite(PLAYER_MAIN_HUD_PATH, true).
	WithCustomRenderer().
	AddSprite(PLAYER_HEART_HUD_PATH, true).
	WithCustomRenderer().
	AddSprite(PLAYER_DEFEAT_SCREEN_PATH, true).
	WithAnimations(animations.DefeatScreenAnim).
	WithCustomRenderer()

var DEFAULT_PLAYER_SND_BUNDLE = client.NewSoundBundle().
	AddSoundFromConfig(sounds.Registry.Characters[characterkeys.BoxHead][sounds.Run]).
	AddSoundFromConfig(sounds.Registry.Characters[characterkeys.BoxHead][sounds.Jump]).
	AddSoundFromConfig(sounds.Registry.Characters[characterkeys.BoxHead][sounds.Land]).
	AddSoundFromConfig(sounds.Registry.Characters[characterkeys.BoxHead][sounds.PrimaryAttack]).
	AddSoundFromConfig(sounds.Registry.Characters[characterkeys.BoxHead][sounds.Hurt]).
	AddSoundFromConfig(sounds.Registry.Characters[characterkeys.BoxHead][sounds.Dodge]).
	AddSoundFromConfig(sounds.DropMoneyAddedSound).
	AddSoundFromConfig(sounds.Registry.Characters[characterkeys.BoxHead][sounds.ObstacleHit]).
	AddSoundFromConfig(sounds.SaveSound).
	AddSoundFromConfig(sounds.SwapSound).
	AddSoundFromConfig(sounds.ExecuteSound)

func NewPlayer(x, y float64, sto warehouse.Storage) (warehouse.Entity, error) {
	playerArchetype, err := sto.NewOrExistingArchetype(
		PlayerComposition...,
	)

	entities, err := playerArchetype.GenerateAndReturnEntity(1,
		spatial.NewPosition(x, y),
		spatial.NewRectangle(18, 58),
		combatdata.HurtBoxes[characterkeys.BoxHead],
		motion.NewDynamics(10),
		spatial.NewDirectionRight(),
		input.ActionBuffer{ReceiverIndex: 0},
		client.CameraIndex(0),
		DEFAULT_PLAYER_SND_BUNDLE,
		DEFAULT_PLAYER_SPR_BUNDLE,
		characterkeys.BoxHead,
		combat.Health{Value: PLAYER_STARTING_HEALTH},
	)
	if err != nil {
		return nil, err
	}
	return entities[0], nil
}
