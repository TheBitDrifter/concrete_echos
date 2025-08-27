package clientsystems

import (
	"math"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	muteMusic bool
	musicVol  float64
)

type MusicSystem struct {
	Volume float64
}

func (sys MusicSystem) Run(lc coldbrew.LocalClient, scene coldbrew.Scene) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		muteMusic = !muteMusic
	}
	if muteMusic {
		musicVol = sys.Volume
		sys.Volume = 0
	}
	if !muteMusic && musicVol != 0 {
		sys.Volume = musicVol
	}

	currentTick := scene.CurrentTick()
	musicQuery := warehouse.Factory.NewQuery().And(components.MusicPlaylistComponent)
	cursor := scene.NewCursor(musicQuery)

	for range cursor.Next() {
		list := components.MusicPlaylistComponent.GetFromCursor(cursor)
		bundle := client.Components.SoundBundle.GetFromCursor(cursor)

		if sounds.ActivePlaylistID != list.Collection.ID {
			list.CurrentSongIndex = 0
			sounds.ActiveSongIndex = 0

			list.IsFading = true

			list.FadeStartTimeTick = currentTick
			sounds.ActivePlaylistID = list.Collection.ID

		} else {
			list.CurrentSongIndex = sounds.ActiveSongIndex
		}

		if len(list.Collection.Sounds) == 0 {
			continue
		}

		if list.FadeDurationTicks <= 0 {
			list.FadeDurationTicks = 120
		}

		activeSongConfig := list.Collection.Sounds[list.CurrentSongIndex]
		activeSound, err := coldbrew.MaterializeSound(bundle, activeSongConfig)
		if err != nil {
			return err
		}
		activePlayer := activeSound.GetAny()

		if list.IsFading {

			if !activePlayer.IsPlaying() {
				activePlayer.SetVolume(0)
				activePlayer.Rewind()
				activePlayer.Play()
				list.FadeStartTimeTick = currentTick
			}
			elapsedTicks := currentTick - list.FadeStartTimeTick

			progress := math.Min(float64(elapsedTicks)/float64(list.FadeDurationTicks), 1.0)
			activePlayer.SetVolume(sys.Volume * progress)

			if progress >= 1.0 {
				list.IsFading = false
			}
		} else {
			if !activePlayer.IsPlaying() {
				list.CurrentSongIndex = (list.CurrentSongIndex + 1) % len(list.Collection.Sounds)
				sounds.ActiveSongIndex = list.CurrentSongIndex

				list.IsFading = true
				list.FadeStartTimeTick = currentTick
			} else {
				activePlayer.SetVolume(sys.Volume)
			}
		}
	}
	return nil
}
