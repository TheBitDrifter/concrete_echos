package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/scenes"
)

const DROP_DELAY_TICKS = 60

type DropSystem struct{}

func (s DropSystem) Run(scene blueprint.Scene, dt float64) error {
	query := warehouse.Factory.NewQuery().And(
		combat.Components.Defeat,
		components.DropComponent,
	)
	cursor := scene.NewCursor(query)

	playerQuery := warehouse.Factory.NewQuery().And(components.PlayerTag)
	playerCursor := scene.NewCursor(playerQuery)

	var callback func(warehouse.Storage) error

	for range cursor.Next() {
		drop := components.DropComponent.GetFromCursor(cursor)
		if drop.Opened {
			continue
		}

		defeat := combat.Components.Defeat.GetFromCursor(cursor)
		if scene.CurrentTick() != defeat.StartTick+DROP_DELAY_TICKS {
			continue
		}
		drop.TickDropped = scene.CurrentTick()
		drop.Opened = true
		for range playerCursor.Next() {
			wallet := components.WalletComponent.GetFromCursor(playerCursor)
			wallet.Money += drop.MoneyDrop
			health := combat.Components.Health.GetFromCursor(cursor)
			health.Value += drop.HealthDrop
			if health.Value > 10 {
				health.Value = 10
			}
		}

		callback = scenes.ChestCallbacksRegistry[scenes.ChestCallbackEnum(drop.CustomSpawnCallbackKey)]
	}

	if callback != nil {
		err := callback(scene.Storage())
		if err != nil {
			return err
		}
	}

	return nil
}
