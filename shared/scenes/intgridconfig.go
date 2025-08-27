package scenes

import (
	"github.com/TheBitDrifter/bappa/blueprint/ldtk"
)

const (
	blockTerrainIntGrid = 1
	platformIntGrid     = 2
	obstacleIntGrid     = 3
)

var intGridConfigs = map[int]ldtk.IntGridLayerConfig{}

func init() {
	intGridConfigs = map[int]ldtk.IntGridLayerConfig{
		blockTerrainIntGrid: {
			Composition: BlockTerrainComposition,
			// Add any default component values here if needed
			DefaultValues: []any{},
		},
		platformIntGrid: {
			Composition:   PlatformComposition,
			DefaultValues: []any{},
		},
		obstacleIntGrid: {
			Composition:   DefaultObstacleComposition,
			DefaultValues: []any{},
			Padding:       5,
		},
	}
}
