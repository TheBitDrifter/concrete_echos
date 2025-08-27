package clientsystems

import (
	"math"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
)

const SAVE_DELAY_TICKS = 30

type PlayerSoundSystem struct {
	Volume float64
}

func (sys PlayerSoundSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	if scene.TicksSinceSelected() < 30 {
		return nil
	}

	sys.handleHurt(cli, scene, sys)
	sys.handleDodge(scene, sys)

	query := warehouse.Factory.NewQuery().And(
		client.Components.SoundBundle,
		components.PlayerTag,
		motion.Components.Dynamics,
	)

	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		soundBundle := client.Components.SoundBundle.GetFromCursor(cursor)
		dyn := motion.Components.Dynamics.GetFromCursor(cursor)
		onGround, onGroundExists := components.OnGroundComponent.GetFromCursorSafe(cursor)
		jumpState := components.JumpStateComponent.GetFromCursor(cursor)
		currentTick := scene.CurrentTick()
		ckey := *components.CharacterKeyComponent.GetFromCursor(cursor)

		sr, ok := components.SoftResetComponent.GetFromCursorSafe(cursor)
		defeated, okDef := combat.Components.Defeat.GetFromCursorSafe(cursor)
		if ok || okDef {
			if ok && sr.StartedTick != currentTick {
				continue
			}
			if okDef && defeated.StartTick != currentTick {
				continue
			}

			srSound, err := coldbrew.MaterializeSound(soundBundle, sounds.Registry.Characters[ckey][sounds.ObstacleHit])
			if err != nil {
				return err
			}
			player := srSound.GetAny()

			if !player.IsPlaying() {
				player.SetVolume(sys.Volume / 3)
				player.Rewind()
				player.Play()
			}

		}

		// attack sound
		if attack, ok := combat.Components.Attack.GetFromCursorSafe(cursor); ok {
			if attack.StartTick == scene.CurrentTick() {
				attackSound, err := coldbrew.MaterializeSound(soundBundle, sounds.Registry.Characters[ckey][sounds.PrimaryAttack])
				if err != nil {
					return err
				}
				player := attackSound.GetAny()

				if !player.IsPlaying() {
					player.SetVolume(sys.Volume)
					player.Rewind()
					player.Play()
					continue
				}

			}
		}

		// swap sound
		if swap, isSwapping := components.TeleportSwapComponent.GetFromCursorSafe(cursor); isSwapping {
			if swap.StartTick == scene.CurrentTick()-1 {
				swapSound, err := coldbrew.MaterializeSound(soundBundle, sounds.SwapSound)
				if err != nil {
					return err
				}
				player := swapSound.GetAny()

				if !player.IsPlaying() {
					player.SetVolume(sys.Volume)
					player.Rewind()
					player.Play()
					continue
				}

			}
		}

		// exec sound
		if exec, isExec := components.PlayerIsExecutingComponent.GetFromCursorSafe(cursor); isExec {
			if exec.StartTick == scene.CurrentTick()-1 {
				execSound, err := coldbrew.MaterializeSound(soundBundle, sounds.ExecuteSound)
				if err != nil {
					return err
				}
				player := execSound.GetAny()

				if !player.IsPlaying() {
					player.SetVolume(sys.Volume)
					player.Rewind()
					player.Play()
					continue
				}

			}
		}

		// save sound
		if saving, isSaving := components.IsSavingComponent.GetFromCursorSafe(cursor); isSaving {
			if saving.StartedTick+SAVE_DELAY_TICKS == scene.CurrentTick() {
				saveSound, err := coldbrew.MaterializeSound(soundBundle, sounds.SaveSound)
				if err != nil {
					return err
				}
				player := saveSound.GetAny()

				if !player.IsPlaying() {
					player.SetVolume(sys.Volume)
					player.Rewind()
					player.Play()
					continue
				}

			}
		}

		// Landed sound
		if onGroundExists && onGround.Landed == currentTick {
			landingSound, _ := coldbrew.MaterializeSound(soundBundle, sounds.Registry.Characters[ckey][sounds.Land])
			player := landingSound.GetAny()

			// A hack to prevent landing sound artifacts between scenes
			// In a more robust setup, we might track if a player has recently changed scenes via a component
			// Such a component would be helpful here
			sceneRecentlySelected := scene.CurrentTick()-scene.LastSelectedTick() < 30

			if !player.IsPlaying() && !sceneRecentlySelected {
				player.SetVolume(sys.Volume)
				player.Rewind()
				player.Play()
			}
		}

		// Jump sound
		//
		//
		if jumpState.LastJump == currentTick {

			jumpSound, _ := coldbrew.MaterializeSound(soundBundle, sounds.Registry.Characters[ckey][sounds.Jump])
			player := jumpSound.GetAny()

			if !player.IsPlaying() {
				player.SetVolume(sys.Volume)
				player.Rewind()
				player.Play()
			}
		}

		// Cam Border Thud Sound
		//
		if pcm, ok := components.PlayerCamStateComponent.GetFromCursorSafe(cursor); ok {
			// DISABLED
			if pcm.LastBorderHitTick == currentTick-1 && false {
				pcmSound, _ := coldbrew.MaterializeSound(soundBundle, sounds.Registry.Characters[ckey][sounds.Hurt])
				player := pcmSound.GetAny()

				if !player.IsPlaying() {
					player.SetVolume(sys.Volume)
					player.Rewind()
					player.Play()
				}
			}
		}

		// Run Sound
		const minMovementSpeed = 20.0
		if math.Abs(dyn.Vel.X) <= minMovementSpeed {
			continue
		}

		// Ensure onGround is not just available for coyote timer
		touchedGroundThisTick := onGroundExists && onGround.LastTouch == currentTick
		if !touchedGroundThisTick {
			continue
		}

		runSound, err := coldbrew.MaterializeSound(soundBundle, sounds.Registry.Characters[ckey][sounds.Run])
		if err != nil {
			return err
		}
		player := runSound.GetAny()

		if !player.IsPlaying() {
			player.SetVolume(sys.Volume)
			player.Rewind()
			player.Play()
		}

	}

	return nil
}

