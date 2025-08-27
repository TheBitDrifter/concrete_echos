package rendersystems

import (
	"log"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/coldbrew/coldbrew_rendersystems"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/persistence"
	"github.com/TheBitDrifter/concrete_echos/sharedclient/fontdata"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type FastTravelSceneRenderSystem struct {
	CheckPointCount int
	ButtonPositions []spatial.Position
	ButtonShapes    []spatial.Shape
	PaddingX        float64
	PaddingY        float64
	SpacingY        float64
	ButtonSqSize    float64
	FONT_FACE       *text.GoTextFace
}

func (sys FastTravelSceneRenderSystem) Render(scene coldbrew.Scene, screen coldbrew.Screen, c coldbrew.LocalClient) {
	count := 0
	for range persistence.State.FtCheckpoints {
		count++
	}
	if count != sys.CheckPointCount {
		sys.ButtonPositions = []spatial.Position{}
		sys.ButtonShapes = []spatial.Shape{}
		i := 0
		for range persistence.State.FtCheckpoints {
			sys.ButtonPositions = append(sys.ButtonPositions, spatial.NewPosition(sys.PaddingX, sys.PaddingY+(float64(i)*sys.SpacingY)))
			sys.ButtonShapes = append(sys.ButtonShapes, spatial.NewRectangle(sys.ButtonSqSize, sys.ButtonSqSize))
			i++
		}
	}

	for _, cam := range c.ActiveCamerasFor(scene) {
		query := warehouse.Factory.NewQuery().And(client.Components.SpriteBundle)
		cursor := scene.NewCursor(query)

		var spr coldbrew.Sprite

		for range cursor.Next() {
			bundle := client.Components.SpriteBundle.GetFromCursor(cursor)
			if components.IsFastTravelButtonTag.CheckCursor(cursor) {
				found, err := coldbrew.MaterializeSprite(bundle, 0)
				spr = found
				if err != nil {
					log.Println(err)
					return
				}
				break
			}
		}

		if spr == nil {
			return
		}

		i := 0.0
		for _, checkPoint := range persistence.State.FtCheckpoints {

			textOpts := &text.DrawOptions{}
			lineSpacing := sys.FONT_FACE.Size + 2
			textOpts.LineSpacing = lineSpacing
			textToDisplay := string(checkPoint.FastTravelCheckpointName)

			textWidth, textHeight := text.Measure(textToDisplay, sys.FONT_FACE, lineSpacing)
			textPos := vector.Two{X: sys.PaddingX, Y: sys.PaddingY + (i * sys.SpacingY)}

			cam.DrawTextStatic(
				textToDisplay,
				textOpts,
				sys.FONT_FACE,
				textPos,
			)

			coldbrew_rendersystems.RenderSprite(
				spr,
				vector.Two{X: textPos.X + textWidth + 10, Y: textPos.Y - textHeight*3},
				0,
				vector.Two{X: 1, Y: 1},
				vector.Two{X: -32, Y: -32},
				spatial.NewDirectionRight(),
				true,
				cam,
			)
			i++
		}

		textOpts := &text.DrawOptions{}
		cam.DrawTextStatic(
			"Crappy Fast Travel System!",
			textOpts,
			fontdata.TITLE_FONT_FACE,
			vector.Two{X: 180, Y: 20},
		)

		cam.PresentToScreen(screen, coldbrew.ClientConfig.CameraBorderSize())
	}
}
