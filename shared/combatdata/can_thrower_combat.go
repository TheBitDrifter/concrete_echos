package combatdata

import (
	"log"

	"github.com/TheBitDrifter/bappa/combat"
)

var canThrowerAttacks = func() characterAttacks {
	pathPrimarySeq := "can_thrower_primary_attack.json"
	attacksPrimary, err := combat.NewAttacksFromJSON(CombatFS, pathPrimarySeq, "../shared/combatdata/")
	if err != nil {
		log.Fatal(err)
	}
	atkSeqPrimary := combat.NewAttackSeq(attacksPrimary...)

	return characterAttacks{
		PrimarySeq: *atkSeqPrimary,
	}
}()

var canThrowerHurtBoxes = func() combat.HurtBoxes {
	path := "can_thrower_hurtboxes.json"
	boxes, err := combat.NewHurtBoxesFromJSON(CombatFS, path, "../shared/combatdata/")
	if err != nil {
		log.Fatal(err)
	}

	return boxes
}()
