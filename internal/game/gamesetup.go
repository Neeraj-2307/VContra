package game

import (
	"image/color"
	"vcontra/internal/animation"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Animations   map[string]animation.SpriteSheet
	Playlist     []string // Keeps track of the order of animations (e.g., ["walk.png", "idle.png", "jump.png"])
	CurrentAsset int      // Index pointing to the active animation in our playlist
	CurrentFrame int
	TickCounter  int
}

// GetActiveSheet is a clean helper to grab whichever sprite sheet is playing right now
func (g *Game) GetActiveSheet() animation.SpriteSheet {
	activeName := g.Playlist[g.CurrentAsset]
	return g.Animations[activeName]
}

func (g *Game) Update() error {
	g.TickCounter++

	// Advance the frame index every 10 ticks (6 FPS)
	if g.TickCounter >= 10 {
		g.TickCounter = 0
		g.CurrentFrame++

		currentSheet := g.GetActiveSheet()

		// If the current animation has completed a full cycle through all its frames,
		// reset the frame count and cycle to the NEXT image asset in our JSON list!
		if g.CurrentFrame >= currentSheet.TotalFrames {
			g.CurrentFrame = 0
			g.CurrentAsset = (g.CurrentAsset + 1) % len(g.Playlist)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Transparent)

	currentSheet := g.GetActiveSheet()
	subRect := currentSheet.GetFrameBounds(g.CurrentFrame)
	frameImage := currentSheet.EbitenImage.SubImage(subRect).(*ebiten.Image)

	opts := &ebiten.DrawImageOptions{}
	screen.DrawImage(frameImage, opts)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Dynamically adjust the canvas box size based on whatever asset is active
	currentSheet := g.GetActiveSheet()
	return currentSheet.FrameWidth, currentSheet.FrameHeight
}
