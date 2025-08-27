package combatdata

import (
	"log"

	"github.com/TheBitDrifter/bappa/combat"
)

var LazySkullyAttackKeys = struct {
	PrimaryCombo   int
	SecondaryCombo int
	TertCombo      int
}{
	PrimaryCombo:   0,
	SecondaryCombo: 1,
	TertCombo:      2,
}

var LazySkullyAttacks = func() []combat.Attack {
	pathPrimarySeq := "lazy_skully_attacks.json"
	attacks, err := combat.NewAttacksFromJSON(CombatFS, pathPrimarySeq, "../shared/combatdata/")
	if err != nil {
		log.Fatal(err)
	}
	return attacks
}()
