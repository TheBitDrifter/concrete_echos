package actions

import (
	"github.com/TheBitDrifter/bappa/blueprint/input"
)

var (
	Left                 = input.NewAction()
	Right                = input.NewAction()
	Jump                 = input.NewAction()
	Down                 = input.NewAction()
	PrimaryAttack        = input.NewAction()
	Dodge                = input.NewAction()
	VectorTwoMovement    = input.NewAction()
	VectorTwoCamMovement = input.NewAction()
	CameraRight          = input.NewAction()
	CameraLeft           = input.NewAction()
	CameraUp             = input.NewAction()
	CameraDown           = input.NewAction()
	AttackDown           = input.NewAction()
	Interact             = input.NewAction()
	Cancel               = input.NewAction()
	JumpReleased         = input.NewAction()
	Up                   = input.NewAction()
	ShiftTeleTargetLeft  = input.NewAction()
	ShiftTeleTargetRight = input.NewAction()
	ShiftTeleTargetNear  = input.NewAction()
	TeleSwap             = input.NewAction()
)
