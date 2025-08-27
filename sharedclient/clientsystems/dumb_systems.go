package clientsystems

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/dialogue"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/coldbrew/coldbrew_clientsystems"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/animations"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/persistence"
	"github.com/TheBitDrifter/concrete_echos/shared/scenes"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
	"github.com/TheBitDrifter/concrete_echos/sharedclient/fontdata"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type DumbHomeScreenSystem struct{}

// Home Screen -> base
func (DumbHomeScreenSystem) Run(lc coldbrew.LocalClient, scene coldbrew.Scene) error {
	cli := lc.(coldbrew.Client)
	if inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
		old := scene

		allScenes := cli.Cache().All()
		for _, s := range allScenes {
			s.Reset()
		}
		cli.ActivateSceneByName(scenes.INTRO_CUTSCENE.Name)
		cli.DeactivateScene(old)

	}

	return nil
}

type DumbDefeatTransferSystem struct{}

// Regular -> Defeat Screen
func (DumbDefeatTransferSystem) Run(lc coldbrew.LocalClient, scene coldbrew.Scene) error {
	cli := lc.(coldbrew.Client)

	query := warehouse.Factory.NewQuery().And(combat.Components.Defeat, components.PlayerTag)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		defeatedData := combat.Components.Defeat.GetFromCursor(cursor)
		if scene.CurrentTick()-defeatedData.StartTick < 180 {
			continue
		}

		err := clearNonDefaultPL(scene)
		if err != nil {
			return err
		}

		playerEn, err := cursor.CurrentEntity()
		if err != nil {
			return err
		}
		err = scene.Storage().EnqueueDestroyEntities(playerEn)
		if err != nil {
			return err
		}

		for _, cam := range cli.ActiveCamerasFor(scene) {
			_, localPos := cam.Positions()
			localPos.X = 0
			localPos.Y = 0
		}
		old := scene
		cli.ActivateSceneByName(scenes.DEFEAT_SCREEN_SCENE.Name)
		cli.DeactivateScene(old)

		break
	}

	return nil
}

type DumbDefeatScreenSystem struct{}

// Defeat Screen ->  Save or begining
func (DumbDefeatScreenSystem) Run(lc coldbrew.LocalClient, scene coldbrew.Scene) error {
	cli := lc.(coldbrew.Client)
	if inpututil.IsKeyJustReleased(ebiten.KeyEnter) {

		old := scene

		newSceneName := persistence.State.LastScene
		if newSceneName == "" {
			newSceneName = persistence.SceneName(scenes.OUTSIDE_HOME_SCENE.Name)
		}

		cli.ActivateSceneByName(string(newSceneName))
		cli.DeactivateScene(old)

		for _, s := range cli.Cache().All() {
			s.Reset()
		}

	}
	return nil
}

type DumbCutSceneSystem struct {
	step int
}

func (sys *DumbCutSceneSystem) Run(lc coldbrew.LocalClient, scene coldbrew.Scene) error {
	query := warehouse.Factory.NewQuery().And(
		components.CutsceneTag,
	)
	cursor := scene.NewCursor(query)

	limit := len(animations.CutSceneAnimations) - 1
	for range cursor.Next() {
		bundle := client.Components.SpriteBundle.GetFromCursor(cursor)

		bps := &bundle.Blueprints[0]

		activeAnimI := bps.Config.ActiveAnimIndex
		activeAnim := bps.Animations[activeAnimI]

		if sys.step < limit && activeAnim.IsFinished(scene.CurrentTick()) {
			sys.step++
			nextAnimation := animations.CutSceneAnimations[sys.step]
			bps.TryAnimation(nextAnimation)
		}
	}

	return nil
}

type DialogueAutoStepperSystem struct{}

func (sys DialogueAutoStepperSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	query := warehouse.Factory.NewQuery().And(
		dialogue.Components.Conversation,
		components.DialogueAutoStepperComponent,
	)
	cursor := scene.NewCursor(query)
	currentTick := scene.CurrentTick()

	for range cursor.Next() {
		convo := dialogue.Components.Conversation.GetFromCursor(cursor)

		if !convo.IsRevealing {
			coldbrew_clientsystems.InitConversation(scene, convo, fontdata.DEFAULT_FONT_FACE, float64(fontdata.MAX_LINE_WIDTH))
		}

		isRevealFinished := convo.AnimationState.FinalUpdateTick > 0
		if !isRevealFinished {
			continue
		}

		const WAIT_IN_TICKS = 30
		isWaitTimeOver := currentTick >= convo.AnimationState.FinalUpdateTick+WAIT_IN_TICKS
		if !isWaitTimeOver {
			continue
		}

		slides := dialogue.SlidesRegistry[dialogue.SlidesEnum(convo.SlidesID)]
		isLastSlide := convo.ActiveSlideIndex >= len(slides)-1

		if isLastSlide {
			cli.ActivateSceneByName(scenes.OUTSIDE_HOME_SCENE_NAME)
			cli.(coldbrew.Client).DeactivateScene(scene)
			return nil
		} else {
			convo.ActiveSlideIndex++
			convo.IsRevealing = false // Mark for re-initialization by the next loop.
			convo.AnimationState.FinalUpdateTick = 0
			convo.RevealStartTick = scene.CurrentTick()
		}
	}
	return nil
}

