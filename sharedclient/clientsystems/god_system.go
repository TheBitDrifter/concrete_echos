package clientsystems

import (
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/environment"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GodSystem struct{}

func (s GodSystem) Run(lc coldbrew.LocalClient, scene coldbrew.Scene) error {
	if environment.IsProd() {
		return nil
	}
	query := warehouse.Factory.NewQuery().And(components.PlayerTag)
	cursor := scene.NewCursor(query)
	for _, c := range lc.ActiveCamerasFor(scene) {

		_, local := c.Positions()

		for range cursor.Next() {
			if inpututil.IsKeyJustPressed(ebiten.Key9) {
				health := combat.Components.Health.GetFromCursor(cursor)
				health.Value += 10
			}

			if inpututil.IsKeyJustPressed(ebiten.Key8) {
				pos := spatial.Components.Position.GetFromCursor(cursor)
				nx, ny := ebiten.CursorPosition()

				pos.X = float64(nx) + local.X
				pos.Y = float64(ny) + local.Y

			}
		}
	}

	return nil
}
