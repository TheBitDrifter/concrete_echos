package components

import "github.com/TheBitDrifter/bappa/tteokbokki/spatial"

type MobDemonAnt struct {
	LastAttack int

	AttackVisionJabRadius float64
	LastJabAttack         int
	JabAttackDelay        int

	LastSlashAttack  int
	SlashAttackDelay int

	DodgeVisionRadius float64
	DodgeDuration     int
	DodgeSpeed        float64
	DodgeDelay        int
	LastDodged        int
	DodgeDirection    spatial.Direction
}