func (s PlayerSoundSystem) handleHurt(cli coldbrew.LocalClient, scene coldbrew.Scene, sys PlayerSoundSystem) error {
	query := warehouse.Factory.NewQuery().And(
		components.PlayerTag,
		combat.Components.Hurt,
	)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		soundBundle := client.Components.SoundBundle.GetFromCursor(cursor)
		hurt := combat.Components.Hurt.GetFromCursor(cursor)
		ckey := characterkeys.Default
		if components.CharacterKeyComponent.CheckCursor(cursor) {
			ckey = *components.CharacterKeyComponent.GetFromCursor(cursor)
		}

		if hurt.StartTick == scene.CurrentTick() {

			hurtSound, err := coldbrew.MaterializeSound(soundBundle, sounds.Registry.Characters[ckey][sounds.Hurt])
			if err != nil {
				return err
			}
			player := hurtSound.GetAny()

			if !player.IsPlaying() {
				player.SetVolume(sys.Volume)
				player.Rewind()
				player.Play()
			}

		}
	}
	return nil
}

func (s PlayerSoundSystem) handleDodge(scene coldbrew.Scene, sys PlayerSoundSystem) error {
	query := warehouse.Factory.NewQuery().And(
		components.PlayerTag,
		components.DodgeComponent,
	)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		soundBundle := client.Components.SoundBundle.GetFromCursor(cursor)
		dodge := components.DodgeComponent.GetFromCursor(cursor)
		if dodge.StartTick == scene.CurrentTick() {

			dodgeSound, err := coldbrew.MaterializeSound(soundBundle, sounds.Registry.Characters[characterkeys.Default][sounds.Dodge])
			if err != nil {
				return err
			}
			player := dodgeSound.GetAny()

			if !player.IsPlaying() {
				player.SetVolume(sys.Volume / 3)
				player.Rewind()
				player.Play()
			}

		}
	}
	return nil
}
