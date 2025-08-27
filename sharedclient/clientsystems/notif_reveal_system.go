package clientsystems

import (
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/coldbrew/text"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/actions"
	"github.com/TheBitDrifter/concrete_echos/shared/components" // Assuming this is your components path
	"github.com/TheBitDrifter/concrete_echos/sharedclient/fontdata"
)

const NOTIF_REVEAL_DELAY_IN_TICKS = 3

type NotifTextRevealSystem struct {
	REVEAL_START_DELAY int
}

func (sys NotifTextRevealSystem) Run(lc coldbrew.LocalClient, scene coldbrew.Scene) error {
	query := warehouse.Factory.NewQuery().And(components.SimpleNotificationComponent)
	cursor := scene.NewCursor(query)
	currentTick := scene.CurrentTick()

	for range cursor.Next() {
		if sys.closeRequested(scene) {
			en, err := cursor.CurrentEntity()
			if err != nil {
				return err
			}
			scene.Storage().EnqueueDestroyEntities(en)
		}

		noti := components.SimpleNotificationComponent.GetFromCursor(cursor)
		if !noti.Wrapped {
			noti.Title = text.WrapText(noti.Title, fontdata.TITLE_FONT_FACE, noti.TitleMaxWidth)
			noti.Body = text.WrapText(noti.Body, fontdata.UNLOCK_BODY_FONT_FACE, noti.BodyMaxWidth)
		}

		if noti.StartedTick+sys.REVEAL_START_DELAY > scene.CurrentTick() {
			continue
		}

		if noti.IsFinished {
			continue
		}

		if noti.StartedTick == 0 {
			noti.StartedTick = scene.CurrentTick()
		}

		if noti.RevealStarted == 0 {
			noti.RevealStarted = currentTick
		}

		if !noti.IsTitleFinished {
			finished, count := text.CurrentIndexInTextReveal(
				noti.RevealStarted+sys.REVEAL_START_DELAY,
				currentTick,
				NOTIF_REVEAL_DELAY_IN_TICKS,
				noti.Title,
			)

			noti.DisplayedTitle = noti.Title[:count]

			if finished {
				noti.IsTitleFinished = true
				noti.RevealStarted = currentTick
			}
			continue
		}

		finished, count := text.CurrentIndexInTextReveal(
			noti.RevealStarted,
			currentTick,
			NOTIF_REVEAL_DELAY_IN_TICKS,
			noti.Body,
		)

		noti.DisplayedBody = noti.Body[:count]

		if finished {
			noti.IsFinished = true
		}
	}

	return nil
}

func (sys NotifTextRevealSystem) closeRequested(scene coldbrew.Scene) bool {
	query := warehouse.Factory.NewQuery().And(input.Components.ActionBuffer)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		actionsBuffer := input.Components.ActionBuffer.GetFromCursor(cursor)
		if ok := actionsBuffer.HasAction(actions.Interact); ok {
			return true
		}

	}
	return false
}
