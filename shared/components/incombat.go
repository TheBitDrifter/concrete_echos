package components

type LastCombat struct {
	StartTick int
}

func (js LastCombat) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}
