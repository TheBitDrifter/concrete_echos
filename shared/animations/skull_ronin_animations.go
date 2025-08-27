package animations

import (
	"log"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
)

var skullRoninAnims = func() *client.AnimationCollection {
	anims, err := client.LoadAnimationsFromJSON("skull_ronin_animations.json", "../shared/animations/", AnimFS)
	if err != nil {
		log.Fatal(err)
	}
	return anims
}()

var _ = func() error {
	Registry.Characters[characterkeys.SkullRonin] = map[AnimEnum]client.AnimationData{}
	Registry.Characters[characterkeys.SkullRonin][Idle] = skullRoninAnims.Animations[0]
	Registry.Characters[characterkeys.SkullRonin][Hurt] = skullRoninAnims.Animations[1]
	Registry.Characters[characterkeys.SkullRonin][PrimaryAttack] = skullRoninAnims.Animations[2]
	Registry.Characters[characterkeys.SkullRonin][Defeat] = skullRoninAnims.Animations[3]
	Registry.Characters[characterkeys.SkullRonin][Run] = skullRoninAnims.Animations[4]
	return nil
}()
