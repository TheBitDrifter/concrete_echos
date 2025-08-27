package components

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

type MusicPlaylist struct {
	ActiveSong        client.SoundConfig
	CurrentSongIndex  int
	Collection        sounds.SoundCollection
	SongCollectionID  int
	IsFading          bool
	FadeDurationTicks int
	FadeStartTimeTick int
}
