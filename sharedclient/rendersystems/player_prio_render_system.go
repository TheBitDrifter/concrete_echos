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

type PlayerRenderer struct{}

func (PlayerRenderer) Render(scene coldbrew.Scene, screen coldbrew.Screen, c coldbrew.LocalClient) {
	netCli, isNet := c.(coldbrew.NetworkClient)
	if isNet {
		id, ok := netCli.AssociatedEntityID()
		if !ok {
			return
		}

		pEn, err := scene.Storage().Entity(id)
		if err != nil {
			return
		}
		for _, cam := range netCli.ActiveCamerasFor(scene) {
			// If it ain't ready chill out!
			if !netCli.Ready(cam) {
				continue
			}
			bundle := client.Components.SpriteBundle.GetFromEntity(pEn)
			spr := coldbrew.MaterializeSprites(bundle)[0]
			direction := *spatial.Components.Direction.GetFromEntity(pEn)

			if combat.Components.Attack.Check(pEn.Table()) {
				atk := combat.Components.Attack.GetFromEntity(pEn)
				direction = atk.LRDirection
			}

			coldbrew_rendersystems.RenderEntity(
				spatial.Components.Position.GetFromEntity(pEn).Two,
				0,
				vector.Two{1, 1},
				direction,
				spr,
				&bundle.Blueprints[0],
				cam,
				scene.CurrentTick(),
			)
			cam.PresentToScreen(screen, coldbrew.ClientConfig.CameraBorderSize())
		}

	} else {
		query := warehouse.Factory.NewQuery().And(components.PlayerTag)
		cursor := scene.NewCursor(query)

		for range cursor.Next() {
			for _, cam := range c.ActiveCamerasFor(scene) {
				if !c.Ready(cam) {
					continue
				}
				bundle := client.Components.SpriteBundle.GetFromCursor(cursor)
				spr := coldbrew.MaterializeSprites(bundle)[0]
				direction := *spatial.Components.Direction.GetFromCursor(cursor)

				if atk, atkOK := combat.Components.Attack.GetFromCursorSafe(cursor); atkOK {
					direction = atk.LRDirection
				}

				coldbrew_rendersystems.RenderEntity(
					spatial.Components.Position.GetFromCursor(cursor).Two,
					0,
					vector.Two{1, 1},
					direction,
					spr,
					&bundle.Blueprints[0],
					cam,
					scene.CurrentTick(),
				)
				cam.PresentToScreen(screen, coldbrew.ClientConfig.CameraBorderSize())

			}
		}
	}
}
