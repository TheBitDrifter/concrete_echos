package scenes

import (
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/ldtk"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

const OUTSIDE_HOME_SCENE_NAME = "OutsideHome"

var OUTSIDE_HOME_SCENE = Scene{
	Name:    OUTSIDE_HOME_SCENE_NAME,
	Plan:    outsideHomePlan,
	Width:   ldtk.DATA.WidthFor(OUTSIDE_HOME_SCENE_NAME),
	Height:  ldtk.DATA.HeightFor(OUTSIDE_HOME_SCENE_NAME),
	Preload: *DEFAULT_PRELOAD,
}

func outsideHomePlan(width, height int, sto warehouse.Storage) error {
	err := ldtk.DATA.LoadTiles(OUTSIDE_HOME_SCENE_NAME, sto)
	if err != nil {
		return err
	}

	err = ldtk.DATA.LoadIntGridFromConfig(OUTSIDE_HOME_SCENE_NAME, sto, intGridConfigs)
	if err != nil {
		return err
	}

	err = ldtk.DATA.LoadEntities(OUTSIDE_HOME_SCENE_NAME, sto, entityRegistry)
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

	err = NewEasternDistrictGreenSkyBackground(sto, 0, 0)
	if err != nil {
		return err
	}
	err = NewEasternOutskirtsBG(sto, 0, -100)
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
	return loadRelevantState(sto, OUTSIDE_HOME_SCENE_NAME)
}
