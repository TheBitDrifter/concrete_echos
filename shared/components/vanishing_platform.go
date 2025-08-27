package components

type VanishingPlatform struct {
	LastVanished     int
	LiveDuration     int
	RespawnDelay     int
	TimerStarted     bool
	TimerStartedTick int
	Vanished         bool
}
