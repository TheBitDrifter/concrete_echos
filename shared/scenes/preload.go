package scenes

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

var DEFAULT_PRELOAD = client.NewPreLoadBlueprint().
	// Box-Head
	AddSprite(PLAYER_SPRITE_SHEET_PATH).
	AddSprite(PLAYER_MAIN_HUD_PATH).
	AddSprite(PLAYER_HEART_HUD_PATH).
	AddSprite(PLAYER_DEFEAT_SCREEN_PATH).
	AddSound(sounds.Registry.Characters[characterkeys.BoxHead][sounds.Run]).
	AddSound(sounds.Registry.Characters[characterkeys.BoxHead][sounds.Jump]).
	AddSound(sounds.Registry.Characters[characterkeys.BoxHead][sounds.Land]).
	AddSound(sounds.Registry.Characters[characterkeys.BoxHead][sounds.PrimaryAttack]).
	AddSound(sounds.Registry.Characters[characterkeys.BoxHead][sounds.Hurt]).
	AddSound(sounds.Registry.Characters[characterkeys.BoxHead][sounds.Dodge]).
	AddSound(sounds.Registry.Characters[characterkeys.BoxHead][sounds.ObstacleHit]).
	AddSound(sounds.Registry.Characters[characterkeys.Default][sounds.Jump]).
	AddSound(sounds.Registry.Characters[characterkeys.Default][sounds.Land]).
	AddSound(sounds.Registry.Characters[characterkeys.Default][sounds.PrimaryAttack]).
	AddSound(sounds.Registry.Characters[characterkeys.Default][sounds.Hurt]).
	AddSound(sounds.DropMoneyAddedSound).
	AddSound(sounds.SaveSound).
	AddSound(sounds.SwapSound).
	AddSound(sounds.ExecuteSound).
	AddSound(sounds.BossSong).
	AddSound(sounds.PostBossSong).
	// Rest Misc
	AddSound(sounds.DialoguePing).
	AddSprite("images/projectiles/trash_projectile.png").
	AddSprite("images/portraits/1.png").
	AddSprite("images/portraits/2.png").
	AddSprite("images/portraits/3.png").
	AddSprite("images/dialogue_box_sheet.png").
	AddSprite("images/dialogue_next_sheet.png").
	AddSprite("images/unlock_abil_sheet.png").
	AddSprite("images/boss_defeat_sheet.png").
	AddSprite("images/mobs/chest_sheet.png").
	AddSprite("images/mobs/chest_sheet_alt.png")
