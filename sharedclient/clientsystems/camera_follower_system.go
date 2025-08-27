package clientsystems

import (
	"math"
	"math/rand/v2"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/actions"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

const (
	// Lol was supposed to be peek, but yolo
	MAX_PEE                 = 150
	HURT_SHAKE_MAGNITUDE    = 4.0
	ATK_HIT_SHAKE_MAGNITUDE = 2.0
	UNLOCK_DISTANCE         = 132.0
	STOP_VELOCITY_THRESHOLD = 0.8
	CAMERA_LERP_SPEED       = 0.04
	PEEK_LERP_SPEED         = 0.06
	BORDER_RECOIL_AMOUNT    = 30.0
	CAMERA_SNAP_THRESHOLD   = 300.0
	RECOIL_DURATION_TICKS   = 45.0 // Made float for calculations
)

type CameraFollowerSystem struct {
	LockedPositions [coldbrew.MaxSplit]*spatial.Position
	peekOffset      vector.Two

	// State-based recoil tracking
	RecoilingX      [coldbrew.MaxSplit]bool
	RecoilingY      [coldbrew.MaxSplit]bool
	recoilStartTick [coldbrew.MaxSplit]vector.Two
	recoilSign      [coldbrew.MaxSplit]vector.Two
}

func (sys *CameraFollowerSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	_, isNet := cli.(coldbrew.NetworkClient)
	if isNet {
		return nil
	}

	playerQuery := warehouse.Factory.NewQuery().And(components.PlayerTag, components.PlayerCamStateComponent)
	playerCursor := scene.NewCursor(playerQuery)

	for range playerCursor.Next() {
		camIndex := int(*client.Components.CameraIndex.GetFromCursor(playerCursor))
		cam := cli.Cameras()[camIndex]
		playerDir := spatial.Components.Direction.GetFromCursor(playerCursor)
		currentTick := scene.CurrentTick()

		baseTarget := sys.calculateFollowerTarget(
			*spatial.Components.Position.GetFromCursor(playerCursor),
			*motion.Components.Dynamics.GetFromCursor(playerCursor),
			*playerDir,
			cam,
			&sys.LockedPositions[camIndex],
		)

		peekInput := sys.getPeekInput(playerCursor)
		sys.peekOffset.X = lerp(sys.peekOffset.X, peekInput.X, PEEK_LERP_SPEED)
		sys.peekOffset.Y = lerp(sys.peekOffset.Y, peekInput.Y, PEEK_LERP_SPEED)

		var bounceOffset vector.Two

		// Handle Recoil State Machine
		if sys.RecoilingX[camIndex] {
			elapsed := float64(currentTick) - sys.recoilStartTick[camIndex].X
			if elapsed > RECOIL_DURATION_TICKS {
				sys.RecoilingX[camIndex] = false
			} else {
				sys.peekOffset.X = 0 // Suppress peek input

				// Calculate bounce using an ease-out sine curve
				progress := elapsed / RECOIL_DURATION_TICKS
				easeOut := math.Sin(progress * math.Pi) // Creates an arc from 0 -> 1 -> 0
				bounceOffset.X = BORDER_RECOIL_AMOUNT * easeOut * sys.recoilSign[camIndex].X
			}
		}
		if sys.RecoilingY[camIndex] {
			elapsed := float64(currentTick) - sys.recoilStartTick[camIndex].Y
			if elapsed > RECOIL_DURATION_TICKS {
				sys.RecoilingY[camIndex] = false
			} else {
				sys.peekOffset.Y = 0 // Suppress peek input

				progress := elapsed / RECOIL_DURATION_TICKS
				easeOut := math.Sin(progress * math.Pi)
				bounceOffset.Y = BORDER_RECOIL_AMOUNT * easeOut * sys.recoilSign[camIndex].Y
			}
		}

		desiredTarget := vector.Two{
			X: baseTarget.X + sys.peekOffset.X,
			Y: baseTarget.Y + sys.peekOffset.Y,
		}

		peekClampedTarget, _ := sys.clampPeekDistance(baseTarget, desiredTarget, peekInput)
		finalClampedTarget, sceneHitX, sceneHitY := sys.clampToSceneBoundaries(peekClampedTarget, cam, scene, peekInput)

		_, cameraPos := cam.Positions()

		// Trigger Recoil State
		if sceneHitX && !sys.RecoilingX[camIndex] {
			sys.peekOffset.X = 0 // Reset peek momentum on impact
			sys.RecoilingX[camIndex] = true
			sys.recoilStartTick[camIndex].X = float64(currentTick)
			sys.recoilSign[camIndex].X = -1.0
			if peekInput.X < 0 {
				sys.recoilSign[camIndex].X = 1.0
			}
			playerCamState := components.PlayerCamStateComponent.GetFromCursor(playerCursor)
			playerCamState.LastBorderHitTick = currentTick
		}

		if sceneHitY && !sys.RecoilingY[camIndex] {
			sys.peekOffset.Y = 0 // Reset peek momentum on impact
			sys.RecoilingY[camIndex] = true
			sys.recoilStartTick[camIndex].Y = float64(currentTick)
			sys.recoilSign[camIndex].Y = -1.0
			if peekInput.Y < 0 {
				sys.recoilSign[camIndex].Y = 1.0
			}
			playerCamState := components.PlayerCamStateComponent.GetFromCursor(playerCursor)
			playerCamState.LastBorderHitTick = currentTick
		}

		finalTarget := vector.Two{
			X: finalClampedTarget.X + bounceOffset.X,
			Y: finalClampedTarget.Y + bounceOffset.Y,
		}

		distanceToTarget := math.Hypot(cameraPos.X-finalTarget.X, cameraPos.Y-finalTarget.Y)
		if distanceToTarget > CAMERA_SNAP_THRESHOLD {
			cameraPos.X = finalTarget.X
			cameraPos.Y = finalTarget.Y
		} else {
			cameraPos.X = lerp(cameraPos.X, finalTarget.X, CAMERA_LERP_SPEED)
			cameraPos.Y = lerp(cameraPos.Y, finalTarget.Y, CAMERA_LERP_SPEED)
		}

		isHurt := combat.Components.Hurt.CheckCursor(playerCursor)
		if isHurt {
			sys.shake(cam, HURT_SHAKE_MAGNITUDE)
		}
		atk, isHitting := combat.Components.Attack.GetFromCursorSafe(playerCursor)
		if isHitting && math.Abs(float64(currentTick)-float64(atk.LastHitTick)) < 30 {
			sys.shake(cam, ATK_HIT_SHAKE_MAGNITUDE)
		}
		softReset, softOK := components.SoftResetComponent.GetFromCursorSafe(playerCursor)
		if softOK && math.Abs(float64(currentTick-softReset.StartedTick)) < 60 {
			sys.shake(cam, HURT_SHAKE_MAGNITUDE*1.2)
		}
	}
	return nil
}

func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

func (sys *CameraFollowerSystem) calculateFollowerTarget(
	playerPos spatial.Position,
	playerDyn motion.Dynamics,
	dir spatial.Direction,
	cam coldbrew.Camera,
	lockedState **spatial.Position,
) vector.Two {
	_, cameraScenePosition := cam.Positions()

	playerIsStopped := math.Abs(playerDyn.Vel.X) < STOP_VELOCITY_THRESHOLD &&
		math.Abs(playerDyn.Vel.Y) < STOP_VELOCITY_THRESHOLD

	yMovement := 0
	if playerDyn.Vel.Y > 5 {
		yMovement = 1
	}
	if playerDyn.Vel.Y < -5 {
		yMovement = -1
	}
	centerX := float64(cam.Surface().Bounds().Dx())/2 - (dir.AsFloat() * 100)
	centerY := float64(cam.Surface().Bounds().Dy())/2 - float64(yMovement*100)

	isPeeking := sys.peekOffset.X != 0 || sys.peekOffset.Y != 0

	if *lockedState != nil {
		lockedFocalPointX := (*lockedState).X + centerX
		lockedFocalPointY := (*lockedState).Y + centerY
		distFromLock := math.Hypot(playerPos.X-lockedFocalPointX, playerPos.Y-lockedFocalPointY)

		if distFromLock > UNLOCK_DISTANCE {
			*lockedState = nil
		}
	} else if playerIsStopped && !isPeeking {
		newPos := spatial.NewPosition(cameraScenePosition.X, cameraScenePosition.Y)
		*lockedState = &newPos
	}

	if *lockedState != nil && !isPeeking {
		return (*lockedState).Two
	}

	centeredPlayerX := playerPos.X
	centeredPlayerY := playerPos.Y
	centeredCameraX := cameraScenePosition.X + centerX
	centeredCameraY := cameraScenePosition.Y + centerY

	diffX := centeredPlayerX - centeredCameraX
	diffY := centeredPlayerY - centeredCameraY

	deadzoneX := 60.0
	deadzoneY := 60.0

	targetX := cameraScenePosition.X
	targetY := cameraScenePosition.Y

	if math.Abs(diffX) > deadzoneX {
		if diffX > 0 {
			targetX = centeredPlayerX - centerX - deadzoneX
		} else {
			targetX = centeredPlayerX - centerX + deadzoneX
		}
	}
	if math.Abs(diffY) > deadzoneY {
		if diffY > 0 {
			targetY = centeredPlayerY - centerY - deadzoneY
		} else {
			targetY = centeredPlayerY - centerY + deadzoneY
		}
	}
	return vector.Two{X: targetX, Y: targetY}
}

func (sys *CameraFollowerSystem) getPeekInput(cursor *warehouse.Cursor) vector.Two {
	buffer := input.Components.ActionBuffer.GetFromCursor(cursor)
	targetOffset := vector.Two{X: 0, Y: 0}

	if buffer.HasAction(actions.CameraRight) {
		targetOffset.X = MAX_PEE
	}
	if buffer.HasAction(actions.CameraLeft) {
		targetOffset.X = -MAX_PEE
	}
	if buffer.HasAction(actions.CameraDown) {
		targetOffset.Y = MAX_PEE
	}
	if buffer.HasAction(actions.CameraUp) {
		targetOffset.Y = -MAX_PEE
	}

	buffer.ConsumeAction(actions.CameraRight)
	buffer.ConsumeAction(actions.CameraLeft)
	buffer.ConsumeAction(actions.CameraDown)
	buffer.ConsumeAction(actions.CameraUp)

	return targetOffset
}

func (sys *CameraFollowerSystem) clampPeekDistance(base, desired, peekInput vector.Two) (vector.Two, bool) {
	offset := vector.Two{X: desired.X - base.X, Y: desired.Y - base.Y}
	dist := math.Hypot(offset.X, offset.Y)
	var hit bool

	if math.Abs(dist-MAX_PEE) < 1 {
		hit = true
	}
	if dist > MAX_PEE {
		ratio := MAX_PEE / dist
		clampedOffset := vector.Two{X: offset.X * ratio, Y: offset.Y * ratio}

		if peekInput.X != 0 || peekInput.Y != 0 {
			return vector.Two{X: base.X + clampedOffset.X, Y: base.Y + clampedOffset.Y}, hit
		}
		return vector.Two{X: base.X + clampedOffset.X, Y: base.Y + clampedOffset.Y}, hit
	}

	return desired, hit
}

func (CameraFollowerSystem) clampToSceneBoundaries(target vector.Two, cam coldbrew.Camera, scene coldbrew.Scene, peekInput vector.Two) (vector.Two, bool, bool) {
	sceneWidth := scene.Width()
	sceneHeight := scene.Height()
	camWidth, camHeight := cam.Dimensions()
	maxX := float64(sceneWidth - camWidth)
	maxY := float64(sceneHeight - camHeight)

	clampedTarget := target
	hitX := false
	hitY := false

	if clampedTarget.X > maxX {
		clampedTarget.X = maxX
		if peekInput.X > 0 {
			hitX = true
		}
	}
	if clampedTarget.X < 0 {
		clampedTarget.X = 0
		if peekInput.X < 0 {
			hitX = true
		}
	}
	if clampedTarget.Y > maxY {
		clampedTarget.Y = maxY
		if peekInput.Y > 0 {
			hitY = true
		}
	}
	if clampedTarget.Y < 0 {
		clampedTarget.Y = 0
		if peekInput.Y < 0 {
			hitY = true
		}
	}

	return clampedTarget, hitX, hitY
}

func (sys *CameraFollowerSystem) shake(cam coldbrew.Camera, mag float64) {
	_, camPos := cam.Positions()
	offsetX := (rand.Float64()*2 - 1) * mag
	offsetY := (rand.Float64()*2 - 1) * mag
	camPos.X += offsetX
	camPos.Y += offsetY / 2
}
