package components

import (
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
)

type JumpState struct {
	LastJump          int
	LastWallJump      int
	LastJumpRelease   int
	Locked            bool
	WallJumpDirection spatial.Direction
}

func (js JumpState) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}
