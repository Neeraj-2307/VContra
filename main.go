package main

import (
	"encoding/json"
	"log"
	"os"
	"vcontra/internal/game"
	"vcontra/internal/service"

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

	// Initialize state management (which loads all animations and sets up initial state)
	stateManager := service.NewStateManagementService(configPath, database)

	gameInstance := &game.Game{
		Animations:       stateManager.GetAnimations(),
		CurrentAssetInfo: stateManager.CurrentState.CurrentAssetName(),
		CurrentFrame:     0,
		StateManager:     stateManager,
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
