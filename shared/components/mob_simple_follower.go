package components

import "github.com/TheBitDrifter/bappa/blueprint/vector"

type MobFollower struct {
	VisionRadius float64
	Speed        float64
	StopRadius   float64
	Offset       vector.Two
}
