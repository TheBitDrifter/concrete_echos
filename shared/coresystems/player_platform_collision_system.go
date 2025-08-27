package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/concrete_echos/shared/components"

	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
)

// PlayerPlatformCollisionSystem handles collisions between players and one-way platforms.
// It tracks historical player positions to determine if the player approached from above.
type PlayerPlatformCollisionSystem struct {
	// Map of player entity ID to their position history
	playerPositionHistory map[uint64][]vector.Two
	maxPositionsToTrack   int // Number of positions to track
}

// NewPlayerPlatformCollisionSystem creates a new collision system with initialized position tracking.
// It uses a pointer because the system is not pure and must retain its state.
func NewPlayerPlatformCollisionSystem() *PlayerPlatformCollisionSystem {
	const TRACK_COUNT = 5 // higher count == more tunneling protection == higher cost
	return &PlayerPlatformCollisionSystem{
		playerPositionHistory: make(map[uint64][]vector.Two),
		maxPositionsToTrack:   TRACK_COUNT,
	}
}

func (s *PlayerPlatformCollisionSystem) Run(scene blueprint.Scene, dt float64) error {
	handleVanishTimers(scene)

	platformTerrainQuery := warehouse.Factory.NewQuery().And(components.PlatformTag)
	platformCursor := scene.NewCursor(platformTerrainQuery)
	playerCursor := scene.NewCursor(blueprint.Queries.ActionBuffer)

	for range playerCursor.Next() {
		playerEntity, err := playerCursor.CurrentEntity()
		if err != nil {
			return err
		}
		playerID := uint64(playerEntity.ID())

		for range platformCursor.Next() {
			err = s.resolve(scene, platformCursor, playerCursor, playerID)
			if err != nil {
				return err
			}
		}

		playerPos := spatial.Components.Position.GetFromCursor(playerCursor)
		s.trackPosition(playerID, playerPos.Two)
	}
	return nil
}

