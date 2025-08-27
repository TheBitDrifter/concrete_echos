package scenes

import (
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/ldtk"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

const OUTER_EAST_DISTRICT_SCENE_NAME = "OuterEastDistrict"

var OUTER_EAST_DISTRICT_SCENE = Scene{
	Name:    OUTER_EAST_DISTRICT_SCENE_NAME,
	Plan:    outerEastDistrictPlan,
	Width:   ldtk.DATA.WidthFor(OUTER_EAST_DISTRICT_SCENE_NAME),
	Height:  ldtk.DATA.HeightFor(OUTER_EAST_DISTRICT_SCENE_NAME),
	Preload: *DEFAULT_PRELOAD,
}

func outerEastDistrictPlan(width, height int, sto warehouse.Storage) error {
	err := ldtk.DATA.LoadTiles(OUTER_EAST_DISTRICT_SCENE_NAME, sto)
	if err != nil {
		return err
	}

	err = ldtk.DATA.LoadIntGridFromConfig(OUTER_EAST_DISTRICT_SCENE_NAME, sto, intGridConfigs)
	if err != nil {
		return err
	}

	err = ldtk.DATA.LoadEntities(OUTER_EAST_DISTRICT_SCENE_NAME, sto, entityRegistry)
	if err != nil {
		return err
	}

	err = NewEasternDistrictGreenSkyBackground(sto, 0, 0)
	if err != nil {
		return err
	}
	err = NewEasternDistrictBG(sto, 0, -250)
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

	return loadRelevantState(sto, OUTER_EAST_DISTRICT_SCENE_NAME)
}
