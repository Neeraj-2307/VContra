package main

import (
	"encoding/json"
	"log"
	"os"
	"vcontra/internal/animation"
	"vcontra/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	configPath := "res/assets/assetConfig.json"

	// 1. Read the JSON keys dynamically so we don't have to hardcode individual load calls
	jsonFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read json config: %v", err)
	}

	// TODO: Remove the generic interface{} and replace it with a proper struct type for better type safety
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

	gameInstance := &game.Game{
		Animations:   loadedAnimations,
		Playlist:     playlist,
		CurrentAsset: 0,
		CurrentFrame: 0,
	}

	// Grab dimensions of our very first asset to initialize the window parameters
	firstSheet := gameInstance.GetActiveSheet()
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowSize(firstSheet.FrameWidth*2, firstSheet.FrameHeight*2) // Upscaled 2x
	ebiten.SetWindowTitle("vcontra engine loop alpha")

	ebiten.SetTPS(30)
	ebiten.SetVsyncEnabled(true)

	options := &ebiten.RunGameOptions{
		ScreenTransparent: true,
	}

	if err := ebiten.RunGameWithOptions(gameInstance, options); err != nil {
		log.Fatal(err)
	}
}
