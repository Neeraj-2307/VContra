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
	movingStateAssets     [][]string // Pre-organized MOVING asset names by soldier number
	idleStateAssets       [][]string // Pre-organized IDLE asset names by soldier number
	transitionStateAssets [][]string // Pre-organized TRANSITION asset names by soldier number
}

// NewStateManagementService loads all animations and initializes with the first one
func NewStateManagementService(configPath string, database map[string]interface{}) *StateManagementService {
	animations := make(map[string]animation.SpriteSheet)
	var movingAssets, idleAssets, transitionAssets [][]string
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

		// Organize by state kind for fast lookup later, indexed by soldier number
		soldierNum := sheet.SoldierNumber
		// Ensure the slice is large enough for this soldier number
		for soldierNum >= len(movingAssets) {
			movingAssets = append(movingAssets, []string{})
			idleAssets = append(idleAssets, []string{})
			transitionAssets = append(transitionAssets, []string{})
		}

		switch sheet.StateKind {
		case "MOVING":
			movingAssets[soldierNum] = append(movingAssets[soldierNum], imageName)
		case "IDLE":
			idleAssets[soldierNum] = append(idleAssets[soldierNum], imageName)
		case "TRANSITION":
			transitionAssets[soldierNum] = append(transitionAssets[soldierNum], imageName)
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
	soldierNum := s.CurrentState.SoldierNumber()

	// Pick the right asset list based on current state kind and soldier number
	switch s.CurrentState.StateKind() {
	case models.MovingState:
		// While moving, stay in moving/idle states of the same soldier
		if soldierNum < len(s.movingStateAssets) {
			assetList = s.movingStateAssets[soldierNum]
			assetList = append(assetList, s.idleStateAssets[soldierNum]...)
		}
	case models.IdleState:
		// While idle, can go to moving/idle/transition of same soldier
		if soldierNum < len(s.idleStateAssets) {
			assetList = s.idleStateAssets[soldierNum]
			assetList = append(assetList, s.movingStateAssets[soldierNum]...)
			assetList = append(assetList, s.transitionStateAssets[soldierNum]...)
		}
	case models.TransitionState:
		// Load MOVING and IDLE assets for all soldiers except the current one
		// This transitions away from the current soldier to another one
		for i := 0; i < len(s.movingStateAssets); i++ {
			if i != soldierNum {
				assetList = append(assetList, s.movingStateAssets[i]...)
			}
		}
		for i := 0; i < len(s.idleStateAssets); i++ {
			if i != soldierNum {
				assetList = append(assetList, s.idleStateAssets[i]...)
			}
		}
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
