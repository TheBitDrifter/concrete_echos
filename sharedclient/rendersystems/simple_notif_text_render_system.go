package rendersystems

import (
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type SimpleNotificationRenderSystem struct {
	TITLE_FONT_FACE *text.GoTextFace
	BODY_FONT_FACE  *text.GoTextFace

	BODY_TEXT_PADDING_X int
	BODY_TEXT_PADDING_Y int
	REVEAL_START_DELAY  int
}

func (sys SimpleNotificationRenderSystem) Render(scene coldbrew.Scene, screen coldbrew.Screen, cli coldbrew.LocalClient) {
	query := warehouse.Factory.NewQuery().And(
		components.SimpleNotificationComponent,
	)
	cursor := scene.NewCursor(query)
	cameras := cli.ActiveCamerasFor(scene)

	for _, c := range cameras {
		for range cursor.Next() {

			pos := spatial.Components.Position.GetFromCursor(cursor)
			noti := components.SimpleNotificationComponent.GetFromCursor(cursor)
			if noti.StartedTick+sys.REVEAL_START_DELAY > scene.CurrentTick() {
				continue
			}

			title := noti.DisplayedTitle
			if title != "" {
				textOpts := &text.DrawOptions{}

				textOpts.LineSpacing = float64(sys.TITLE_FONT_FACE.Size)
				texPos := vector.Two{X: pos.X + noti.PaddingX, Y: pos.Y + noti.PaddingY}

				c.DrawTextStatic(
					title,
					textOpts,
					sys.TITLE_FONT_FACE,
					texPos,
				)
			}

			body := noti.DisplayedBody
			if body != "" {
				textOpts := &text.DrawOptions{}

				textOpts.LineSpacing = float64(sys.TITLE_FONT_FACE.Size)
				texPos := vector.Two{X: pos.X + float64(sys.BODY_TEXT_PADDING_X), Y: pos.Y + float64(sys.BODY_TEXT_PADDING_Y)}

				c.DrawTextStatic(
					body,
					textOpts,
					sys.BODY_FONT_FACE,
					texPos,
				)
			}

			textOpts := &text.DrawOptions{}

			textOpts.LineSpacing = float64(sys.BODY_FONT_FACE.Size)

			if noti.IsFinished || true {
				var textToDraw string
				currentTick := scene.CurrentTick()

				cycle := currentTick % 90

				if cycle < 30 {
					textToDraw = "Press Enter To Close ."
				} else if cycle < 60 {
					textToDraw = "Press Enter To Close .."
				} else {
					textToDraw = "Press Enter To Close ..."
				}

				// Show the skip option
				someExtra := 180
				if noti.StartedTick+sys.REVEAL_START_DELAY+someExtra <= scene.CurrentTick() {
					textOpts := &text.DrawOptions{}
					textOpts.LineSpacing = float64(sys.BODY_FONT_FACE.Size)
					texPos := vector.Two{X: pos.X + float64(sys.BODY_TEXT_PADDING_X+230), Y: pos.Y + float64(sys.BODY_TEXT_PADDING_Y) + 116}

					c.DrawTextStatic(
						textToDraw,
						textOpts,
						sys.BODY_FONT_FACE,
						texPos,
					)
				}

			}
		}
		c.PresentToScreen(screen, coldbrew.ClientConfig.CameraBorderSize())
	}
}
