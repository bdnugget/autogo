package main

import (
	"fmt"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SQUARE_SIZE  int32 = 200
	FPS_TARGET   int32 = 60
	GARAGE_COUNT int32 = 5
)

type Car struct {
	position  rl.Vector2
	size      rl.Vector2
	color     rl.Color
	direction int32
}

type Garage struct {
	position rl.Vector2
	size     rl.Vector2
	color    rl.Color
}

type GameState struct {
	score         int32
	car           Car
	framesCounter int32
	gameOver      bool
	pause         bool
	allowMove     bool
	garages       []Garage
}

var (
	screenWidth  int32 = 1600
	screenHeight int32 = 1000

	autoTexture   rl.Texture2D
	autoHappy     rl.Sound
	autoSad       rl.Sound
	garageColors        = []rl.Color{rl.Red, rl.Green, rl.Blue, rl.Yellow, rl.Purple}
	gameState     GameState
	offset        rl.Vector2
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Dikke Vette Cargame voor Milo")
	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()
	defer rl.CloseWindow()

	autoHappy = rl.LoadSound("resources/auto_happy_vob.ogg")
	autoSad = rl.LoadSound("resources/auto_sad_vob.ogg")
	autoTexture = rl.LoadTexture("resources/car_200px.png")
	defer rl.UnloadSound(autoHappy)
	defer rl.UnloadSound(autoSad)
	defer rl.UnloadTexture(autoTexture)

	offset = rl.Vector2{X: float32(screenWidth % SQUARE_SIZE), Y: float32(screenHeight % SQUARE_SIZE)}

	InitGame()

	rl.SetTargetFPS(FPS_TARGET)

	for !rl.WindowShouldClose() {
		UpdateDrawFrame()
	}
}

func InitGame() {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(int(GARAGE_COUNT))

	gameState.framesCounter = 0
	gameState.gameOver = false
	gameState.pause = false
	gameState.allowMove = false
	gameState.car = Car{
		position:  rl.Vector2{X: offset.X / 2, Y: offset.Y/2 + 2*float32(SQUARE_SIZE)},
		size:      rl.Vector2{X: float32(SQUARE_SIZE), Y: float32(SQUARE_SIZE)},
		color:     garageColors[randomIndex],
		direction: 0,
	}
	gameState.garages = make([]Garage, GARAGE_COUNT)

	InitGarages()
}

func UpdateDrawFrame() {
	UpdateGame()
	DrawGame()
}

func UpdateGame() {
	if !gameState.gameOver {
		HandlePause()
		HandleMovement()
		CheckCollisions()
		IncrementFrameCounter()
	} else {
		HandleRestart()
	}
}

func HandlePause() {
	if rl.IsKeyPressed(int32('P')) {
		gameState.pause = !gameState.pause
	}
}

func HandleMovement() {
	if !gameState.pause {
		if rl.IsKeyPressed(rl.KeyUp) && gameState.allowMove && gameState.car.position.Y > 0 {
			gameState.car.direction = -1
			gameState.allowMove = false
		}
		if (rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyRight)) &&
			gameState.allowMove &&
			gameState.car.position.Y < float32(screenHeight-SQUARE_SIZE) {
			gameState.car.direction = 1
			gameState.allowMove = false
		}

		if gameState.framesCounter%(FPS_TARGET/60) == 0 {
			MoveCar()
		}
	}
}

func MoveCar() {
	gameState.car.position.X += float32(SQUARE_SIZE) / float32(FPS_TARGET)
	gameState.car.position.Y += float32(gameState.car.direction * SQUARE_SIZE)
	gameState.allowMove = true
	gameState.car.direction = 0
}

func CheckCollisions() {
	for i := int32(0); i < GARAGE_COUNT; i++ {
		if rl.CheckCollisionRecs(
			rl.Rectangle{X: gameState.car.position.X, Y: gameState.car.position.Y, Width: gameState.car.size.X, Height: gameState.car.size.Y},
			rl.Rectangle{X: gameState.garages[i].position.X, Y: gameState.garages[i].position.Y, Width: gameState.garages[i].size.X, Height: gameState.garages[i].size.Y},
		) {
			HandleCollision(i)
		}
	}
}

func HandleCollision(index int32) {
	if rl.ColorToInt(gameState.car.color) == rl.ColorToInt(gameState.garages[index].color) {
		rl.PlaySound(autoHappy)
		gameState.score++
		InitGame()
	} else {
		rl.PlaySound(autoSad)
		gameState.gameOver = true
	}
}

func IncrementFrameCounter() {
	gameState.framesCounter++
}

func HandleRestart() {
	if rl.IsKeyPressed(rl.KeyEnter) {
		gameState.score = 0
		InitGame()
	}
}

func DrawGame() {
	rl.BeginDrawing()

	rl.ClearBackground(rl.DarkGray)

	if !gameState.gameOver {
		for i := int32(1); i < screenHeight/SQUARE_SIZE+1; i++ {
			for j := int32(0); j < screenWidth/(SQUARE_SIZE/2); j += 2 {
				rl.DrawRectangle(
					int32(j)*((SQUARE_SIZE/2)+int32(offset.X)/2),
					int32(i*SQUARE_SIZE+int32(offset.Y)/2),
					SQUARE_SIZE/2,
					SQUARE_SIZE/8,
					rl.RayWhite,
				)
			}
		}

		// Draw garages
		for i := int32(0); i < GARAGE_COUNT; i++ {
			rl.DrawRectangleRec(
				rl.Rectangle{X: gameState.garages[i].position.X, Y: gameState.garages[i].position.Y, Width: gameState.garages[i].size.X, Height: gameState.garages[i].size.Y},
				gameState.garages[i].color,
			)
		}

		rl.DrawTextureV(autoTexture, gameState.car.position, gameState.car.color)

		// Draw score
		rl.DrawText(fmt.Sprintf("Score: %d", gameState.score), 20, 20, 20, rl.RayWhite)

		if gameState.pause {
			rl.DrawText(
				"GAME PAUSED",
				screenWidth/2-rl.MeasureText("GAME PAUSED", 40)/2,
				screenHeight/2-40,
				40,
				rl.Gray,
			)
		}
	} else {
		rl.DrawText(
			"A WINRAR IS NOT YOU :(\nPRESS [ENTER] TO PLAY AGAIN",
			screenWidth/2-rl.MeasureText("A WINRAR IS NOT YOU :(\nPRESS [ENTER] TO PLAY AGAIN", 20)/2,
			screenHeight/2-50,
			20,
			rl.Gray,
		)
	}

	rl.EndDrawing()
}

func InitGarages() {
	for i := int32(0); i < GARAGE_COUNT; i++ {
		gameState.garages[i].size = rl.Vector2{X: float32(SQUARE_SIZE), Y: float32(SQUARE_SIZE)}
		gameState.garages[i].position = rl.Vector2{X: float32(screenWidth - SQUARE_SIZE), Y: float32(i * SQUARE_SIZE)}
		gameState.garages[i].color = garageColors[i]
	}
}

