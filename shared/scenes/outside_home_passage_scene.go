package scenes

import (
	"log"

	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/ldtk"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

const OUTSIDE_HOME_PASSAGE_NAME = "OutsideHomePassage"

var OUTSIDE_HOME_PASSAGE_SCENE = Scene{
	Name:    OUTSIDE_HOME_PASSAGE_NAME,
	Plan:    outsideHomePassagePlan,
	Width:   ldtk.DATA.WidthFor(OUTSIDE_HOME_PASSAGE_NAME),
	Height:  ldtk.DATA.HeightFor(OUTSIDE_HOME_PASSAGE_NAME),
	Preload: *DEFAULT_PRELOAD,
}

func outsideHomePassagePlan(width, height int, sto warehouse.Storage) error {
	// Load the image tiles
	err := ldtk.DATA.LoadTiles(OUTSIDE_HOME_PASSAGE_NAME, sto)
	if err != nil {
		return err
	}

	err = ldtk.DATA.LoadIntGridFromConfig(OUTSIDE_HOME_PASSAGE_NAME, sto, intGridConfigs)
	if err != nil {
		return err
	}

	// Load custom LDTK entities
	err = ldtk.DATA.LoadEntities(OUTSIDE_HOME_PASSAGE_NAME, sto, entityRegistry)
	if err != nil {
		log.Printf("Error loading entities: %v", err)
		return err
	}

	err = NewEasternDistrictGreenSkyBackground(sto, 0, 0)
	if err != nil {
		return err
	}
	err = NewEasternOutskirtsBG(sto, 0, 40)
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

	return loadRelevantState(sto, OUTSIDE_HOME_PASSAGE_NAME)
}
