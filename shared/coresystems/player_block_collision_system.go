package coresystems

import (
	"math"

	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

type PlayerBlockCollisionSystem struct{}

func (s PlayerBlockCollisionSystem) Run(scene blueprint.Scene, dt float64) error {
	blockTerrainQuery := warehouse.Factory.NewQuery().And(components.BlockTerrainTag)
	blockTerrainCursor := scene.NewCursor(blockTerrainQuery)
	playerCursor := scene.NewCursor(blueprint.Queries.ActionBuffer)

	for range blockTerrainCursor.Next() {
		for range playerCursor.Next() {
			// Delegate to helper
			err := s.resolve(scene, blockTerrainCursor, playerCursor)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (PlayerBlockCollisionSystem) resolve(scene blueprint.Scene, blockCursor, playerCursor *warehouse.Cursor) error {
	playerPosition := spatial.Components.Position.GetFromCursor(playerCursor)
	playerDir := spatial.Components.Direction.GetFromCursor(playerCursor)
	playerShape := spatial.Components.Shape.GetFromCursor(playerCursor)
	playerDynamics := motion.Components.Dynamics.GetFromCursor(playerCursor)
	onGround, playerAlreadyGrounded := components.OnGroundComponent.GetFromCursorSafe(playerCursor)
	onWall, playerAlreadyWallGrounded := components.OnWallComponent.GetFromCursorSafe(playerCursor)

	blockPosition := spatial.Components.Position.GetFromCursor(blockCursor)
	blockShape := spatial.Components.Shape.GetFromCursor(blockCursor)
	blockDynamics := motion.Components.Dynamics.GetFromCursor(blockCursor)
	jumpState := components.JumpStateComponent.GetFromCursor(playerCursor)

	// Check for a collision
	if ok, collisionResult := spatial.Detector.Check(
		*playerShape, *blockShape, playerPosition.Two, blockPosition.Two,
	); ok {

		// Determine if ground is sloped
		n := collisionResult.Normal
		horizontal := n.X == 0 && n.Y == 1 || n.X == 0 && n.Y == -1
		vertical := n.X == -1 && n.Y == 0 || n.X == 1 && n.Y == 0
		isSloped := !horizontal && !vertical

		if playerAlreadyGrounded && onGround.LastTouch != scene.CurrentTick()-1 {
			playerAlreadyGrounded = false
		}

		// Prevents snapping on AAB corner transitions/collisions
		playerOnTopOfBlock := collisionResult.IsTopB()
		blockOnTopOfPlayer := collisionResult.IsTop()

		hasWallJump := components.WallJumpUnlockedTag.CheckCursor(playerCursor)

		if !playerOnTopOfBlock && !blockOnTopOfPlayer {
			if playerAlreadyWallGrounded {
				onWall.LastTouch = scene.CurrentTick()
				jumpState.WallJumpDirection = getWallJumpDirection(collisionResult.Normal, *playerDir)
			} else {
				playerEn, _ := playerCursor.CurrentEntity()
				playerEn.EnqueueAddComponentWithValue(components.OnWallComponent, components.OnWall{LastTouch: scene.CurrentTick()})
				jumpState.WallJumpDirection = getWallJumpDirection(collisionResult.Normal, *playerDir)
			}
			if playerDynamics.Vel.Y > 0 && hasWallJump {
				playerDynamics.Vel.Y = 40
			}
		}

		if playerOnTopOfBlock && playerDynamics.Vel.Y < 0 && !isSloped {
			return nil
		}
		if blockOnTopOfPlayer && playerDynamics.Vel.Y > 0 {
			return nil
		}

		// Vertical resolver to prevent positional sliding on slopes
		if isSloped {
			motion.VerticalResolver.Resolve(
				&playerPosition.Two,
				&blockPosition.Two,
				playerDynamics,
				blockDynamics,
				collisionResult,
			)
		} else {
			// Otherwise resolve as normal
			motion.Resolver.Resolve(
				&playerPosition.Two,
				&blockPosition.Two,
				playerDynamics,
				blockDynamics,
				collisionResult,
			)
		}

		// Ensure the player is on top of the terrain before marking them as grounded
		if !playerOnTopOfBlock {
			return nil
		}

		currentTick := scene.CurrentTick()

		// Update onGround accordingly (create or update)
		if !playerAlreadyGrounded {
			playerEntity, err := playerCursor.CurrentEntity()
			if err != nil {
				return err
			}
			// We cannot mutate during a cursor iteration, so we use the enqueue API
			err = playerEntity.EnqueueAddComponentWithValue(
				components.OnGroundComponent,
				components.OnGround{LastTouch: currentTick, Landed: currentTick, SlopeNormal: collisionResult.Normal},
			)
			if err != nil {
				return err
			}
		} else {
			onGround.LastTouch = scene.CurrentTick()
			onGround.SlopeNormal = collisionResult.Normal
		}

	}
	return nil
}

func getWallJumpDirection(collisionNormal vector.Two, currentDir spatial.Direction) spatial.Direction {
	if math.Abs(collisionNormal.X) > math.Abs(collisionNormal.Y) {
		if collisionNormal.X > 0 {
			return spatial.NewDirectionRight()
		} else {
			return spatial.NewDirectionLeft()
		}
	}

	return currentDir
}
