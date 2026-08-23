package service

import (
	"log"
	"math/rand"
	"vcontra/internal/animation"
	"vcontra/internal/models"
)

type StateManagementService struct {
	CurrentState          models.AssetState
	animations            map[string]animation.SpriteSheet
	movingStateAssets     []string // Pre-organized MOVING asset names
	idleStateAssets       []string // Pre-organized IDLE asset names
	transitionStateAssets []string // Pre-organized TRANSITION asset names
}

// NewStateManagementService loads all animations and initializes with the first one
func NewStateManagementService(configPath string, database map[string]interface{}) *StateManagementService {
	animations := make(map[string]animation.SpriteSheet)
	var movingAssets, idleAssets, transitionAssets []string
	var firstAssetName string

	// Load all animations from database and organize by state kind
	for imageName := range database {
		sheet, err := animation.LoadAnimation(configPath, imageName)
		if err != nil {
			log.Printf("Skipping asset %s due to loading error: %v", imageName, err)
			continue
		}
		animations[imageName] = sheet
		if firstAssetName == "" {
			firstAssetName = imageName
		}
		log.Printf("Loaded asset into cycle sequence: %s", imageName)

		// Organize by state kind for fast lookup later
		switch sheet.StateKind {
		case "MOVING":
			movingAssets = append(movingAssets, imageName)
		case "IDLE":
			idleAssets = append(idleAssets, imageName)
		case "TRANSITION":
			transitionAssets = append(transitionAssets, imageName)
		}
	}

	if len(animations) == 0 {
		log.Fatalf("No valid sprite sheet assets were loaded successfully.")
	}

	// Initialize with the first asset
	firstSheet := animations[firstAssetName]
	initialState := models.NewStateFromKind(firstSheet.StateKind, firstSheet.SoldierNumber, firstAssetName)

	return &StateManagementService{
		CurrentState:          initialState,
		animations:            animations,
		movingStateAssets:     movingAssets,
		idleStateAssets:       idleAssets,
		transitionStateAssets: transitionAssets,
	}
}

// GetAnimations returns the loaded animations map
func (s *StateManagementService) GetAnimations() map[string]animation.SpriteSheet {
	return s.animations
}

// GetNextState returns the next state based on current state logic
func (s *StateManagementService) GetNextState() models.AssetState {
	var assetList []string

	// Pick the right asset list based on current state kind
	switch s.CurrentState.StateKind() {
	case models.MovingState:
		assetList = s.movingStateAssets
	case models.IdleState:
		assetList = s.idleStateAssets
	case models.TransitionState:
		assetList = s.transitionStateAssets
	}

	// If no assets available for this state kind, stay in current state
	if len(assetList) == 0 {
		return s.CurrentState
	}

	// Pick a random asset from the list
	randomAssetName := assetList[rand.Intn(len(assetList))]
	sheet := s.animations[randomAssetName]
	s.CurrentState = models.NewStateFromKind(sheet.StateKind, sheet.SoldierNumber, randomAssetName)

	return s.CurrentState
}
