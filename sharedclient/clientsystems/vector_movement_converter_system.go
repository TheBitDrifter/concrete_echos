package clientsystems

import (
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/warehouse"

	"github.com/TheBitDrifter/concrete_echos/shared/actions"
)

type VectorMovementConverterSytem struct{}

func (VectorMovementConverterSytem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	query := warehouse.Factory.NewQuery().And(input.Components.ActionBuffer)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {

		buffer := input.Components.ActionBuffer.GetFromCursor(cursor)
		stampedVectorMovement, existed := buffer.ConsumeAction(actions.VectorTwoMovement)

		if stampedVectorMovement.Tick != scene.CurrentTick() || !existed {
			continue
		}

		dz := 40

		dzyD := -70

		dzyU := 55

		if stampedVectorMovement.X > dz {
			newStamped := input.StampedAction{
				Tick: scene.CurrentTick() + 1,
				Val:  actions.Right,
			}

			buffer.Add(newStamped)
		}

		if stampedVectorMovement.X < -dz {
			newStamped := input.StampedAction{
				Tick: scene.CurrentTick() + 1,
				Val:  actions.Left,
			}

			buffer.Add(newStamped)
		}

		if stampedVectorMovement.Y < dzyD {
			if !combat.Components.Attack.CheckCursor(cursor) {
				newStamped := input.StampedAction{
					Tick: scene.CurrentTick() + 1,
					Val:  actions.Down,
				}

				buffer.Add(newStamped)
			}

			newStamped := input.StampedAction{
				Tick: scene.CurrentTick() + 1,
				Val:  actions.AttackDown,
			}

			buffer.Add(newStamped)
		}

		if stampedVectorMovement.Y > dzyU {
			if !combat.Components.Attack.CheckCursor(cursor) {
				newStamped := input.StampedAction{
					Tick: scene.CurrentTick() + 1,
					Val:  actions.Up,
				}

				buffer.Add(newStamped)
			}
		}

	}

	cursor = scene.NewCursor(query)
	for range cursor.Next() {

		buffer := input.Components.ActionBuffer.GetFromCursor(cursor)
		stampedVectorMovement, existed := buffer.ConsumeAction(actions.VectorTwoCamMovement)

		if stampedVectorMovement.Tick != scene.CurrentTick() || !existed {
			continue
		}

		dz := 40
		dzy := 70

		if stampedVectorMovement.X > dz {
			newStamped := input.StampedAction{
				Tick: scene.CurrentTick(),
				Val:  actions.CameraRight,
			}
			buffer.Add(newStamped)
		}

		if stampedVectorMovement.X < -dz {
			newStamped := input.StampedAction{
				Tick: scene.CurrentTick(),
				Val:  actions.CameraLeft,
			}

			buffer.Add(newStamped)
		}
		if stampedVectorMovement.Y < -dzy {
			newStamped := input.StampedAction{
				Tick: scene.CurrentTick(),
				Val:  actions.CameraDown,
			}

			buffer.Add(newStamped)
		}

		if stampedVectorMovement.Y > -dzy {
			newStamped := input.StampedAction{
				Tick: scene.CurrentTick(),
				Val:  actions.CameraUp,
			}

			buffer.Add(newStamped)
		}
	}
	return nil
}
