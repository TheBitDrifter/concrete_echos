package rendersystems

import (
	"fmt"

	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/actions"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type KeyRebindingRenderSystem struct {
	FONT_FACE *text.GoTextFace
}

func (sys KeyRebindingRenderSystem) validate() {
	if sys.FONT_FACE == nil {
		panic("missing font face for KeyRebindingRenderSystem")
	}
}

func (sys KeyRebindingRenderSystem) Render(scene coldbrew.Scene, screen coldbrew.Screen, cli coldbrew.LocalClient) {
	sys.validate()

	query := warehouse.Factory.NewQuery().And(components.KeyRebindingStateComponent)
	cursor := scene.NewCursor(query)
	var state *components.KeyRebindingState

	if cursor.TotalMatched() == 0 {
		return
	}

	for range cursor.Next() {
		state = components.KeyRebindingStateComponent.GetFromCursor(cursor)
	}

	if state == nil || !state.Active {
		return
	}

	cameras := cli.ActiveCamerasFor(scene)
	for _, c := range cameras {
		prompt := fmt.Sprintf("Press a button for: %s", actionToString(state.CurrentAction))

		textOpts := &text.DrawOptions{}
		lineSpacing := sys.FONT_FACE.Size + 2
		textOpts.LineSpacing = lineSpacing

		textWidth, textHeight := text.Measure(prompt, sys.FONT_FACE, lineSpacing)
		screenWidth, screenHeight := coldbrew.ClientConfig.Resolution()

		finalX := (float64(screenWidth) - textWidth) / 2
		finalY := (float64(screenHeight) - textHeight) / 2

		pos := vector.Two{X: finalX, Y: finalY}

		c.DrawTextStatic(
			prompt,
			textOpts,
			sys.FONT_FACE,
			pos,
		)
		c.PresentToScreen(screen, coldbrew.ClientConfig.CameraBorderSize())
	}
}

func actionToString(a input.Action) string {
	switch a {
	case actions.Jump:
		return "Jump"
	case actions.PrimaryAttack:
		return "Primary Attack"
	case actions.Dodge:
		return "Dodge"
	case actions.Interact:
		return "Interact"
	case actions.Cancel:
		return "Cancel"
	case actions.TeleSwap:
		return "Teleport Swap(if unlocked)"
	case actions.ShiftTeleTargetRight:
		return "Shift Teleport Target Right(if unlocked)"
	case actions.ShiftTeleTargetLeft:
		return "Shift Teleport Target Left(if unlocked)"
	case actions.ShiftTeleTargetNear:
		return "Shift Teleport Target Near(if unlocked)"
	default:
		return fmt.Sprintf("Action %d", a)
	}
}
