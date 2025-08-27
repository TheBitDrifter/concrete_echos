package rendersystems

import (
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/coldbrew/coldbrew_rendersystems"
	"github.com/TheBitDrifter/bappa/coldbrew/combat_rendersystems"
	"github.com/TheBitDrifter/concrete_echos/shared/dialoguedata"
	"github.com/TheBitDrifter/concrete_echos/sharedclient/fontdata"
)

var DefaultRenderSystems = []coldbrew.RenderSystem{
	PlayerRenderer{},
	combat_rendersystems.HitBoxRenderSystem{},
	combat_rendersystems.HurtBoxRenderSystem{},
	FlashRenderSystem{},
	NewFogShaderRenderSystem(),
	PlayerUIRenderSystem{},
	coldbrew_rendersystems.DefaultDialogueRenderSystem{
		FONT_FACE:                fontdata.DEFAULT_FONT_FACE,
		PORTRAIT_FONT_FACE:       fontdata.SMALLER_FONT_FACE,
		PADDING_X:                dialoguedata.DEFAULT_TEXT_PADDING.X,
		PADDING_Y:                dialoguedata.DEFAULT_TEXT_PADDING.Y,
		PORTRAIT_PADDING_X:       dialoguedata.DEFAULT_PORTRAIT_PADDING.X,
		PORTRAIT_PADDING_Y:       dialoguedata.DEFAULT_PORTRAIT_PADDING.Y,
		PORTRAIT_TEXT_PADDING_X:  dialoguedata.DEFAULT_PORTRAIT_TEXT_PADDING.X,
		PORTRAIT_TEXT_PADDING_Y:  dialoguedata.DEFAULT_PORTRAIT_TEXT_PADDING.Y,
		PORTRAIT_NAME_BOX_WIDTH:  54,
		NEXT_INDICATOR_PADDING_X: 380,
		NEXT_INDICATOR_PADDING_Y: 63,
		NEXT_MIN_DELAY:           15,
	},
	SceneTitleRenderSystem{
		FONT_FACE:           fontdata.TITLE_FONT_FACE,
		TicksPerCharacter:   10,
		HoldDurationInTicks: 40,
	},
	SimpleNotificationRenderSystem{
		TITLE_FONT_FACE:     fontdata.TITLE_FONT_FACE,
		BODY_FONT_FACE:      fontdata.UNLOCK_BODY_FONT_FACE,
		BODY_TEXT_PADDING_X: 145,
		BODY_TEXT_PADDING_Y: 115,
		REVEAL_START_DELAY:  60,
	},
	TeleportRenderSystem{},
	MoneyRenderSystem{},
	InteractionMarkerRenderSystem{},
	KeyRebindingRenderSystem{
		FONT_FACE: fontdata.DEFAULT_FONT_FACE,
	},
}

var DefaultFastTravelSystems = []coldbrew.RenderSystem{
	&FastTravelSceneRenderSystem{
		PaddingX:     50,
		PaddingY:     100,
		SpacingY:     60,
		ButtonSqSize: 32,
		FONT_FACE:    fontdata.DEFAULT_FONT_FACE,
	},
}

type PlayerSpriteBundleIndex int

const (
	PlayerIndexAnimations PlayerSpriteBundleIndex = iota
	PlayerIndexPortraitHUD
	PlayerIndexHearts
	PlayerIndexDefeatScreen
)
