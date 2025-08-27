package rendersystems

import (
	"fmt"

	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/sharedclient/fontdata"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type MoneyRenderSystem struct{}

func (MoneyRenderSystem) Render(scene coldbrew.Scene, screen coldbrew.Screen, c coldbrew.LocalClient) {
	query := warehouse.Factory.NewQuery().And(components.WalletComponent)
	cursor := scene.NewCursor(query)

	const rightAlignX = 600

	for _, cam := range c.ActiveCamerasFor(scene) {
		for range cursor.Next() {
			wallet := components.WalletComponent.GetFromCursor(cursor)
			textStr := fmt.Sprintf("Cash Money: %.0f", wallet.Money)

			textOpts := &text.DrawOptions{}
			textOpts.LineSpacing = float64(8)

			textWidth, _ := text.Measure(textStr, fontdata.DEFAULT_FONT_FACE, textOpts.LineSpacing)

			newX := rightAlignX - textWidth

			texPos := vector.Two{X: newX, Y: 10}
			cam.DrawTextStatic(
				textStr, // Use the formatted string
				textOpts,
				fontdata.DEFAULT_FONT_FACE,
				texPos,
			)

			cam.PresentToScreen(screen, coldbrew.ClientConfig.CameraBorderSize())
		}
	}
}
