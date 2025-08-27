package scenes

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

const FAST_TRAVEL_SCENE_NAME = "FT"

var FAST_TRAVEL_SCENE = Scene{
	Name:    FAST_TRAVEL_SCENE_NAME,
	Plan:    ftPlan,
	Width:   640,
	Height:  360,
	Preload: *DEFAULT_PRELOAD,
}

func ftPlan(width, height int, sto warehouse.Storage) error {
	err := NewAmbientWindNoise(sto)
	if err != nil {
		return err
	}

	err = AddPlaylist(sto, 0, sounds.DefaultSoundCollection)
	if err != nil {
		return err
	}
	composition := []warehouse.Component{
		client.Components.SpriteBundle,
		spatial.Components.Position,
		client.Components.CameraIndex,
	}

	arche, err := sto.NewOrExistingArchetype(composition...)
	if err != nil {
		return err
	}

	err = arche.Generate(1,
		spatial.NewPosition(0, 0),
		client.NewSpriteBundle().
			AddSprite("images/fast_travel_screen.png", true).WithOffset(vector.Two{X: 0, Y: 0}),
	)
	if err != nil {
		return nil
	}

	composition = []warehouse.Component{
		client.Components.SpriteBundle,
		spatial.Components.Position,
		client.Components.CameraIndex,
		components.IsFastTravelButtonTag,
	}
	arche, err = sto.NewOrExistingArchetype(composition...)
	if err != nil {
		return err
	}

	err = arche.Generate(1,
		spatial.NewPosition(0, 0),
		client.NewSpriteBundle().
			AddSprite("images/fast_travel_btn.png", true).WithCustomRenderer(),
	)
	if err != nil {
		return nil
	}
	return nil
}
