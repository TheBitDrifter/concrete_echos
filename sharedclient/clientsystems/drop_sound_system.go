package clientsystems

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

type DropSoundSystem struct {
	Volume float64
}

func (sys DropSoundSystem) Run(lc coldbrew.LocalClient, scene coldbrew.Scene) error {
	var soundBundle *client.SoundBundle

	playerQuery := warehouse.Factory.NewQuery().And(components.PlayerTag)
	pCursor := scene.NewCursor(playerQuery)

	for range pCursor.Next() {
		soundBundle = client.Components.SoundBundle.GetFromCursor(pCursor)
		break
	}

	query := warehouse.Factory.NewQuery().And(components.DropComponent)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		drop := components.DropComponent.GetFromCursor(cursor)
		if drop.Opened || drop.TickDropped == 0 {
			continue
		}

		if drop.TickDropped != scene.CurrentTick() {
			continue
		}

		sound, err := coldbrew.MaterializeSound(soundBundle, sounds.DropMoneyAddedSound)
		if err != nil {
			return err
		}
		player := sound.GetAny()
		player.SetVolume(sys.Volume)

		if !player.IsPlaying() {
			player.Rewind()
			player.Play()
		}
	}
	return nil
}
