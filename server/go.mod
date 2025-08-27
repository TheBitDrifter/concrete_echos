module github.com/TheBitDrifter/concrete_echos/server

go 1.24.1

replace github.com/TheBitDrifter/concrete_echos/shared => ../shared/

require (
	github.com/TheBitDrifter/bappa/blueprint v0.0.0-20250827172029-d96a17f7aaf4
	github.com/TheBitDrifter/bappa/drip v0.0.0-20250805033540-6c492b297684
	github.com/TheBitDrifter/bappa/warehouse v0.0.0-20250827171242-3f6179875b16
	github.com/TheBitDrifter/concrete_echos/shared v0.0.0-00010101000000-000000000000
)

require (
	github.com/TheBitDrifter/bappa/combat v0.0.0-20250805033540-6c492b297684 // indirect
	github.com/TheBitDrifter/bappa/environment v0.0.0-20250805064827-5d0801e0a7d5 // indirect
	github.com/TheBitDrifter/bappa/table v0.0.0-20250827171242-3f6179875b16 // indirect
	github.com/TheBitDrifter/bappa/tteokbokki v0.0.0-20250805064827-5d0801e0a7d5 // indirect
	github.com/TheBitDrifter/bark v0.0.0-20250302175939-26104a815ed9 // indirect
	github.com/TheBitDrifter/mask v0.0.1-early-alpha.1 // indirect
	github.com/TheBitDrifter/util v0.0.0-20241102212109-342f4c0a810e // indirect
)