type DumbSingletonOptionalBossMusicChangeSystem struct {
	ChangeTick int
}

func (sys *DumbSingletonOptionalBossMusicChangeSystem) Run(lc coldbrew.LocalClient, scene coldbrew.Scene) error {
	query := warehouse.Factory.NewQuery().And(
		components.FriendlyAgroComponent,
		components.IsBossTag,
	)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		fa := components.FriendlyAgroComponent.GetFromCursor(cursor)
		if fa.StartTick != 0 && sys.ChangeTick != fa.StartTick {
			sys.ChangeTick = scene.CurrentTick()
		}

	}

	query = warehouse.Factory.NewQuery().And(
		components.MusicPlaylistComponent,
	)

	cursor = scene.NewCursor(query)

	createNew := false

	for range cursor.Next() {
		bundle := client.Components.SoundBundle.GetFromCursor(cursor)
		list := components.MusicPlaylistComponent.GetFromCursor(cursor)
		activeSongConfig := list.Collection.Sounds[list.CurrentSongIndex]
		activeSound, err := coldbrew.MaterializeSound(bundle, activeSongConfig)
		if err != nil {
			return err
		}
		activePlayer := activeSound.GetAny()

		if sys.ChangeTick == scene.CurrentTick() && sys.ChangeTick != 0 {
			createNew = true
			list.SongCollectionID = sounds.BossSoundCollection.ID
			list.Collection = sounds.BossSoundCollection
			activePlayer.Pause()

			en, err := cursor.CurrentEntity()
			if err != nil {
				return err
			}
			scene.Storage().EnqueueDestroyEntities(en)
		}

	}

	if createNew {
		scenes.AddPlaylist(scene.Storage(), 0, sounds.BossSoundCollection)
	}

	query = warehouse.Factory.NewQuery().And(combat.Components.Defeat, components.IsBossTag)
	cursor = scene.NewCursor(query)
	createNew = false

	for range cursor.Next() {
		defeatTick := combat.Components.Defeat.GetFromCursor(cursor).StartTick
		if scene.CurrentTick() == defeatTick+120 {
			err := clearBossPL(scene)
			if err != nil {
				return err
			}
			createNew = true

		}
	}

	if createNew {
		scenes.AddPlaylist(scene.Storage(), 0, sounds.PostBossSoundCollection)
	}
	return nil
}

func clearNonDefaultPL(scene coldbrew.Scene) error {
	queryTunes := warehouse.Factory.NewQuery().And(components.MusicPlaylistComponent)
	cursorTunes := scene.NewCursor(queryTunes)
	for range cursorTunes.Next() {
		bundle := client.Components.SoundBundle.GetFromCursor(cursorTunes)
		pl := components.MusicPlaylistComponent.GetFromCursor(cursorTunes)

		if pl.Collection.ID != sounds.DefaultSoundCollection.ID {

			activeSongConfig := pl.Collection.Sounds[pl.CurrentSongIndex]
			activeSound, err := coldbrew.MaterializeSound(bundle, activeSongConfig)
			if err != nil {
				return err
			}
			activePlayer := activeSound.GetAny()
			activePlayer.Pause()
			activePlayer.Rewind()

			en, err := cursorTunes.CurrentEntity()
			if err != nil {
				return err
			}
			scene.Storage().EnqueueDestroyEntities(en)
		}

	}
	return nil
}

func clearBossPL(scene coldbrew.Scene) error {
	queryTunes := warehouse.Factory.NewQuery().And(components.MusicPlaylistComponent)
	cursorTunes := scene.NewCursor(queryTunes)
	for range cursorTunes.Next() {
		bundle := client.Components.SoundBundle.GetFromCursor(cursorTunes)
		pl := components.MusicPlaylistComponent.GetFromCursor(cursorTunes)

		if pl.Collection.ID == sounds.BossSoundCollection.ID {

			activeSongConfig := pl.Collection.Sounds[pl.CurrentSongIndex]
			activeSound, err := coldbrew.MaterializeSound(bundle, activeSongConfig)
			if err != nil {
				return err
			}
			activePlayer := activeSound.GetAny()
			activePlayer.Pause()
			activePlayer.Rewind()

			en, err := cursorTunes.CurrentEntity()
			if err != nil {
				return err
			}
			scene.Storage().EnqueueDestroyEntities(en)
		}

	}
	return nil
}