func (s *PlayerPlatformCollisionSystem) resolve(scene blueprint.Scene, platformCursor, playerCursor *warehouse.Cursor, playerID uint64) error {
	playerEntity, err := playerCursor.CurrentEntity()
	if err != nil {
		return err
	}

	playerShape := spatial.Components.Shape.GetFromCursor(playerCursor)
	playerPosition := spatial.Components.Position.GetFromCursor(playerCursor)
	playerDynamics := motion.Components.Dynamics.GetFromCursor(playerCursor)

	platformShape := spatial.Components.Shape.GetFromCursor(platformCursor)
	platformPosition := spatial.Components.Position.GetFromCursor(platformCursor)
	platformRotation := float64(*spatial.Components.Rotation.GetFromCursor(platformCursor))
	platformDynamics := motion.Components.Dynamics.GetFromCursor(platformCursor)

	if ok, collisionResult := spatial.Detector.Check(
		*playerShape, *platformShape, playerPosition.Two, platformPosition.Two,
	); ok {

		if _, isDropping := components.IsDroppingThroughPlatformTag.GetFromCursorSafe(playerCursor); isDropping {
			// The player is trying to drop. Since we have a collision, it means they are still
			// inside the platform. We add a temporary tag to signal this.
			playerEntity.EnqueueAddComponent(components.StillCollidingWithPlatformTag)

			// By returning here, we skip all the normal collision resolution logic,
			// effectively ignoring the platform.
			return nil
		}

		// Check if any of the past player positions indicate the player was above the platform
		platformTop := platformShape.Polygon.WorldVertices[0].Y // Just using top left vert for non rotated rect platforms
		var playerWasAbove bool

		// Checking for 'above' is much easier when the edge is flat (fixed y value)
		if platformRotation == 0 {
			playerWasAbove = s.checkAnyPlayerPositionWasAbove(playerID, platformTop, playerShape.LocalAAB.Height)

			// Rotation check is more complicated using vector math to determine if player 'cleared top'
		} else {
			playerWasAbove = s.checkAnyPlayerPositionWasAboveAdvanced(
				playerID,
				// The top edge for a rotated triangle platform is always 0,1
				[]vector.Two{
					platformShape.Polygon.WorldVertices[0],
					platformShape.Polygon.WorldVertices[1],
				},
				playerShape.LocalAAB.Width, playerShape.LocalAAB.Height,
			)
		}

		if playerDynamics.Vel.Y > 0 && collisionResult.IsTopB() && playerWasAbove {

			if vp, vanishOK := components.VanishingPlatformComponent.GetFromCursorSafe(platformCursor); vanishOK {
				if !vp.TimerStarted {
					vp.TimerStarted = true
					vp.TimerStartedTick = scene.CurrentTick()
				}
			}

			// Use a vertical resolver since we can't collide with the sides
			motion.VerticalResolver.Resolve(
				&playerPosition.Two,
				&platformPosition.Two,
				playerDynamics,
				platformDynamics,
				collisionResult,
			)

			// Ground state handling
			currentTick := scene.CurrentTick()

			onGround, playerAlreadyGrounded := components.OnGroundComponent.GetFromCursorSafe(playerCursor)

			if !playerAlreadyGrounded {
				playerEntity, _ := playerCursor.CurrentEntity()
				err := playerEntity.EnqueueAddComponentWithValue(
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
	}
	return nil
}

func (s *PlayerPlatformCollisionSystem) trackPosition(playerID uint64, pos vector.Two) {
	// Initialize the position history for this player if it doesn't exist
	if _, exists := s.playerPositionHistory[playerID]; !exists {
		s.playerPositionHistory[playerID] = make([]vector.Two, 0, s.maxPositionsToTrack)
	}

	// Add the new position to this player's history
	s.playerPositionHistory[playerID] = append(s.playerPositionHistory[playerID], pos)

	// If we've exceeded our max, remove the oldest position
	if len(s.playerPositionHistory[playerID]) > s.maxPositionsToTrack {
		s.playerPositionHistory[playerID] = s.playerPositionHistory[playerID][1:]
	}
}

// checkAnyPlayerPositionWasAbove checks if the player was above a non-rotated platform in any historical position
func (s *PlayerPlatformCollisionSystem) checkAnyPlayerPositionWasAbove(playerID uint64, platformTop float64, playerHeight float64) bool {
	positions, exists := s.playerPositionHistory[playerID]
	if !exists || len(positions) == 0 {
		return false
	}

	// Check all stored positions to see if the player was above in any of them
	for _, pos := range positions {
		playerBottom := pos.Y + playerHeight/2
		if playerBottom <= platformTop {
			return true // Found at least one position where player was above
		}
	}

	return false // No positions found where player was above
}

// checkAnyPlayerPositionWasAboveAdvanced checks if the player was above a rotated platform's top edge
func (s *PlayerPlatformCollisionSystem) checkAnyPlayerPositionWasAboveAdvanced(
	playerID uint64,
	platformTopVerts []vector.Two,
	playerWidth, playerHeight float64,
) bool {
	positions, exists := s.playerPositionHistory[playerID]
	if !exists || len(positions) == 0 {
		return false
	}
	v1 := platformTopVerts[0]
	v2 := platformTopVerts[1]

	edgeVector := v2.Sub(v1)
	edgeLength := edgeVector.Mag()
	if edgeLength < 0.001 {
		return false
	}

	edgeNormalized := edgeVector.Norm()
	edgeNormal := vector.Two{X: -edgeNormalized.Y, Y: edgeNormalized.X}

	worldUp := vector.Two{X: 0, Y: -1}
	if edgeNormal.ScalarProduct(worldUp) < 0 {
		edgeNormal = edgeNormal.Scale(-1)
	}

	for _, historicalPos := range positions {
		halfHeight := playerHeight / 2
		halfWidth := playerWidth / 2
		checkPoints := []vector.Two{
			{X: historicalPos.X, Y: historicalPos.Y + halfHeight},
			{X: historicalPos.X - halfWidth, Y: historicalPos.Y + halfHeight},
			{X: historicalPos.X + halfWidth, Y: historicalPos.Y + halfHeight},
		}

		for _, point := range checkPoints {
			v1ToPoint := point.Sub(v1)
			distanceAlongNormal := v1ToPoint.ScalarProduct(edgeNormal)
			projectionOnEdge := v1ToPoint.ScalarProduct(edgeNormalized)

			const margin = 10.0
			const minAbove = 1.0
			const maxAbove = 75.0

			isAbove := distanceAlongNormal >= minAbove &&
				distanceAlongNormal < maxAbove &&
				projectionOnEdge >= -margin &&
				projectionOnEdge <= edgeLength+margin

			if isAbove {
				return true
			}
		}
	}

	return false
}

func handleVanishTimers(scene blueprint.Scene) {
	query := warehouse.Factory.NewQuery().And(components.VanishingPlatformComponent)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		platformEntity, _ := cursor.CurrentEntity()

		vp := components.VanishingPlatformComponent.GetFromCursor(cursor)

		if vp.TimerStarted && vp.TimerStartedTick+vp.LiveDuration <= scene.CurrentTick() && !vp.Vanished {
			platformEntity.EnqueueRemoveComponent(components.PlatformTag)
			vp.TimerStartedTick = scene.CurrentTick()
			vp.Vanished = true
		}

		if vp.Vanished && vp.TimerStartedTick+vp.RespawnDelay <= scene.CurrentTick() {
			platformEntity.EnqueueAddComponent(components.PlatformTag)
			vp.Vanished = false
			vp.TimerStarted = false
		}

	}
}
