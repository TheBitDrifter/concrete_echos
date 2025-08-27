package animations

import (
	"embed"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
)

//go:embed *
var AnimFS embed.FS

type AnimEnum int

const (
	Idle AnimEnum = iota
	Run
	Jump
	Fall
	PrimaryAttack
	Hurt
	Dodge
	Aerial
	AerialDownSmash
	AerialDownSmashLanding
	Defeat
	SecondaryAttack
	PrimaryRanged
	Teleport
	TeleportMarker
	TeleportEffect
	ConvoStart
	InConvo
	IsSaving
	UpAttack
	UpAerial
	Execute
)

type registry struct {
	Characters map[characterkeys.CharEnum]map[AnimEnum]client.AnimationData
}

var Registry = registry{
	Characters: make(map[characterkeys.CharEnum]map[AnimEnum]client.AnimationData),
}
