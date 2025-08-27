package scenes

import (
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/ldtk"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

const UPPER_EAST_DISTRICT_SCENE_NAME = "UpperEastDistrict"

var UPPER_EAST_DISTRICT_SCENE = Scene{
	Name:    UPPER_EAST_DISTRICT_SCENE_NAME,
	Plan:    upperEastDistrictPlan,
	Width:   ldtk.DATA.WidthFor(UPPER_EAST_DISTRICT_SCENE_NAME),
	Height:  ldtk.DATA.HeightFor(UPPER_EAST_DISTRICT_SCENE_NAME),
	Preload: *DEFAULT_PRELOAD,
}

func upperEastDistrictPlan(width, height int, sto warehouse.Storage) error {
	err := ldtk.DATA.LoadTiles(UPPER_EAST_DISTRICT_SCENE_NAME, sto)
	if err != nil {
		return err
	}

	err = ldtk.DATA.LoadIntGridFromConfig(UPPER_EAST_DISTRICT_SCENE_NAME, sto, intGridConfigs)
	if err != nil {
		return err
	}

	err = ldtk.DATA.LoadEntities(UPPER_EAST_DISTRICT_SCENE_NAME, sto, entityRegistry)
	if err != nil {
		return err
	}

	err = NewEasternDistrictGreenSkyBackground(sto, 0, 0)
	if err != nil {
		return err
	}
	err = NewEasternDistrictBGAlt(sto, 0, 1670)
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

	return loadRelevantState(sto, UPPER_EAST_DISTRICT_SCENE_NAME)
}
