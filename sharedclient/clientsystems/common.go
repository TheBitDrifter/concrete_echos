package clientsystems

import (
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/coldbrew/coldbrew_clientsystems"
	"github.com/TheBitDrifter/concrete_echos/shared/dialoguedata"
	"github.com/TheBitDrifter/concrete_echos/sharedclient/fontdata"
)

const (
	DIALOGUE_PING_VOL = 0.135
	MUSIC_VOL         = 0.075
	AMBIENT_NOISE_VOL = 0.07
	PLAYER_SOUNDS_VOL = 0.45
	MOB_SOUNDS_VOL    = 0.4
)

var perfMonitor = NewPerformanceMonitorSystem(55.0, 60, 120)

var DefaultClientSystems = []coldbrew.ClientSystem{
	perfMonitor,
	VectorMovementConverterSytem{},
	DebugPositionSystem{},
	PlayerSoundSystem{Volume: PLAYER_SOUNDS_VOL},
	MobSoundSystem{Volume: MOB_SOUNDS_VOL},
	MobLazyRoninSoundSystem{Volume: MOB_SOUNDS_VOL},
	MusicSystem{Volume: MUSIC_VOL},
	AmbientNoiseSystem{Volume: AMBIENT_NOISE_VOL},
	PlayerAnimationSystem{},
	&CameraFollowerSystem{},
	&coldbrew_clientsystems.BackgroundScrollSystem{},
	PlayerSpawnSystem{},
	CollisionPlayerTransferSystem{},
	MobAnimationSystem{},
	DumbDefeatTransferSystem{},
	coldbrew_clientsystems.DialogueTextSystem{
		TEXT_REVEAL_DELAY_IN_TICKS: dialoguedata.TEXT_REVEAL_DELAY_IN_TICKS,
		FONT_SIZE:                  fontdata.DEFAULT_FONT_SIZE,
		FONT_FACE:                  fontdata.DEFAULT_FONT_FACE,
		MAX_LINE_WIDTH:             fontdata.MAX_LINE_WIDTH,
	},
	coldbrew_clientsystems.DialogueSoundSystem{Volume: DIALOGUE_PING_VOL, TEXT_REVEAL_DELAY_IN_TICKS: dialoguedata.TEXT_REVEAL_DELAY_IN_TICKS, SoundOnWord: true},
	DialogueAutoStepperSystem{},
	DialogueManualStepperSystem{},
	VanishingPlatformAnimationSystem{},
	GodSystem{},
	TrapDoorAnimationStateSystem{},
	TrapDoorSoundSystem{Volume: PLAYER_SOUNDS_VOL},
	DropSoundSystem{Volume: PLAYER_SOUNDS_VOL},
	DialogueManualStepperActivationSystem{},
	DialogueManualStepperClearingSystem{},
	SaveSystem{
		MinSaveTicks: 80,
	},

	SavingClearingSystem{},
	NotifTextRevealSystem{
		REVEAL_START_DELAY: 60,
	},
	FastTravelSceneActivationSystem{},
	&DumbSingletonOptionalBossMusicChangeSystem{},
	InteractRemoverSystem{},
	KeyRebindingSystem{},
}

var DefaultFastTravelSystems = []coldbrew.ClientSystem{
	FastTravelSceneSystem{
		PaddingX:     100,
		PaddingY:     100,
		SpacingY:     60,
		ButtonSqSize: 32,
		FONT_FACE:    fontdata.DEFAULT_FONT_FACE,
	},
}

var CutSceneClientSystems = []coldbrew.ClientSystem{
	&DumbCutSceneSystem{},
	coldbrew_clientsystems.DialogueTextSystem{
		TEXT_REVEAL_DELAY_IN_TICKS: dialoguedata.TEXT_REVEAL_DELAY_IN_TICKS_ALT,
		FONT_SIZE:                  fontdata.DEFAULT_FONT_SIZE,
		FONT_FACE:                  fontdata.DEFAULT_FONT_FACE,
		MAX_LINE_WIDTH:             fontdata.MAX_LINE_WIDTH,
	},
	MusicSystem{Volume: MUSIC_VOL},

	coldbrew_clientsystems.DialogueSoundSystem{Volume: DIALOGUE_PING_VOL, TEXT_REVEAL_DELAY_IN_TICKS: dialoguedata.TEXT_REVEAL_DELAY_IN_TICKS_ALT, SoundOnWord: true},
	AmbientNoiseSystem{Volume: AMBIENT_NOISE_VOL},
	DialogueAutoStepperSystem{},
}

var DefaultClientSystemsNetworked = []coldbrew.ClientSystem{
	VectorMovementConverterSytem{},
	PlayerSoundSystem{Volume: PLAYER_SOUNDS_VOL},
	MobSoundSystem{Volume: MOB_SOUNDS_VOL},
	MusicSystem{Volume: MUSIC_VOL},
	AmbientNoiseSystem{Volume: AMBIENT_NOISE_VOL},

	coldbrew_clientsystems.DialogueSoundSystem{Volume: DIALOGUE_PING_VOL},
	PlayerAnimationSystem{},
	&CameraFollowerSystem{},
	&coldbrew_clientsystems.BackgroundScrollSystem{},
	coldbrew_clientsystems.DialogueTextSystem{
		TEXT_REVEAL_DELAY_IN_TICKS: dialoguedata.TEXT_REVEAL_DELAY_IN_TICKS,
		FONT_SIZE:                  fontdata.DEFAULT_FONT_SIZE,
		FONT_FACE:                  fontdata.DEFAULT_FONT_FACE,
		MAX_LINE_WIDTH:             fontdata.MAX_LINE_WIDTH,
	},
	DialogueAutoStepperSystem{},

	&DumbSingletonOptionalBossMusicChangeSystem{},
}
