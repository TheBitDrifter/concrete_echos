package scenes

import (
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/ldtk"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

const RISING_PASSAGE_SCENE_NAME = "RisingPassage"

var RISING_PASSAGE_SCENE = Scene{
	Name:    RISING_PASSAGE_SCENE_NAME,
	Plan:    risingPassagePlan,
	Width:   ldtk.DATA.WidthFor(RISING_PASSAGE_SCENE_NAME),
	Height:  ldtk.DATA.HeightFor(RISING_PASSAGE_SCENE_NAME),
	Preload: *DEFAULT_PRELOAD,
}

func risingPassagePlan(width, height int, sto warehouse.Storage) error {
	err := ldtk.DATA.LoadTiles(RISING_PASSAGE_SCENE_NAME, sto)
	if err != nil {
		return err
	}

	err = ldtk.DATA.LoadIntGridFromConfig(RISING_PASSAGE_SCENE_NAME, sto, intGridConfigs)
	if err != nil {
		return err
	}

	err = ldtk.DATA.LoadEntities(RISING_PASSAGE_SCENE_NAME, sto, entityRegistry)
	if err != nil {
		return err
	}

	err = NewEasternDistrictGreenSkyBackground(sto, 0, 0)
	if err != nil {
		return err
	}
	err = NewEasternOutskirtsBG(sto, -40, 400)
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

	return loadRelevantState(sto, RISING_PASSAGE_SCENE_NAME)
}
