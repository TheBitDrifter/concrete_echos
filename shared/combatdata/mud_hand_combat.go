package combatdata

import (
	"log"

	"github.com/TheBitDrifter/bappa/combat"
)

var mudHandHurtboxes = func() combat.HurtBoxes {
	path := "mud_hand_hurtboxes.json"
	boxes, err := combat.NewHurtBoxesFromJSON(CombatFS, path, "../shared/combatdata/")
	if err != nil {
		log.Fatal(err)
	}

	return boxes
}()
