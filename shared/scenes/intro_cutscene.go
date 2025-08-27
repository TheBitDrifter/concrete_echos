package scenes

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/animations"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/dialoguedata"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

const INTRO_CUTSCENE_NAME = "INTROCS"

var INTRO_CUTSCENE = Scene{
	Name:   HOME_SCREEN_NAME,
	Plan:   intro_cutscene_plan,
	Width:  640,
	Height: 360,
}

func intro_cutscene_plan(width, height int, sto warehouse.Storage) error {
	_, err := NewDialogueEntityAutoStepper(sto, dialoguedata.IntroCutsceneDialogue)
	if err != nil {
		return err
	}

	// Cutscene

	composition := []warehouse.Component{
		client.Components.SpriteBundle,
		spatial.Components.Position,
		client.Components.CameraIndex,
		components.CutsceneTag,
	}

	arche, err := sto.NewOrExistingArchetype(composition...)
	if err != nil {
		return err
	}

	err = NewAmbientWindNoise(sto)
	if err != nil {
		return err
	}

	err = AddPlaylist(sto, 0, sounds.DefaultSoundCollection)
	if err != nil {
		return err
	}

	return arche.Generate(1,
		spatial.NewPosition(0, 0),
		client.NewSpriteBundle().
			AddSprite("images/intro_cutscene_sheet.png", true).
			WithAnimations(
				animations.CutSceneAnimations...,
			),
	)
}
