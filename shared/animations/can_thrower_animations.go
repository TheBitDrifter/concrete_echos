package animations

import (
	"log"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
)

var canThrowerAnims = func() *client.AnimationCollection {
	anims, err := client.LoadAnimationsFromJSON("can_thrower_animations.json", "../shared/animations/", AnimFS)
	if err != nil {
		log.Fatal(err)
	}
	return anims
}()

var _ = func() error {
	Registry.Characters[characterkeys.CanThrower] = map[AnimEnum]client.AnimationData{}
	Registry.Characters[characterkeys.CanThrower][Idle] = canThrowerAnims.Animations[0]
	Registry.Characters[characterkeys.CanThrower][PrimaryRanged] = canThrowerAnims.Animations[1]
	Registry.Characters[characterkeys.CanThrower][Defeat] = canThrowerAnims.Animations[2]
	Registry.Characters[characterkeys.CanThrower][Hurt] = canThrowerAnims.Animations[3]
	return nil
}()

var canStraightThrowerAnims = func() *client.AnimationCollection {
	anims, err := client.LoadAnimationsFromJSON("can_thrower_animations.json", "../shared/animations/", AnimFS)
	if err != nil {
		log.Fatal(err)
	}
	return anims
}()

var _ = func() error {
	Registry.Characters[characterkeys.CanStraightThrower] = map[AnimEnum]client.AnimationData{}
	Registry.Characters[characterkeys.CanStraightThrower][Idle] = canStraightThrowerAnims.Animations[0]
	Registry.Characters[characterkeys.CanStraightThrower][PrimaryRanged] = canStraightThrowerAnims.Animations[1]
	Registry.Characters[characterkeys.CanStraightThrower][Defeat] = canStraightThrowerAnims.Animations[2]
	Registry.Characters[characterkeys.CanStraightThrower][Hurt] = canStraightThrowerAnims.Animations[3]
	return nil
}()
