package scenes

import (
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/ldtk"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

const JUNCTION_SCENE_NAME = "Junction"

var JUNCTION_SCENE = Scene{
	Name:    JUNCTION_SCENE_NAME,
	Plan:    junctionScenePlan,
	Width:   ldtk.DATA.WidthFor(JUNCTION_SCENE_NAME),
	Height:  ldtk.DATA.HeightFor(JUNCTION_SCENE_NAME),
	Preload: *DEFAULT_PRELOAD,
}

func junctionScenePlan(width, height int, sto warehouse.Storage) error {
	err := ldtk.DATA.LoadTiles(JUNCTION_SCENE_NAME, sto)
	if err != nil {
		return err
	}

	err = ldtk.DATA.LoadIntGridFromConfig(JUNCTION_SCENE_NAME, sto, intGridConfigs)
	if err != nil {
		return err
	}

	err = ldtk.DATA.LoadEntities(JUNCTION_SCENE_NAME, sto, entityRegistry)
	if err != nil {
		return err
	}

	err = NewEasternDistrictGreenSkyBackground(sto, 0, 0)
	if err != nil {
		return err
	}

	err = NewEasternOutskirtsJunctionBG(sto, 0, 0)
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

	return loadRelevantState(sto, JUNCTION_SCENE_NAME)
}
