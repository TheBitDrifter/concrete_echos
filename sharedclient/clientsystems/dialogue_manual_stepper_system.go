package clientsystems

import (
	"log"

	"github.com/TheBitDrifter/bappa/blueprint/dialogue"
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/coldbrew/coldbrew_clientsystems"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/actions"

	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/scenes"
	"github.com/TheBitDrifter/concrete_echos/sharedclient/fontdata"
)

type DialogueManualStepperSystem struct{}

func (sys DialogueManualStepperSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	nextSlideRequested := false
	cancelConvo := false

	actionsQuery := warehouse.Factory.NewQuery().And(input.Components.ActionBuffer)
	actionsCursor := scene.NewCursor(actionsQuery)

	query := warehouse.Factory.NewQuery().And(
		dialogue.Components.Conversation,
		components.DialogueManualStepperComponent,
	)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {

		for range actionsCursor.Next() {
			actionsBuffer := input.Components.ActionBuffer.GetFromCursor(actionsCursor)
			nextSlideRequested = actionsBuffer.HasAction(actions.Interact)
			_, cancelConvo = actionsBuffer.ConsumeAction(actions.Cancel)

		}

		convo := dialogue.Components.Conversation.GetFromCursor(cursor)
		stepper := components.DialogueManualStepperComponent.GetFromCursor(cursor)

		if !convo.IsRevealing {
			coldbrew_clientsystems.InitConversation(scene, convo, fontdata.DEFAULT_FONT_FACE, float64(fontdata.MAX_LINE_WIDTH))
		}

		if cancelConvo {
			en, err := cursor.CurrentEntity()
			if err != nil {
				return err
			}

			references := components.EntityReferencesComponent.GetFromCursor(cursor)

			toClear := references.AllActive(scene.Storage())

			for _, clearMe := range toClear {
				convo := components.InConversationComponent.GetFromEntity(clearMe)
				convo.EndedTick = scene.CurrentTick()
			}
			scene.Storage().EnqueueDestroyEntities(en)

		}

		slides := dialogue.SlidesRegistry[dialogue.SlidesEnum(convo.SlidesID)]
		if !nextSlideRequested {
			continue
		}

		if scene.CurrentTick() < convo.FinalUpdateTick+stepper.MinDelayInTicks {
			// TODO: ENABLE DELAY EVENTUALLY!
			// continue
		}

		if convo.ActiveSlideIndex < len(slides)-1 {

			convo.ActiveSlideIndex++
			convo.IsRevealing = false
			convo.RevealStartTick = scene.CurrentTick()
		} else {
			en, err := cursor.CurrentEntity()
			if err != nil {
				return err
			}

			cb, okCB := dialogue.CallbackRegistry[convo.CallbackID]
			if okCB {
				err = cb(scene)
				if err != nil {
					return err
				}
			} else {
				log.Println("No CB for", convo.CallbackID)
			}

			references := components.EntityReferencesComponent.GetFromCursor(cursor)
			toClear := references.AllActive(scene.Storage())

			for _, clearMe := range toClear {
				convo := components.InConversationComponent.GetFromEntity(clearMe)
				convo.EndedTick = scene.CurrentTick()
			}
			scene.Storage().EnqueueDestroyEntities(en)
			return nil
		}
	}
	return nil
}

type DialogueManualStepperActivationSystem struct{}

func (sys DialogueManualStepperActivationSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	activate := false
	actionsQuery := warehouse.Factory.NewQuery().And(input.Components.ActionBuffer)
	actionsCursor := scene.NewCursor(actionsQuery)

	for range actionsCursor.Next() {
		actionsBuffer := input.Components.ActionBuffer.GetFromCursor(actionsCursor)
		activate = actionsBuffer.HasAction(actions.Interact)

	}

	queryP := warehouse.Factory.NewQuery().And(
		components.PlayerTag,
	)
	cursorP := scene.NewCursor(queryP)

	queryF := warehouse.Factory.NewQuery().And(
		components.DialogueActivationComponent,
	)
	cursorF := scene.NewCursor(queryF)

	var playerEn warehouse.Entity

	var friendlyEns []warehouse.Entity

	convosToCreate := []dialogue.SlidesEnum{}

	convoCallbacks := []dialogue.CallbackEnum{}

	for range cursorP.Next() {
		playerPos := spatial.Components.Position.GetFromCursor(cursorP)
		playerEn, _ = cursorP.CurrentEntity()
		if components.InConversationComponent.CheckCursor(cursorP) {
			continue
		}

		for range cursorF.Next() {
			if components.InConversationComponent.CheckCursor(cursorF) {
				continue
			}

			en, error := cursorF.CurrentEntity()
			if error != nil {
				return error
			}

			friendlyMobPos := spatial.Components.Position.GetFromCursor(cursorF)
			currentDistSq := playerPos.Two.Sub(friendlyMobPos.Two).MagSquared()
			diaActiv := components.DialogueActivationComponent.GetFromCursor(cursorF)

			inRange := currentDistSq <= diaActiv.Range*diaActiv.Range
			if inRange {
				inRange = currentDistSq > diaActiv.MinRange*diaActiv.MinRange
			}
			if inRange && diaActiv.MustBeRight {
				inRange = playerPos.X > friendlyMobPos.X
			}
			if inRange && diaActiv.MustBeLeft {
				inRange = playerPos.X < friendlyMobPos.X
			}

			if inRange && activate {
				da := components.DialogueActivationComponent.GetFromCursor(cursorF)
				friendlyEns = append(friendlyEns, en)

				convosToCreate = append(convosToCreate, da.SlidesID)
				convoCallbacks = append(convoCallbacks, da.ConvoCallbackID)

				playerEn.EnqueueAddComponentWithValue(components.InConversationComponent, components.InConversation{StartedTick: scene.CurrentTick()})
				en.EnqueueAddComponentWithValue(components.InConversationComponent, components.InConversation{StartedTick: scene.CurrentTick()})
			}
		}
	}

	entitiesToRef := append(friendlyEns, playerEn)
	for i, convoKey := range convosToCreate {
		_, err := scenes.NewDialogueEntityManualStepper(scene.Storage(), convoKey, convoCallbacks[i], 60, entitiesToRef...)
		if err != nil {
			return err
		}
	}
	return nil
}

type DialogueManualStepperClearingSystem struct{}

func (sys DialogueManualStepperClearingSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	query := warehouse.Factory.NewQuery().And(components.InConversationComponent)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		inConvo := components.InConversationComponent.GetFromCursor(cursor)

		if scene.CurrentTick() == inConvo.EndedTick+30 {
			en, err := cursor.CurrentEntity()
			if err != nil {
				return err
			}

			en.EnqueueRemoveComponent(components.InConversationComponent)
		}
	}

	return nil
}
