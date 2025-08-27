package animations

import (
	"log"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
)

var batHandAnims = func() *client.AnimationCollection {
	anims, err := client.LoadAnimationsFromJSON("bat_animations.json", "../shared/animations/", AnimFS)
	if err != nil {
		log.Fatal(err)
	}
	return anims
}()

var _ = func() error {
	Registry.Characters[characterkeys.Bat] = map[AnimEnum]client.AnimationData{}
	Registry.Characters[characterkeys.Bat][Idle] = batHandAnims.Animations[0]
	Registry.Characters[characterkeys.Bat][Hurt] = batHandAnims.Animations[1]
	Registry.Characters[characterkeys.Bat][Defeat] = batHandAnims.Animations[2]
	return nil
}()
