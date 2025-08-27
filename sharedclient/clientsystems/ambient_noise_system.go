package clientsystems

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

type AmbientNoiseSystem struct {
	Volume float64
}

func (sys AmbientNoiseSystem) Run(lc coldbrew.LocalClient, scene coldbrew.Scene) error {
	query := warehouse.Factory.NewQuery().And(components.AmbientNoiseTag)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		soundBundle := client.Components.SoundBundle.GetFromCursor(cursor)
		sound, err := coldbrew.MaterializeSound(soundBundle, sounds.AmbientWindNoise)
		if err != nil {
			return err
		}
		player := sound.GetAny()

		if !player.IsPlaying() {

			player.SetVolume(sys.Volume)
			player.Rewind()
			player.Play()
		}
	}
	return nil
}
