package components

type Drop struct {
	Opened                 bool
	TickDropped            int `json:"-"`
	MoneyDrop              float64
	HealthDrop             int
	CustomSpawnCallbackKey int
	// func(sto warehouse.Storage) error `json:"-"`
}
