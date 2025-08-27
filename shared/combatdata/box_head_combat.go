package combatdata

import (
	"log"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/concrete_echos/shared/animations"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
)

type characterAttacks struct {
	PrimarySeq   combat.AttackSequence
	SecondarySeq combat.AttackSequence

	UpSeq combat.AttackSequence

	AerialSeq   combat.AttackSequence
	AerialUpSeq combat.AttackSequence

	AerialDownSmash        combat.Attack
	AerialDownSmashLanding combat.Attack
}

var boxHeadAttacks = func() characterAttacks {
	primaryAttacks, err := combat.NewAttacksFromJSON(CombatFS, "box_head_primary_attack_seq.json", "../shared/combatdata/")
	if err != nil {
		log.Fatal(err)
	}
	primaryAttackSeq := combat.NewAttackSeq(primaryAttacks...)

	aerialAttacks, err := combat.NewAttacksFromJSON(CombatFS, "box_head_aerial_attack_seq.json", "../shared/combatdata/")
	if err != nil {
		log.Fatal(err)
	}
	aerialAttackSeq := combat.NewAttackSeq(aerialAttacks...)

	upAttacks, err := combat.NewAttacksFromJSON(CombatFS, "box_head_up_attack_seq.json", "../shared/combatdata/")
	if err != nil {
		log.Fatal(err)
	}
	upAttackSeq := combat.NewAttackSeq(upAttacks...)

	aerialUpAttacks, err := combat.NewAttacksFromJSON(CombatFS, "box_head_aerial_up_attack_seq.json", "../shared/combatdata/")
	if err != nil {
		log.Fatal(err)
	}
	aerialUpAttackSeq := combat.NewAttackSeq(aerialUpAttacks...)

	// Down smash is weird
	pathAerialDownSmash := "box_head_aerial_attack_down_smash.json"
	attacksDownSmash, err := combat.NewAttacksFromJSON(CombatFS, pathAerialDownSmash, "../shared/combatdata/")
	if err != nil {
		log.Fatal(err)
	}

	aerSeq := combat.NewAttackSeq(attacksDownSmash...)
	atkDownSmash := aerSeq.First()
	atkDownSmashLanding := aerSeq.Attacks[1]
	// ----

	return characterAttacks{
		PrimarySeq:             *primaryAttackSeq,
		AerialSeq:              *aerialAttackSeq,
		UpSeq:                  *upAttackSeq,
		AerialUpSeq:            *aerialUpAttackSeq,
		AerialDownSmash:        atkDownSmash,
		AerialDownSmashLanding: atkDownSmashLanding,
	}
}()

var BoxHeadAttackAnimationMapping = map[int]client.AnimationData{
	boxHeadAttacks.PrimarySeq.First().ID:  animations.Registry.Characters[characterkeys.BoxHead][animations.PrimaryAttack],
	boxHeadAttacks.AerialSeq.First().ID:   animations.Registry.Characters[characterkeys.BoxHead][animations.Aerial],
	boxHeadAttacks.AerialUpSeq.First().ID: animations.Registry.Characters[characterkeys.BoxHead][animations.UpAerial],
	boxHeadAttacks.UpSeq.First().ID:       animations.Registry.Characters[characterkeys.BoxHead][animations.UpAttack],

	boxHeadAttacks.AerialDownSmash.ID:        animations.Registry.Characters[characterkeys.BoxHead][animations.AerialDownSmash],
	boxHeadAttacks.AerialDownSmashLanding.ID: animations.Registry.Characters[characterkeys.BoxHead][animations.AerialDownSmashLanding],
	boxHeadAttacks.AerialDownSmashLanding.ID: animations.Registry.Characters[characterkeys.BoxHead][animations.AerialDownSmashLanding],
}

var boxHeadHurtBoxes = func() combat.HurtBoxes {
	path := "box_head_hurtboxes.json"
	boxes, err := combat.NewHurtBoxesFromJSON(CombatFS, path, "../shared/combatdata/")
	if err != nil {
		log.Fatal(err)
	}

	return boxes
}()
