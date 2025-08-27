package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/scenes"
)

type DefeatSystem struct{}

func (DefeatSystem) Run(scene blueprint.Scene, dt float64) error {
	query := warehouse.Factory.NewQuery().And(combat.Components.Health)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		healthData := combat.Components.Health.GetFromCursor(cursor)
		if healthData.Value <= 0 {
			en, err := cursor.CurrentEntity()
			if err != nil {
				return err
			}

			if !combat.Components.Defeat.CheckCursor(cursor) {
				en.EnqueueAddComponentWithValue(combat.Components.Defeat, combat.Defeat{StartTick: scene.CurrentTick()})
			}
		}
	}

	spawnBossDefeat := false
	query = warehouse.Factory.NewQuery().And(combat.Components.Defeat, components.IsBossTag)
	cursor = scene.NewCursor(query)

	for range cursor.Next() {
		defeatTick := combat.Components.Defeat.GetFromCursor(cursor).StartTick
		if scene.CurrentTick() == defeatTick+120 {
			spawnBossDefeat = true
		}
	}

	if spawnBossDefeat {
		err := scenes.NewBossDefeat(scene.Storage(), scene.CurrentTick())
		if err != nil {
			return err
		}
	}

	query = warehouse.Factory.NewQuery().And(components.BossDefeatedComponent)
	cursor = scene.NewCursor(query)

	for range cursor.Next() {
		defeatTick := components.BossDefeatedComponent.GetFromCursor(cursor).StartTick
		if scene.CurrentTick() == defeatTick+120 {

			en, err := cursor.CurrentEntity()
			if err != nil {
				return err
			}
			scene.Storage().EnqueueDestroyEntities(en)
		}
	}
	return nil
}
