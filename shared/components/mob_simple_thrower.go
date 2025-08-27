package components

import "github.com/TheBitDrifter/bappa/blueprint/vector"

type MobSimpleThrower struct {
	LastStartedThrow int
	IsWindingUp      bool
	LastFiredThrow   int
	VisionRadius     float64
	Delay            int

	// If the 'throw' takes 60 ticks (including recovery frames a.k.a ThrowDuration)
	// an example release could be at 45 ticks ...
	//
	// Usually the release happens towards the end so the projectile comes out at the apex of the
	// throwing motion, but before the follow through/recovery frames
	ThrowRelease  int
	ThrowDuration int

	SpawnOffset vector.Two
}

type MobSimpleHoriShooter struct {
	LastStartedShot int
	IsWindingUp     bool
	LastFiredShot   int
	VisionRadius    float64
	Delay           int
	ShootRelease    int
	ShootDuration   int
	SpawnOffset     vector.Two
}
