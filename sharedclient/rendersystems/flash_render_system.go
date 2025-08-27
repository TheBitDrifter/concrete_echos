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
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

type FlashRenderSystem struct{}

func (FlashRenderSystem) Render(scene coldbrew.Scene, screen coldbrew.Screen, c coldbrew.LocalClient) {
	tick := scene.CurrentTick()

	for _, cam := range c.ActiveCamerasFor(scene) {
		queryDodge := warehouse.Factory.NewQuery().And(client.Components.SpriteBundle, components.DodgeComponent)
		cursorDodge := scene.NewCursor(queryDodge)
		for range cursorDodge.Next() {
			dodge := components.DodgeComponent.GetFromCursor(cursorDodge)

			bundle := client.Components.SpriteBundle.GetFromCursor(cursorDodge)
			spr := coldbrew.MaterializeSprites(bundle)[0]
			pos := spatial.Components.Position.GetFromCursor(cursorDodge)
			dir := spatial.Components.Direction.GetFromCursor(cursorDodge)

			flashPurple := (tick-dodge.StartTick)%4 == 0
			if !flashPurple {
				continue
			}

			var purpMatrix colorm.ColorM
			purpMatrix.Scale(0, 0, 0, 1)
			purpMatrix.SetElement(0, 4, 1)
			purpMatrix.SetElement(2, 4, 1)

			coldbrew_rendersystems.RenderSpriteSheetAnimation(
				spr,
				&bundle.Blueprints[0],
				bundle.Blueprints[0].Config.ActiveAnimIndex,
				pos.Two,
				0,
				vector.Two{X: 1, Y: 1},
				*dir,
				bundle.Blueprints[0].Config.Offset,
				false,
				cam,
				tick,
				nil,
				&purpMatrix,
			)
		}

		queryHurt := warehouse.Factory.NewQuery().And(client.Components.SpriteBundle, combat.Components.Hurt)
		cursorHurt := scene.NewCursor(queryHurt)
		for range cursorHurt.Next() {
			hurt := combat.Components.Hurt.GetFromCursor(cursorHurt)
			hurtTick := hurt.StartTick

			bundle := client.Components.SpriteBundle.GetFromCursor(cursorHurt)
			spr := coldbrew.MaterializeSprites(bundle)[0]
			pos := spatial.Components.Position.GetFromCursor(cursorHurt)
			dir := spatial.Components.Direction.GetFromCursor(cursorHurt)

			flashRed := (tick-hurtTick)%8 == 0
			if !flashRed {
				continue
			}

			var redMatrix colorm.ColorM
			redMatrix.Scale(0, 0, 0, 1)
			redMatrix.SetElement(0, 4, 1)

			coldbrew_rendersystems.RenderSpriteSheetAnimation(
				spr,
				&bundle.Blueprints[0],
				bundle.Blueprints[0].Config.ActiveAnimIndex,
				pos.Two,
				0,
				vector.Two{X: 1, Y: 1},
				*dir,
				bundle.Blueprints[0].Config.Offset,
				false,
				cam,
				tick,
				nil,
				&redMatrix,
			)
		}

		queryInvincible := warehouse.Factory.NewQuery().And(
			client.Components.SpriteBundle, combat.Components.Invincible,
			warehouse.Factory.NewQuery().Not(combat.Components.Hurt, components.DodgeComponent),
		)
		cursorInvincible := scene.NewCursor(queryInvincible)
		for range cursorInvincible.Next() {
			invincible := combat.Components.Invincible.GetFromCursor(cursorInvincible)
			invincibleTick := invincible.StartTick

			bundle := client.Components.SpriteBundle.GetFromCursor(cursorInvincible)
			spr := coldbrew.MaterializeSprites(bundle)[0]
			pos := spatial.Components.Position.GetFromCursor(cursorInvincible)
			dir := spatial.Components.Direction.GetFromCursor(cursorInvincible)

			flashWhite := (tick-invincibleTick)%8 == 0
			if !flashWhite {
				continue
			}

			var whiteMatrix colorm.ColorM
			whiteMatrix.Scale(0, 0, 0, 1)
			whiteMatrix.SetElement(0, 4, 1)
			whiteMatrix.SetElement(1, 4, 1)
			whiteMatrix.SetElement(2, 4, 1)

			coldbrew_rendersystems.RenderSpriteSheetAnimation(
				spr,
				&bundle.Blueprints[0],
				bundle.Blueprints[0].Config.ActiveAnimIndex,
				pos.Two,
				0,
				vector.Two{X: 1, Y: 1},
				*dir,
				bundle.Blueprints[0].Config.Offset,
				false,
				cam,
				tick,
				nil,
				&whiteMatrix,
			)
		}
		// -----
		querySwapVurn := warehouse.Factory.NewQuery().And(
			components.SwapVulnerableComponent,
			warehouse.Factory.NewQuery().Not(components.WarpTotemComponent),
		)
		cursorSwapVurn := scene.NewCursor(querySwapVurn)
		for range cursorSwapVurn.Next() {

			bundle := client.Components.SpriteBundle.GetFromCursor(cursorSwapVurn)
			spr := coldbrew.MaterializeSprites(bundle)[0]
			pos := spatial.Components.Position.GetFromCursor(cursorSwapVurn)
			dir := spatial.Components.Direction.GetFromCursor(cursorSwapVurn)

			flash := (tick)%4 == 0
			if !flash {
				continue
			}

			var yellowMatrix colorm.ColorM
			yellowMatrix.Scale(0, 0, 0, 1)
			yellowMatrix.SetElement(0, 4, 1)
			yellowMatrix.SetElement(1, 4, 1)

			coldbrew_rendersystems.RenderSpriteSheetAnimation(
				spr,
				&bundle.Blueprints[0],
				bundle.Blueprints[0].Config.ActiveAnimIndex,
				pos.Two,
				0,
				vector.Two{X: 1, Y: 1},
				*dir,
				bundle.Blueprints[0].Config.Offset,
				false,
				cam,
				tick,
				nil,
				&yellowMatrix,
			)

		}

		cam.PresentToScreen(screen, coldbrew.ClientConfig.CameraBorderSize())
	}
}
