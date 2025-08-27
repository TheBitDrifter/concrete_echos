package scenes

import (
	"github.com/TheBitDrifter/bappa/blueprint"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/ldtk"
)

var entityRegistry = ldtk.NewLDtkEntityRegistry()

// Local scene object makes it easier to organize scene plans
type Scene struct {
	Name          string
	Plan          blueprint.Plan
	Width, Height int
	// Optional: manual assets to preload (usually for entities that get added 'dynamically', while the scene is running)
	Preload client.PreLoadAssetBundle
}
