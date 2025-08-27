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

type TeleportRenderSystem struct{}

func (TeleportRenderSystem) Render(scene coldbrew.Scene, screen coldbrew.Screen, c coldbrew.LocalClient) {
	const cd = 240
	query := warehouse.Factory.NewQuery().And(components.TeleportSwapComponent)
	cursor := scene.NewCursor(query)

	// This is pretty yolo but its sorta contained in this scope so effit
	var bundleRef *client.SpriteBundle
	var sprRef coldbrew.Sprite

	for _, cam := range c.ActiveCamerasFor(scene) {
		for range cursor.Next() {

			bundle := client.Components.SpriteBundle.GetFromCursor(cursor)
			bundleRef = bundle

			sprites := coldbrew.MaterializeSprites(bundle)
			sprRef = sprites[0]

			teleState := components.TeleportSwapComponent.GetFromCursor(cursor)
			if !teleState.HasTarget {
				continue
			}

			if scene.CurrentTick() < teleState.StartTick+cd {
				continue
			}

			if !teleState.ActiveTarget.Valid() {
				return
			}

			if combat.Components.Defeat.Check(teleState.ActiveTarget.Table()) {
				return
			}

			pos := spatial.Components.Position.GetFromEntity(teleState.ActiveTarget)

			coldbrew_rendersystems.RenderSpriteSheetAnimation(
				sprites[0],
				&bundle.Blueprints[0],
				11,
				vector.Two{X: pos.X, Y: pos.Y},
				0,
				vector.Two{X: 1, Y: 1},
				spatial.NewDirectionRight(),
				vector.Two{X: -77, Y: -101},
				false,
				cam,
				scene.CurrentTick(),
				nil,
				nil,
			)

		}
		query = warehouse.Factory.NewQuery().And(components.SwapVulnerableComponent)
		cursor = scene.NewCursor(query)

		for range cursor.Next() {
			pos := spatial.Components.Position.GetFromCursor(cursor)
			coldbrew_rendersystems.RenderSpriteSheetAnimation(
				sprRef,
				&bundleRef.Blueprints[0],
				12,
				vector.Two{X: pos.X, Y: pos.Y},
				0,
				vector.Two{X: 1, Y: 1},
				spatial.NewDirectionRight(),
				vector.Two{X: -77, Y: -101},
				false,
				cam,
				scene.CurrentTick(),
				nil,
				nil,
			)
		}

		cam.PresentToScreen(screen, coldbrew.ClientConfig.CameraBorderSize())

	}
}
