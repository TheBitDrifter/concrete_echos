package clientsystems

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
	"github.com/TheBitDrifter/concrete_echos/shared/combatdata"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

type MobSoundSystem struct {
	Volume float64
}

func (sys MobSoundSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	if scene.TicksSinceSelected() < 30 {
		return nil
	}
	query := warehouse.Factory.NewQuery().And(
		components.MobTag,
	)
	cursor := scene.NewCursor(query)
	for range cursor.Next() {
		cKey := components.CharacterKeyComponent.GetFromCursor(cursor)

		soundBundle := client.Components.SoundBundle.GetFromCursor(cursor)
		hurt, isHurt := combat.Components.Hurt.GetFromCursorSafe(cursor)
		if isHurt && hurt.StartTick == scene.CurrentTick() {

			hurtSoundConfig, hasCustomSound := sounds.Registry.Characters[*cKey][sounds.Hurt]
			var hurtSound coldbrew.Sound

			if hasCustomSound {
				hs, err := coldbrew.MaterializeSound(soundBundle, hurtSoundConfig)
				if err != nil {
					return err
				}
				hurtSound = hs
			} else {

				defaultHurtConfig := sounds.Registry.Characters[characterkeys.Default][sounds.Hurt]
				hs, err := coldbrew.MaterializeSound(soundBundle, defaultHurtConfig)
				if err != nil {
					return err
				}
				hurtSound = hs
			}

			player := hurtSound.GetAny()

			if !player.IsPlaying() {
				player.SetVolume(sys.Volume)
				player.Rewind()
				player.Play()
			}

		}

		defeated, isDef := combat.Components.Defeat.GetFromCursorSafe(cursor)
		if isDef && defeated.StartTick == scene.CurrentTick() {
			customDefeatSoundConfig, hasCustomDefeatSound := sounds.Registry.Characters[*cKey][sounds.Defeat]
			var defeatSound coldbrew.Sound

			if hasCustomDefeatSound {
				ds, err := coldbrew.MaterializeSound(soundBundle, customDefeatSoundConfig)
				if err != nil {
					return err
				}
				defeatSound = ds
			} else {
				defaultDefeatSoundConfig := sounds.Registry.Characters[characterkeys.Default][sounds.Defeat]

				ds, err := coldbrew.MaterializeSound(soundBundle, defaultDefeatSoundConfig)
				if err != nil {
					return err
				}
				defeatSound = ds
			}

			player := defeatSound.GetAny()

			if !player.IsPlaying() {
				player.SetVolume(sys.Volume)
				player.Rewind()
				player.Play()
			}
		}

		attack, isAttacking := combat.Components.Attack.GetFromCursorSafe(cursor)

		if isAttacking {
			firstBoxIndex := attack.FirstActiveBoxIndex()
			soundTick := attack.Speed*firstBoxIndex + attack.StartTick
			if scene.CurrentTick() != soundTick {
				continue
			}

			primaryID := 0
			primarySeq, okPrim := combatdata.PrimarySeqs[*cKey]
			if okPrim {
				primaryID = primarySeq.First().ID
			}

			secondaryID := 0
			secSeq, okSec := combatdata.SecondarySeqs[*cKey]
			if okSec {
				secondaryID = secSeq.First().ID
			}

			isPrim := primaryID == attack.ID
			isSec := secondaryID == attack.ID

			var soundConfig client.SoundConfig

			regis := sounds.Registry.Characters[*cKey]

			if isPrim {
				soundConfig = sounds.Registry.Characters[characterkeys.Default][sounds.PrimaryAttack]
				customPrimAtkSoundConfig, hasCustomPrimAttack := regis[sounds.PrimaryAttack]

				if hasCustomPrimAttack {
					soundConfig = customPrimAtkSoundConfig
				}
			}

			if isSec {
				soundConfig = sounds.Registry.Characters[characterkeys.Default][sounds.SecondaryAttack]
				customSecAtkSoundConfig, hasCustomSecAtk := regis[sounds.SecondaryAttack]

				if hasCustomSecAtk {
					soundConfig = customSecAtkSoundConfig
				}
			}

			attackSound, err := coldbrew.MaterializeSound(soundBundle, soundConfig)
			if err != nil {
				// log.Println("noAttackSound!")
				return nil
			}

			player := attackSound.GetAny()

			if !player.IsPlaying() {
				player.SetVolume(sys.Volume)
				player.Rewind()
				player.Play()
			}
		}
	}
	return nil
}
