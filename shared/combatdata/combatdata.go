package combatdata

import (
	"embed"

	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
)

//go:embed *.json
var CombatFS embed.FS

var PrimarySeqs = map[characterkeys.CharEnum]combat.AttackSequence{
	characterkeys.BoxHead:    boxHeadAttacks.PrimarySeq,
	characterkeys.SkullRonin: skullRoninAttacks.PrimarySeq,
	characterkeys.CanThrower: canThrowerAttacks.PrimarySeq,
	characterkeys.DemonAnt:   demonAntAttacks.PrimarySeq,
}

var SecondarySeqs = map[characterkeys.CharEnum]combat.AttackSequence{
	characterkeys.DemonAnt: demonAntAttacks.SecondarySeq,
}

var UpSeqs = map[characterkeys.CharEnum]combat.AttackSequence{
	characterkeys.BoxHead: boxHeadAttacks.UpSeq,
}

var AerialSeqs = map[characterkeys.CharEnum]combat.AttackSequence{
	characterkeys.BoxHead: boxHeadAttacks.AerialSeq,
}

var AerialUpSeq = map[characterkeys.CharEnum]combat.AttackSequence{
	characterkeys.BoxHead: boxHeadAttacks.AerialUpSeq,
}

var AerialDownSmashes = map[characterkeys.CharEnum]combat.Attack{
	characterkeys.BoxHead: boxHeadAttacks.AerialDownSmash,
}

var AerialDownSmashLandings = map[characterkeys.CharEnum]combat.Attack{
	characterkeys.BoxHead: boxHeadAttacks.AerialDownSmashLanding,
}

var HurtBoxes = map[characterkeys.CharEnum]combat.HurtBoxes{
	characterkeys.BoxHead:            boxHeadHurtBoxes,
	characterkeys.SkullRonin:         skullRoninHurtBoxes,
	characterkeys.MudHand:            mudHandHurtboxes,
	characterkeys.CanThrower:         canThrowerHurtBoxes,
	characterkeys.CanStraightThrower: canThrowerHurtBoxes,
	characterkeys.DemonAnt:           demonAntHurtBoxes,
	characterkeys.Chest:              chestHurtBoxes,
}
