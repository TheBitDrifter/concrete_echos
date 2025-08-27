package coresystems

import (
	"math"

	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/actions"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
	"github.com/TheBitDrifter/concrete_echos/shared/combatdata"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

const (
	SPEED_X_DEFAULT   = 120.0
	HURT_KNOCK_BACK_X = 120

	SNAP_FORCE = 40.0
	JUMP_FORCE = 320.0

	WALL_JUMP_FORCE_Y = 245.0
	WALL_JUMP_FORCE_X = 140.0
	WALL_JUMP_DELAY_A = 4
	WALL_JUMP_DELAY_B = 8
)

var INPUT_DISABLED_COMPONENTS = []warehouse.Component{
	components.SoftResetComponent,
	components.InConversationComponent,
	components.IsSavingComponent,
	combat.Components.Defeat,
	components.PlayerIsExecutingComponent,
}

// PlayerMovementSystem handles all player movement mechanics including horizontal
// movement on flat ground and slopes, jumping with coyote time + early jump buffering,
// and platform drop-through functionality.
type PlayerMovementSystem struct{}

func (sys PlayerMovementSystem) Run(scene blueprint.Scene, dt float64) error {
	jumpReleaseQuery := warehouse.Factory.NewQuery().And(components.JumpStateComponent, input.Components.ActionBuffer)
	jumpReleaseCursor := scene.NewCursor(jumpReleaseQuery)

	for range jumpReleaseCursor.Next() {
		jumpState := components.JumpStateComponent.GetFromCursor(jumpReleaseCursor)
		actionBuffer := input.Components.ActionBuffer.GetFromCursor(jumpReleaseCursor)
		if _, ok := actionBuffer.ConsumeAction(actions.JumpReleased); ok {
			jumpState.LastJumpRelease = scene.CurrentTick()
			jumpState.Locked = false
		}
	}

	err := sys.handleDown(scene)
	if err != nil {
		return err
	}
	sys.handleHorizontal(scene)
	sys.handleJump(scene)
	sys.handleWallJump(scene)
	sys.handleDodge(scene)
	return nil
}

// handleHorizontal processes left/right movement with different behaviors for:
// - Air movement
// - Flat ground movement
// - Uphill/downhill slope movement with proper tangent calculations
func (PlayerMovementSystem) handleHorizontal(scene blueprint.Scene) {
	query := warehouse.Factory.NewQuery().And(
		components.PlayerTag,
		warehouse.Factory.NewQuery().Not(
			INPUT_DISABLED_COMPONENTS,
		),
	)
	cursor := scene.NewCursor(query)
	currentTick := scene.CurrentTick()

	for range cursor.Next() {
		dyn := motion.Components.Dynamics.GetFromCursor(cursor)
		incomingActions := input.Components.ActionBuffer.GetFromCursor(cursor)
		direction := spatial.Components.Direction.GetFromCursor(cursor)
		hurt, isHurt := combat.Components.Hurt.GetFromCursorSafe(cursor)
		jumpState := components.JumpStateComponent.GetFromCursor(cursor)
		atk, isAttacking := combat.Components.Attack.GetFromCursorSafe(cursor)

		// Check ground status
		onGround, isGroundComponentPresent := components.OnGroundComponent.GetFromCursorSafe(cursor)
		isGrounded := isGroundComponentPresent && currentTick-1 == onGround.LastTouch

		if (isAttacking && isGrounded) || (isAttacking && atk.ID != combatdata.AerialSeqs[characterkeys.BoxHead].First().ID) {
			return
		}

		var SPEED_X_ADJUSTED float64
		var KNOCK_BACK_X float64

		if isHurt {
			SPEED_X_ADJUSTED = SPEED_X_DEFAULT / 8
			KNOCK_BACK_X = HURT_KNOCK_BACK_X * hurt.Direction.X
		} else {
			SPEED_X_ADJUSTED = SPEED_X_DEFAULT
		}

		if scene.CurrentTick()-jumpState.LastWallJump <= WALL_JUMP_DELAY_A && jumpState.LastWallJump != 0 && !isHurt {
			dyn.Vel.X = -WALL_JUMP_FORCE_X * jumpState.WallJumpDirection.AsFloat()
			return
		}
		if scene.CurrentTick()-jumpState.LastWallJump <= WALL_JUMP_DELAY_B && jumpState.LastWallJump != 0 && !isHurt {
			KNOCK_BACK_X = WALL_JUMP_FORCE_X * jumpState.WallJumpDirection.AsFloat()
		}

		_, pressedLeft := incomingActions.ConsumeAction(actions.Left)
		if pressedLeft && !isHurt {
			direction.SetLeft()
		}

		_, pressedRight := incomingActions.ConsumeAction(actions.Right)
		if pressedRight && !isHurt {
			direction.SetRight()
		}

		isMovingHorizontal := (pressedLeft || pressedRight)

		if !isMovingHorizontal && KNOCK_BACK_X != 0 {
			dyn.Vel.X = -KNOCK_BACK_X
		}

		// Horizontal airborne movement
		if !isGrounded {
			if isMovingHorizontal {
				dyn.Vel.X = (SPEED_X_ADJUSTED * direction.AsFloat()) - KNOCK_BACK_X
			}
			// Skip grounded horizontal movement
			continue
		}

		// Apply small downward force to keep player attached to slopes when grounded
		dyn.Vel.Y = math.Max(dyn.Vel.Y, SNAP_FORCE)

		// Horizontal flat movement
		flat := onGround.SlopeNormal.X == 0 && onGround.SlopeNormal.Y == 1
		if flat {
			if isMovingHorizontal {
				dyn.Vel.X = SPEED_X_ADJUSTED*direction.AsFloat() - KNOCK_BACK_X
			}
			// Skip slope horizontal movement
			continue
		}

		// Horizontal sloped movement
		if isMovingHorizontal {
			// Calculate tangent vector along the slope
			tangent := onGround.SlopeNormal.Perpendicular()

			isUphill := (direction.AsFloat() * onGround.SlopeNormal.X) > 0

			slopeDir := tangent.Scale(direction.AsFloat())

			if isUphill {
				// When going uphill, only set X velocity and let physics handle Y
				dyn.Vel.X = slopeDir.X * SPEED_X_ADJUSTED
			} else {
				// When going downhill, help player follow the slope with both X and Y velocities
				dyn.Vel.X = slopeDir.X * SPEED_X_ADJUSTED
				dyn.Vel.Y = slopeDir.Y * SPEED_X_ADJUSTED
			}
		}
	}
}

// handleJump processes jump inputs with coyote time and input buffering features
// Coyote time: Player can jump shortly after leaving a platform
// Input buffering: Jump inputs are remembered and applied when landing
func (PlayerMovementSystem) handleJump(scene blueprint.Scene) {
	playersEligibleToJumpQuery := warehouse.Factory.NewQuery().
		And(
			components.OnGroundComponent,
			input.Components.ActionBuffer,
			warehouse.Factory.NewQuery().Not(INPUT_DISABLED_COMPONENTS),
		)

	cursor := scene.NewCursor(playersEligibleToJumpQuery)
	currentTick := scene.CurrentTick()

	for range cursor.Next() {
		// Get required components
		dyn := motion.Components.Dynamics.GetFromCursor(cursor)
		incomingActions := input.Components.ActionBuffer.GetFromCursor(cursor)
		jumpState := components.JumpStateComponent.GetFromCursor(cursor)

		onGround := components.OnGroundComponent.GetFromCursor(cursor)

		if stampedAction, actionReceived := incomingActions.ConsumeAction(actions.Jump); actionReceived {

			playerHasNotJumpedSinceGroundTouch := jumpState.LastJump <= onGround.LastTouch
			actionAfterRelease := stampedAction.Tick > jumpState.LastJumpRelease
			canJump := playerHasNotJumpedSinceGroundTouch && actionAfterRelease && !jumpState.Locked

			if canJump {
				dyn.Vel.Y = -JUMP_FORCE
				dyn.Accel.Y = -JUMP_FORCE
				jumpState.LastJump = currentTick
				jumpState.Locked = true
			}
		}
	}
}

func (PlayerMovementSystem) handleDown(scene blueprint.Scene) error {
	query := warehouse.Factory.NewQuery()
	query.And(components.OnGroundComponent, components.PlayerTag,
		warehouse.Factory.NewQuery().Not(INPUT_DISABLED_COMPONENTS),
	)

	cursor := scene.NewCursor(query)
	currentTick := scene.CurrentTick()

	for range cursor.Next() {

		incomingActions := input.Components.ActionBuffer.GetFromCursor(cursor)
		if stampedAction, inputReceived := incomingActions.PeekLatestOfType(actions.Down); inputReceived {
			if currentTick != stampedAction.Tick {
				continue
			}
			incomingActions.ConsumeAction(actions.Down)

			playerEntity, _ := cursor.CurrentEntity()

			playerEntity.EnqueueAddComponent(components.IsDroppingThroughPlatformTag)

			// playerEntity.EnqueueRemoveComponent(components.OnGroundComponent)
		}
	}
	return nil
}

func (PlayerMovementSystem) handleWallJump(scene blueprint.Scene) {
	const cd = 30
	playersEligibleToWallJumpQuery := warehouse.Factory.NewQuery().
		And(
			components.OnWallComponent,
			input.Components.ActionBuffer,
			components.WallJumpUnlockedTag,
			warehouse.Factory.NewQuery().Not(INPUT_DISABLED_COMPONENTS),
		)

	cursor := scene.NewCursor(playersEligibleToWallJumpQuery)
	currentTick := scene.CurrentTick()

	for range cursor.Next() {
		// Get required components
		dyn := motion.Components.Dynamics.GetFromCursor(cursor)
		incomingActions := input.Components.ActionBuffer.GetFromCursor(cursor)
		jumpState := components.JumpStateComponent.GetFromCursor(cursor)

		if scene.CurrentTick()-jumpState.LastWallJump <= cd {
			return
		}

		onWall := components.OnWallComponent.GetFromCursor(cursor)

		if stampedAction, actionReceived := incomingActions.ConsumeAction(actions.Jump); actionReceived {

			playerHasNotJumpedSinceGroundTouch := jumpState.LastJump <= onWall.LastTouch
			actionAfterRelease := stampedAction.Tick > jumpState.LastJumpRelease
			canJump := playerHasNotJumpedSinceGroundTouch && actionAfterRelease && !jumpState.Locked

			if canJump {
				dyn.Vel.Y = -WALL_JUMP_FORCE_Y
				dyn.Accel.Y = -WALL_JUMP_FORCE_Y
				jumpState.LastJump = currentTick
				jumpState.LastWallJump = currentTick
				jumpState.Locked = true

			}
		}
	}
}

var PlayerDodge = combat.NewIR()

func (PlayerMovementSystem) handleDodge(scene blueprint.Scene) error {
	// for yolo removal (should be isolated system...)
	playersThatAreDodging := warehouse.Factory.NewQuery().And(
		components.DodgeComponent,
		components.PlayerTag,
	)
	playersThatAreDodgingCursor := scene.NewCursor(playersThatAreDodging)

	playersThatAreNotDodging := warehouse.Factory.NewQuery().And(
		input.Components.ActionBuffer,
		warehouse.Factory.NewQuery().Not(
			components.DodgeComponent,
			combat.Components.Hurt,
			combat.Components.Attack,
			INPUT_DISABLED_COMPONENTS,
		),
	)

	for range playersThatAreDodgingCursor.Next() {
		dodge := components.DodgeComponent.GetFromCursor(playersThatAreDodgingCursor)
		if scene.CurrentTick()-dodge.StartTick > 15 {
			en, err := playersThatAreDodgingCursor.CurrentEntity()
			if err != nil {
				return err
			}
			err = en.EnqueueRemoveComponent(components.DodgeComponent)
			if err != nil {
				return err
			}
		}

		dyn := motion.Components.Dynamics.GetFromCursor(playersThatAreDodgingCursor)
		direction := spatial.Components.Direction.GetFromCursor(playersThatAreDodgingCursor)
		dyn.Vel.X = direction.AsFloat() * 350

	}

	notDodgingCursor := scene.NewCursor(playersThatAreNotDodging)
	currentTick := scene.CurrentTick()

	for range notDodgingCursor.Next() {
		buffer := input.Components.ActionBuffer.GetFromCursor(notDodgingCursor)
		movementCDs := components.MovementCooldownComponent.GetFromCursor(notDodgingCursor)

		if combat.Components.Hurt.CheckCursor(notDodgingCursor) {
			continue
		}

		if stampedAction, inputReceived := buffer.ConsumeAction(actions.Dodge); inputReceived {
			if currentTick != stampedAction.Tick || !movementCDs.Dodge.Available(scene.CurrentTick()) {
				continue
			}

			playerEntity, err := notDodgingCursor.CurrentEntity()
			if err != nil {
				return err
			}

			err = playerEntity.EnqueueAddComponentWithValue(components.DodgeComponent, components.Dodge{StartTick: scene.CurrentTick()})
			if err != nil {
				return err
			}
			err = playerEntity.EnqueueAddComponentWithValue(combat.Components.Invincible, combat.Invincible{StartTick: scene.CurrentTick(), Reason: PlayerDodge})
			if err != nil {
				return err
			}

			movementCDs.Dodge.StartTick = scene.CurrentTick()
			movementCDs.Dodge.Duration = 60

		}
	}

	return nil
}
