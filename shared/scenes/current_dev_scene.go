package scenes

import (
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/ldtk"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

const CURRENT_DEV_SCENE_NAME = "CurrentDev"

var CURRENT_DEV_SCENE = Scene{
	Name:    CURRENT_DEV_SCENE_NAME,
	Plan:    currentDevPlan,
	Width:   ldtk.DATA.WidthFor(CURRENT_DEV_SCENE_NAME),
	Height:  ldtk.DATA.HeightFor(CURRENT_DEV_SCENE_NAME),
	Preload: *DEFAULT_PRELOAD,
}

func currentDevPlan(width, height int, sto warehouse.Storage) error {
	err := ldtk.DATA.LoadTiles(CURRENT_DEV_SCENE_NAME, sto)
	if err != nil {
		return err
	}

	err = ldtk.DATA.LoadIntGridFromConfig(CURRENT_DEV_SCENE_NAME, sto, intGridConfigs)
	if err != nil {
		return err
	}

	err = NewEasternDistrictGreenSkyBackground(sto, 0, 0)
	if err != nil {
		return err
	}
	err = NewEasternDistrictBLDBackground(sto, 0, 1830)
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

	return ldtk.DATA.LoadEntities(CURRENT_DEV_SCENE_NAME, sto, entityRegistry)
}
