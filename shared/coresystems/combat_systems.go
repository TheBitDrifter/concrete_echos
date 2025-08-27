package coresystems

import (
	"log"
	"math"

	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/combat"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"

	"github.com/TheBitDrifter/concrete_echos/shared/actions"
	"github.com/TheBitDrifter/concrete_echos/shared/characterkeys"
	"github.com/TheBitDrifter/concrete_echos/shared/combatdata"
	"github.com/TheBitDrifter/concrete_echos/shared/components"
	"github.com/TheBitDrifter/concrete_echos/shared/scenes"
)

type PlayerAttackSystem struct{}

func (PlayerAttackSystem) Run(scene blueprint.Scene, dt float64) error {
	atkQuery := warehouse.Factory.NewQuery().And(
		combat.Components.Attack, spatial.Components.Direction,
	)
	attackCursor := scene.NewCursor(atkQuery)

	for range attackCursor.Next() {
		cKey := components.CharacterKeyComponent.GetFromCursor(attackCursor)
		hasAtk := combat.Components.Attack.CheckCursor(attackCursor)

		// Remove stale attacks
		if hasAtk {
			atk := combat.Components.Attack.GetFromCursor(attackCursor)
			start := atk.StartTick
			duration := atk.Speed * atk.Length
			if (scene.CurrentTick() - duration) > start {
				en, err := attackCursor.CurrentEntity()
				if err != nil {
					return err
				}
				en.EnqueueRemoveComponent(combat.Components.Attack)
			} else {
				// Handle DownSmash
				og, grounded := components.OnGroundComponent.GetFromCursorSafe(attackCursor)
				if grounded {
					grounded = scene.CurrentTick() == og.LastTouch
				}

				if grounded {
					matchedAtk, ok := combatdata.AerialDownSmashes[*cKey]
					if ok && matchedAtk.ID == atk.ID {
						smashAtkLanding := combatdata.AerialDownSmashLandings[*cKey]
						*atk = smashAtkLanding

						atk.StartTick = scene.CurrentTick()
					}

				}
			}

		}
	}

	query := warehouse.Factory.NewQuery().And(components.PlayerTag,
		warehouse.Factory.NewQuery().Not(components.SoftResetComponent, components.InConversationComponent, combat.Components.Defeat),
	)
	cursor := scene.NewCursor(query)
	for range cursor.Next() {

		cKey := components.CharacterKeyComponent.GetFromCursor(cursor)
		hasAtk := combat.Components.Attack.CheckCursor(cursor)
		actionBuffer := input.Components.ActionBuffer.GetFromCursor(cursor)

		_, attackingOK := actionBuffer.PeekLatest()

		if !attackingOK {
			continue
		}

		stamptedAtkAction, _ := actionBuffer.ConsumeAction(actions.PrimaryAttack)

		stampedUp, upOK := actionBuffer.ConsumeAction(actions.Up)

		og, grounded := components.OnGroundComponent.GetFromCursorSafe(cursor)
		if grounded {
			grounded = scene.CurrentTick() == og.LastTouch
		}

		if stamptedAtkAction.Tick == scene.CurrentTick() {
			if hasAtk {
				continue
			}
			en, err := cursor.CurrentEntity()
			if err != nil {
				return err
			}

			attackCooldown := components.AttackCooldownComponent.GetFromCursor(cursor)

			if grounded && upOK && stampedUp.Tick == scene.CurrentTick() {
				attack := combatdata.UpSeqs[*cKey].First()
				attack.StartTick = scene.CurrentTick()
				attack.LRDirection = *spatial.Components.Direction.GetFromCursor(cursor)
				en.EnqueueAddComponentWithValue(combat.Components.Attack, attack)
				continue
			}

			if grounded {
				attack := combatdata.PrimarySeqs[*cKey].First()
				attack.StartTick = scene.CurrentTick()
				attack.LRDirection = *spatial.Components.Direction.GetFromCursor(cursor)
				en.EnqueueAddComponentWithValue(combat.Components.Attack, attack)
			} else {

				stampedDown, downOK := actionBuffer.PeekLatestOfType(actions.AttackDown)
				downOK = stampedDown.Tick == scene.CurrentTick()

				if attackCooldown.Aerial.Available(scene.CurrentTick()) && !downOK && !upOK {
					attack := combatdata.AerialSeqs[*cKey].First()
					attack.StartTick = scene.CurrentTick()
					attackCooldown.Aerial.StartTick = scene.CurrentTick()
					attackCooldown.Aerial.Duration = 60
					attack.LRDirection = *spatial.Components.Direction.GetFromCursor(cursor)
					en.EnqueueAddComponentWithValue(combat.Components.Attack, attack)
				} else if attackCooldown.Aerial.Available(scene.CurrentTick()) && downOK {
					actionBuffer.ConsumeAction(actions.Down)
					attack := combatdata.AerialDownSmashes[*cKey]
					attack.StartTick = scene.CurrentTick()
					attackCooldown.Aerial.StartTick = scene.CurrentTick()
					attackCooldown.Aerial.Duration = 60

					attack.LRDirection = *spatial.Components.Direction.GetFromCursor(cursor)
					en.EnqueueAddComponentWithValue(combat.Components.Attack, attack)
				} else if attackCooldown.Aerial.Available(scene.CurrentTick()) && upOK {

					attack := combatdata.AerialUpSeq[*cKey].First()
					attack.StartTick = scene.CurrentTick()
					attackCooldown.Aerial.StartTick = scene.CurrentTick()
					attackCooldown.Aerial.Duration = 60

					attack.LRDirection = *spatial.Components.Direction.GetFromCursor(cursor)
					en.EnqueueAddComponentWithValue(combat.Components.Attack, attack)
				}

			}

		}
	}
	return nil
}

