package scenes

import (
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/ldtk"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

const EASTERN_GARDENS_SCENE_NAME = "EasternGardens"

var EASTERN_GARDENS_SCENE = Scene{
	Name:    EASTERN_GARDENS_SCENE_NAME,
	Plan:    easternGardensScenePlan,
	Width:   ldtk.DATA.WidthFor(EASTERN_GARDENS_SCENE_NAME),
	Height:  ldtk.DATA.HeightFor(EASTERN_GARDENS_SCENE_NAME),
	Preload: *DEFAULT_PRELOAD,
}

func easternGardensScenePlan(width, height int, sto warehouse.Storage) error {
	err := ldtk.DATA.LoadTiles(EASTERN_GARDENS_SCENE_NAME, sto)
	if err != nil {
		return err
	}

	err = ldtk.DATA.LoadIntGridFromConfig(EASTERN_GARDENS_SCENE_NAME, sto, intGridConfigs)
	if err != nil {
		return err
	}

	err = ldtk.DATA.LoadEntities(EASTERN_GARDENS_SCENE_NAME, sto, entityRegistry)
	if err != nil {
		return err
	}

	err = NewEasternDistrictGreenSkyBackground(sto, 0, 0)
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

	return loadRelevantState(sto, EASTERN_GARDENS_SCENE_NAME)
}
