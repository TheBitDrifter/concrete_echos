package combatdata

import (
	"log"

	"github.com/TheBitDrifter/bappa/combat"
)

var skullRoninAttacks = func() characterAttacks {
	pathPrimarySeq := "skull_ronin_primary_attack_seq.json"
	attacksPrimary, err := combat.NewAttacksFromJSON(CombatFS, pathPrimarySeq, "../shared/combatdata/")
	if err != nil {
		log.Fatal(err)
	}
	atkSeqPrimary := combat.NewAttackSeq(attacksPrimary...)

	return characterAttacks{
		PrimarySeq: *atkSeqPrimary,
	}
}()

var skullRoninHurtBoxes = func() combat.HurtBoxes {
	path := "skull_ronin_hurtboxes.json"
	boxes, err := combat.NewHurtBoxesFromJSON(CombatFS, path, "../shared/combatdata/")
	if err != nil {
		log.Fatal(err)
	}

	return boxes
}()
