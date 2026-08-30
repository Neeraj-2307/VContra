package animation

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/png" // Registers the PNG decoder
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// SpriteConfig matches the structure of our JSON entries
type SpriteConfig struct {
	TotalFrames   int    `json:"total_frames"`
	StateKind     string `json:"state_kind"`
	SoldierNumber int    `json:"soldier_number"`
}

// SpriteSheet holds the raw pixel grid and calculated dimensions.
// All field names MUST be capitalized so main.go can access them!
type SpriteSheet struct {
	EbitenImage   *ebiten.Image
	TotalFrames   int
	FrameWidth    int
	FrameHeight   int
	StateKind     string
	SoldierNumber int
}

// LoadAnimation reads the JSON config, opens the image, and computes dimensions
func LoadAnimation(jsonPath string, imageName string) (SpriteSheet, error) {
	// 1. Read and parse the JSON config file
	jsonFile, err := os.ReadFile(jsonPath)
	if err != nil {
		return SpriteSheet{}, fmt.Errorf("failed to read json: %v", err)
	}

	// Parse into a generic map string -> SpriteConfig
	var database map[string]SpriteConfig
	if err := json.Unmarshal(jsonFile, &database); err != nil {
		return SpriteSheet{}, fmt.Errorf("failed to parse json: %v", err)
	}

	// Look up our specific image settings
	config, exists := database[imageName]
	if !exists {
		return SpriteSheet{}, fmt.Errorf("image %s not found in config", imageName)
	}

	// 2. Open and decode the image asset file
	imgFile, err := os.Open(imageName)
	if err != nil {
		return SpriteSheet{}, fmt.Errorf("failed to open image: %v", err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return SpriteSheet{}, fmt.Errorf("failed to decode image: %v", err)
	}

	// 3. Do the dynamic division math
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()
	calculatedFrameWidth := imgWidth / config.TotalFrames

	// Convert standard Go image into an optimized Ebitengine image surface
	ebitenImg := ebiten.NewImageFromImage(img)

	return SpriteSheet{
		EbitenImage:   ebitenImg,
		TotalFrames:   config.TotalFrames,
		FrameWidth:    calculatedFrameWidth,
		FrameHeight:   imgHeight,
		StateKind:     config.StateKind,
		SoldierNumber: config.SoldierNumber,
	}, nil
}

// GetFrameBounds returns the exact pixel bounding box for any given frame index
func (s *SpriteSheet) GetFrameBounds(frameIndex int) image.Rectangle {
	// Simple safety loop boundary protection
	index := frameIndex % s.TotalFrames

	x0 := index * s.FrameWidth
	y0 := 0
	x1 := x0 + s.FrameWidth
	y1 := s.FrameHeight

	return image.Rect(x0, y0, x1, y1)
}
