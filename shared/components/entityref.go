package components

import "github.com/TheBitDrifter/bappa/warehouse"

type EntityRef struct {
	ID, Recycled int
}

type EntityReferences struct {
	Refs   [10]EntityRef
	Active [10]bool
}

func (e EntityReferences) AllActive(sto warehouse.Storage) []warehouse.Entity {
	res := []warehouse.Entity{}
	for i, ref := range e.Refs {
		if !e.Active[i] {
			continue
		}

		en, err := sto.Entity(ref.ID)
		if err == nil && en.Recycled() == ref.Recycled {
			res = append(res, en)
		}
	}

	return res
}
