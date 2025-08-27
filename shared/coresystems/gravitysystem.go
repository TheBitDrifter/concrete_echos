package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
	"github.com/TheBitDrifter/concrete_echos/shared/combatdata"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

const (
	DEFAULT_GRAVITY  = 9.8
	PIXELS_PER_METER = 50.0
)

type GravitySystem struct{}

func (GravitySystem) Run(scene blueprint.Scene, dt float64) error {
	query := warehouse.Factory.NewQuery().And(
		motion.Components.Dynamics,
		warehouse.Factory.NewQuery().Not(components.NoGravityTag, components.SwapVulnerableComponent),
	)

	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		dyn := motion.Components.Dynamics.GetFromCursor(cursor)

		if combat.Components.Attack.CheckCursor(cursor) {
			atk := combat.Components.Attack.GetFromCursor(cursor)
			cKey := components.CharacterKeyComponent.GetFromCursor(cursor)

			matchedAtk, ok := combatdata.AerialDownSmashes[*cKey]
			if ok && matchedAtk.ID == atk.ID {
				dyn.Vel.Y = 450
			}
		}

		if atk, atkOK := combat.Components.Attack.GetFromCursorSafe(cursor); atkOK && dyn.Vel.Y > 0 {
			if atk.ID == combatdata.AerialSeqs[characterkeys.BoxHead].First().ID {
				continue
			}
		}

		mass := 1 / dyn.InverseMass

		gravity := motion.Forces.Generator.NewGravityForce(mass, DEFAULT_GRAVITY, PIXELS_PER_METER)

		motion.Forces.AddForce(dyn, gravity)
	}
	return nil
}
