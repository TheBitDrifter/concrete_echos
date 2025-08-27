package animations

import (
	"log"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
)

var chestAnims = func() *client.AnimationCollection {
	anims, err := client.LoadAnimationsFromJSON("chest_animations.json", "../shared/animations/", AnimFS)
	if err != nil {
		log.Fatal(err)
	}
	return anims
}()

var _ = func() error {
	Registry.Characters[characterkeys.Chest] = map[AnimEnum]client.AnimationData{}
	Registry.Characters[characterkeys.Chest][Idle] = chestAnims.Animations[0]
	Registry.Characters[characterkeys.Chest][Defeat] = chestAnims.Animations[1]
	Registry.Characters[characterkeys.Chest][Hurt] = chestAnims.Animations[2]
	return nil
}()
