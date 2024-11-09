package core

import (
	"math/rand"
	"time"

	"platformer-game/gameobjects"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	camera     rl.Camera2D
	background rl.Texture2D
	zombies    []*gameobjects.Zombie // Slice to hold pointers to all zombies
)

const (
	worldWidth   = 5000
	worldHeight  = 1200
	screenWidth  = 800
	screenHeight = 450
)

const (
	miniMapWidth  = 200
	miniMapHeight = 150
	miniMapX      = screenWidth - miniMapWidth - 10
	miniMapY      = 10
	deadZoneWidth = 200
)

func InitGame(worldWidth, worldHeight int) {
	background = rl.LoadTexture("assets/levelonebg.png")

	// Initializing  player
	gameobjects.InitPlayer(worldWidth, worldHeight)

	// Initializing zombies
	initZombies(5) // Spawning 5 zombies

	// Initializing camera
	camera = rl.Camera2D{
		Target: gameobjects.PlayerInstance.Position,
		Offset: rl.NewVector2(float32(screenWidth)/2, float32(screenHeight)/2),
		Zoom:   1.0,
	}
}

// Initializing zombies with random positions
func initZombies(numZombies int) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < numZombies; i++ {
		// Randomize position within world boundaries
		x := float32(rand.Intn(worldWidth-100) + 50) // Keep zombies within world bounds
		y := float32(worldHeight - 50)               // Spawn zombies at ground level

		// Initialize zombie with random position and append to the zombies slice
		zombie := gameobjects.InitZombie(x, y-50, 1)
		zombies = append(zombies, &zombie)
	}
}

func UpdateGame(worldHeight int) {
	// Updating player and call Shoot to check for zombie hits
	gameobjects.PlayerInstance.Update(worldHeight, worldWidth, zombies)
	gameobjects.PlayerInstance.Shoot() // Call Shoot to check for zombie hits

	playerPosition := gameobjects.PlayerInstance.Position

	// Updating each zombie in the zombies slice
	for i := len(zombies) - 1; i >= 0; i-- {
		zombies[i].Update(worldWidth, playerPosition)
		if !zombies[i].IsAlive && zombies[i].State == gameobjects.ZombieDead && zombies[i].CurrentFrame == len(zombies[i].DeadFrames)-1 {
			// Remove zombie once dead animation completes
			zombies = append(zombies[:i], zombies[i+1:]...)
		}
	}

	// For more smooth camera transition using dead zones
	playerX := gameobjects.PlayerInstance.Position.X
	if playerX > camera.Target.X+float32(screenWidth)/2-deadZoneWidth {
		camera.Target.X = playerX - float32(screenWidth)/2 + deadZoneWidth
	} else if playerX < camera.Target.X-float32(screenWidth)/2+deadZoneWidth {
		camera.Target.X = playerX + float32(screenWidth)/2 - deadZoneWidth
	}

	// Keeping camera within world bounds
	camera.Target.X = clampFloat(camera.Target.X, float32(screenWidth)/2, float32(worldWidth)-float32(screenWidth)/2)
	camera.Target.Y = clampFloat(camera.Target.Y, float32(screenHeight)/2, float32(worldHeight)-float32(screenHeight)/2)
}

func DrawMiniMap() {
	rl.DrawRectangle(miniMapX, miniMapY, miniMapWidth, miniMapHeight, rl.LightGray)

	// Calculate scaling factors for mini-map
	scaleX := float32(miniMapWidth) / float32(worldWidth)
	scaleY := float32(miniMapHeight) / float32(worldHeight)

	// Drawing world boundary on mini-map
	rl.DrawRectangleLines(miniMapX, miniMapY, miniMapWidth, miniMapHeight, rl.DarkGray)

	// Drawing camera view on mini-map
	viewX := miniMapX + int((camera.Target.X-float32(screenWidth)/2)*scaleX)
	viewY := miniMapY + int((camera.Target.Y-float32(screenHeight)/2)*scaleY)
	viewWidth := int(float32(screenWidth) * scaleX * 0.8)
	viewHeight := int(float32(screenHeight) * scaleY * 0.8)

	viewX = clamp(viewX, miniMapX, miniMapX+miniMapWidth-viewWidth)
	viewY = clamp(viewY, miniMapY, miniMapY+miniMapHeight-viewHeight)
	rl.DrawRectangleLines(int32(viewX), int32(viewY), int32(viewWidth), int32(viewHeight), rl.Red)
}

func DrawGame() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	// Drawing game world with camera
	rl.BeginMode2D(camera)
	rl.DrawTexture(background, 0, 0, rl.White)
	gameobjects.PlayerInstance.Draw()

	// Draw each zombie in the zombies slice
	for _, zombie := range zombies {
		zombie.Draw()
	}
	rl.EndMode2D()

	
	DrawMiniMap()

	rl.EndDrawing()
}

// Utility function to clamp an integer within a range
func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Utility function to clamp a float32 within a range
func clampFloat(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
