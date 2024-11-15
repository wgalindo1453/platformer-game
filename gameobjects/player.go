package gameobjects

import (
	"fmt"
	"platformer-game/rendering"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PlayerState int

const (
	Idle PlayerState = iota
	Walking
	Running
	Shooting
	Sitting
	SittingShooting
	Jumping
	Resting
	Sleeping
	Dying
)
const (
	jumpVelocity = -5.0 // Initial upward velocity for jumping
	gravity      = 800  // Gravity value, pulling player down each frame
	groundYPos   = 0    // The ground level, adjust to your world height
)

type Player struct {
	Position              rl.Vector2
	Speed                 rl.Vector2
	Acceleration          rl.Vector2
	Width, Height         float32
	Color                 rl.Color
	FacingRight           bool           // Direction the player is facing
	CurrentFrame          int            // Current frame index for animation
	FrameCounter          int            // Counter to control frame switch timing
	State                 PlayerState    // Current animation state
	IdleTimer             time.Time      // Timer for idle state
	RestTimer             time.Time      // Timer for resting state
	WalkFrames            []rl.Texture2D // Frames for walking animation
	RunFrames             []rl.Texture2D // Frames for running animation
	IdleFrames            []rl.Texture2D // Frames for idle animation
	ShootFrames           []rl.Texture2D // Frames for shooting animation
	SittingFrames         []rl.Texture2D // Frames for sitting animation
	SittingShootingFrames []rl.Texture2D // Frames for sitting shooting animation
	JumpFrames            []rl.Texture2D // Frames for jumping animation
	RestingFrames         []rl.Texture2D // Frames for resting animation
	SleepingFrames        []rl.Texture2D // Frames for sleeping animation
	DyingFrames           []rl.Texture2D // Frames for dying animation
	Bullets               []*Bullet      // Add bullets slice
	switchDown            bool           // Indicates when to start descending

	// Sounds
	WalkSound  rl.Sound
	RunSound   rl.Sound
	ShootSound rl.Sound

	// New attributes
	Health    float64 // Player health
	MaxHealth float64 // Maximum health to keep track for the health bar
	Inventory  Inventory
	HeldItem Item // The currently held item
}
func (p *Player) UpdateHeldItem() {
    if p.Inventory.Slots[p.Inventory.SelectedSlot].Type != Other {
        p.HeldItem = p.Inventory.Slots[p.Inventory.SelectedSlot]
    } else {
        p.HeldItem = Item{} // No item held if slot is empty
    }
}

func (p *Player) Shoot() {
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		bulletPosition := p.Position
		bulletPosition.Y += p.Height / 2 // Adjust to shoot from the middle
		newBullet := NewBullet(bulletPosition.X, bulletPosition.Y, 10, p.FacingRight)
		p.Bullets = append(p.Bullets, newBullet)

		// Play shoot sound
		if !rl.IsSoundPlaying(p.ShootSound) {
			rl.PlaySound(p.ShootSound)
		}
	}
}
func (p *Player) IsGameOver() bool {

	return p.Health <= 0
}

func (p *Player) Unload() {
	for _, frame := range p.WalkFrames {
		rl.UnloadTexture(frame)
	}
	for _, frame := range p.RunFrames {
		rl.UnloadTexture(frame)
	}
	for _, frame := range p.IdleFrames {
		rl.UnloadTexture(frame)
	}
	for _, frame := range p.ShootFrames {
		rl.UnloadTexture(frame)
	}
	// Unload sounds
	rl.UnloadSound(p.WalkSound)
	rl.UnloadSound(p.RunSound)
	rl.UnloadSound(p.ShootSound)
}

var PlayerInstance Player

