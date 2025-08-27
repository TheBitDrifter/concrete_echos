package clientsystems

import (
	"log"

	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type DebugPositionSystem struct{}

func (sys DebugPositionSystem) Run(lc coldbrew.LocalClient, scene coldbrew.Scene) error {
	query := warehouse.Factory.NewQuery().And(input.Components.ActionBuffer, spatial.Components.Position)
	cursor := scene.NewCursor(query)

	i := 1
	for range cursor.Next() {
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			pos := spatial.Components.Position.GetFromCursor(cursor)
			log.Println("Position", pos, i)
		}
		i++
	}
	return nil
}
