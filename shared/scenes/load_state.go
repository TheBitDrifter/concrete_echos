package scenes

import (
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/persistence"
)

func loadRelevantState(sto warehouse.Storage, sceneName string) error {
	chestQuery := warehouse.Factory.NewQuery().And(components.ChestTag, components.PersistenceComponent)
	chestCursor := warehouse.Factory.NewCursor(chestQuery, sto)

	chests := []warehouse.Entity{}
	for range chestCursor.Next() {
		chestEn, err := chestCursor.CurrentEntity()
		if err != nil {
			continue
		}
		chests = append(chests, chestEn)
	}

	for _, chestEn := range chests {
		chestPersist := components.PersistenceComponent.GetFromEntity(chestEn)
		seriChest, ok := persistence.State.Get(sceneName, chestPersist.EntityType, chestPersist.PersistID)
		if !ok {
			continue
		}
		_, err := sto.ForceSerializedEntityWithID(seriChest, int(chestEn.ID()))
		if err != nil {
			return err
		}
		err = seriChest.SetValue(chestEn)
		if err != nil {
			return err
		}
	}

	// TODO: YO THIS WILL BREAK IF THERE IS MORE THEN ONE TD PER SCENE (NOT YET), BUT WE CAN fIX BY EXTRACTING THE TD ENTITIES
	// INTO A SLICE
	// AND LOOPINNG THRU THE SLICE INSTEAD OF USING THE CURSOR (ITS WHAT WERE DOING FOR CHESTS ABOVE!!!)
	trapDoorQuery := warehouse.Factory.NewQuery().And(components.TrapDoorComponent, components.PersistenceComponent)
	trapDoorCursor := warehouse.Factory.NewCursor(trapDoorQuery, sto)
	for range trapDoorCursor.Next() {
		trapDoorEn, err := trapDoorCursor.CurrentEntity()
		if err != nil {
			return err
		}
		trapDoorPersist := components.PersistenceComponent.GetFromCursor(trapDoorCursor)

		seriTrapDoor, ok := persistence.State.Get(sceneName, trapDoorPersist.EntityType, trapDoorPersist.PersistID)

		if !ok {
			continue
		}

		_, err = sto.ForceSerializedEntityWithID(seriTrapDoor, int(trapDoorEn.ID()))
		if err != nil {
			return err
		}
		err = seriTrapDoor.SetValue(trapDoorEn)
		if err != nil {
			return err
		}

	}

	return nil
}
