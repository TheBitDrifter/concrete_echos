package callbacks

import (
	"log"

	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

const (
	_ components.TrapDoorEnum = iota
	OuterEastWarpUnlock
	FinalBldUnlock
)

var TrapDoorCallbacks = map[components.TrapDoorEnum]func(sto warehouse.Storage, width, height float64) bool{
	OuterEastWarpUnlock: outerEastWarpUnlockCallback,
	FinalBldUnlock:      finalBldUnlockCallback,
}

func CheckTrapDoor(sto warehouse.Storage, width, height float64, td components.TrapDoor) bool {
	callback, ok := TrapDoorCallbacks[td.IsOpenCallbackID]
	if !ok {
		log.Println("Warning! No callback", td.IsOpenCallbackID)
		return false
	}
	return callback(sto, width, height)
}

func outerEastWarpUnlockCallback(sto warehouse.Storage, width, height float64) bool {
	query := warehouse.Factory.NewQuery().And(components.MobTag)
	cursor := warehouse.Factory.NewCursor(query, sto)
	for range cursor.Next() {
		if !combat.Components.Defeat.CheckCursor(cursor) && !components.ChestTag.CheckCursor(cursor) {
			return false
		}
	}
	return true
}

func finalBldUnlockCallback(sto warehouse.Storage, width, height float64) bool {
	query := warehouse.Factory.NewQuery().And(components.IsBossTag)
	cursor := warehouse.Factory.NewCursor(query, sto)
	for range cursor.Next() {
		return combat.Components.Defeat.CheckCursor(cursor)
	}
	return false
}
