package scenes

import (
	"fmt"

	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

const PLAYER_EVIL_THRESHOLD = 20

type ChestCallbackEnum int

const (
	_ ChestCallbackEnum = iota
	WarpSwapDrop
	WallJumpDrop
	ProtoEnd
)

var ChestCallbacksRegistry = map[ChestCallbackEnum]func(sto warehouse.Storage) error{
	WarpSwapDrop: warpSwapDrop,
	WallJumpDrop: wallJumpDrop,
	ProtoEnd:     protoEndDrop,
}

func warpSwapDrop(sto warehouse.Storage) error {
	playerQuery := warehouse.Factory.NewQuery().And(components.PlayerTag)
	playerCursor := warehouse.Factory.NewCursor(playerQuery, sto)
	for range playerCursor.Next() {
		playerEN, err := playerCursor.CurrentEntity()
		if err != nil {
			return err
		}
		playerEN.EnqueueAddComponent(components.WarpSwapUnlockedTag)
	}
	err := NewSimpleNotification(
		sto, "New Ability Unlocked!",
		"You have unlocked the Warp Swap Ability! Press the K-Key(or right trigger) to swap places with a highlighted enemy. You may shift the target using the U/I-Keys or D-Pad(L-R). The O-Key and D-Pad Up will target the nearest mob.",
		222,
		85,
	)
	if err != nil {
		return err
	}
	return nil
}

func wallJumpDrop(sto warehouse.Storage) error {
	playerQuery := warehouse.Factory.NewQuery().And(components.PlayerTag)
	playerCursor := warehouse.Factory.NewCursor(playerQuery, sto)
	for range playerCursor.Next() {
		playerEN, err := playerCursor.CurrentEntity()
		if err != nil {
			return err
		}
		playerEN.EnqueueAddComponent(components.WallJumpUnlockedTag)
	}
	err := NewSimpleNotification(
		sto, "New Ability Unlocked!",
		"You have unlocked the Wall Jump Ability! You can now jump on walls. It's like jumping on floors, but walls!",
		222,
		85,
	)
	if err != nil {
		return err
	}
	return nil
}

func protoEndDrop(sto warehouse.Storage) error {
	playerQuery := warehouse.Factory.NewQuery().And(components.PlayerTag)
	playerCursor := warehouse.Factory.NewCursor(playerQuery, sto)

	playerIsEvil := false
	count := 0

	for range playerCursor.Next() {
		playerEN, err := playerCursor.CurrentEntity()
		if err != nil {
			return err
		}
		execCount := components.PlayerExecutionCountComponent.GetFromCursor(playerCursor)
		count = execCount.Count

		if execCount.Count >= PLAYER_EVIL_THRESHOLD {
			playerIsEvil = true
		}
		playerEN.EnqueueAddComponent(components.WallJumpUnlockedTag)
	}
	body := fmt.Sprintf(
		"You only executed (absorbed) %d mobs. So the Lazy Skully boss was willing to work things out without violence. You can optionally hit him if you're looking for a fight!",
		count,
	)
	if playerIsEvil {
		body = fmt.Sprintf(
			"You executed (absorbed) %d mobs. So the Lazy Skully boss was not willing to work things out without violence. If you had absorbed less mobs, he might of been willing to talk it out!",
			count,
		)
	}

	err := NewSimpleNotification(
		sto, "You Beat the Prototype!",
		body,
		210,
		85,
	)
	if err != nil {
		return err
	}
	return nil
}
