package clientsystems

import (
	"log"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/combatdata"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

type MobLazyRoninSoundSystem struct {
	Volume float64
}

func (sys MobLazyRoninSoundSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	query := warehouse.Factory.NewQuery().And(
		components.MobLazySkullyComponent,
	)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {

		soundBundle := client.Components.SoundBundle.GetFromCursor(cursor)
		currAttack, isAttacking := combat.Components.Attack.GetFromCursorSafe(cursor)
		currentTick := scene.CurrentTick()
		isDefeated := combat.Components.Defeat.CheckCursor(cursor)

		if isAttacking && !isDefeated {
			elapsedTicks := currentTick - currAttack.StartTick
			idx := elapsedTicks / currAttack.Speed

			if idx == 0 {
				continue
			}

			boxesCurr := currAttack.Boxes[idx]

			currActive := false
			for _, b := range boxesCurr {
				if b.LocalAAB.Height != 0 && b.LocalAAB.Width != 0 {
					currActive = true
					break
				}
			}

			boxesPrev := currAttack.Boxes[idx-1]
			prevActive := false
			for _, b := range boxesPrev {
				if b.LocalAAB.Height != 0 && b.LocalAAB.Width != 0 {
					prevActive = true
					break
				}
			}

			var soundConfig client.SoundConfig

			if currActive && !prevActive &&
				currAttack.ID != combatdata.LazySkullyAttacks[combatdata.LazySkullyAttackKeys.TertCombo].ID {
				soundConfig = sounds.LazyRoninSounds.SlashSound
			}

			const exploMainCuttof = 12

			if currAttack.ID == combatdata.LazySkullyAttacks[combatdata.LazySkullyAttackKeys.TertCombo].ID &&
				currAttack.FirstActiveBoxIndex() == idx {
				soundConfig = sounds.LazyRoninSounds.SlashSound
			} else if currAttack.ID == combatdata.LazySkullyAttacks[combatdata.LazySkullyAttackKeys.TertCombo].ID &&
				currActive && !prevActive {
				if idx < exploMainCuttof {
					soundConfig = sounds.LazyRoninSounds.ExplosionSound
				} else {
					soundConfig = sounds.LazyRoninSounds.SmallExplosionSound
				}
			}

			if soundConfig.AudioPlayerCount == 0 {
				continue
			}

			attackSound, err := coldbrew.MaterializeSound(soundBundle, soundConfig)
			if err != nil {
				log.Println("noAttackSound!")
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

	query = warehouse.Factory.NewQuery().And(
		components.MobLazySkullyComponent,
		components.FriendlyAgroComponent,
	)
	cursor = scene.NewCursor(query)
	for range cursor.Next() {

		soundBundle := client.Components.SoundBundle.GetFromCursor(cursor)
		fa := components.FriendlyAgroComponent.GetFromCursor(cursor)

		if fa.StartTick == scene.CurrentTick() {
			soundConfig := sounds.LazyRoninSounds.AgroSound

			snd, err := coldbrew.MaterializeSound(soundBundle, soundConfig)
			if err != nil {
				return nil
			}

			player := snd.GetAny()

			if !player.IsPlaying() {
				player.SetVolume(sys.Volume)
				player.Rewind()
				player.Play()
			}

		}
	}

	return nil
}
