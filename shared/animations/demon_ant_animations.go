package animations

import (
	"log"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
)

var demonAntAnimations = func() *client.AnimationCollection {
	anims, err := client.LoadAnimationsFromJSON("demon_ant_animations.json", "../shared/animations/", AnimFS)
	if err != nil {
		log.Fatal(err)
	}
	return anims
}()

var _ = func() error {
	Registry.Characters[characterkeys.DemonAnt] = map[AnimEnum]client.AnimationData{}
	Registry.Characters[characterkeys.DemonAnt][Idle] = demonAntAnimations.Animations[0]
	Registry.Characters[characterkeys.DemonAnt][PrimaryAttack] = demonAntAnimations.Animations[1]
	Registry.Characters[characterkeys.DemonAnt][Hurt] = demonAntAnimations.Animations[2]
	Registry.Characters[characterkeys.DemonAnt][Dodge] = demonAntAnimations.Animations[3]
	Registry.Characters[characterkeys.DemonAnt][Defeat] = demonAntAnimations.Animations[4]
	Registry.Characters[characterkeys.DemonAnt][SecondaryAttack] = demonAntAnimations.Animations[5]
	Registry.Characters[characterkeys.DemonAnt][Run] = demonAntAnimations.Animations[6]

	return nil
}()
