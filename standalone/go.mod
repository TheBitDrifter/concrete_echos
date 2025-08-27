module github.com/TheBitDrifter/concrete_echos/standalone

go 1.24.1

replace github.com/TheBitDrifter/concrete_echos/shared => ../shared/

replace github.com/TheBitDrifter/concrete_echos/sharedclient => ../sharedclient/

require (
	github.com/TheBitDrifter/bappa/blueprint v0.0.0-20250827172029-d96a17f7aaf4
	github.com/TheBitDrifter/bappa/coldbrew v0.0.0-20250827172029-d96a17f7aaf4
	github.com/TheBitDrifter/bappa/environment v0.0.0-20250805064827-5d0801e0a7d5
	github.com/TheBitDrifter/bappa/warehouse v0.0.0-20250827171242-3f6179875b16
	github.com/TheBitDrifter/concrete_echos/shared v0.0.0-00010101000000-000000000000
	github.com/TheBitDrifter/concrete_echos/sharedclient v0.0.0-00010101000000-000000000000
	github.com/hajimehoshi/ebiten/v2 v2.8.8
)

require (
	github.com/TheBitDrifter/bappa/combat v0.0.0-20250827172029-d96a17f7aaf4 // indirect
	github.com/TheBitDrifter/bappa/drip v0.0.0-20250805064827-5d0801e0a7d5 // indirect
	github.com/TheBitDrifter/bappa/table v0.0.0-20250827171242-3f6179875b16 // indirect
	github.com/TheBitDrifter/bappa/tteokbokki v0.0.0-20250805064827-5d0801e0a7d5 // indirect
	github.com/TheBitDrifter/bark v0.0.0-20250302175939-26104a815ed9 // indirect
	github.com/TheBitDrifter/mask v0.0.1-early-alpha.1 // indirect
	github.com/TheBitDrifter/util v0.0.0-20241102212109-342f4c0a810e // indirect
	github.com/ebitengine/gomobile v0.0.0-20240911145611-4856209ac325 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/oto/v3 v3.3.3 // indirect
	github.com/ebitengine/purego v0.8.0 // indirect
	github.com/go-text/typesetting v0.2.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	golang.org/x/image v0.30.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.28.0 // indirect
)
