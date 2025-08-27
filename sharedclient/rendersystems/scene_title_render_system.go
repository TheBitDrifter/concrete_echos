package rendersystems

import (
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type SceneTitleRenderSystem struct {
	FONT_FACE           *text.GoTextFace
	TicksPerCharacter   int
	HoldDurationInTicks int
}

func (sys SceneTitleRenderSystem) validate() {
	if sys.FONT_FACE == nil {
		panic("missing font face for SceneTitleRenderSystem")
	}
	if sys.TicksPerCharacter < 0 {
		panic("TicksPerCharacter cannot be negative")
	}
}

func (sys SceneTitleRenderSystem) Render(scene coldbrew.Scene, screen coldbrew.Screen, cli coldbrew.LocalClient) {
	sys.validate()

	var startTime int
	if sys.checkForActiveNotif(scene) {
		startTime = scene.CurrentTick()
	} else {
		startTime = scene.LastSelectedTick()
	}

	query := warehouse.Factory.NewQuery().And(components.SceneTitleComponent)
	cursor := scene.NewCursor(query)

	cameras := cli.ActiveCamerasFor(scene)

	for _, c := range cameras {
		for range cursor.Next() {
			title := components.SceneTitleComponent.GetFromCursor(cursor)

			currentTime := scene.CurrentTick()
			elapsedTicks := currentTime - startTime

			revealDuration := len(title.Value) * sys.TicksPerCharacter
			totalVisibleDuration := revealDuration + sys.HoldDurationInTicks

			if elapsedTicks < 0 {
				continue
			}

			if elapsedTicks >= totalVisibleDuration {
				continue
			}
			fullText := title.Value
			var textToDisplay string

			if sys.TicksPerCharacter == 0 {
				textToDisplay = fullText
			} else {
				charsToShow := int(elapsedTicks / sys.TicksPerCharacter)
				if charsToShow > len(fullText) {
					charsToShow = len(fullText)
				}
				textToDisplay = fullText[:charsToShow]
			}

			if textToDisplay != "" {
				textOpts := &text.DrawOptions{}
				lineSpacing := sys.FONT_FACE.Size + 2
				textOpts.LineSpacing = lineSpacing

				// Measure the text to get its dimensions
				textWidth, textHeight := text.Measure(textToDisplay, sys.FONT_FACE, lineSpacing)

				// Get the screen dimensions
				screenWidth, screenHeight := coldbrew.ClientConfig.Resolution()

				// Calculate the centered position
				finalX := (float64(screenWidth) - textWidth) / 2
				finalY := ((float64(screenHeight) - textHeight) / 2) - 80

				pos := vector.Two{X: finalX, Y: finalY}

				c.DrawTextStatic(
					textToDisplay,
					textOpts,
					sys.FONT_FACE,
					pos,
				)
			}
		}

		c.PresentToScreen(screen, coldbrew.ClientConfig.CameraBorderSize())
	}
}

func (sys SceneTitleRenderSystem) checkForActiveNotif(scene coldbrew.Scene) bool {
	query := warehouse.Factory.NewQuery().And(components.SimpleNotificationComponent)
	cursor := scene.NewCursor(query)
	for range cursor.Next() {
		return true
	}
	return false
}
