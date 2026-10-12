package main

import (
    "image/color"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct{}

var ballx,bally = 160,120
var ballvx,ballvy float32 = 2,2

func (g *Game) Update() error {
	_ , cursory := ebiten.CursorPosition()
    ballx = ballx + int(ballvx)
	bally = bally + int(ballvy)

	if bally > 230 || bally < 10 {
		ballvy = 0 - ballvy
		bally = bally + int(ballvy)
	}

	if ballx > 290 {
		ballvx = 0 - ballvx
		ballx = ballx + int(ballvx)
	}

	if bally - 20 < cursory && bally + 20 > cursory && ballx < 30 {
		ballvx = (0 - ballvx) + 0.5
		ballx = ballx + int(ballvx)
	}

	if ballx < 0 {
		return ebiten.Termination
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	_ , cursory := ebiten.CursorPosition()
    // Two lines
    vector.StrokeLine(screen, 20, float32(cursory) -20, 20, float32(cursory) + 20, 5, color.White, false)
    vector.StrokeLine(screen,300, float32(bally) -20, 300, float32(bally) + 20, 5, color.White, false)

    // A dot (small filled circle)
    vector.DrawFilledCircle(screen, float32(ballx), float32(bally), 5, color.White, false)
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
    return 320, 240
}

func main() {
    ebiten.SetWindowSize(640, 480)
    ebiten.SetWindowTitle("Pongo")
    ebiten.RunGame(&Game{})
}