func InitPlayer(worldWidth, worldHeight int) {

	rl.InitAudioDevice() // Initialize audio device
	PlayerInstance = Player{
		Position:     rl.NewVector2(100, float32(worldHeight-50)),
		Speed:        rl.NewVector2(0, 0),
		Acceleration: rl.NewVector2(0, 0.5),
		Width:        113,
		Height:       113,
		Color:        rl.White,
		CurrentFrame: 0,
		FrameCounter: 0,
		State:        Idle,
		FacingRight:  true,
		Health:       100, // Initialize with full health
		MaxHealth:    100, // Set maximum health
		Inventory: NewInventory(10), // Initialize with 10 slots
	}
	// Load sounds
	PlayerInstance.WalkSound = rl.LoadSound("assets/sounds/walking.mp3")
	PlayerInstance.RunSound = rl.LoadSound("assets/sounds/running.mp3")
	PlayerInstance.ShootSound = rl.LoadSound("assets/sounds/machineguneffect.wav")

	// Load sprite sheet
	spriteSheet := rendering.LoadSpriteSheet("assets/sprites/shooterspritesheet.png")
	spriteSheet2 := rendering.LoadSpriteSheet("assets/sprites/shooterspritesheet2.png")

	// Load walking frames
	walkingFrames := []rl.Rectangle{
		{X: 309, Y: 301, Width: 63, Height: 136},  // Frame 1
		{X: 500, Y: 301, Width: 66, Height: 136},  // Frame 2
		{X: 690, Y: 303, Width: 72, Height: 134},  // Frame 3
		{X: 878, Y: 302, Width: 72, Height: 136},  // Frame 4
		{X: 1075, Y: 299, Width: 70, Height: 138}, // Frame 5
	}
	for _, frame := range walkingFrames {
		PlayerInstance.WalkFrames = append(PlayerInstance.WalkFrames, spriteSheet.ImageAt(frame, rl.Blank))
	}

	// Load running frames
	runningFrames := []rl.Rectangle{
		{X: 267, Y: 525, Width: 76, Height: 122},  // Frame 1
		{X: 456, Y: 535, Width: 78, Height: 122},  // Frame 2
		{X: 644, Y: 535, Width: 84, Height: 122},  // Frame 3
		{X: 840, Y: 525, Width: 80, Height: 122},  // Frame 4
		{X: 1042, Y: 533, Width: 68, Height: 124}, // Frame 5
	}
	for _, frame := range runningFrames {
		PlayerInstance.RunFrames = append(PlayerInstance.RunFrames, spriteSheet.ImageAt(frame, rl.Blank))
	}

	// Load idle frames
	idleFrames := []rl.Rectangle{
		{X: 296, Y: 71, Width: 94, Height: 134},  // Frame 1
		{X: 488, Y: 71, Width: 94, Height: 134},  // Frame 2
		{X: 681, Y: 69, Width: 94, Height: 136},  // Frame 3
		{X: 873, Y: 69, Width: 94, Height: 136},  // Frame 4
		{X: 1063, Y: 69, Width: 94, Height: 136}, // Frame 5
		{X: 1256, Y: 71, Width: 93, Height: 134}, // Frame 6
	}
	for _, frame := range idleFrames {
		PlayerInstance.IdleFrames = append(PlayerInstance.IdleFrames, spriteSheet.ImageAt(frame, rl.Blank))
	}

	// Load shooting frames
	shooting1Frames := []rl.Rectangle{
		{X: 294, Y: 739, Width: 95, Height: 130},  // Frame 1
		{X: 487, Y: 739, Width: 108, Height: 130}, // Frame 2
		{X: 677, Y: 739, Width: 125, Height: 130}, // Frame 3
		{X: 869, Y: 739, Width: 102, Height: 131}, // Frame 4
		{X: 300, Y: 951, Width: 102, Height: 130}, // Frame 5
		{X: 492, Y: 951, Width: 111, Height: 130}, // Frame 6
		{X: 683, Y: 951, Width: 130, Height: 130}, // Frame 7
		{X: 877, Y: 951, Width: 106, Height: 130}, // Frame 8
	}
	for _, frame := range shooting1Frames {
		PlayerInstance.ShootFrames = append(PlayerInstance.ShootFrames, spriteSheet.ImageAt(frame, rl.Blank))
	}

	// Load Sitting frames
	sittingFrames := []rl.Rectangle{
		{X: 234, Y: 82, Width: 75, Height: 88}, // Frame 1
		{X: 394, Y: 83, Width: 75, Height: 87}, // Frame 2
		{X: 555, Y: 85, Width: 75, Height: 86}, // Frame 3
	}
	for _, frame := range sittingFrames {
		PlayerInstance.SittingFrames = append(PlayerInstance.SittingFrames, spriteSheet2.ImageAt(frame, rl.Blank))
	}

	//Load Sitting Shooting frames
	sittingShootingFrames := []rl.Rectangle{
		{X: 242, Y: 275, Width: 85, Height: 89},  // Frame 1
		{X: 399, Y: 275, Width: 84, Height: 89},  // Frame 2
		{X: 560, Y: 275, Width: 110, Height: 89}, // Frame 3
	}
	for _, frame := range sittingShootingFrames {
		PlayerInstance.SittingShootingFrames = append(PlayerInstance.SittingShootingFrames, spriteSheet2.ImageAt(frame, rl.Blank))
	}

	//Jumping frames
	jumpingFrames := []rl.Rectangle{
		{X: 240, Y: 444, Width: 78, Height: 103}, // Frame 1
		{X: 401, Y: 450, Width: 80, Height: 96},  // Frame 2
		{X: 561, Y: 434, Width: 79, Height: 113}, // Frame 3
		{X: 722, Y: 444, Width: 78, Height: 98},  // Frame 4
		{X: 1043, Y: 457, Width: 68, Height: 89}, // Frame 5

	}
	for _, frame := range jumpingFrames {
		PlayerInstance.JumpFrames = append(PlayerInstance.JumpFrames, spriteSheet2.ImageAt(frame, rl.Blank))
	}

	//Resting frames
	restingFrames := []rl.Rectangle{
		{X: 240, Y: 621, Width: 78, Height: 102}, // Frame 1
		{X: 400, Y: 626, Width: 78, Height: 97},  // Frame 2
		{X: 559, Y: 644, Width: 71, Height: 79},  // Frame 3
		{X: 686, Y: 651, Width: 87, Height: 72},  // Frame 4
	}

	for _, frame := range restingFrames {
		PlayerInstance.RestingFrames = append(PlayerInstance.RestingFrames, spriteSheet2.ImageAt(frame, rl.Blank))
	}

	//Sleeping frames
	sleepingFrames := []rl.Rectangle{
		{X: 231, Y: 864, Width: 113, Height: 32}, // Frame 1
		{X: 390, Y: 847, Width: 115, Height: 49}, // Frame 2
		{X: 541, Y: 825, Width: 124, Height: 71}, // Frame 3
		{X: 711, Y: 864, Width: 114, Height: 32}, // Frame 4
		{X: 869, Y: 863, Width: 114, Height: 33}, // Frame 5

	}
	for _, frame := range sleepingFrames {
		PlayerInstance.SleepingFrames = append(PlayerInstance.SleepingFrames, spriteSheet2.ImageAt(frame, rl.Blank))
	}

	//Dying frames
	dyingFrames := []rl.Rectangle{
		{X: 315, Y: 952, Width: 92, Height: 128},
		{X: 504, Y: 943, Width: 94, Height: 137},
		{X: 651, Y: 984, Width: 128, Height: 96},
		{X: 814, Y: 1041, Width: 160, Height: 39},
	}

	for _, frame := range dyingFrames {
		PlayerInstance.DyingFrames = append(PlayerInstance.DyingFrames, spriteSheet2.ImageAt(frame, rl.Blank))
	}

}

