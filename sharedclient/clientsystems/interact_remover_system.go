package clientsystems

import (
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/actions"
)

type InteractRemoverSystem struct{}

func (sys InteractRemoverSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	query := warehouse.Factory.NewQuery().And(input.Components.ActionBuffer)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		actionsBuf := input.Components.ActionBuffer.GetFromCursor(cursor)
		actionsBuf.ConsumeAction(actions.Interact)
	}
	return nil
}
