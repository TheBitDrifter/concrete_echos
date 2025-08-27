package scenes

import (
	"github.com/TheBitDrifter/bappa/blueprint/dialogue"
	"github.com/TheBitDrifter/bappa/blueprint/ldtk"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/persistence"
)

// Registering custom LDTK entities
func init() {
	entityRegistry.Register("PlayerStart", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		// Create the player at the position defined in LDtk
		_, err := NewPlayerSpawn(float64(entity.Position[0]), float64(entity.Position[1]), sto)
		if err != nil {
			return err
		}
		return nil
	})

	entityRegistry.Register("Ramp", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		return NewRamp(
			sto,
			float64(entity.Position[0]),
			float64(entity.Position[1]),
		)
	})

	entityRegistry.Register("RotatedPlatformRight", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		return NewPlatformRotated(
			sto,
			float64(entity.Position[0]),
			float64(entity.Position[1]),
			-0.25,
		)
	})
	entityRegistry.Register("SceneTransfer", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		// Extract properties from LDtk entity
		targetScene := entity.StringFieldOr("targetScene", "")
		targetX := entity.FloatFieldOr("targetX", 20.0)
		targetY := entity.FloatFieldOr("targetY", 400.0)
		width := entity.FloatFieldOr("width", 100)
		height := entity.FloatFieldOr("height", 100)

		return NewCollisionPlayerTransfer(
			sto,
			float64(entity.Position[0]),
			float64(entity.Position[1]),
			width,
			height,
			targetX,
			targetY,
			targetScene,
		)
	})

	entityRegistry.Register("MudHandMob", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		isLeft := entity.BoolFieldOr("left", true)
		minX := entity.FloatFieldOr("minX", 0)
		maxX := entity.FloatFieldOr("maxX", 0)
		_, err := NewMudHandMob(
			float64(entity.Position[0]),
			float64(entity.Position[1]),
			sto,
			isLeft,
			minX,
			maxX,
		)
		return err
	})

	entityRegistry.Register("SkullRoninMob", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		isLeft := entity.BoolFieldOr("left", true)
		minX := entity.FloatFieldOr("minX", 0)
		maxX := entity.FloatFieldOr("maxX", 0)
		vision := entity.FloatFieldOr("vision", 0)

		_, err := NewSkullRoninMob(
			float64(entity.Position[0]),
			float64(entity.Position[1]),
			sto,
			isLeft,
			minX,
			maxX,
			vision,
		)
		return err
	})

	entityRegistry.Register("CanThrower", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		isLeft := entity.BoolFieldOr("left", true)
		vision := entity.FloatFieldOr("vision", 0)
		shotDelay := entity.FloatFieldOr("delay", 0)
		_, err := NewCanThrower(
			float64(entity.Position[0]),
			float64(entity.Position[1]),
			sto,
			isLeft,
			vision,
			shotDelay,
		)
		return err
	})

	entityRegistry.Register("CanStraightThrower", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		isLeft := entity.BoolFieldOr("left", true)
		vision := entity.FloatFieldOr("vision", 0)
		shotDelay := entity.FloatFieldOr("delay", 0)

		_, err := NewCanStraightThrower(
			float64(entity.Position[0]),
			float64(entity.Position[1]),
			sto,
			isLeft,
			vision,
			shotDelay,
		)
		return err
	})

	entityRegistry.Register("DemonAnt", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		isLeft := entity.BoolFieldOr("left", true)
		minX := entity.FloatFieldOr("minX", 0)
		maxX := entity.FloatFieldOr("maxX", 0)
		vision := entity.FloatFieldOr("vision", 0)

		_, err := NewDemonAnt(
			float64(entity.Position[0]),
			float64(entity.Position[1]),
			sto,
			isLeft,
			minX,
			maxX,
			vision,
		)
		return err
	})

	entityRegistry.Register("VanishingPlatformSmall", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		err := NewVanishingPlatformSmall(
			sto,
			float64(entity.Position[0]),
			float64(entity.Position[1]),
			0,
		)
		return err
	})

	entityRegistry.Register("VanishingPlatformMed", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		err := NewVanishingPlatformMed(
			sto,
			float64(entity.Position[0]),
			float64(entity.Position[1]),
			0,
		)
		return err
	})
	entityRegistry.Register("VanishingPlatformLarge", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		err := NewVanishingPlatformLg(
			sto,
			float64(entity.Position[0]),
			float64(entity.Position[1]),
			0,
		)
		return err
	})
	entityRegistry.Register("ObstacleBasicMain", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		minX := entity.FloatFieldOr("minX", 0)
		maxX := entity.FloatFieldOr("maxX", 0)

		minY := entity.FloatFieldOr("minY", 0)
		maxY := entity.FloatFieldOr("maxY", 0)

		travelTime := entity.FloatFieldOr("travelTime", 0)
		pauseDuration := entity.FloatFieldOr("pauseDuration", 0)
		easingType := entity.StringFieldOr("easingType", "ease-in-out")

		distX := entity.FloatFieldOr("distX", 0)
		distY := entity.FloatFieldOr("distY", 0)

		_, err := NewObstacleBasic(
			float64(entity.Position[0]),
			float64(entity.Position[1]),
			sto,
			travelTime,
			pauseDuration,
			minX,
			maxX,
			minY,
			maxY,
			easingType,
			distX,
			distY,
		)
		return err
	})

	entityRegistry.Register("SoftResetCheckpoint", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		automaticActive := entity.BoolFieldOr("automaticActive", false)
		w := entity.FloatFieldOr("w", 0)
		h := entity.FloatFieldOr("h", 0)

		err := NewSoftResetCheckpoint(

			sto,
			float64(entity.Position[0]),
			float64(entity.Position[1]),
			w,
			h,
			automaticActive,
		)
		return err
	})

	entityRegistry.Register("TrapDoor", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		isLeft := entity.BoolFieldOr("left", false)

		callbackKey := entity.IntFieldOr("callbackKey", 0)
		persistID := entity.IntFieldOr("persistID", 0)

		err := NewTrapDoor(
			sto,
			float64(entity.Position[0]),
			float64(entity.Position[1]),
			isLeft,
			components.TrapDoorEnum(callbackKey),
			persistence.PersistenceID(persistID),
		)
		return err
	})

	entityRegistry.Register("Chest", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		isLeft := entity.BoolFieldOr("left", false)

		moneyDrop := entity.FloatFieldOr("left", 400)
		persistID := entity.IntFieldOr("persistID", 0)

		callbackDropID := entity.IntFieldOr("callbackDropID", 0)

		useAlt := entity.BoolFieldOr("useAlt", false)

		_, err := NewChest(
			float64(entity.Position[0]),
			float64(entity.Position[1]),
			sto,
			isLeft,
			moneyDrop,
			persistence.PersistenceID(persistID),
			callbackDropID,
			useAlt,
		)
		return err
	})

	entityRegistry.Register("Drifter", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		isLeft := entity.BoolFieldOr("left", false)
		convoEnum := entity.IntFieldOr("convoEnum", 0)

		_, err := NewDrifterFriendly(
			float64(entity.Position[0]),
			float64(entity.Position[1]),
			sto,
			isLeft,
			convoEnum,
		)
		return err
	})

	entityRegistry.Register("SceneTitle", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		title := entity.StringFieldOr("title", "")
		_, err := NewSceneTitle(
			sto,
			title,
		)
		return err
	})

	entityRegistry.Register("SaveBench", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		isLeft := entity.BoolFieldOr("left", false)

		optionalID := entity.IntFieldOr("optionalID", 0)

		err := NewSaveBench(
			sto,
			float64(entity.Position[0]),
			float64(entity.Position[1]),
			isLeft,
			persistence.OptionalSaveID(optionalID),
		)
		return err
	})

	entityRegistry.Register("WarpTotem", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		isLeft := entity.BoolFieldOr("left", false)

		err := NewWarpTotem(
			sto,
			float64(entity.Position[0]),
			float64(entity.Position[1]),
			isLeft,
		)
		return err
	})

	entityRegistry.Register("LazySkullRonin", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		convoEnum := dialogue.SlidesEnum(entity.IntFieldOr("convoEnum", 0))
		convoCallback := dialogue.CallbackEnum(entity.IntFieldOr("convoCallback", 0))

		_, err := NewLazySkullRoninFriendly(
			float64(entity.Position[0]),
			float64(entity.Position[1]),
			sto,
			convoEnum,
			convoCallback,
		)
		return err
	})

	entityRegistry.Register("Bat", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		isLeft := entity.BoolFieldOr("left", false)

		distX := entity.FloatFieldOr("distX", 0)
		distY := entity.FloatFieldOr("distY", 0)

		_, err := NewBatFlyer(
			float64(entity.Position[0]),
			float64(entity.Position[1]),
			sto,
			isLeft,
			distX,
			distY,
		)
		return err
	})
	entityRegistry.Register("FastTravelTotem", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		isLeft := entity.BoolFieldOr("left", false)

		name := entity.StringFieldOr("name", "foo")

		err := NewTravelTotem(
			sto,
			float64(entity.Position[0]),
			float64(entity.Position[1]),
			isLeft,
			name,
		)
		return err
	})
}
