package scenes

import (
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/ldtk"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

const WARP_UNLOCK_SCENE_NAME = "WarpUnlock"

var WARP_UNLOCK_SCENE = Scene{
	Name:    WARP_UNLOCK_SCENE_NAME,
	Plan:    warpUnlock,
	Width:   ldtk.DATA.WidthFor(WARP_UNLOCK_SCENE_NAME),
	Height:  ldtk.DATA.HeightFor(WARP_UNLOCK_SCENE_NAME),
	Preload: *DEFAULT_PRELOAD,
}

func warpUnlock(width, height int, sto warehouse.Storage) error {
	err := ldtk.DATA.LoadTiles(WARP_UNLOCK_SCENE_NAME, sto)
	if err != nil {
		return err
	}

	err = ldtk.DATA.LoadIntGridFromConfig(WARP_UNLOCK_SCENE_NAME, sto, intGridConfigs)
	if err != nil {
		return err
	}

	err = ldtk.DATA.LoadEntities(WARP_UNLOCK_SCENE_NAME, sto, entityRegistry)
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

	return loadRelevantState(sto, WARP_UNLOCK_SCENE_NAME)
}
