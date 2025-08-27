package main

import (
	"errors"
	"log"
	"math/rand/v2"
	"os"
	"runtime/pprof"

	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/coldbrew/coldbrew_clientsystems"
	"github.com/TheBitDrifter/bappa/coldbrew/coldbrew_rendersystems"
	"github.com/TheBitDrifter/bappa/environment"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/concrete_echos/shared/actions"
	"github.com/TheBitDrifter/concrete_echos/shared/coresystems"
	"github.com/TheBitDrifter/concrete_echos/shared/persistence"
	"github.com/TheBitDrifter/concrete_echos/shared/scenes"
	"github.com/TheBitDrifter/concrete_echos/shared/sounds"
	"github.com/TheBitDrifter/concrete_echos/sharedclient"
	"github.com/TheBitDrifter/concrete_echos/sharedclient/assets"
	"github.com/TheBitDrifter/concrete_echos/sharedclient/clientsystems"
	"github.com/TheBitDrifter/concrete_echos/sharedclient/rendersystems"

	"github.com/hajimehoshi/ebiten/v2"
)

var isFromSave bool

func main() {
	if false {
		f, err := os.Create("cpu.prof")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // Make sure to close the file when the function exits.

		// Start the CPU profiler.
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		// Stop the profiler when the function exits. This is essential.
		defer pprof.StopCPUProfile()
	}

	warehouse.MemConfig.Set(1000, 200)
	warehouse.MemConfig.Preallocate()

	client := coldbrew.NewClient(
		sharedclient.RESOLUTION_X,
		sharedclient.RESOLUTION_Y,
		sharedclient.MAX_SPRITES_CACHED,
		sharedclient.MAX_SOUNDS_CACHED,
		sharedclient.MAX_SCENES_CACHED,
		assets.FS,
	)

	client.SetLocalAssetPath("../sharedclient/assets/")

	client.SetTitle("Concrete Echos Prototype")
	client.SetResizable(true)
	client.SetMinimumLoadTime(0)
	client.SetWindowSize(1920, 1080)
	client.SetDebugMode(false)
	ebiten.SetWindowSize(1920, 1080)

	err := persistence.LoadState("ce_save.json")
	if err != nil {
		if errors.Is(err, persistence.ErrStateNotFound) {
			log.Println("No save file found, starting fresh.")
		} else {
			log.Fatalf("Failed to load state with unexpected error: %v", err)
		}
	} else {
		isFromSave = true
		log.Println("Successfully loaded state from save.")

		log.Println("Randomizing Playlist!")
		playlist := sounds.DefaultSoundCollection

		rand.Shuffle(len(playlist.Sounds), func(i, j int) {
			playlist.Sounds[i], playlist.Sounds[j] = playlist.Sounds[j], playlist.Sounds[i]
		})

		log.Println("Playlist after shuffle:", playlist.Sounds)
		sounds.DefaultSoundCollection = playlist

	}

	// log.Println(persistence.State.PlayerSingleton.Components)
	sceneOverride := ""

	if persistence.State.LastScene != "" {
		sceneOverride = string(persistence.State.LastScene)
	}

	homeScreen := true
	currentDev := false
	// TODO: Make this a build tag or command line tag thingy for ezpz dev life?
	// sceneOverride = scenes.OUTSIDE_HOME_SCENE.Name

	if homeScreen {
		log.Println("Registering HomeScreen ...")
		err := client.RegisterScene(
			scenes.HOME_SCREEN_SCENE.Name,
			scenes.HOME_SCREEN_SCENE.Width,
			scenes.HOME_SCREEN_SCENE.Height,
			scenes.HOME_SCREEN_SCENE.Plan,
			rendersystems.DefaultRenderSystems,
			[]coldbrew.ClientSystem{clientsystems.DumbHomeScreenSystem{}},
			coresystems.DefaultCoreSystems,
		)
		if err != nil {
			log.Fatalf("Failed to register %v — %s", err, scenes.HOME_SCREEN_SCENE)
		}

		log.Println("Registering IntroCutscene ...")
		err = client.RegisterScene(
			scenes.INTRO_CUTSCENE.Name,
			scenes.INTRO_CUTSCENE.Width,
			scenes.INTRO_CUTSCENE.Height,
			scenes.INTRO_CUTSCENE.Plan,
			rendersystems.DefaultRenderSystems,
			clientsystems.CutSceneClientSystems,
			[]blueprint.CoreSystem{},
		)
		if err != nil {
			log.Fatalf("Failed to register %v — %s", err, scenes.INTRO_CUTSCENE.Name)
		}
	}

	// LEVELS ----------------

	if currentDev {
		err := client.RegisterScene(
			scenes.CURRENT_DEV_SCENE.Name,
			scenes.CURRENT_DEV_SCENE.Width,
			scenes.CURRENT_DEV_SCENE.Height,
			scenes.CURRENT_DEV_SCENE.Plan,
			rendersystems.DefaultRenderSystems,
			clientsystems.DefaultClientSystems,
			coresystems.DefaultCoreSystems,
			*scenes.DEFAULT_PRELOAD...,
		)
		if err != nil {
			log.Fatalf("Failed to register %v — %s", err, scenes.CURRENT_DEV_SCENE.Name)
		}
	}

	err = client.RegisterScene(
		scenes.OUTSIDE_HOME_SCENE.Name,
		scenes.OUTSIDE_HOME_SCENE.Width,
		scenes.OUTSIDE_HOME_SCENE.Height,
		scenes.OUTSIDE_HOME_SCENE.Plan,
		rendersystems.DefaultRenderSystems,
		clientsystems.DefaultClientSystems,
		coresystems.DefaultCoreSystems,
		*scenes.DEFAULT_PRELOAD...,
	)
	if err != nil {
		log.Fatalf("Failed to register %v — %s", err, scenes.OUTSIDE_HOME_SCENE.Name)
	}

	err = client.RegisterScene(
		scenes.EASTERN_GARDENS_SCENE.Name,
		scenes.EASTERN_GARDENS_SCENE.Width,
		scenes.EASTERN_GARDENS_SCENE.Height,
		scenes.EASTERN_GARDENS_SCENE.Plan,
		rendersystems.DefaultRenderSystems,
		clientsystems.DefaultClientSystems,
		coresystems.DefaultCoreSystems,
		*scenes.DEFAULT_PRELOAD...,
	)
	if err != nil {
		log.Fatalf("Failed to register %v — %s", err, scenes.OUTSIDE_HOME_SCENE.Name)
	}

	err = client.RegisterScene(
		scenes.OUTER_EAST_DISTRICT_SCENE.Name,
		scenes.OUTER_EAST_DISTRICT_SCENE.Width,
		scenes.OUTER_EAST_DISTRICT_SCENE.Height,
		scenes.OUTER_EAST_DISTRICT_SCENE.Plan,
		rendersystems.DefaultRenderSystems,
		clientsystems.DefaultClientSystems,
		coresystems.DefaultCoreSystems,
		scenes.OUTSIDE_HOME_PASSAGE_SCENE.Preload...,
	)
	if err != nil {
		log.Fatalf("Failed to register %v — %s", err, scenes.OUTER_EAST_DISTRICT_SCENE.Name)
	}

	err = client.RegisterScene(
		scenes.UPPER_EAST_DISTRICT_SCENE.Name,
		scenes.UPPER_EAST_DISTRICT_SCENE.Width,
		scenes.UPPER_EAST_DISTRICT_SCENE.Height,
		scenes.UPPER_EAST_DISTRICT_SCENE.Plan,
		rendersystems.DefaultRenderSystems,
		clientsystems.DefaultClientSystems,
		coresystems.DefaultCoreSystems,
		*scenes.DEFAULT_PRELOAD...,
	)
	if err != nil {
		log.Fatalf("Failed to register %v — %s", err, scenes.UPPER_EAST_DISTRICT_SCENE.Name)
	}

	err = client.RegisterScene(
		scenes.JUNCTION_SCENE.Name,
		scenes.JUNCTION_SCENE.Width,
		scenes.JUNCTION_SCENE.Height,
		scenes.JUNCTION_SCENE.Plan,
		rendersystems.DefaultRenderSystems,
		clientsystems.DefaultClientSystems,
		coresystems.DefaultCoreSystems,
		*scenes.DEFAULT_PRELOAD...,
	)
	if err != nil {
		log.Fatalf("Failed to register %v — %s", err, scenes.JUNCTION_SCENE.Name)
	}

	err = client.RegisterScene(
		scenes.RISING_PASSAGE_SCENE.Name,
		scenes.RISING_PASSAGE_SCENE.Width,
		scenes.RISING_PASSAGE_SCENE.Height,
		scenes.RISING_PASSAGE_SCENE.Plan,
		rendersystems.DefaultRenderSystems,
		clientsystems.DefaultClientSystems,
		coresystems.DefaultCoreSystems,
		*scenes.DEFAULT_PRELOAD...,
	)
	if err != nil {
		log.Fatalf("Failed to register %v — %s", err, scenes.RISING_PASSAGE_SCENE.Name)
	}

	err = client.RegisterScene(
		scenes.OUTSIDE_HOME_PASSAGE_SCENE.Name,
		scenes.OUTSIDE_HOME_PASSAGE_SCENE.Width,
		scenes.OUTSIDE_HOME_PASSAGE_SCENE.Height,
		scenes.OUTSIDE_HOME_PASSAGE_SCENE.Plan,
		rendersystems.DefaultRenderSystems,
		clientsystems.DefaultClientSystems,
		coresystems.DefaultCoreSystems,
		*scenes.DEFAULT_PRELOAD...,
	)
	if err != nil {
		log.Fatalf("Failed to register %v — %s", err, scenes.OUTSIDE_HOME_PASSAGE_SCENE.Name)
	}

	err = client.RegisterScene(
		scenes.SKY_GROTTO_SCENE.Name,
		scenes.SKY_GROTTO_SCENE.Width,
		scenes.SKY_GROTTO_SCENE.Height,
		scenes.SKY_GROTTO_SCENE.Plan,
		rendersystems.DefaultRenderSystems,
		clientsystems.DefaultClientSystems,
		coresystems.DefaultCoreSystems,
		*scenes.DEFAULT_PRELOAD...,
	)
	if err != nil {
		log.Fatalf("Failed to register %v — %s", err, scenes.SKY_GROTTO_SCENE.Name)
	}
	err = client.RegisterScene(
		scenes.SKY_GROTTO_REWARD_CUT_SCENE.Name,
		scenes.SKY_GROTTO_REWARD_CUT_SCENE.Width,
		scenes.SKY_GROTTO_REWARD_CUT_SCENE.Height,
		scenes.SKY_GROTTO_REWARD_CUT_SCENE.Plan,
		rendersystems.DefaultRenderSystems,
		clientsystems.DefaultClientSystems,
		coresystems.DefaultCoreSystems,
		*scenes.DEFAULT_PRELOAD...,
	)
	if err != nil {
		log.Fatalf("Failed to register %v — %s", err, scenes.SKY_GROTTO_REWARD_CUT_SCENE.Name)
	}
	// -----------------------------
	err = client.RegisterScene(
		scenes.WARP_UNLOCK_SCENE.Name,
		scenes.WARP_UNLOCK_SCENE.Width,
		scenes.WARP_UNLOCK_SCENE.Height,
		scenes.WARP_UNLOCK_SCENE.Plan,
		rendersystems.DefaultRenderSystems,
		clientsystems.DefaultClientSystems,
		coresystems.DefaultCoreSystems,
		*scenes.DEFAULT_PRELOAD...,
	)
	if err != nil {
		log.Fatalf("Failed to register %v — %s", err, scenes.WARP_UNLOCK_SCENE.Name)
	}

	// -----------------------------
	err = client.RegisterScene(
		scenes.FAST_TRAVEL_SCENE.Name,
		scenes.FAST_TRAVEL_SCENE.Width,
		scenes.FAST_TRAVEL_SCENE.Height,
		scenes.FAST_TRAVEL_SCENE.Plan,
		rendersystems.DefaultFastTravelSystems,
		clientsystems.DefaultFastTravelSystems,
		coresystems.DefaultCoreSystems,
		*scenes.DEFAULT_PRELOAD...,
	)
	if err != nil {
		log.Fatalf("Failed to register %v — %s", err, scenes.FAST_TRAVEL_SCENE.Name)
	}

	// -----------------------------

	log.Println("Registering DefeatScreen ...")
	err = client.RegisterScene(
		scenes.DEFEAT_SCREEN_SCENE.Name,
		scenes.DEFEAT_SCREEN_SCENE.Width,
		scenes.DEFEAT_SCREEN_SCENE.Height,
		scenes.DEFEAT_SCREEN_SCENE.Plan,
		rendersystems.DefaultRenderSystems,
		[]coldbrew.ClientSystem{clientsystems.DumbDefeatScreenSystem{}},
		coresystems.DefaultCoreSystems,
	)
	if err != nil {
		log.Fatalf("Failed to register %v — %s", err, scenes.DEFEAT_SCREEN_NAME)
	}

	log.Println("Registering Global Systems...")
	client.RegisterGlobalRenderSystem(
		coldbrew_rendersystems.GlobalRenderer{},
		&coldbrew_rendersystems.DebugRenderer{},
	)
	client.RegisterGlobalClientSystem(
		coldbrew_clientsystems.InputBufferSystem{},
		&coldbrew_clientsystems.CameraSceneAssignerSystem{},
	)

	log.Println("Activating Camera...")
	_, err = client.ActivateCamera()
	if err != nil {
		log.Fatalf("Failed to activate camera: %v", err)
	}

	log.Println("Activating Input Receiver and Mapping Keys...")
	receiver1, err := client.ActivateReceiver()
	if err != nil {
		log.Fatalf("Failed to activate receiver: %v", err)
	}

	receiver1.RegisterPad(0)

	// Jump
	receiver1.RegisterKey(ebiten.KeySpace, actions.Jump)
	receiver1.RegisterReleasedKey(ebiten.KeySpace, actions.JumpReleased)

	if environment.IsWASM() {
		receiver1.RegisterGamepadButton(ebiten.GamepadButton0, actions.Jump)
		receiver1.RegisterGamepadReleasedButton(ebiten.GamepadButton0, actions.JumpReleased)
	} else {
		receiver1.RegisterGamepadButton(ebiten.GamepadButton1, actions.Jump)
		receiver1.RegisterGamepadReleasedButton(ebiten.GamepadButton1, actions.JumpReleased)
	}

	// Directional Movement
	receiver1.RegisterKey(ebiten.KeyA, actions.Left)
	receiver1.RegisterKey(ebiten.KeyD, actions.Right)
	receiver1.RegisterJustPressedKey(ebiten.KeyS, actions.Down)
	receiver1.RegisterKey(ebiten.KeyW, actions.Up)
	receiver1.RegisterKey(ebiten.KeyS, actions.AttackDown)

	receiver1.RegisterGamepadAxes(true, actions.VectorTwoMovement)

	// Attack
	receiver1.RegisterJustPressedKey(ebiten.KeyJ, actions.PrimaryAttack)
	if environment.IsWASM() {
		receiver1.RegisterGamepadJustPressedButton(ebiten.GamepadButton1, actions.PrimaryAttack)
	} else {
		receiver1.RegisterGamepadJustPressedButton(ebiten.GamepadButton2, actions.PrimaryAttack)
	}

	// Dodge
	receiver1.RegisterJustPressedKey(ebiten.KeyShift, actions.Dodge)
	if environment.IsWASM() {
		receiver1.RegisterGamepadJustPressedButton(ebiten.GamepadButton6, actions.Dodge)
	} else {
		receiver1.RegisterGamepadJustPressedButton(ebiten.GamepadButton6, actions.Dodge)
	}

	// Camera Movement
	receiver1.RegisterKey(ebiten.KeyRight, actions.CameraRight)
	receiver1.RegisterKey(ebiten.KeyLeft, actions.CameraLeft)
	receiver1.RegisterKey(ebiten.KeyUp, actions.CameraUp)
	receiver1.RegisterKey(ebiten.KeyDown, actions.CameraDown)
	receiver1.RegisterGamepadAxes(false, actions.VectorTwoCamMovement)

	// Interactions
	receiver1.RegisterJustPressedKey(ebiten.KeyEnter, actions.Interact)
	receiver1.RegisterGamepadJustPressedButton(ebiten.GamepadButton3, actions.Interact)

	receiver1.RegisterJustPressedKey(ebiten.KeyEscape, actions.Cancel)
	if environment.IsWASM() {
		receiver1.RegisterGamepadJustPressedButton(ebiten.GamepadButton2, actions.Cancel)
	} else {
		receiver1.RegisterGamepadJustPressedButton(ebiten.GamepadButton0, actions.Cancel)
	}

	// Teleport Swap
	receiver1.RegisterJustPressedKey(ebiten.KeyI, actions.ShiftTeleTargetRight)
	receiver1.RegisterJustPressedKey(ebiten.KeyU, actions.ShiftTeleTargetLeft)
	receiver1.RegisterJustPressedKey(ebiten.KeyO, actions.ShiftTeleTargetNear)

	receiver1.RegisterJustPressedKey(ebiten.KeyK, actions.TeleSwap)
	receiver1.RegisterGamepadJustPressedButton(ebiten.GamepadButton7, actions.TeleSwap)

	receiver1.RegisterGamepadJustPressedButton(ebiten.GamepadButton15, actions.ShiftTeleTargetRight)
	if environment.IsWASM() {
		receiver1.RegisterGamepadJustPressedButton(ebiten.GamepadButton14, actions.ShiftTeleTargetLeft)
	} else {
		receiver1.RegisterGamepadJustPressedButton(ebiten.GamepadButton17, actions.ShiftTeleTargetLeft)
	}

	if environment.IsWASM() {
		receiver1.RegisterGamepadJustPressedButton(ebiten.GamepadButton12, actions.ShiftTeleTargetNear)
	} else {
		receiver1.RegisterGamepadJustPressedButton(ebiten.GamepadButton14, actions.ShiftTeleTargetNear)
	}

	log.Println("Starting Ebiten game loop (blocking)...")

	if sceneOverride != "" {
		for s := range client.ActiveScenes() {
			client.DeactivateScene(s)
		}
		client.ActivateSceneByName(sceneOverride)
	}

	const prePlan = true

	if prePlan && !environment.IsWASM() {
		err = client.PreExecAllPlans()
		if err != nil {
			log.Fatalf("Pre exec err %v", err)
		}

	}

	if prePlan && environment.IsWASM() {
		client.PreExecSceneByName(scenes.EASTERN_GARDENS_SCENE.Name)
	}

	if err := client.Start(); err != nil {
		log.Fatalf("Client exited with error: %v", err)
	}
}
