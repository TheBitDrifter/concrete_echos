package scenes

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/dialogue"
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
)

// These are slices of common component compositions for various archetypes.
// They only include/represent the initial and static components of archetype
// Components can still be added or removed dynamically at runtime
//
// These slices are especially useful for creating starting entities, via archetypes, inside plan functions

var PlayerComposition = []warehouse.Component{
	// Hacky unlocky
	//	components.WallJumpUnlockedTag,
	//	components.WarpSwapUnlockedTag,

	spatial.Components.Position,
	combat.Components.HurtBoxes,
	combat.Components.Health,
	client.Components.SpriteBundle,
	spatial.Components.Direction,
	input.Components.ActionBuffer,
	client.Components.CameraIndex,
	spatial.Components.Shape,
	motion.Components.Dynamics,
	client.Components.SoundBundle,
	components.JumpStateComponent,
	components.CharacterKeyComponent,
	components.MovementCooldownComponent,
	components.AttackCooldownComponent,
	components.PlayerTag,
	components.TeleportSwapComponent,
	components.WalletComponent,
	components.PlayerExecutionCountComponent,
	components.LastCombatComponent,
	components.PlayerCamStateComponent,
}

var BlockTerrainComposition = []warehouse.Component{
	components.BlockTerrainTag,
	spatial.Components.Shape,
	spatial.Components.Position,
	motion.Components.Dynamics,
}

var PlatformComposition = []warehouse.Component{
	components.PlatformTag,
	spatial.Components.Rotation,
	client.Components.SpriteBundle,
	spatial.Components.Shape,
	spatial.Components.Position,
	motion.Components.Dynamics,
}

var MusicComposition = []warehouse.Component{
	client.Components.SoundBundle,
	components.MusicTag,
}

var AmbientNoiseComposition = []warehouse.Component{
	client.Components.SoundBundle,
	components.AmbientNoiseTag,
}

var CollisionPlayerTransferComposition = []warehouse.Component{
	spatial.Components.Position,
	spatial.Components.Shape,
	components.PlayerSceneTransferComponent,
}

var DefaultMobComposition = []warehouse.Component{
	spatial.Components.Position,
	client.Components.SpriteBundle,
	client.Components.SoundBundle,
	spatial.Components.Direction,
	spatial.Components.Shape,
	motion.Components.Dynamics,
	combat.Components.HurtBoxes,
	combat.Components.Health,
	components.CharacterKeyComponent,
	components.ContactComponent,
	components.MobTag,
	components.DropComponent,
}

var DefaultObstacleComposition = []warehouse.Component{
	spatial.Components.Position,
	client.Components.SpriteBundle,
	client.Components.SoundBundle,
	spatial.Components.Direction,
	spatial.Components.Shape,
	motion.Components.Dynamics,
	components.ObstacleComponent,
}

var DialogueComposition = []warehouse.Component{
	spatial.Components.Position,
	client.Components.SpriteBundle,
	client.Components.SoundBundle,
	spatial.Components.Direction,
	dialogue.Components.Conversation,
	client.Components.CameraIndex,
}

var ProjectileComposition = []warehouse.Component{
	spatial.Components.Position,
	spatial.Components.Direction,
	spatial.Components.Shape,
	motion.Components.Dynamics,
	client.Components.SoundBundle,
	client.Components.SpriteBundle,
	components.ProjectileTag,
	components.EntityRefComponent,
}

var TrapDoorComposition = []warehouse.Component{
	components.TrapDoorComponent,
	spatial.Components.Position,
	spatial.Components.Direction,
	spatial.Components.Shape,
	motion.Components.Dynamics,
	client.Components.SoundBundle,
	client.Components.SpriteBundle,
	components.PersistenceComponent,
}
