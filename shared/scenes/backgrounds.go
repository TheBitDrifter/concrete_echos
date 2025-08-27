package scenes

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/warehouse"
)

func NewEasternDistrictBG(sto warehouse.Storage, x, y float64) error {
	return blueprint.NewParallaxBackgroundBuilder(sto).
		AddLayer("images/backgrounds/intro/far.png", 0.025, 0.05).
		AddLayer("images/backgrounds/intro/mid.png", 0.1, 0.1).
		AddLayer("images/backgrounds/intro/near.png", 0.2, 0.2).
		WithOffset(vector.Two{X: x, Y: y}).
		Build()
}

func NewEasternDistrictBGAlt(sto warehouse.Storage, x, y float64) error {
	return blueprint.NewParallaxBackgroundBuilder(sto).
		AddLayer("images/backgrounds/intro/far.png", 0.025, 0.05).
		AddLayer("images/backgrounds/intro/mid.png", 0.1, 0.1).
		WithOffset(vector.Two{X: x, Y: y}).
		Build()
}

func NewEasternDistrictGreenSkyBackground(sto warehouse.Storage, x, y float64) error {
	return blueprint.NewParallaxBackgroundBuilder(sto).
		AddLayer("images/backgrounds/eastern_district/sky.png", 0.05, 0.05).
		WithOffset(vector.Two{X: x, Y: y}).
		Build()
}

func NewEasternDistrictBLDBackground(sto warehouse.Storage, x, y float64) error {
	return blueprint.NewParallaxBackgroundBuilder(sto).
		AddLayer("images/backgrounds/eastern_district/back.png", 0.1, 0.1).
		AddLayer("images/backgrounds/eastern_district/front.png", 0.2, 0.2).
		WithOffset(vector.Two{X: x, Y: y}).
		Build()
}

func NewEasternOutskirtsBG(sto warehouse.Storage, x, y float64) error {
	return blueprint.NewParallaxBackgroundBuilder(sto).
		AddLayer("images/backgrounds/eastern_outskirts/3.png", 0.025, 0.025).
		AddLayer("images/backgrounds/eastern_outskirts/2.png", 0.050, 0.05).
		AddLayer("images/backgrounds/eastern_outskirts/1.png", 0.1, 0.1).
		AddLayer("images/backgrounds/eastern_outskirts/0.png", 0.13, 0.13).
		WithOffset(vector.Two{X: x, Y: y}).
		Build()
}

func NewEasternOutskirtsJunctionBG(sto warehouse.Storage, x, y float64) error {
	err := blueprint.NewParallaxBackgroundBuilder(sto).
		AddLayer("images/backgrounds/eastern_outskirts/3.png", 0.1, 0.1).
		WithOffset(vector.Two{X: x - 600, Y: y}).
		Build()
	if err != nil {
		return err
	}

	return blueprint.NewParallaxBackgroundBuilder(sto).
		AddLayer("images/backgrounds/eastern_outskirts/2.png", 0.2, 0.2).
		WithOffset(vector.Two{X: x, Y: y + 100}).
		Build()
}
