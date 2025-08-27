package clientsystems

import (
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
	"github.com/TheBitDrifter/concrete_echos/shared/combatdata"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/persistence"
	"github.com/TheBitDrifter/concrete_echos/shared/scenes"
)

type PlayerSpawnSystem struct{}

func (s PlayerSpawnSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	query := warehouse.Factory.NewQuery().And(components.PlayerTag)
	playerCursor := scene.NewCursor(query)
	playerCount := playerCursor.TotalMatched()
	foundSpawn := false

	if playerCount != 0 {
		return nil
	}

	spawnQuery := warehouse.Factory.NewQuery().And(components.PlayerSpawnComponent)
	spawnCursor := scene.NewCursor(spawnQuery)

	var spawn components.PlayerSpawn
	for range spawnCursor.Next() {
		found := components.PlayerSpawnComponent.GetFromCursor(spawnCursor)
		spawn = *found
		foundSpawn = true
		break
	}

	if foundSpawn {
		_, err := loadPlayerAt(spawn.X, spawn.Y, scene)
		if err != nil {
			return err
		}
		return nil
	}

	saveCheckpointQuery := warehouse.Factory.NewQuery().And(components.SaveActivationComponent)
	saveCheckpointCursor := scene.NewCursor(saveCheckpointQuery)
	cx := 0.0
	cy := 0.0

	for range saveCheckpointCursor.Next() {
		saveActivation := components.SaveActivationComponent.GetFromCursor(saveCheckpointCursor)
		checkPointPos := spatial.Components.Position.GetFromCursor(saveCheckpointCursor)
		cx = checkPointPos.X
		cy = checkPointPos.Y

		if persistence.State.LastOptionalSaveID == int(saveActivation.OptionalID) {
			break
		}
	}

	_, err := loadPlayerAt(cx, cy, scene)
	if err != nil {
		return err
	}
	return nil
}

func loadPlayerAt(x, y float64, scene coldbrew.Scene) (warehouse.Entity, error) {
	en, err := scenes.NewPlayer(x, y, scene.Storage())
	if err != nil {
		return nil, err
	}
	if persistence.State.PlayerSingleton != nil {
		en, err = scene.Storage().ForceSerializedEntityWithID(*persistence.State.PlayerSingleton, int(en.ID()))
		if err != nil {
			return nil, err
		}

		err = persistence.State.PlayerSingleton.SetValue(en)
		if err != nil {
			return nil, err
		}
		emptyHurtBoxes := combat.Components.HurtBoxes.GetFromEntity(en)

		properBoxes := combatdata.HurtBoxes[characterkeys.BoxHead]
		for i := range emptyHurtBoxes {
			emptyHurtBoxes[i] = properBoxes[i]
		}

	}

	pos := spatial.Components.Position.GetFromEntity(en)
	pos.X = x
	pos.Y = y

	if components.DodgeComponent.Check(en.Table()) {
		en.EnqueueRemoveComponent(components.DodgeComponent)
	}

	return en, nil
}
