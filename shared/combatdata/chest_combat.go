package combatdata

import (
	"log"

	"github.com/TheBitDrifter/bappa/combat"
)

var chestHurtBoxes = func() combat.HurtBoxes {
	path := "chest_hurtboxes.json"
	boxes, err := combat.NewHurtBoxesFromJSON(CombatFS, path, "../shared/combatdata/")
	if err != nil {
		log.Fatal(err)
	}

	return boxes
}()
