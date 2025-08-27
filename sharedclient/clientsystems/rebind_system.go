package clientsystems

import (
	"fmt"
	"log"

	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/actions"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type KeyRebindingSystem struct{}

func (KeyRebindingSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	query := warehouse.Factory.NewQuery().And(components.KeyRebindingStateComponent)
	cursor := scene.NewCursor(query)
	var state *components.KeyRebindingState

	if cursor.TotalMatched() == 0 {
		arche, err := scene.Storage().NewOrExistingArchetype(components.KeyRebindingStateComponent)
		if err != nil {
			return err
		}
		arche.Generate(1)
	} else {
		for range cursor.Next() {
			state = components.KeyRebindingStateComponent.GetFromCursor(cursor)
			break
		}
	}

	receiver := cli.Receiver(0)
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {

		receiver.ResetPadButtonMapping()

		if !state.Active {
			state.Active = true
			state.ActionsToBind = []input.Action{
				actions.Jump,
				actions.PrimaryAttack,
				actions.Dodge,
				actions.Interact,
				actions.Cancel,
				actions.ShiftTeleTargetRight,
				actions.ShiftTeleTargetLeft,
				actions.ShiftTeleTargetNear,
				actions.TeleSwap,
			}
			state.CurrentAction = state.ActionsToBind[0]
			fmt.Println("Key rebinding activated.")
		} else {
			state.Active = false
			fmt.Println("Key rebinding cancelled.")
		}
	}

	if state != nil && !state.Active {
		return nil
	}

	const gamepadID = 0

	if !receiver.PadActive() {
		return nil
	}

	for i := 0; i <= int(ebiten.GamepadButtonMax); i++ {
		button := ebiten.GamepadButton(i)
		if inpututil.IsGamepadButtonJustPressed(gamepadID, button) {

			fmt.Printf("Bound action %s to button %d\n", actionToString(state.CurrentAction), button)

			log.Println(button, "button")

			if state.CurrentAction == actions.Jump {
				receiver.RegisterGamepadReleasedButton(button, actions.JumpReleased)
				receiver.RegisterGamepadButton(button, state.CurrentAction)
			} else {
				receiver.RegisterGamepadJustPressedButton(button, state.CurrentAction)
			}

			if len(state.ActionsToBind) > 1 {
				state.ActionsToBind = state.ActionsToBind[1:]
				state.CurrentAction = state.ActionsToBind[0]
			} else {
				state.Active = false
				fmt.Println("Key rebinding finished.")
			}
			break
		}
	}

	return nil
}

func actionToString(a input.Action) string {
	switch a {
	case actions.Jump:
		return "Jump"
	case actions.PrimaryAttack:
		return "Primary Attack"
	case actions.Dodge:
		return "Dodge"
	case actions.Interact:
		return "Interact"
	case actions.Cancel:
		return "Cancel Interaction"
	case actions.TeleSwap:
		return "Swap Teleport(if unlocked)"
	case actions.ShiftTeleTargetRight:
		return "Shift Teleport Target Right(if unlocked)"
	case actions.ShiftTeleTargetLeft:
		return "Shift Teleport Target Left(if unlocked)"
	case actions.ShiftTeleTargetNear:
		return "Shift Teleport Target Near(if unlocked)"

	default:
		return fmt.Sprintf("Action %d", a)
	}
}
