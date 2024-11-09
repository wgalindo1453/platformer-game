package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"platformer-game/core"
)

const (
	screenWidth  = 800
	screenHeight = 450
	worldWidth   = 5000
	worldHeight  = 1200
)

var gameOver bool

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Platformer Game")
	defer rl.CloseWindow()

	core.InitGame(worldWidth, worldHeight)

	for !rl.WindowShouldClose() && !gameOver {
		core.UpdateGame(worldHeight) //need to pass worldHeight to update zombies 
		core.DrawGame()
	}
}