/***********************************STATES*********************************************** */

func (p *Player) setState(state PlayerState) {
	if p.State != state {
		p.State = state
		p.CurrentFrame = 0
		p.FrameCounter = 0
	}

	// Reset timers when changing to idle, resting, or sleeping states
	if state == Idle {
		p.IdleTimer = time.Now()
		p.RestTimer = time.Time{}
	} else if state == Resting {
		p.RestTimer = time.Now()
	} else {
		p.IdleTimer = time.Time{}
		p.RestTimer = time.Time{}
	}
}

/***********************************UPDATE*********************************************** */

func (p *Player) Update(worldHeight int, worldWidth int, zombies []*Zombie) {
	// fmt.Println("players starting out y position: ", p.Position.Y)

	// Update bullets
	for _, bullet := range p.Bullets {
		if bullet.IsActive {
			bullet.Update()

			// Here we are checking if bullet hits any zombie
			for _, zombie := range zombies {
				if zombie.IsAlive && rl.CheckCollisionPointCircle(bullet.Position, zombie.Position, zombie.Width/2) {
					zombie.TakeDamage(20) // Adjust damage as needed
					bullet.IsActive = false
					break
				}
			}

			// Deactivate bullet if it goes out of bounds
			if bullet.Position.X < 0 || bullet.Position.X > float32(worldWidth) {
				bullet.IsActive = false
			}
		}
	}

	//print inventory
	//if inventory is not empty then print the item inside
	if ( len(p.Inventory.Slots) != 0){
		// fmt.Println("Inventory:", p.Inventory)

	}

	// Filter out inactive bullets
	activeBullets := p.Bullets[:0]
	for _, bullet := range p.Bullets {
		if bullet.IsActive {
			activeBullets = append(activeBullets, bullet)
		}
	}
	p.Bullets = activeBullets
	// Check if player is on the ground
	onGround := p.Position.Y >= float32(worldHeight)-p.Height

	// Apply gravity and handle jumping
	if !onGround || p.State == Jumping {
		// Print player's position for debugging
		// fmt.Println("Player Y Position:", p.Position.Y)

		// Apply gravity effect based on ascending or descending state
		if p.Speed.Y < 0 && !p.switchDown { // Ascending
			fmt.Println("Ascending")

			// Switch to descending if near the apex
			if p.Speed.Y >= -0.5 { // Lower threshold for more gradual transition
				fmt.Println("Switching down")
				p.switchDown = true
			}
			p.Speed.Y += gravity * 0.000001 // Reduce gravity effect while ascending
		} else { // Descending
			p.Speed.Y += gravity * 0.005 // Normal gravity effect for descent
		}

		// Update the player's vertical position with the adjusted speed
		p.Position.Y += p.Speed.Y


	}

	// If player is grounded and was jumping, reset to Idle and reset switchDown
	if p.Position.Y >= float32(worldHeight)-p.Height {
		p.Position.Y = float32(worldHeight) - p.Height
		p.Speed.Y = 0
		p.switchDown = false // Reset switchDown for the next jump
		if p.State == Jumping {
			p.setState(Idle) // Reset to Idle after landing
		}
	}

	// Player state logic based on key inputs, prioritizing crouching
	switch {
	case rl.IsKeyDown(rl.KeyLeftControl):
		// Crouching has priority, halts forward movement
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			fmt.Println("Sitting and shooting")
			p.setState(SittingShooting)
			p.Shoot() // Call shoot when sitting and shooting
			//call shoot method simul
			p.Speed.X = 0 // Halt horizontal movement
			rl.StopSound(p.WalkSound)

			if !rl.IsSoundPlaying(p.ShootSound) {
				rl.PlaySound(p.ShootSound)
			}
			//stop walking sound
			rl.StopSound(p.WalkSound)

		} else {
			p.setState(Sitting)
			p.Speed.X = 0 // Halt horizontal movement
			rl.StopSound(p.WalkSound)
		}

		// When initiating the jump, set a lower initial speed
	case rl.IsKeyPressed(rl.KeySpace) && onGround:
		// Jump initiation
		p.setState(Jumping)
		p.Speed.Y = -2.0 // Lower initial jump speed for a shorter jump
		fmt.Println("Jumping")

		// Apply gravity and handle jumping
		if !onGround || p.State == Jumping {
			// Print player's position for debugging
			fmt.Println(p.Position.Y)

			if p.Speed.Y < 0 && !p.switchDown { // Ascending
				fmt.Println("Ascending")
				if p.Speed.Y > -0.1 {
					fmt.Println("Switching down")
					p.switchDown = true
				}
				p.Speed.Y += p.Acceleration.Y * 0.001 // Maintain slow upward deceleration
			} else { // Descending
				p.Speed.Y += p.Acceleration.Y * 0.1 // Slightly faster but controlled descent
			}
			p.Position.Y += p.Speed.Y
		}

		// If player lands on the ground, reset to Idle and reset switchDown
		if p.Position.Y >= float32(worldHeight)-p.Height {
			p.Position.Y = float32(worldHeight) - p.Height
			p.Speed.Y = 0
			p.switchDown = false // Reset switchDown for the next jump
			if p.State == Jumping {
				p.setState(Idle) // Reset to Idle after landing
			}
		}

	case rl.IsMouseButtonDown(rl.MouseLeftButton) && p.State != Sitting && p.State != SittingShooting:
		// Shooting (no horizontal movement)
		p.setState(Shooting)
		p.Speed.X = 0
		if !rl.IsSoundPlaying(p.ShootSound) {
			rl.PlaySound(p.ShootSound)
		}
		//stop walking sound
		rl.StopSound(p.WalkSound)
		//stop running sound
		rl.StopSound(p.RunSound)

	case rl.IsKeyDown(rl.KeyD) && rl.IsKeyDown(rl.KeyLeftShift) && p.State != Shooting && p.State != Sitting:
		// Running (right) if not shooting or crouching
		p.setState(Running)
		p.FacingRight = true
		p.Speed.X = 0.2
		if !rl.IsSoundPlaying(p.RunSound) {
			rl.PlaySound(p.RunSound)
		}
		rl.StopSound(p.WalkSound)

	case rl.IsKeyDown(rl.KeyD) && p.State != Shooting && p.State != Sitting && p.State != SittingShooting:
		// Walking (right) if not shooting or crouching
		p.setState(Walking)
		p.FacingRight = true
		p.Speed.X = 0.05
		if !rl.IsSoundPlaying(p.WalkSound) {
			rl.PlaySound(p.WalkSound)
		}
		rl.StopSound(p.RunSound)

	case rl.IsKeyDown(rl.KeyA) && rl.IsKeyDown(rl.KeyLeftShift) && p.State != Shooting && p.State != Sitting:
		// Running (left) if not shooting or crouching
		p.setState(Running)
		p.FacingRight = false
		p.Speed.X = -0.2
		if !rl.IsSoundPlaying(p.RunSound) {
			rl.PlaySound(p.RunSound)
		}
		rl.StopSound(p.WalkSound)

	case rl.IsKeyDown(rl.KeyA) && p.State != Shooting && p.State != Sitting && p.State != SittingShooting:
		// Walking (left) if not shooting or crouching
		p.setState(Walking)
		p.FacingRight = false
		p.Speed.X = -0.05
		if !rl.IsSoundPlaying(p.WalkSound) {
			rl.PlaySound(p.WalkSound)
		}
		rl.StopSound(p.RunSound)

	case onGround && p.State != Resting && p.State != Sleeping:
		// Idle if no movement
		p.setState(Idle)
		p.Speed.X = 0
		rl.StopSound(p.WalkSound)
		rl.StopSound(p.RunSound)
		rl.StopSound(p.ShootSound)
	}

	if !rl.IsMouseButtonDown(rl.MouseLeftButton) {
		rl.StopSound(p.ShootSound)
	}

	// Update horizontal position
	p.Position.X += p.Speed.X

	// this is to constrain player within screen bounds (X-axis)
	if p.Position.X < 0 {
		p.Position.X = 0
	} else if p.Position.X > float32(worldWidth)-p.Width {
		p.Position.X = float32(worldWidth) - p.Width
	}

	// this is to Ensure player doesn't sink below ground level (Y-axis)
	if p.Position.Y >= float32(worldHeight) {
		p.Position.Y = float32(worldHeight)
		p.Speed.Y = 0
	}

	// Updating animation frames based on state of the player
	p.FrameCounter++
	var frames []rl.Texture2D
	frameDelay := 300
	switch p.State {
		case Walking:
			frames = p.WalkFrames
		case Running:
			frames = p.RunFrames
		case Shooting:
			frames = p.ShootFrames
		case Sitting:
			frames = p.SittingFrames
		case SittingShooting:
			frames = p.SittingShootingFrames
		case Jumping:
			frames = p.JumpFrames
			frameDelay = 500
		case Resting:
			frames = p.RestingFrames
			frameDelay = 5000
		case Sleeping:
			frames = p.SleepingFrames
			frameDelay = 5000
		case Dying:
			frames = p.DyingFrames
			frameDelay = 10000
		default:
			frames = p.IdleFrames
	}

	// Only update frame based on delay
	if len(frames) > 0 && p.FrameCounter >= frameDelay {
		p.CurrentFrame = (p.CurrentFrame + 1) % len(frames)
		p.FrameCounter = 0
	}
}

