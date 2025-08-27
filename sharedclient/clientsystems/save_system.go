package clientsystems

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/actions"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/persistence"
)

type SaveSystem struct {
	MinSaveTicks int
}

func (sys SaveSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	activate := false
	actionsQuery := warehouse.Factory.NewQuery().And(input.Components.ActionBuffer)
	actionsCursor := scene.NewCursor(actionsQuery)
	shouldResetCurr := false

	for range actionsCursor.Next() {
		actionsBuffer := input.Components.ActionBuffer.GetFromCursor(actionsCursor)
		activate = actionsBuffer.HasAction(actions.Interact)

	}

	queryP := warehouse.Factory.NewQuery().And(
		components.PlayerTag,
	)
	cursorP := scene.NewCursor(queryP)

	queryF := warehouse.Factory.NewQuery().And(
		components.SaveActivationComponent,
	)
	cursorF := scene.NewCursor(queryF)

	var playerEn warehouse.Entity

	for range cursorP.Next() {
		playerPos := spatial.Components.Position.GetFromCursor(cursorP)
		playerEn, _ = cursorP.CurrentEntity()

		saving, savingOK := components.IsSavingComponent.GetFromCursorSafe(cursorP)
		if savingOK {
			if scene.CurrentTick() >= saving.StartedTick+sys.MinSaveTicks && saving.EndedTick == 0 {
				saving.EndedTick = scene.CurrentTick()
				shouldResetCurr = true
				continue
			}
		}

		for range cursorF.Next() {
			savePos := spatial.Components.Position.GetFromCursor(cursorF)
			saveCheckPointData := components.SaveActivationComponent.GetFromCursor(cursorF)

			currentDistSq := playerPos.Two.Sub(savePos.Two).MagSquared()
			inRange := currentDistSq <= saveCheckPointData.Range*saveCheckPointData.Range
			if inRange && activate {
				if !savingOK {
					playerDir := spatial.Components.Direction.GetFromCursor(cursorP)
					playerDir.SetRight()
					playerEn.EnqueueAddComponentWithValue(components.IsSavingComponent, components.IsSaving{StartedTick: scene.CurrentTick()})

					persistence.State.LastOptionalSaveID = int(saveCheckPointData.OptionalID)
					persistence.State.LastScene = persistence.SceneName(scene.Name())

					health := combat.Components.Health.GetFromCursor(cursorP)
					health.Value = 100 // One day we will have a max health value or something tracked.

					err := seriPlayer(cursorP)
					if err != nil {
						return err
					}
					for _, cScene := range cli.(coldbrew.Client).Cache().All() {
						err := saveRelevantState(cScene)
						if err != nil {
							return err
						}

						if scene.Name() != cScene.Name() {
							err = cScene.Reset()
							if err != nil {
								return err
							}

						}
					}
				}
			}
		}

	}
	if shouldResetCurr {
		err := clearNonDefaultPL(scene)
		if err != nil {
			return err
		}
		err = scene.Reset()
		if err != nil {
			return err
		}
		err = persistence.SaveState("ce_save.json")
		if err != nil {
			return err
		}
	}

	return nil
}

func saveRelevantState(scene coldbrew.Scene) error {
	if scene.LastActivatedTick() == 0 && !scene.IsLoaded() {
		return nil
	}

	chestQuery := warehouse.Factory.NewQuery().And(components.ChestTag, components.PersistenceComponent)
	chestCursor := scene.NewCursor(chestQuery)

	for range chestCursor.Next() {

		chestEn, err := chestCursor.CurrentEntity()
		if err != nil {
			return err
		}
		chestPersist := components.PersistenceComponent.GetFromCursor(chestCursor)
		chestDefeat, okDef := combat.Components.Defeat.GetFromCursorSafe(chestCursor)
		if okDef {
			chestDefeat.StartTick = 0
		}

		seriChest := chestEn.Serialize()
		persistence.State.Set(scene.Name(), chestPersist.EntityType, chestPersist.PersistID, seriChest)
	}

	trapDoorQuery := warehouse.Factory.NewQuery().And(components.TrapDoorComponent, components.PersistenceComponent)
	trapDoorCursor := scene.NewCursor(trapDoorQuery)

	for range trapDoorCursor.Next() {
		trapDoorEn, err := trapDoorCursor.CurrentEntity()
		if err != nil {
			return err
		}
		seriTrapDoor := trapDoorEn.Serialize()
		trapDoorPersist := components.PersistenceComponent.GetFromCursor(trapDoorCursor)

		persistence.State.Set(scene.Name(), trapDoorPersist.EntityType, trapDoorPersist.PersistID, seriTrapDoor)
	}
	return nil
}

type SavingClearingSystem struct{}

func (sys SavingClearingSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	query := warehouse.Factory.NewQuery().And(components.IsSavingComponent)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		isSaving := components.IsSavingComponent.GetFromCursor(cursor)

		if scene.CurrentTick() == isSaving.EndedTick+30 {
			en, err := cursor.CurrentEntity()
			if err != nil {
				return err
			}

			en.EnqueueRemoveComponent(components.IsSavingComponent)
		}
	}

	return nil
}

func seriPlayer(cursor *warehouse.Cursor) error {
	playerEn, err := cursor.CurrentEntity()
	if err != nil {
		return err
	}
	actions := input.Components.ActionBuffer.GetFromCursor(cursor)
	actions.Clear()

	lc := components.LastCombatComponent.GetFromCursor(cursor)
	lc.StartTick = 0

	spriteBundle := client.Components.SpriteBundle.GetFromCursor(cursor)

	for i := range spriteBundle.Blueprints {
		for k := range spriteBundle.Blueprints {
			spriteBundle.Blueprints[i].Animations[k].StartTick = 0
		}
	}

	dyn := motion.Components.Dynamics.GetFromEntity(playerEn)
	dyn.Vel.X = 0
	dyn.Vel.Y = 0

	atk, atkOK := combat.Components.Attack.GetFromCursorSafe(cursor)
	if atkOK {
		atk.StartTick = 0
	}

	seriPlayer := playerEn.Serialize()
	persistence.State.PlayerSingleton = &seriPlayer

	return nil
}
