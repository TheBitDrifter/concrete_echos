package animations

import (
	"log"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
)

var mudHandAnims = func() *client.AnimationCollection {
	anims, err := client.LoadAnimationsFromJSON("mud_hand_animations.json", "../shared/animations/", AnimFS)
	if err != nil {
		log.Fatal(err)
	}
	return anims
}()

var _ = func() error {
	Registry.Characters[characterkeys.MudHand] = map[AnimEnum]client.AnimationData{}
	Registry.Characters[characterkeys.MudHand][Idle] = mudHandAnims.Animations[0]
	Registry.Characters[characterkeys.MudHand][Hurt] = mudHandAnims.Animations[1]
	Registry.Characters[characterkeys.MudHand][Defeat] = mudHandAnims.Animations[2]
	return nil
}()
