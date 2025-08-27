package rendersystems

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/coldbrew/coldbrew_rendersystems"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

type InteractionMarkerRenderSystem struct{}

func (sys InteractionMarkerRenderSystem) Render(scene coldbrew.Scene, screen coldbrew.Screen, c coldbrew.LocalClient) {
	queryP := warehouse.Factory.NewQuery().And(components.PlayerTag)
	cursorP := scene.NewCursor(queryP)

	queryF := warehouse.Factory.NewQuery().Or(
		components.DialogueActivationComponent,
		components.SaveActivationComponent,
		components.MobExecuteComponent,
		components.FastTravelActivationComponent,
	)
	cursorF := scene.NewCursor(queryF)

	for _, cam := range c.ActiveCamerasFor(scene) {
		for range cursorP.Next() {
			for range cursorF.Next() {
				if components.InConversationComponent.CheckCursor(cursorF) {
					continue
				}
				if components.IsSavingComponent.CheckCursor(cursorP) {
					continue
				}

				bundle := client.Components.SpriteBundle.GetFromCursor(cursorF)
				sprites := coldbrew.MaterializeSprites(bundle)
				pos := spatial.Components.Position.GetFromCursor(cursorF)
				// dir := spatial.Components.Direction.GetFromCursor(cursorF)
				offset := bundle.Blueprints[1].Config.Offset

				playerPos := spatial.Components.Position.GetFromCursor(cursorP)
				playerDir := spatial.Components.Direction.GetFromCursor(cursorP)
				currentDistSq := pos.Two.Sub(playerPos.Two).MagSquared()

				rangeDist := 20.0

				if saveActive, ok := components.SaveActivationComponent.GetFromCursorSafe(cursorF); ok {
					rangeDist = saveActive.Range
				}

				diaActive, diaOK := components.DialogueActivationComponent.GetFromCursorSafe(cursorF)
				if diaOK {
					rangeDist = diaActive.Range
				}

				ftActive, ftOK := components.FastTravelActivationComponent.GetFromCursorSafe(cursorF)
				if ftOK {
					rangeDist = ftActive.Range
				}

				inRange := currentDistSq <= rangeDist*rangeDist

				if diaOK && inRange {
					minDist := diaActive.MinRange
					inRange = currentDistSq > float64(minDist*minDist)
				}

				if inRange && diaOK && diaActive.MustBeRight {
					inRange = playerPos.X > pos.X
				}

				if inRange && diaOK && diaActive.MustBeLeft {
					inRange = playerPos.X < pos.X
				}

				_, okExec := components.MobExecuteComponent.GetFromCursorSafe(cursorF)
				if okExec {

					playerFacingRightWay := (playerDir.IsLeft() && playerPos.X > pos.X) || (playerDir.IsRight() && playerPos.X < pos.X)
					inRange = currentDistSq <= 45*45 && currentDistSq > 15*15 &&
						(components.LastCombatComponent.GetFromCursor(cursorP).StartTick+120 <= scene.CurrentTick()) &&
						playerFacingRightWay
				}

				if len(sprites) < 2 {
					continue
				}
				if inRange {
					coldbrew_rendersystems.RenderSpriteSheetAnimation(
						sprites[1],
						&bundle.Blueprints[1],
						0,
						vector.Two{X: pos.X, Y: pos.Y},
						0,
						vector.Two{X: 1, Y: 1},
						spatial.NewDirectionRight(),
						offset,
						false,
						cam,
						scene.CurrentTick(),
						nil,
						nil,
					)
				}

			}
		}

		cam.PresentToScreen(screen, coldbrew.ClientConfig.CameraBorderSize())

	}
}
