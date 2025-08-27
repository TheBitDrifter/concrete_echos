package sounds

import "github.com/TheBitDrifter/bappa/blueprint/client"

var P1 = client.SoundConfig{
	Path:             "sounds/tunes/p1.wav",
	AudioPlayerCount: 1,
}

var P2 = client.SoundConfig{
	Path:             "sounds/tunes/p2.wav",
	AudioPlayerCount: 1,
}

var P3 = client.SoundConfig{
	Path:             "sounds/tunes/p3.wav",
	AudioPlayerCount: 1,
}

var P4 = client.SoundConfig{
	Path:             "sounds/tunes/p4.wav",
	AudioPlayerCount: 1,
}

var P5 = client.SoundConfig{
	Path:             "sounds/tunes/p5.wav",
	AudioPlayerCount: 1,
}

var BossSong = client.SoundConfig{
	Path:             "sounds/tunes/boss.wav",
	AudioPlayerCount: 1,
}

var PostBossSong = client.SoundConfig{
	Path:             "sounds/tunes/post_boss.wav",
	AudioPlayerCount: 1,
}

type SoundCollection struct {
	ID     int
	Sounds []client.SoundConfig
}

var DefaultSoundCollection = SoundCollection{
	ID:     1,
	Sounds: defaultPlaylistCollectionSounds,
}

var defaultPlaylistCollectionSounds = []client.SoundConfig{
	P1,
	P2,
	P3,
	P4,
	P5,
}

var BossSoundCollection = SoundCollection{
	ID:     2,
	Sounds: bossPlaylistCollectionSounds,
}

var bossPlaylistCollectionSounds = []client.SoundConfig{
	BossSong,
}

var PostBossSoundCollection = SoundCollection{
	ID:     3,
	Sounds: postBossPlaylistCollectionSounds,
}

var postBossPlaylistCollectionSounds = []client.SoundConfig{
	PostBossSong,
	P1,
	P2,
	P3,
	P4,
}

var ActivePlaylistID = 0

var ActiveSongIndex = 0
