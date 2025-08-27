package dialoguedata

import (
	"embed"

	"github.com/TheBitDrifter/bappa/blueprint/dialogue"
)

//go:embed *
var DialogueFS embed.FS

const TEXT_REVEAL_DELAY_IN_TICKS = 2

const TEXT_REVEAL_DELAY_IN_TICKS_ALT = 7

type vec2Int struct {
	X int
	Y int
}

// TinyGo Shenanigans(cant assert a global variable from int/float directly)
func (v vec2Int) AsFloats() (x, y float64) {
	vintX := v.X
	vintY := v.Y

	return float64(vintX), float64(vintY)
}

var (
	DEFAULT_BOX_POS               = vec2Int{X: 120, Y: 270}
	DEFAULT_TEXT_PADDING          = vec2Int{X: 70, Y: 19}
	DEFAULT_PORTRAIT_PADDING      = vec2Int{X: 0, Y: 11}
	DEFAULT_PORTRAIT_TEXT_PADDING = vec2Int{X: 22, Y: 1}
)

// Slides Registration ---------------------------------------------------------------------------------------------
const (
	_ dialogue.SlidesEnum = iota
	IntroCutsceneDialogue
	DrifterFirstMeetingDialogue
	PeacefulEchoesEncounter
)

var _ = func() error {
	dialogue.SlidesRegistry = map[dialogue.SlidesEnum]dialogue.Slides{}
	dialogue.SlidesRegistry[IntroCutsceneDialogue] = introCutsceneDia
	dialogue.SlidesRegistry[DrifterFirstMeetingDialogue] = drifterFirstMeetingDia
	dialogue.SlidesRegistry[PeacefulEchoesEncounter] = peacefulEchosEncounter
	return nil
}()

// Slides Registration End ---------------------------------------------------------------------------------------------

// Portraits && Names ---------------------------------------------------------------------------------------------
const (
	_ dialogue.PortraitEnum = iota
	BoxHead
	Drifter
	LazySkully
)

const (
	BoxHeadName    = "Box Head"
	DrifterName    = "Smokey Guy"
	LazySkullyName = "Lazy Skully"
)

const (
	LazySkullySpeed = 6
)

// Portraits && Names END ---------------------------------------------------------------------------------------------

// SLIDES DATA ---------------------------------------------------------------------------------------------
var introCutsceneDia = dialogue.Slides{
	dialogue.Slide{
		OwnerName:  BoxHeadName,
		PortraitID: BoxHead,
		Text:       "There it is again.",
	},
	dialogue.Slide{
		OwnerName:  BoxHeadName,
		PortraitID: BoxHead,
		Text:       "Same feeling as back then.",
	},
	dialogue.Slide{
		OwnerName:  BoxHeadName,
		PortraitID: BoxHead,
		Text:       "It's just a picture in my head. What am I supposed to do, just start walking?",
	},
}

var drifterFirstMeetingDia = dialogue.Slides{
	dialogue.Slide{
		OwnerName:  DrifterName,
		PortraitID: Drifter,
		Text:       "Alright, I gotta ask. What's the deal with the box?",
	},
	dialogue.Slide{
		OwnerName:  BoxHeadName,
		PortraitID: BoxHead,
		Text:       "It's a filter. Helps me focus.",
	},
	dialogue.Slide{
		OwnerName:  DrifterName,
		PortraitID: Drifter,
		Text:       "A filter. For what? The smog? The government mind-control rays?",
	},
	dialogue.Slide{
		OwnerName:  BoxHeadName,
		PortraitID: BoxHead,
		Text:       "The background noise. From... people. It gets loud.",
	},
	dialogue.Slide{
		OwnerName:  DrifterName,
		PortraitID: Drifter,
		Text:       "Speaking of weird... You alright there? You're kinda... squatting.",
	},
	dialogue.Slide{
		OwnerName:  BoxHeadName,
		PortraitID: BoxHead,
		Text:       "*sigh*... Right. Of course that's what it looks like.",
	},
	dialogue.Slide{
		OwnerName:  BoxHeadName,
		PortraitID: BoxHead,
		Text:       "Well, this was interesting. See you around.",
	},
	dialogue.Slide{
		OwnerName:  DrifterName,
		PortraitID: Drifter,
		Text:       "Yeah... you too. Stay weird, Box Head.",
	},
}

var peacefulEchosEncounter = dialogue.Slides{
	dialogue.Slide{
		OwnerName:  BoxHeadName,
		PortraitID: BoxHead,
		Text:       "You're not attacking me. That's new.",
	},
	dialogue.Slide{
		OwnerName:        LazySkullyName,
		PortraitID:       LazySkully,
		Text:             "You're not... dangerous... so I'm not... bothered.",
		CustomSpeedTicks: LazySkullySpeed,
	},
	dialogue.Slide{
		OwnerName:  BoxHeadName,
		PortraitID: BoxHead,
		Text:       "'Dangerous?' I've been fighting this whole time.",
	},
	dialogue.Slide{
		OwnerName:        LazySkullyName,
		PortraitID:       LazySkully,
		Text:             "Fighting is just... fighting... Taking what's left... that's different.",
		CustomSpeedTicks: LazySkullySpeed,
	},
	dialogue.Slide{
		OwnerName:  BoxHeadName,
		PortraitID: BoxHead,
		Text:       "You mean absorbing? I thought I was just cleaning up.",
	},
	dialogue.Slide{
		OwnerName:        LazySkullyName,
		PortraitID:       LazySkully,
		Text:             "What you take... sticks to you... It makes you... brighter.",
		CustomSpeedTicks: LazySkullySpeed,
	},
	dialogue.Slide{
		OwnerName:        LazySkullyName,
		PortraitID:       LazySkully,
		Text:             "And the really hungry things... are drawn to the brightest lights... zzzzzz...",
		CustomSpeedTicks: LazySkullySpeed,
	},
	dialogue.Slide{
		OwnerName:  BoxHeadName,
		PortraitID: BoxHead,
		Text:       "...",
	},
}
