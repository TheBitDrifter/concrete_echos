package callbacks

import (
	"log"
	"math"

	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/blueprint/dialogue"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/scenes"
)

const (
	_ dialogue.CallbackEnum = iota
	PeacefulEchoesEncounterCallback
)

var _ = func() error {
	dialogue.CallbackRegistry[PeacefulEchoesEncounterCallback] = peacefulEchoesEncounterCallback
	return nil
}()

var peacefulEchoesEncounterCallback = func(scene blueprint.Scene) error {
	log.Println("firing peaceful callback yo!")
	query := warehouse.Factory.NewQuery().And(components.TrapDoorComponent)
	cursor := scene.NewCursor(query)

	execCountQuery := warehouse.Factory.NewQuery().And(components.PlayerExecutionCountComponent)
	execCursor := scene.NewCursor(execCountQuery)

	highCount := 0.0

	for range execCursor.Next() {
		counter := components.PlayerExecutionCountComponent.GetFromCursor(execCursor)
		highCount = math.Max(float64(counter.Count), highCount)
	}

	foundSomeEvilDudePlayer := highCount >= scenes.PLAYER_EVIL_THRESHOLD
	if foundSomeEvilDudePlayer {
		return nil
	}

	for range cursor.Next() {
		td := components.TrapDoorComponent.GetFromCursor(cursor)
		trapDoorMatched := td.IsOpenCallbackID == FinalBldUnlock

		if !trapDoorMatched {
			continue
		}

		td.LastChangedTick = scene.CurrentTick()
		td.Open = true

	}
	return nil
}
