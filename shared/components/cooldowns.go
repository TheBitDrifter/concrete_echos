package components

type Cooldown struct {
	StartTick, Duration int `json:"-"`
}

func (cd Cooldown) Available(currentTick int) bool {
	expiresOn := cd.StartTick + cd.Duration
	return expiresOn < currentTick
}

type MovementCooldowns struct {
	Dodge Cooldown
}

type AttackCooldowns struct {
	Aerial Cooldown
}
