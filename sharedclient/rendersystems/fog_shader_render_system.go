package rendersystems

import (
	"math"

	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/sharedclient/shaders"
	"github.com/hajimehoshi/ebiten/v2"
)

// lerp performs linear interpolation between two points.
func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

// Moved the tempSurfaces map and time into the struct to hold state correctly.
type FogShaderRenderSystem struct {
	tempSurfaces  map[int]*ebiten.Image
	time          float32
	lastTargetPos vector.Two // Store the last frame's target position for smooth interpolation.
}

func NewFogShaderRenderSystem() *FogShaderRenderSystem {
	return &FogShaderRenderSystem{
		tempSurfaces: make(map[int]*ebiten.Image),
	}
}

func (s *FogShaderRenderSystem) Render(scene coldbrew.Scene, screen coldbrew.Screen, c coldbrew.LocalClient) {
	s.time += 0.01

	query := warehouse.Factory.NewQuery().And(components.PlayerTag)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		for _, cam := range c.ActiveCamerasFor(scene) {
			if !c.Ready(cam) {
				continue
			}

			camW, camH := cam.Surface().Bounds().Dx(), cam.Surface().Bounds().Dy()
			temp, ok := s.tempSurfaces[cam.Index()]

			if !ok || temp.Bounds().Dx() != camW || temp.Bounds().Dy() != camH {
				if temp != nil {
					temp.Dispose()
				}
				temp = ebiten.NewImage(camW, camH)
				s.tempSurfaces[cam.Index()] = temp
			}
			temp.Clear()

			// --- CONDITIONAL AVERAGING LOGIC ---

			// 1. Get positions in world space.
			playerWorldPos := spatial.Components.Position.GetFromCursor(cursor).Two
			_, cameraWorldPos := cam.Positions()
			cameraCenterWorldPos := vector.Two{
				X: cameraWorldPos.X + float64(camW)/2.0,
				Y: cameraWorldPos.Y + float64(camH)/2.0,
			}

			// 2. Calculate the distance between the player and the camera's center.
			dist := math.Hypot(playerWorldPos.X-cameraCenterWorldPos.X, playerWorldPos.Y-cameraCenterWorldPos.Y)

			// 3. Define the final target position for the fog.
			var currentTargetPos vector.Two
			const PEEK_THRESHOLD = 40.0 // Only start shifting the fog when the camera is 40 pixels away.

			if dist > PEEK_THRESHOLD {
				// If peeking, the target is halfway between the player and the camera's center.
				currentTargetPos.X = (playerWorldPos.X + cameraCenterWorldPos.X) / 2.0
				currentTargetPos.Y = (playerWorldPos.Y + cameraCenterWorldPos.Y) / 2.0
			} else {
				// Otherwise, the target is just the player's position.
				currentTargetPos = playerWorldPos
			}

			// 4. Smoothly interpolate or snap to the new target position.
			const LERP_SPEED = 0.08      // Adjust for faster or slower smoothing.
			const SNAP_THRESHOLD = 200.0 // If the distance is greater than this, snap instantly.

			// Initialize lastTargetPos on the first frame to prevent a snap from (0,0).
			if s.lastTargetPos.X == 0 && s.lastTargetPos.Y == 0 {
				s.lastTargetPos = currentTargetPos
			}

			distanceToTarget := math.Hypot(s.lastTargetPos.X-currentTargetPos.X, s.lastTargetPos.Y-currentTargetPos.Y)
			if distanceToTarget > SNAP_THRESHOLD {
				s.lastTargetPos = currentTargetPos // Snap directly to the target.
			} else {
				// Otherwise, smoothly interpolate.
				s.lastTargetPos.X = lerp(s.lastTargetPos.X, currentTargetPos.X, LERP_SPEED)
				s.lastTargetPos.Y = lerp(s.lastTargetPos.Y, currentTargetPos.Y, LERP_SPEED)
			}

			// 5. Localize the final interpolated position to get the correct coordinate for the shader.
			playerShape := spatial.Components.Shape.GetFromCursor(cursor)
			finalTargetLocal := cam.Localize(s.lastTargetPos)
			finalTargetCenterX := finalTargetLocal.X + (playerShape.LocalAAB.Width / 2.0)
			finalTargetCenterY := finalTargetLocal.Y + (playerShape.LocalAAB.Height / 2.0)

			uniforms := map[string]any{
				"PlayerPos": []float32{float32(finalTargetCenterX), float32(finalTargetCenterY)},
				"Time":      s.time,
			}

			opts := &ebiten.DrawRectShaderOptions{
				Uniforms: uniforms,
				Images:   [4]*ebiten.Image{cam.Surface()},
			}

			temp.DrawRectShader(camW, camH, shaders.FogShader, opts)
			cam.Surface().DrawImage(temp, nil)
		}
	}
}
