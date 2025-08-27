package scenes

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/animations"
)

const HOME_SCREEN_NAME = "HOME_SCREEN"

var HOME_SCREEN_SCENE = Scene{
	Name:   HOME_SCREEN_NAME,
	Plan:   home_screen_plan,
	Width:  640,
	Height: 360,
}

func home_screen_plan(width, height int, sto warehouse.Storage) error {
	composition := []warehouse.Component{
		client.Components.SpriteBundle,
		spatial.Components.Position,
		client.Components.CameraIndex,
	}

	arche, err := sto.NewOrExistingArchetype(composition...)
	if err != nil {
		return err
	}

	return arche.Generate(1,
		spatial.NewPosition(0, 0),
		client.NewSpriteBundle().
			AddSprite("images/home_screen_sheet.png", true).
			WithAnimations(animations.HomeScreenAnim),
	)
}

const DEFEAT_SCREEN_NAME = "DEFEAT_SCREEN"

var DEFEAT_SCREEN_SCENE = Scene{
	Name:   DEFEAT_SCREEN_NAME,
	Plan:   defeat_screen_plan,
	Width:  640,
	Height: 360,
}

func defeat_screen_plan(width, height int, sto warehouse.Storage) error {
	composition := []warehouse.Component{
		client.Components.SpriteBundle, spatial.Components.Position,
		client.Components.CameraIndex,
	}

	arche, err := sto.NewOrExistingArchetype(composition...)
	if err != nil {
		return err
	}

	return arche.Generate(1,
		spatial.NewPosition(0, 0),
		client.NewSpriteBundle().
			AddSprite("images/defeat_screen_sheet.png", true).
			WithAnimations(animations.DefeatScreenAnim),
	)
}
