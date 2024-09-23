package main

import (
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1600
	screenHeight = 800
	numLanes     = 5
	numGarages   = 5
	carSpeed     = 5
	laneHeight   = screenHeight / numLanes
)

var (
	garageColors  = [numGarages]rl.Color{rl.Red, rl.Blue, rl.Green, rl.Yellow, rl.Purple}
	autoHappy     rl.Sound
	autoSad       rl.Sound
	autoTexture   rl.Texture2D
	garageTexture rl.Texture2D
	highScore     int
)

type Car struct {
	texture rl.Texture2D
	lane    int
	color   rl.Color
	xPos    int
}

type Garage struct {
	lane  int
	color rl.Color
}

type Game struct {
	car      Car
	garages  [numGarages]Garage
	score    int
	gameOver bool
}

func NewCar() Car {
	return Car{
		texture: autoTexture,
		lane:    rand.Intn(numLanes),
		color:   randomColor(),
		xPos:    25,
	}
}

func NewGame() *Game {
	game := &Game{
		car:      NewCar(),
		score:    0,
		gameOver: false,
	}

	// Initialize garages
	for i := range game.garages {
		game.garages[i] = Garage{
			lane:  i,
			color: garageColors[i],
		}
	}

	return game
}

func randomColor() rl.Color {
	return garageColors[rand.Intn(len(garageColors))]
}

func (g *Game) draw() {
	// Draw stripes on the road
	for i := 0; i < numLanes; i++ {
		rl.DrawRectangle(0, int32(i*laneHeight), screenWidth, 2, rl.White)
	}

	// Draw garages
	for _, garage := range g.garages {
		rl.DrawRectangle(int32(screenWidth-150), int32(garage.lane*laneHeight), 150, int32(laneHeight), garage.color)
		rl.DrawTexture(garageTexture, int32(screenWidth-150), int32(garage.lane*laneHeight), garage.color)
	}

	// Draw car
	rl.DrawTexture(g.car.texture, int32(g.car.xPos), int32(g.car.lane*laneHeight), g.car.color)

	rl.DrawText(fmt.Sprintf("Score: %d\nHigh Score: %d", g.score, highScore), 10, 10, 32, rl.DarkGreen)
	if g.gameOver {
		rl.DrawText("You're not an WinRAR :(\nPress Enter to Restart", int32(screenWidth/2-150), int32(screenHeight/2-20), 32, rl.Red)
	}
}

func (g *Game) update() {
	if g.gameOver {
		if rl.IsKeyPressed(rl.KeyEnter) {
			*g = *NewGame()
		}
		return
	}

	g.car.xPos += carSpeed

	// Check if the car has reached the garages
	if int32(g.car.xPos)+g.car.texture.Width >= screenWidth-100 {
		garage := g.garages[g.car.lane]
		if g.car.color == garage.color {
			rl.PlaySound(autoHappy)
			g.score++
			if g.score > highScore {
				highScore = g.score
			}
			g.car = NewCar()
		} else {
			rl.PlaySound(autoSad)
			g.gameOver = true
		}
	}

	// Move car
	if rl.IsKeyPressed(rl.KeyUp) && g.car.lane > 0 {
		g.car.lane--
	} else if (rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyRight)) && g.car.lane < numLanes-1 {
		g.car.lane++
	}
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Dikke Vette Cargame voor Milo")
	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	// Load assets
	autoHappy = rl.LoadSound("resources/auto_happy_vob.ogg")
	autoSad = rl.LoadSound("resources/auto_sad_vob.ogg")
	autoTexture = rl.LoadTexture("resources/car_200px.png")
	garageTexture = rl.LoadTexture("resources/garage_200px.png")
	defer rl.UnloadTexture(autoTexture)
	defer rl.UnloadSound(autoHappy)
	defer rl.UnloadSound(autoSad)

	game := NewGame()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkGray)

		game.update()
		game.draw()

		rl.EndDrawing()
	}
}
