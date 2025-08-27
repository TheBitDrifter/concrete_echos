package coresystems

import (
	"math"

	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

type ObstacleMovementSystem struct{}

func (s ObstacleMovementSystem) Run(scene blueprint.Scene, dt float64) error {
	obsQuery := warehouse.Factory.NewQuery().And(components.ObstacleComponent)
	obsCursor := scene.NewCursor(obsQuery)

	for range obsCursor.Next() {

		obsState := components.ObstacleComponent.GetFromCursor(obsCursor)
		if obsState.TravelTime == 0 {
			continue
		}
		obsPos := spatial.Components.Position.GetFromCursor(obsCursor)

		// Update the timer based on the current state
		obsState.Timer += dt

		switch obsState.State {
		case components.StateMoving:
			// --- Handle Movement ---
			progress := obsState.Timer / obsState.TravelTime
			if progress >= 1.0 {
				progress = 1.0 // Clamp to 1.0 to ensure it lands perfectly

				// If there's a pause, switch to the paused state
				if obsState.PauseDuration > 0 {
					obsState.State = components.StatePaused
					obsState.Timer = 0 // Reset timer for pausing
				} else {
					// No pause, just reverse direction immediately
					obsState.IsReversingX = !obsState.IsReversingX
					obsState.IsReversingY = !obsState.IsReversingY
					obsState.Timer = 0 // Reset timer for the next movement cycle
				}
			}

			// Apply easing to the progress
			easedProgress := applyEasing(progress, obsState.EasingType)

			// Update position using linear interpolation (Lerp)
			if obsState.MaxX != obsState.MinX {
				if obsState.IsReversingX {
					obsPos.X = obsState.MaxX - (obsState.MaxX-obsState.MinX)*easedProgress
				} else {
					obsPos.X = obsState.MinX + (obsState.MaxX-obsState.MinX)*easedProgress
				}
			}
			if obsState.MaxY != obsState.MinY {
				if obsState.IsReversingY {
					obsPos.Y = obsState.MaxY - (obsState.MaxY-obsState.MinY)*easedProgress
				} else {
					obsPos.Y = obsState.MinY + (obsState.MaxY-obsState.MinY)*easedProgress
				}
			}

		case components.StatePaused:
			// --- Handle Pausing ---
			if obsState.Timer >= obsState.PauseDuration {
				// Pause is over, switch back to moving
				obsState.State = components.StateMoving

				obsState.Timer = 0 // Reset timer for movement

				// Reverse direction for the next movement
				obsState.IsReversingX = !obsState.IsReversingX
				obsState.IsReversingY = !obsState.IsReversingY
			}
		}
	}
	return nil
}

func applyEasing(t float64, easingType string) float64 {
	switch easingType {
	case "ease-in-out":
		// This formula creates a smooth acceleration and deceleration (Sine curve)
		return -(math.Cos(math.Pi*t) - 1) / 2
	default: // "linear"
		return t // No easing, constant speed
	}
}