type AttackHurtBoxCollisionSystem struct{}

func (s AttackHurtBoxCollisionSystem) Run(scene blueprint.Scene, dt float64) error {
	hurtBoxesQuery := warehouse.Factory.NewQuery().And(
		combat.Components.HurtBoxes, spatial.Components.Position,
		warehouse.Factory.NewQuery().Not(combat.Components.Invincible, combat.Components.Defeat),
	)
	hurtBoxCursor := scene.NewCursor(hurtBoxesQuery)

	attackQuery := warehouse.Factory.NewQuery().And(
		combat.Components.Attack, spatial.Components.Position,
		warehouse.Factory.NewQuery().Not(combat.Components.Defeat),
	)
	attackCursor := scene.NewCursor(attackQuery)

	for range hurtBoxCursor.Next() {
		for range attackCursor.Next() {

			if (components.MobTag.CheckCursor(hurtBoxCursor) && !components.SwapVulnerableComponent.CheckCursor(hurtBoxCursor)) &&
				(components.MobTag.CheckCursor(attackCursor) && !components.SwapVulnerableComponent.CheckCursor(attackCursor)) {
				continue
			}
			err := s.resolve(hurtBoxCursor, attackCursor, scene.CurrentTick())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (AttackHurtBoxCollisionSystem) resolve(hurtBoxCursor, attackCursor *warehouse.Cursor, tick int) error {
	aCurr, _ := attackCursor.CurrentEntity()
	hCurr, _ := hurtBoxCursor.CurrentEntity()

	if aCurr.ID() == hCurr.ID() { // Prevent self-harm
		return nil
	}

	var hitOccurred bool

	hurtCenterPos := spatial.Components.Position.GetFromCursor(hurtBoxCursor)
	hurtBoxesComponent := combat.Components.HurtBoxes.GetFromCursor(hurtBoxCursor)

	atkPos := spatial.Components.Position.GetFromCursor(attackCursor)
	atk := combat.Components.Attack.GetFromCursor(attackCursor)

	if atk.Length == 0 { // Attack has no active frames
		return nil
	}
	index := ((tick - atk.StartTick) / atk.Speed) % atk.Length

	// Validate attack data for the current frame
	if index >= len(atk.Boxes) || index >= len(atk.BoxesPositionOffsets) {
		// log.Printf("Attack %s frame %d out of bounds for boxes/offsets", atk.Name, index)
		return nil
	}
	if index < 0 {
		log.Println("bad index for attack(liekly from save)")
		en, err := attackCursor.CurrentEntity()
		if err != nil {
			return err
		}
		en.EnqueueRemoveComponent(combat.Components.Attack)
		return nil
	}
	atkBoxesInFrame := atk.Boxes[index]
	atkBoxOffsetsInFrame := atk.BoxesPositionOffsets[index]

HitCheckLoop:
	for atkBoxIdx, currentAtkBox := range atkBoxesInFrame {
		if currentAtkBox.LocalAAB.Height == 0 { // Skip inactive/empty hitboxes
			continue
		}
		if atkBoxIdx >= len(atkBoxOffsetsInFrame) { // Missing offset for this box
			continue
		}

		worldAtkBoxX := atkPos.X + atkBoxOffsetsInFrame[atkBoxIdx].X
		worldAtkBoxY := atkPos.Y + atkBoxOffsetsInFrame[atkBoxIdx].Y

		attackerDirPtr, hasAttackerDir := spatial.Components.Direction.GetFromCursorSafe(attackCursor)
		attackerDir := *attackerDirPtr

		if atk.LRDirection.Valid() {
			attackerDir = atk.LRDirection
			hasAttackerDir = true
		}

		if hasAttackerDir && attackerDir.IsLeft() {
			worldAtkBoxX = atkPos.X - atkBoxOffsetsInFrame[atkBoxIdx].X
		}

		atkOffsetPos := vector.Two{X: worldAtkBoxX, Y: worldAtkBoxY}

		hurtEntityDir, hasHurtEntityDir := spatial.Components.Direction.GetFromCursorSafe(hurtBoxCursor)

		for currentHurtBox := range hurtBoxesComponent.Active() {
			adjustedHurtBoxPos := currentHurtBox.RelativePos
			if hasHurtEntityDir && hurtEntityDir.IsLeft() {
				// Assumes RelativePos.X is positive for base right-facing pose
				adjustedHurtBoxPos.X = -currentHurtBox.RelativePos.X
			}
			finalHurtBoxWorldPos := hurtCenterPos.Two.Add(adjustedHurtBoxPos)

			collisionDetected, _ := spatial.Detector.Check(
				spatial.Shape(currentHurtBox.Shape), spatial.Shape(currentAtkBox), finalHurtBoxWorldPos, atkOffsetPos,
			)

			if collisionDetected {
				hitOccurred = true
				break HitCheckLoop // Single hit registered, proceed to effects
			}
		}
	}

	if hitOccurred {
		if !combat.Components.Invincible.CheckCursor(hurtBoxCursor) {
			err := hCurr.EnqueueAddComponentWithValue(combat.Components.Invincible, combat.Invincible{StartTick: tick})
			if err != nil {
				return err
			}
		}

		hurtEntityPos := spatial.Components.Position.GetFromCursor(hurtBoxCursor)
		attackerEntityPos := spatial.Components.Position.GetFromCursor(attackCursor)
		health := combat.Components.Health.GetFromCursor(hurtBoxCursor)

		var hitDirection vector.Two
		if hurtEntityPos.X > attackerEntityPos.X {
			hitDirection = vector.Two{X: -1, Y: 0}
		} else {
			hitDirection = vector.Two{X: 1, Y: 0}
		}

		if existingHurt, hasHurt := combat.Components.Hurt.GetFromCursorSafe(hurtBoxCursor); !hasHurt {
			err := hCurr.EnqueueAddComponentWithValue(combat.Components.Hurt, combat.Hurt{StartTick: tick, Direction: hitDirection})
			if err != nil {
				return err
			}
		} else {
			existingHurt.StartTick = tick // Re-trigger/extend hurt state
			existingHurt.Direction = hitDirection
		}

		health.Value -= atk.Damage
		atk.LastHitTick = tick

		if lc, ok := components.LastCombatComponent.GetFromCursorSafe(attackCursor); ok {
			lc.StartTick = tick
		}
		if lc, ok := components.LastCombatComponent.GetFromCursorSafe(hurtBoxCursor); ok {
			lc.StartTick = tick
		}

		return nil
	}

	return nil
}

type IframeRemoverSystem struct{}

func (IframeRemoverSystem) Run(scene blueprint.Scene, dt float64) error {
	query := warehouse.Factory.NewQuery().And(combat.Components.Invincible)
	cursor := scene.NewCursor(query)

	IFRAME_DURATION_IN_TICKS := 25

	for range cursor.Next() {
		if components.PlayerTag.CheckCursor(cursor) {
			IFRAME_DURATION_IN_TICKS = 30
		}
		iFrame := combat.Components.Invincible.GetFromCursor(cursor)

		if (iFrame.StartTick+IFRAME_DURATION_IN_TICKS)-scene.CurrentTick() < 0 {
			en, err := cursor.CurrentEntity()
			if err != nil {
				return err
			}
			err = en.EnqueueRemoveComponent(combat.Components.Invincible)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type hurtDurationRegistry map[characterkeys.CharEnum]int

var HurtDurationRegistry = map[characterkeys.CharEnum]int{
	characterkeys.BoxHead: 15,
}

type HurtRemoverSystem struct{}

func (HurtRemoverSystem) Run(scene blueprint.Scene, dt float64) error {
	query := warehouse.Factory.NewQuery().And(
		combat.Components.Hurt,
		components.CharacterKeyComponent,
	)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {

		hurt := combat.Components.Hurt.GetFromCursor(cursor)
		cKey := components.CharacterKeyComponent.GetFromCursor(cursor)
		hurtDura, ok := HurtDurationRegistry[*cKey]
		if !ok {
			// TODO: Log level based warning
			// log.Println("hey no value for hit duration for ckey (defaulting 45 ticks)", *cKey)
			hurtDura = 45
		}

		if (hurt.StartTick+hurtDura)-scene.CurrentTick() < 0 {
			en, err := cursor.CurrentEntity()
			if err != nil {
				return err
			}
			err = en.EnqueueRemoveComponent(combat.Components.Hurt)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type ContactDamageSystem struct{}

func (ContactDamageSystem) Run(scene blueprint.Scene, dt float64) error {
	playerQuery := warehouse.Factory.NewQuery().And(
		components.PlayerTag,
		warehouse.Factory.NewQuery().Not(
			combat.Components.Defeat,
		),
	)
	playerCursor := scene.NewCursor(playerQuery)
	mobQuery := warehouse.Factory.NewQuery().And(
		scenes.DefaultMobComposition,
		warehouse.Factory.NewQuery().Not(
			combat.Components.Defeat,
			components.PlayerTag,
			components.DodgeComponent,
			components.IgnoreContactDamageTag,
		),
	)
	mobCursor := scene.NewCursor(mobQuery)

	for range playerCursor.Next() {
		playerDirection, playerDirectionOK := spatial.Components.Direction.GetFromCursorSafe(playerCursor)
		for range mobCursor.Next() {

			if agro, friAgroOK := components.FriendlyAgroComponent.GetFromCursorSafe(mobCursor); friAgroOK {
				if agro.StartTick == 0 {
					continue
				}
			}

			mobDirection, mobDirectionOK := spatial.Components.Direction.GetFromCursorSafe(mobCursor)

			playerPos := spatial.Components.Position.GetFromCursor(playerCursor)
			playerBoxes := combat.Components.HurtBoxes.GetFromCursor(playerCursor)

			mobPos := spatial.Components.Position.GetFromCursor(mobCursor)
			mobHurtBoxes := combat.Components.HurtBoxes.GetFromCursor(mobCursor)

			for phb := range playerBoxes.Active() {
				for mhb := range mobHurtBoxes.Active() {

					adjP := phb.RelativePos
					adjM := mhb.RelativePos

					if playerDirectionOK {
						adjP.X = adjP.X * playerDirection.AsFloat()
					}

					if mobDirectionOK {
						adjM.X = adjM.X * mobDirection.AsFloat()
					}

					if ok, _ := spatial.Detector.Check(
						phb.Shape, mhb.Shape, playerPos.Add(adjP), mobPos.Add(adjM),
					); ok {

						en, err := playerCursor.CurrentEntity()
						tick := scene.CurrentTick()
						if err != nil {
							return err
						}
						isIv := combat.Components.Invincible.CheckCursor(playerCursor)
						isHurt := combat.Components.Hurt.CheckCursor(playerCursor)

						if !isHurt && !isIv {
							health := combat.Components.Health.GetFromCursor(playerCursor)
							health.Value -= 10

							var direction vector.Two

							if playerPos.X < mobPos.X {
								direction = vector.Two{X: 1, Y: 0}
							} else {
								direction = vector.Two{X: -1, Y: 0}
							}

							en.EnqueueAddComponentWithValue(combat.Components.Hurt, combat.Hurt{StartTick: tick, Direction: direction})
							en.EnqueueAddComponentWithValue(combat.Components.Invincible, combat.Invincible{StartTick: tick})

							contactData := components.ContactComponent.GetFromCursor(mobCursor)
							contactData.LastHit = tick

							if lc, ok := components.LastCombatComponent.GetFromCursorSafe(playerCursor); ok {
								lc.StartTick = tick
							}

							return nil
						}

					}
				}
			}

			// temp disable
			if cnpt, cnptOK := components.CannotDashThroughComponent.GetFromCursorSafe(mobCursor); cnptOK && false {
				blocker := cnpt.Blocker
				if blocker.LocalAAB.Height == 0 {
					return nil
				}
				playerShape := spatial.Components.Shape.GetFromCursor(playerCursor)

				if ok, col := spatial.Detector.Check(
					*playerShape, blocker, playerPos, mobPos,
				); ok {
					playerDyn := motion.Components.Dynamics.GetFromCursor(playerCursor)
					mobDyn := motion.NewDynamics(0)
					motion.Resolver.Resolve(&playerPos.Two, &mobPos.Two, playerDyn, &mobDyn, col)
				}
			}
		}
	}

	return nil
}

// ----------------------
type ProjectileCollisionSystem struct{}

func (s ProjectileCollisionSystem) Run(scene blueprint.Scene, dt float64) error {
	hurtBoxesQuery := warehouse.Factory.NewQuery().Or(
		warehouse.Factory.NewQuery().And(
			combat.Components.HurtBoxes,
			components.PlayerTag,
			warehouse.Factory.NewQuery().Not(
				combat.Components.Invincible,
				combat.Components.Defeat,
				components.SoftResetComponent,
			),
		),
		warehouse.Factory.NewQuery().And(
			combat.Components.HurtBoxes,
			components.SwapVulnerableComponent,
			warehouse.Factory.NewQuery().Not(
				combat.Components.Invincible,
				combat.Components.Defeat,
			),
		),
	)
	hurtBoxesCursor := scene.NewCursor(hurtBoxesQuery)

	projectileQuery := warehouse.Factory.NewQuery().And(components.ProjectileTag)
	projectileCursor := scene.NewCursor(projectileQuery)

	tick := scene.CurrentTick()

	for range hurtBoxesCursor.Next() {

		hitOccurred := false
		hitVel := vector.Two{}

		hCurr, _ := hurtBoxesCursor.CurrentEntity()
		hBoxesCenterPos := spatial.Components.Position.GetFromCursor(hurtBoxesCursor)
		hBoxes := combat.Components.HurtBoxes.GetFromCursor(hurtBoxesCursor)
		hurtDir := spatial.Components.Direction.GetFromCursor(hurtBoxesCursor)

		for range projectileCursor.Next() {
			projectilePos := spatial.Components.Position.GetFromCursor(projectileCursor)
			projectileShape := spatial.Components.Shape.GetFromCursor(projectileCursor)
			projectileDyn := motion.Components.Dynamics.GetFromCursor(projectileCursor)

			for hb := range hBoxes.Active() {
				hbPosAdjustement := hb.RelativePos

				if hurtDir.IsLeft() {
					hbPosAdjustement.X = -hbPosAdjustement.X
				}

				hurtWorldPos := hBoxesCenterPos.Two.Add(hbPosAdjustement)

				collisionDetected, _ := spatial.Detector.Check(
					hb.Shape, *projectileShape, hurtWorldPos, projectilePos,
				)

				if collisionDetected {
					hitOccurred = true
					hitVel = projectileDyn.Vel
					pCurr, err := projectileCursor.CurrentEntity()
					err = scene.Storage().EnqueueDestroyEntities(pCurr)
					if err != nil {
						return err
					}
				}
			}

		}

		if hitOccurred {
			if !combat.Components.Invincible.CheckCursor(hurtBoxesCursor) {
				err := hCurr.EnqueueAddComponentWithValue(combat.Components.Invincible, combat.Invincible{StartTick: tick})
				if err != nil {
					return err
				}
			}

			health := combat.Components.Health.GetFromCursor(hurtBoxesCursor)

			var hitDirection vector.Two
			if hitVel.X > 0 {
				hitDirection = vector.Two{X: -1, Y: 0}
			} else {
				hitDirection = vector.Two{X: 1, Y: 0}
			}

			if existingHurt, hasHurt := combat.Components.Hurt.GetFromCursorSafe(hurtBoxesCursor); !hasHurt {
				err := hCurr.EnqueueAddComponentWithValue(combat.Components.Hurt, combat.Hurt{StartTick: tick, Direction: hitDirection})
				if err != nil {
					return err
				}
			} else {
				existingHurt.StartTick = tick
				existingHurt.Direction = hitDirection
			}

			health.Value -= 10

			return nil
		}

	}
	return nil
}

type SwapVurnRemoverSystem struct{}

func (SwapVurnRemoverSystem) Run(scene blueprint.Scene, dt float64) error {
	query := warehouse.Factory.NewQuery().And(components.SwapVulnerableComponent)
	cursor := scene.NewCursor(query)

	DURATION := 90

	for range cursor.Next() {

		sv := components.SwapVulnerableComponent.GetFromCursor(cursor)

		if (sv.StartTick+DURATION)-scene.CurrentTick() < 0 {
			en, err := cursor.CurrentEntity()
			if err != nil {
				return err
			}
			err = en.EnqueueRemoveComponent(components.SwapVulnerableComponent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type SwapVurnStunSystem struct{}

func (SwapVurnStunSystem) Run(scene blueprint.Scene, dt float64) error {
	query := warehouse.Factory.NewQuery().And(components.SwapVulnerableComponent, motion.Components.Dynamics)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		dyn := motion.Components.Dynamics.GetFromCursor(cursor)
		dyn.Vel.X = 0
		dyn.Vel.Y = 5
	}
	return nil
}

type MobKBSystem struct{}

func (MobKBSystem) Run(scene blueprint.Scene, dt float64) error {
	query := warehouse.Factory.NewQuery().And(components.MobTag, combat.Components.Hurt)
	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		hurt := combat.Components.Hurt.GetFromCursor(cursor)
		diff := scene.CurrentTick() - hurt.StartTick
		if math.Abs(float64(diff)) < 10 {
			dyn := motion.Components.Dynamics.GetFromCursor(cursor)
			dyn.Vel.X = hurt.Direction.Norm().X * -30
		}

	}
	return nil
}

type PlayerAerialInterruptSystem struct{}

func (PlayerAerialInterruptSystem) Run(scene blueprint.Scene, dt float64) error {
	query := warehouse.Factory.NewQuery().And(components.PlayerTag, combat.Components.Attack, components.OnGroundComponent)
	cursor := scene.NewCursor(query)

	addPrimaryAtkBack := map[int]bool{}

	for range cursor.Next() {
		atk := combat.Components.Attack.GetFromCursor(cursor)
		isAerial := atk.ID == combatdata.AerialSeqs[characterkeys.BoxHead].First().ID
		onGround := components.OnGroundComponent.GetFromCursor(cursor)
		if onGround.LastTouch == scene.CurrentTick() && isAerial {

			en, err := cursor.CurrentEntity()
			if err != nil {
				return err
			}
			if math.Abs(float64(scene.CurrentTick()-atk.StartTick)) <= 15 {
				addPrimaryAtkBack[int(en.ID())] = true
			}
			en.EnqueueRemoveComponent(combat.Components.Attack)
		}
	}

	query = warehouse.Factory.NewQuery().And(components.PlayerTag)
	cursor = scene.NewCursor(query)

	for range cursor.Next() {
		en, err := cursor.CurrentEntity()
		if err != nil {
			return err
		}

		_, ok := addPrimaryAtkBack[int(en.ID())]

		if ok {
			atk := combatdata.PrimarySeqs[characterkeys.BoxHead].First()
			atk.LRDirection = *spatial.Components.Direction.GetFromCursor(cursor)
			atk.StartTick = scene.CurrentTick()
			en.EnqueueAddComponentWithValue(combat.Components.Attack, atk)
		}
	}

	return nil
}
