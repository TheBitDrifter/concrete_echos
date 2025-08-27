package clientsystems

import (
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/actions"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/persistence"
	"github.com/TheBitDrifter/concrete_echos/shared/scenes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type FastTravelSceneActivationSystem struct{}

func (sys FastTravelSceneActivationSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	activate := false
	actionsQuery := warehouse.Factory.NewQuery().And(input.Components.ActionBuffer)
	actionsCursor := scene.NewCursor(actionsQuery)

	for range actionsCursor.Next() {
		actionsBuffer := input.Components.ActionBuffer.GetFromCursor(actionsCursor)
		activate = actionsBuffer.HasAction(actions.Interact)

	}

	queryP := warehouse.Factory.NewQuery().And(
		components.PlayerTag,
	)
	cursorP := scene.NewCursor(queryP)

	queryF := warehouse.Factory.NewQuery().And(
		components.FastTravelActivationComponent,
	)
	cursorF := scene.NewCursor(queryF)
	transfer := false

	for range cursorP.Next() {
		playerPos := spatial.Components.Position.GetFromCursor(cursorP)

		for range cursorF.Next() {
			checkpoint := components.FastTravelActivationComponent.GetFromCursor(cursorF)
			rangeDist := checkpoint.Range
			unlockRangeDist := checkpoint.UnlockRange
			if unlockRangeDist == 0 {
				unlockRangeDist = 150
			}

			ffTotemPos := spatial.Components.Position.GetFromCursor(cursorF)

			currentDistSq := playerPos.Two.Sub(ffTotemPos.Two).MagSquared()
			inRange := currentDistSq <= rangeDist*rangeDist
			inUnlockRange := currentDistSq <= unlockRangeDist*unlockRangeDist

			if inUnlockRange {
				_, ok := persistence.State.FtCheckpointsMap[persistence.FastTravelCheckpointName(checkpoint.Name)]
				if !ok {
					newCp := persistence.FastTravelCheckpoint{
						FastTravelCheckpointName: persistence.FastTravelCheckpointName(checkpoint.Name),
						SceneName:                persistence.SceneName(scene.Name()),
						DropOff:                  *ffTotemPos,
					}
					persistence.State.FtCheckpointsMap[persistence.FastTravelCheckpointName(checkpoint.Name)] = newCp
					persistence.State.FtCheckpoints = append(persistence.State.FtCheckpoints, newCp)
					err := seriPlayer(cursorP)
					if err != nil {
						return err
					}

				}
			}

			if inRange && activate {
				playerEn, err := cursorP.CurrentEntity()
				if err != nil {
					return err
				}
				scene.Storage().EnqueueDestroyEntities(playerEn)
				transfer = true

			}
		}
	}
	if transfer {
		err := clearNonDefaultPL(scene)
		if err != nil {
			return err
		}

		_, err = cli.ActivateSceneByName(scenes.FAST_TRAVEL_SCENE.Name)
		if err != nil {
			return err
		}

		for _, cam := range cli.ActiveCamerasFor(scene) {
			_, local := cam.Positions()
			local.X = 0
			local.Y = 0
		}
		cli.(coldbrew.Client).DeactivateScene(scene)
	}

	return nil
}

type FastTravelSceneSystem struct {
	CheckPointCount int
	ButtonPositions []spatial.Position
	ButtonShapes    []spatial.Shape

	PaddingX     float64
	PaddingY     float64
	SpacingY     float64
	ButtonSqSize float64

	FONT_FACE *text.GoTextFace
}

func (sys FastTravelSceneSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	count := 0
	for range persistence.State.FtCheckpoints {
		count++
	}
	if count != sys.CheckPointCount {
		sys.ButtonPositions = []spatial.Position{}
		sys.ButtonShapes = []spatial.Shape{}

		i := 0
		for range persistence.State.FtCheckpoints {
			sys.ButtonPositions = append(sys.ButtonPositions, spatial.NewPosition(sys.PaddingX, sys.PaddingY+(float64(i)*sys.SpacingY)))
			sys.ButtonShapes = append(sys.ButtonShapes, spatial.NewRectangle(sys.ButtonSqSize, sys.ButtonSqSize))
			i++
		}

	}

	for _, c := range cli.ActiveCamerasFor(scene) {
		_, local := c.Positions()
		nx, ny := ebiten.CursorPosition()

		nx = nx + int(local.X)
		ny = ny + int(local.Y)
		clickPos := spatial.NewPosition(float64(nx), float64(ny))

		for i, btnShape := range sys.ButtonShapes {
			btnPos := sys.ButtonPositions[i]

			textWidth, _ := text.Measure(string(persistence.State.FtCheckpoints[i].FastTravelCheckpointName), sys.FONT_FACE, 0)

			btnPosAdj := spatial.NewPosition(btnPos.X+textWidth, btnPos.Y)
			clickShape := spatial.NewRectangle(5, 5)

			clicked, _ := spatial.Detector.Check(
				clickShape,
				btnShape,
				clickPos,
				btnPosAdj,
			)

			if clicked && inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
				targetName := string(persistence.State.FtCheckpoints[i].FastTravelCheckpointName)
				dest := persistence.State.FtCheckpointsMap[persistence.FastTravelCheckpointName(targetName)]

				en, err := loadPlayerAt(dest.DropOff.X, dest.DropOff.Y, scene)
				if err != nil {
					return err
				}

				_, err = cli.ActivateSceneByName(string(dest.SceneName), en)
				if err != nil {
					return err
				}
				cli.(coldbrew.Client).DeactivateScene(scene)
				for _, cam := range cli.ActiveCamerasFor(scene) {
					_, local := cam.Positions()
					local.X = dest.DropOff.X - 320
					local.Y = dest.DropOff.Y - 180

				}
			}

		}

	}
	return nil
}
