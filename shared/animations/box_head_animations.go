package animations

import (
	"log"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
)

var boxHeadAnimations = func() *client.AnimationCollection {
	anims, err := client.LoadAnimationsFromJSON("box_head_animations.json", "../shared/animations/", AnimFS)
	if err != nil {
		log.Fatal(err)
	}
	return anims
}()

var _ = func() error {
	Registry.Characters[characterkeys.BoxHead] = map[AnimEnum]client.AnimationData{}
	Registry.Characters[characterkeys.BoxHead][Idle] = boxHeadAnimations.Animations[0]
	Registry.Characters[characterkeys.BoxHead][Run] = boxHeadAnimations.Animations[1]
	Registry.Characters[characterkeys.BoxHead][Jump] = boxHeadAnimations.Animations[2]
	Registry.Characters[characterkeys.BoxHead][Fall] = boxHeadAnimations.Animations[3]
	Registry.Characters[characterkeys.BoxHead][PrimaryAttack] = boxHeadAnimations.Animations[4]
	Registry.Characters[characterkeys.BoxHead][Hurt] = boxHeadAnimations.Animations[6]
	Registry.Characters[characterkeys.BoxHead][Dodge] = boxHeadAnimations.Animations[7]
	Registry.Characters[characterkeys.BoxHead][Aerial] = boxHeadAnimations.Animations[8]
	Registry.Characters[characterkeys.BoxHead][AerialDownSmash] = boxHeadAnimations.Animations[9]
	Registry.Characters[characterkeys.BoxHead][AerialDownSmashLanding] = boxHeadAnimations.Animations[10]
	Registry.Characters[characterkeys.BoxHead][Teleport] = boxHeadAnimations.Animations[11]
	Registry.Characters[characterkeys.BoxHead][TeleportMarker] = boxHeadAnimations.Animations[12]
	Registry.Characters[characterkeys.BoxHead][TeleportEffect] = boxHeadAnimations.Animations[13]
	Registry.Characters[characterkeys.BoxHead][Defeat] = boxHeadAnimations.Animations[14]
	Registry.Characters[characterkeys.BoxHead][InConvo] = boxHeadAnimations.Animations[15]
	Registry.Characters[characterkeys.BoxHead][IsSaving] = boxHeadAnimations.Animations[16]
	Registry.Characters[characterkeys.BoxHead][UpAttack] = boxHeadAnimations.Animations[17]
	Registry.Characters[characterkeys.BoxHead][UpAerial] = boxHeadAnimations.Animations[18]
	Registry.Characters[characterkeys.BoxHead][Execute] = boxHeadAnimations.Animations[19]

	return nil
}()
