package scenes

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/animations"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

var OBS_BUNDLE_MAIN = client.NewSpriteBundle().
	AddSprite("images/obstacles/main_sheet.png", true).
	WithAnimations(
		animations.ObstacleMainAnim,
	).
	WithOffset(vector.Two{X: -32, Y: -32}).
	WithPriority(20)

func NewObstacleBasic(x, y float64, sto warehouse.Storage, travelTime, pauseDuration, minX, maxX, minY, maxY float64, easingType string, distX, distY float64) (warehouse.Entity, error) {
	comps := []warehouse.Component{}
	comps = append(comps, DefaultObstacleComposition...)

	arche, err := sto.NewOrExistingArchetype(comps...)
	if err != nil {
		return nil, err
	}

	if distX != 0 {
		minX = x
		maxX = minX + distX
	}

	if distY != 0 {
		minY = y
		maxY = minY + distY
	}

	dir := spatial.NewDirectionRight()
	dyn := motion.NewDynamics(0)

	entities, err := arche.GenerateAndReturnEntity(1,
		spatial.NewPosition(x, y),
		spatial.NewRegularPolygon(20, 8),
		dyn,
		OBS_BUNDLE_MAIN,
		dir,
		components.Obstacle{
			Damage:        10,
			TravelTime:    travelTime,
			PauseDuration: (pauseDuration),
			MinX:          minX,
			MaxX:          maxX,
			MinY:          minY,
			MaxY:          maxY,
			EasingType:    easingType,
		},
	)
	if err != nil {
		return nil, err
	}
	return entities[0], nil
}
