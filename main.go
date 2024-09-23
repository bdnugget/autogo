package main

import (
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth      = 1300
	screenHeight     = 800
	numLanes         = 5
	laneHeight       = screenHeight / numLanes
	carTextureWidth  = 200
	carTextureHeight = 200
	garageWidth      = 200
	garageHeight     = 200
	scorePosx        = 10
	scorePosy        = 10
	scoreFontSize    = 32
	gameOverPosx     = screenWidth/2 - 150
	gameOverPosy     = screenHeight/2 - 20
	gameOverFontSize = 32
)

var (
	garageColors   = []rl.Color{rl.Red, rl.Blue, rl.Green, rl.Yellow, rl.Purple}
	autoHappy      rl.Sound
	autoSad        rl.Sound
	autoTexture    rl.Texture2D
	garageTexture  rl.Texture2D
	asphaltTexture rl.Texture2D
	highScore      int
)

type Car struct {
	texture rl.Texture2D
	lane    int
	color   rl.Color
	xPos    int
	speed   int
}

func NewCar() Car {
	return Car{
		texture: autoTexture,
		lane:    rand.Intn(numLanes),
		color:   randomColor(),
		xPos:    25,
		speed:   1,
	}
}

type Garage struct {
	lane  int
	color rl.Color
}

type Game struct {
	car      Car
	garages  [numLanes]Garage
	score    int
	gameOver bool
	debugMsg string
}

func NewGame() Game {
	game := Game{
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

func (g Game) draw() {
	for i := 0; i <= screenHeight/laneHeight; i++ {
		for j := 0; j < screenWidth/int(asphaltTexture.Width); j++ {
			rl.DrawTexturePro(asphaltTexture, rl.NewRectangle(0, 0, float32(asphaltTexture.Width), float32(asphaltTexture.Height)), rl.NewRectangle(float32(j)*float32(asphaltTexture.Width), float32(i*laneHeight-laneHeight/2), float32(asphaltTexture.Width), float32(asphaltTexture.Height)), rl.Vector2{}, 0, rl.White)
		}
	}

	// Draw garages
	for _, garage := range g.garages {
		rl.DrawTexture(garageTexture, int32(screenWidth-garageWidth), int32(garage.lane*laneHeight), garage.color)
	}

	// Draw car
	rl.DrawTexture(g.car.texture, int32(g.car.xPos), int32(g.car.lane*laneHeight), g.car.color)

	rl.DrawText(fmt.Sprintf("Score: %d\nHigh Score: %d", g.score, highScore), scorePosx, scorePosy, scoreFontSize, rl.DarkGreen)
	if g.gameOver {
		rl.DrawText("You're not an WinRAR :(\nPress Enter to Restart", gameOverPosx, gameOverPosy, gameOverFontSize, rl.Red)
	}
}

func (g *Game) update() {
	if rl.IsKeyPressed(rl.KeyF) {
		rl.ToggleFullscreen()
	}

	if g.gameOver {
		if rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeyKpEnter) || rl.IsGestureDetected(rl.GestureTap) {
			*g = NewGame()
		}
		return
	}

	g.car.xPos += g.car.speed

	// Check if the car has reached the garages
	if int32(g.car.xPos)+g.car.texture.Width >= screenWidth-garageWidth {
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

	g.car.speed = 3 + g.score/5

	// Move car
	if rl.IsGestureDetected(rl.GestureTap) {
		touchY := int(rl.GetTouchY())
		lane := touchY / (screenHeight / numLanes)
		g.debugMsg = fmt.Sprintf("Y: %d, Lane: %d", touchY, lane)
		g.car.lane = lane
	}

	if rl.IsKeyPressed(rl.KeyUp) && g.car.lane > 0 {
		g.car.lane--
	} else if (rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyRight)) && g.car.lane < numLanes-1 {
		g.car.lane++
	}

	rl.DrawText(g.debugMsg, 20, 100, 32, rl.Red)
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
	asphaltTexture = rl.LoadTexture("resources/asphalt.png")
	defer rl.UnloadTexture(autoTexture)
	defer rl.UnloadTexture(garageTexture)
	defer rl.UnloadTexture(asphaltTexture)
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
