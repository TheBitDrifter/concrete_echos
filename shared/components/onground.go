package components

import "github.com/TheBitDrifter/bappa/blueprint/vector"

type OnGround struct {
	LastTouch   int        `json:"-"`
	Landed      int        `json:"-"`
	SlopeNormal vector.Two `json:"-"`
}
