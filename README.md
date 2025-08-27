# Concrete Echos Prototype

![ceprotogif](https://github.com/user-attachments/assets/1016de2e-9f92-4f69-b5e4-d996a419db94)
<img width="" height="225" alt="Screenshot 2025-08-27 at 11 41 23 AM" src="https://github.com/user-attachments/assets/ed46659e-5260-45c8-851a-819e6c729e5a" />

## About The Game

Concrete Echos is a free prototype for a 2D action-platformer concept I've been developing. It's an early, personal hobby kind of project.

You play as "Box Head," a character navigating a strange somewhat modern world. The game features platforming, combat, and a unique "swap" mechanic instead of a traditional parry.

As a young project, it's still rough around the edges, but my main goal was to create a specific mood and to test if the core ideas

## About The Project

This prototype is written in Go and stands on the shoulders of the following open-source projects:

- [Ebitengine:](https://ebitengine.org/) A 2D game library that provides the foundational tools for rendering, input, and audio.

- [LDtk:](https://ldtk.io/) A modern and flexible level editor, with custom JSON export.

The game's architecture, including the Entity Component System (ECS), physics, and the code that integrates with these tools, is handled by [Bappa](https://bappa.net/), a custom framework I wrote specifically for this and future projects (maybe lol).

Bringing all these elements together has been a huge learning process. I expected to be humbled in terms of programming, but it turns out game development comes with a plethora of other skills that I was lacking. Without even considering the code (which is a bit al dente), everything from the actual game design—like literally placing stuff in the level—to creating/procuring art and sounds has been an exhausting challenge. Truly an intense hobby.

## Key Features

- Swap Parry: A unique defensive mechanic where you teleport and swap places with an attacker, turning enemies against each other.

- Choice & Consequence: Your actions have a direct impact on the game world. Absorbing souls may seem harmless, but it will attract the attention of dangerous foe(s). Currently only one.

- Action Platforming: Explore a mysterious world with tight controls, challenging combat encounters, and difficult platforming sections.

## How to Play

You can play the prototype directly in your browser or download a native build for your operating system [here](https://thebitdrifter.itch.io/concrete-echos).

## Feedback Wanted

This is a prototype, and your feedback is incredibly valuable! After playing, please share your thoughts. I'm especially interested in what you think of the swap mechanic and the overall game feel.

Feedback: [Here](https://thebitdrifter.itch.io/concrete-echos/community)


## For Developers (A Note on Assets)

Please note that this repository contains the complete source code for the game logic, but it does not include the game's assets (art, sound, music).

This is due to licensing restrictions on some third-party assets. While all the assets used are free, some licenses request that they not be redistributed. As a result, the project cannot be built directly from this repository at this time.

My long-term goal is to create a fully original set of assets, at which point I will update the repository to be completely self-contained and compilable. In the meantime, you are welcome to explore the code.

## Some Commands

```bash
  # From ./standalone

  # Run Local
  go run -tags="m256 unsafe" .
  go run github.com/hajimehoshi/wasmserve@latest -tags="m256 unsafe" .

  # For Unix/Linux/macOS shells (cross-compiling for Windows)
  GOOS=windows GOARCH=amd64 go build -o build/concrete_echos_windows.exe -tags="m256 unsafe" -ldflags="-X 'github.com/TheBitDrifter/bappa/environment.Environment=production'" .

  # WASM
  GOOS=js GOARCH=wasm go build -o build/concrete_echos.wasm -tags="m256 unsafe" -ldflags="-X 'github.com/TheBitDrifter/bappa/environment.Environment=production'" .

  # Mac

  # ARM
  GOOS=darwin GOARCH=arm64 go build -o build/concrete_echos_mac_arm64 -tags="m256 unsafe" -ldflags="-X 'github.com/TheBitDrifter/bappa/environment.Environment=production'" .

  # INTEL
  CGO_ENABLED=1 CC=clang CGO_CFLAGS="-arch x86_64" CGO_LDFLAGS="-arch x86_64" \
  GOOS=darwin GOARCH=amd64 go build -o build/concrete_echos_mac_amd64 \
  -tags="m256 unsafe" -ldflags="-X 'github.com/TheBitDrifter/bappa/environment.Environment=production'" .

  ## Univeral MAC (using both binaries above)
  lipo -create -output build/concrete_echos_mac build/concrete_echos_mac_arm64 build/concrete_echos_mac_amd64
```
