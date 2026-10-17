package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct{}

var ballx, bally, cursorx, cursory = 160, 120, 0, 0
var ballvx, ballvy float32 = 2, 2
var score int
var state bool

func (g *Game) Update() error {
	cursorx, cursory = ebiten.CursorPosition()
	if state {
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

		if bally-20 < cursory && bally+20 > cursory && ballx < 30 {
			ballvx = (0 - ballvx) + 0.5
			ballx = ballx + int(ballvx)
			score++
		}

		if ballx < 0 {
			state = false
		}
	} else {
		if cursorx > 100 && cursorx < 220 && cursory > 160 && cursory < 240 && ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
			state = true
			ballvx, ballvy = 2, 2
			ballx, bally, score = 160, 120, 0
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if state {
		vector.StrokeLine(screen, 20, float32(cursory)-20, 20, float32(cursory)+20, 5, color.White, false)
		vector.StrokeLine(screen, 300, float32(bally)-20, 300, float32(bally)+20, 5, color.White, false)

		vector.FillCircle(screen, float32(ballx), float32(bally), 5, color.White, false)
		ebitenutil.DebugPrint(screen, fmt.Sprintf("%s %d", "Score:", score))
	} else {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("%s %d", "Score:", score))
		vector.StrokeLine(screen, 100, 200, 220, 200, 40, color.White, false)
		vector.StrokeLine(screen, 105, 200, 215, 200, 30, color.Black, false)
		ebitenutil.DebugPrintAt(screen, "Play", 110, 190)
		ebitenutil.DebugPrintAt(screen, "Pongo", 120, 50)
	}
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Pongo")
	ebiten.RunGame(&Game{})
}
