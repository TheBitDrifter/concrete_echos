package animations

import (
	"log"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
)

var drifterAnims = func() *client.AnimationCollection {
	anims, err := client.LoadAnimationsFromJSON("drifter_animations.json", "../shared/animations/", AnimFS)
	if err != nil {
		log.Fatal(err)
	}
	return anims
}()

var _ = func() error {
	Registry.Characters[characterkeys.Drifter] = map[AnimEnum]client.AnimationData{}
	Registry.Characters[characterkeys.Drifter][Idle] = drifterAnims.Animations[0]
	Registry.Characters[characterkeys.Drifter][ConvoStart] = drifterAnims.Animations[1]
	Registry.Characters[characterkeys.Drifter][InConvo] = drifterAnims.Animations[2]
	return nil
}()
