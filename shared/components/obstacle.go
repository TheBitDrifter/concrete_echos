package components

type ObstacleMovementState int

const (
	StateMoving ObstacleMovementState = iota
	StatePaused
)

// The updated Obstacle component
type Obstacle struct {
	Damage     int
	MinX, MaxX float64 // Use float64 for better precision with dt
	MinY, MaxY float64

	// New Properties
	TravelTime    float64 // Time in seconds to move from min to max
	PauseDuration float64 // Time in seconds to pause at each end
	EasingType    string  // e.g., "linear", "ease-in-out"

	// Internal State (managed by the system)
	State        ObstacleMovementState
	Timer        float64 // A multi-purpose timer for movement and pausing
	IsReversingX bool    // Tracks the current direction on the X-axis
	IsReversingY bool    // Tracks the current direction on the Y-axis
}
