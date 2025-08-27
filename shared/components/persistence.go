package components

import "github.com/TheBitDrifter/concrete_echos/shared/persistence"

type Persistence struct {
	EntityType persistence.EntityEnum
	PersistID  persistence.PersistenceID
}
