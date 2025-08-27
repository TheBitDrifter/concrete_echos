package clientsystems

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

type TrapDoorSoundSystem struct {
	Volume float64
}

func (sys TrapDoorSoundSystem) Run(lc coldbrew.LocalClient, scene coldbrew.Scene) error {
	query := warehouse.Factory.NewQuery().And(components.TrapDoorComponent)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		td := components.TrapDoorComponent.GetFromCursor(cursor)

		if td.LastChangedTick == 0 {
			continue
		}
		if td.LastChangedTick != scene.CurrentTick() {
			continue
		}

		soundBundle := client.Components.SoundBundle.GetFromCursor(cursor)
		sound, err := coldbrew.MaterializeSound(soundBundle, sounds.TrapDoorOpenSound)
		if err != nil {
			return err
		}

		player := sound.GetAny()
		player.SetVolume(sys.Volume)

		// Loop if needed
		if !player.IsPlaying() {
			player.Rewind()
			player.Play()
		}
	}
	return nil
}