/***********************************DRAW*********************************************** */

func (p *Player) Draw() {

	var frame rl.Texture2D
	switch p.State {
	case Walking:
		frame = p.WalkFrames[p.CurrentFrame]
	case Running:
		frame = p.RunFrames[p.CurrentFrame]
	case Shooting:
		frame = p.ShootFrames[p.CurrentFrame]
	case Sitting:
		frame = p.SittingFrames[p.CurrentFrame]
	case SittingShooting:
		frame = p.SittingShootingFrames[p.CurrentFrame]
	case Jumping:
		frame = p.JumpFrames[p.CurrentFrame]
	case Resting:
		frame = p.RestingFrames[p.CurrentFrame]
	case Sleeping:
		frame = p.SleepingFrames[p.CurrentFrame]
	case Dying:
		frame = p.DyingFrames[p.CurrentFrame]
	default:
		frame = p.IdleFrames[p.CurrentFrame]
	}

	if p.HeldItem.Type != Other && p.HeldItem.Image.ID != 0 {
        heldX := p.Position.X  -10 // Adjust for desired position relative to player
        heldY := p.Position.Y - 10 // Adjust for desired position relative to player
        rl.DrawTextureEx(p.HeldItem.Image, rl.Vector2{X: heldX, Y: heldY}, 0, 0.5, rl.White) // Scale to desired size
    }

	// Source rectangle starts normally
	sourceRect := rl.Rectangle{X: 0, Y: 0, Width: float32(frame.Width), Height: float32(frame.Height)}
	// Flip the source width to achieve a horizontal flip
	if !p.FacingRight {
		sourceRect.Width = -sourceRect.Width
	}

	// Destination rectangle keeps player position and scale
	destinationRect := rl.Rectangle{
		X:      p.Position.X,
		Y:      p.Position.Y,
		Width:  p.Width,
		Height: p.Height,
	}

	// Draw the current frame with adjusted sourceRect for flipping
	if frame.ID != 0 { // Ensure the frame texture is loaded
		rl.DrawTexturePro(
			frame,
			sourceRect,      // Flipped if FacingRight is false
			destinationRect, // Destination position and size on the screen
			rl.Vector2{X: p.Width / 2, Y: p.Height / 2}, // Origin remains centered
			0,       // No rotation
			p.Color, // Tint color
		)
	}

	// Drawing bullets
	for _, bullet := range p.Bullets {
		bullet.Draw()
	}
}
