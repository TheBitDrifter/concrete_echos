package sounds

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
)

var DialoguePing = client.SoundConfig{
	Path:             "sounds/dialogue_ping.wav",
	AudioPlayerCount: 1,
}

var AmbientWindNoise = client.SoundConfig{
	Path:             "sounds/wind.wav",
	AudioPlayerCount: 1,
}

var TrapDoorOpenSound = client.SoundConfig{
	Path:             "sounds/trap_door_open.wav",
	AudioPlayerCount: 1,
}

var DropMoneyAddedSound = client.SoundConfig{
	Path:             "sounds/get_money.wav",
	AudioPlayerCount: 1,
}

var SaveSound = client.SoundConfig{
	Path:             "sounds/notifications/save.wav",
	AudioPlayerCount: 1,
}

var SwapSound = client.SoundConfig{
	Path:             "sounds/box_head/swap.wav",
	AudioPlayerCount: 1,
}

var ExecuteSound = client.SoundConfig{
	Path:             "sounds/box_head/execute.wav",
	AudioPlayerCount: 1,
}

type SoundEnum int

const (
	Idle SoundEnum = iota
	Run
	Jump
	Land
	PrimaryAttack
	Hurt
	Dodge
	ObstacleHit
	Defeat
	SecondaryAttack
)

type registry struct {
	Characters map[characterkeys.CharEnum]map[SoundEnum]client.SoundConfig
}

var Registry = registry{
	Characters: make(map[characterkeys.CharEnum]map[SoundEnum]client.SoundConfig),
}

// Defaults -----

var runDefault = client.SoundConfig{
	Path:             "sounds/default/run.wav",
	AudioPlayerCount: 1,
}

var jumpDefault = client.SoundConfig{
	Path:             "sounds/default/jump.wav",
	AudioPlayerCount: 1,
}

var landDefault = client.SoundConfig{
	Path:             "sounds/default/land.wav",
	AudioPlayerCount: 1,
}

var hurtDefault = client.SoundConfig{
	Path:             "sounds/default/hurt.wav",
	AudioPlayerCount: 1,
}

var defeatDefault = client.SoundConfig{
	Path:             "sounds/default/defeat.wav",
	AudioPlayerCount: 1,
}

var dodgeDefault = client.SoundConfig{
	Path:             "sounds/default/dodge.wav",
	AudioPlayerCount: 4,
}

var primaryAttackDefault = client.SoundConfig{
	Path:             "sounds/default/primary_attack.wav",
	AudioPlayerCount: 3,
}

var secAttackDefault = client.SoundConfig{
	Path:             "sounds/default/secondary_attack.wav",
	AudioPlayerCount: 3,
}

var _ = func() error {
	Registry.Characters[characterkeys.Default] = map[SoundEnum]client.SoundConfig{}
	Registry.Characters[characterkeys.Default][Hurt] = hurtDefault
	Registry.Characters[characterkeys.Default][Defeat] = defeatDefault
	Registry.Characters[characterkeys.Default][Dodge] = dodgeDefault
	Registry.Characters[characterkeys.Default][PrimaryAttack] = primaryAttackDefault
	Registry.Characters[characterkeys.Default][SecondaryAttack] = secAttackDefault
	return nil
}()

// BoxHead -----

var boxHeadHurt = client.SoundConfig{
	Path:             "sounds/box_head/hurt.wav",
	AudioPlayerCount: 1,
}

var boxHeadDefeat = client.SoundConfig{
	Path:             "sounds/box_head/defeat.wav",
	AudioPlayerCount: 1,
}

var boxHeadPrimaryAttack = client.SoundConfig{
	Path:             "sounds/box_head/primary_attack.wav",
	AudioPlayerCount: 1,
}

var boxHeadObstacleHit = client.SoundConfig{
	Path:             "sounds/box_head/obstacle_hit.wav",
	AudioPlayerCount: 1,
}

var _ = func() error {
	Registry.Characters[characterkeys.BoxHead] = map[SoundEnum]client.SoundConfig{}
	Registry.Characters[characterkeys.BoxHead][Run] = runDefault
	Registry.Characters[characterkeys.BoxHead][Hurt] = boxHeadHurt
	Registry.Characters[characterkeys.BoxHead][Jump] = jumpDefault
	Registry.Characters[characterkeys.BoxHead][Land] = landDefault
	Registry.Characters[characterkeys.BoxHead][PrimaryAttack] = boxHeadPrimaryAttack
	Registry.Characters[characterkeys.BoxHead][Dodge] = dodgeDefault
	Registry.Characters[characterkeys.BoxHead][ObstacleHit] = boxHeadObstacleHit
	return nil
}()

// Chest -----

var chestDefeat = client.SoundConfig{
	Path:             "sounds/chest/defeat.wav",
	AudioPlayerCount: 1,
}

var _ = func() error {
	Registry.Characters[characterkeys.Chest] = map[SoundEnum]client.SoundConfig{}
	Registry.Characters[characterkeys.Chest][Defeat] = chestDefeat

	return nil
}()

// DemonAnt -----

var demonAntPrimaryAttack = client.SoundConfig{
	Path:             "sounds/demon_ant/primary_attack.wav",
	AudioPlayerCount: 3,
}

var _ = func() error {
	Registry.Characters[characterkeys.DemonAnt] = map[SoundEnum]client.SoundConfig{}
	Registry.Characters[characterkeys.DemonAnt][PrimaryAttack] = demonAntPrimaryAttack

	return nil
}()

// LazyRonin (customs)

var LazyRoninSounds = struct {
	SlashSound, ExplosionSound, SmallExplosionSound, AgroSound client.SoundConfig
}{
	SlashSound: client.SoundConfig{
		Path:             "sounds/lazy_ronin/slash.wav",
		AudioPlayerCount: 1,
	},
	ExplosionSound: client.SoundConfig{
		Path:             "sounds/lazy_ronin/explosion.wav",
		AudioPlayerCount: 1,
	},
	AgroSound: client.SoundConfig{
		Path:             "sounds/lazy_ronin/agro.wav",
		AudioPlayerCount: 1,
	},
	SmallExplosionSound: client.SoundConfig{
		Path:             "sounds/lazy_ronin/explosion_small.wav",
		AudioPlayerCount: 1,
	},
}
