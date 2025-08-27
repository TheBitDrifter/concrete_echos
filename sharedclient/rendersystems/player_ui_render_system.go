package rendersystems

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/coldbrew/coldbrew_rendersystems"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

type PlayerUIRenderSystem struct{}

func (PlayerUIRenderSystem) Render(scene coldbrew.Scene, screen coldbrew.Screen, c coldbrew.LocalClient) {
	query := warehouse.Factory.NewQuery().And(components.PlayerTag)
	cursor := scene.NewCursor(query)

	for _, cam := range c.ActiveCamerasFor(scene) {
		for range cursor.Next() {

			bundle := client.Components.SpriteBundle.GetFromCursor(cursor)
			hudSpr := coldbrew.MaterializeSprites(bundle)[PlayerIndexPortraitHUD]

			heartSpr := coldbrew.MaterializeSprites(bundle)[PlayerIndexHearts]

			coldbrew_rendersystems.RenderSprite(
				hudSpr,
				vector.Two{X: 10, Y: 10},
				0,
				vector.Two{X: 1, Y: 1},
				vector.Two{X: 0, Y: 0},
				spatial.NewDirectionRight(),
				true,
				cam,
			)

			health := combat.Components.Health.GetFromCursor(cursor)

			heartCount := int(health.Value / 10)

			for i := 0; i < heartCount; i++ {
				x := float64(75 + (i * 22))
				y := 35.0

				coldbrew_rendersystems.RenderSprite(
					heartSpr,
					vector.Two{X: x, Y: y},
					0,
					vector.Two{X: 1, Y: 1},
					vector.Two{X: 0, Y: 0},
					spatial.NewDirectionRight(),
					true,
					cam,
				)
			}

			cam.PresentToScreen(screen, coldbrew.ClientConfig.CameraBorderSize())
		}
	}
}
