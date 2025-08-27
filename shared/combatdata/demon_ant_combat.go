package combatdata

import (
	"log"

	"github.com/TheBitDrifter/bappa/combat"
)

var demonAntAttacks = func() characterAttacks {
	pathPrimarySeq := "demon_ant_primary_attack.json"
	attacksPrimary, err := combat.NewAttacksFromJSON(CombatFS, pathPrimarySeq, "../shared/combatdata/")
	if err != nil {
		log.Fatal(err)
	}
	atkSeqPrimary := combat.NewAttackSeq(attacksPrimary...)

	pathSecondary := "demon_ant_secondary_attack.json"
	attacksSecondary, err := combat.NewAttacksFromJSON(CombatFS, pathSecondary, "../shared/combatdata/")
	if err != nil {
		log.Fatal(err)
	}

	atkSeqSecondary := combat.NewAttackSeq(attacksSecondary...)

	return characterAttacks{
		PrimarySeq:   *atkSeqPrimary,
		SecondarySeq: *atkSeqSecondary,
	}
}()

var demonAntHurtBoxes = func() combat.HurtBoxes {
	path := "demon_ant_hurtboxes.json"
	boxes, err := combat.NewHurtBoxesFromJSON(CombatFS, path, "../shared/combatdata/")
	if err != nil {
		log.Fatal(err)
	}

	return boxes
}()
