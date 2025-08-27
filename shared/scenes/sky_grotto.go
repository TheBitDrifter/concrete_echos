package scenes

import (
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/ldtk"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

const SKY_GROTTO_SCENE_NAME = "SkyGrotto"

var SKY_GROTTO_SCENE = Scene{
	Name:    SKY_GROTTO_SCENE_NAME,
	Plan:    skyGrottoPlan,
	Width:   ldtk.DATA.WidthFor(SKY_GROTTO_SCENE_NAME),
	Height:  ldtk.DATA.HeightFor(SKY_GROTTO_SCENE_NAME),
	Preload: *DEFAULT_PRELOAD,
}

func skyGrottoPlan(width, height int, sto warehouse.Storage) error {
	err := ldtk.DATA.LoadTiles(SKY_GROTTO_SCENE_NAME, sto)
	if err != nil {
		return err
	}

	err = ldtk.DATA.LoadIntGridFromConfig(SKY_GROTTO_SCENE_NAME, sto, intGridConfigs)
	if err != nil {
		return err
	}

	err = ldtk.DATA.LoadEntities(SKY_GROTTO_SCENE_NAME, sto, entityRegistry)
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

	return loadRelevantState(sto, SKY_GROTTO_SCENE_NAME)
}
