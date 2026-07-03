package main

import (
	"encoding/json"
	"image/color"
	"log"
	"os"
	"vcontra/animation"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	animations   map[string]animation.SpriteSheet
	playlist     []string // Keeps track of the order of animations (e.g., ["walk.png", "idle.png", "jump.png"])
	currentAsset int      // Index pointing to the active animation in our playlist
	currentFrame int
	tickCounter  int
}

// GetActiveSheet is a clean helper to grab whichever sprite sheet is playing right now
func (g *Game) GetActiveSheet() animation.SpriteSheet {
	activeName := g.playlist[g.currentAsset]
	return g.animations[activeName]
}

func (g *Game) Update() error {
	g.tickCounter++

	// Advance the frame index every 10 ticks (6 FPS)
	if g.tickCounter >= 10 {
		g.tickCounter = 0
		g.currentFrame++

		currentSheet := g.GetActiveSheet()

		// If the current animation has completed a full cycle through all its frames,
		// reset the frame count and cycle to the NEXT image asset in our JSON list!
		if g.currentFrame >= currentSheet.TotalFrames {
			g.currentFrame = 0
			g.currentAsset = (g.currentAsset + 1) % len(g.playlist)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Transparent)

	currentSheet := g.GetActiveSheet()
	subRect := currentSheet.GetFrameBounds(g.currentFrame)
	frameImage := currentSheet.EbitenImage.SubImage(subRect).(*ebiten.Image)

	opts := &ebiten.DrawImageOptions{}
	screen.DrawImage(frameImage, opts)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Dynamically adjust the canvas box size based on whatever asset is active
	currentSheet := g.GetActiveSheet()
	return currentSheet.FrameWidth, currentSheet.FrameHeight
}

func main() {
	configPath := "res/assets/assetConfig.json"

	// 1. Read the JSON keys dynamically so we don't have to hardcode individual load calls
	jsonFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read json config: %v", err)
	}

	var database map[string]interface{}
	if err := json.Unmarshal(jsonFile, &database); err != nil {
		log.Fatalf("Failed to parse json config: %v", err)
	}

	// 2. Initialize our map and playlist arrays
	loadedAnimations := make(map[string]animation.SpriteSheet)
	var playlist []string

	// Loop through every single entry found inside your JSON file and load it
	for imageName := range database {
		sheet, err := animation.LoadAnimation(configPath, imageName)
		if err != nil {
			log.Printf("Skipping asset %s due to loading error: %v", imageName, err)
			continue
		}
		loadedAnimations[imageName] = sheet
		playlist = append(playlist, imageName)
		log.Printf("Loaded asset into cycle sequence: %s", imageName)
	}

	if len(playlist) == 0 {
		log.Fatalf("No valid sprite sheet assets were loaded successfully.")
	}

	game := &Game{
		animations:   loadedAnimations,
		playlist:     playlist,
		currentAsset: 0,
		currentFrame: 0,
	}

	// Grab dimensions of our very first asset to initialize the window parameters
	firstSheet := game.GetActiveSheet()
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowSize(firstSheet.FrameWidth*2, firstSheet.FrameHeight*2) // Upscaled 2x
	ebiten.SetWindowTitle("vcontra engine loop alpha")

	ebiten.SetTPS(30)
	ebiten.SetVsyncEnabled(true)

	options := &ebiten.RunGameOptions{
		ScreenTransparent: true,
	}

	if err := ebiten.RunGameWithOptions(game, options); err != nil {
		log.Fatal(err)
	}
}
