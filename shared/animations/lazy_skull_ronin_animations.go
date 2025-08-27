package animations

import (
	"log"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
)

var LazySkullRoninAnims = func() *client.AnimationCollection {
	anims, err := client.LoadAnimationsFromJSON("lazy_skull_ronin_animations.json", "../shared/animations/", AnimFS)
	if err != nil {
		log.Fatal(err)
	}
	return anims
}()

var _ = func() error {
	Registry.Characters[characterkeys.LazySkullRonin] = map[AnimEnum]client.AnimationData{}
	Registry.Characters[characterkeys.LazySkullRonin][Idle] = LazySkullRoninAnims.Animations[0]
	Registry.Characters[characterkeys.LazySkullRonin][Hurt] = LazySkullRoninAnims.Animations[6]
	Registry.Characters[characterkeys.LazySkullRonin][Defeat] = LazySkullRoninAnims.Animations[7]
	Registry.Characters[characterkeys.LazySkullRonin][Run] = LazySkullRoninAnims.Animations[5]
	return nil
}()
