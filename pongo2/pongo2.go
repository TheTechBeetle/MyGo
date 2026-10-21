package main

import (
	"bytes"
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"os"

	"github.com/thetechbeetle/mygo/uiengine"
	"github.com/hajimehoshi/ebiten/v2"
	_ "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var data, _ = os.ReadFile("font.ttf")
var r = bytes.NewReader(data)
var src, _ = text.NewGoTextFaceSource(r)

var face = &text.GoTextFace{
	Source:    src,
	Size:      48,
	Direction: text.DirectionLeftToRight,
}

type mousetype struct {
	X  int
	Y  int
	In bool
}

type paddle struct {
	X    int
	Y    int
	Size int
}

type Balltype struct {
	X    float64
	Y    float64
	Xvel float64
	Yvel float64
}

var Ball = Balltype{
	50,
	50,
	2,
	3,
}

var EnemyPaddle = paddle{
	300,
	120,
	50,
}

var PlayerPaddle = paddle{
	20,
	120,
	50,
}

var state bool
var ui []uiengine.UIElement

var score, topscore, enimofset int

var mouse mousetype

type Game struct {
	face text.Face
}

func (g *Game) Update() error {
	mouse.In = ebiten.IsMouseButtonPressed(ebiten.MouseButton0)
	mouse.X, mouse.Y = ebiten.CursorPosition()
	ui = nil
	if state {
		Simulategame()
	} else {
		Simulatemenu()
	}
	switch uiengine.SimulateUi(ui) {
	case "PLY":
		{
			state = true
			score = 0
			Ball = Balltype{
				50,
				50,
				2,
				3,
			}
			enimofset = rand.Intn(70) - 35
		}

	case "EXT":
		{
			return ebiten.Termination
		}

	case "MNU":
		{
			state = false
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if state {
		Drawgame(screen)
	} else {
		Drawmenu(screen)
	}
	uiengine.DrawUi(screen, ui, face)
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Pongo")
	ebiten.RunGame(&Game{})
}

func Drawmenu(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.GeoM.Scale(0.75, 0.75)
	op.GeoM.Translate(90, 60)
	op.ColorScale.SetB(0)
	op.ColorScale.SetG(0)
	text.Draw(screen, "Pongo", face, op)
	op = &text.DrawOptions{}
	op.GeoM.Translate(20, 5)
	op.GeoM.Scale(0.25, 0.25)
	text.Draw(screen, fmt.Sprintf("%s %d", "Score:", score), face, op)
	op = &text.DrawOptions{}
	op.GeoM.Translate(20, 50)
	op.GeoM.Scale(0.25, 0.25)
	text.Draw(screen, fmt.Sprintf("%s %d", "High Score:", topscore), face, op)
}

func Drawgame(screen *ebiten.Image) {
	vector.FillCircle(screen, float32(Ball.X), float32(Ball.Y), 5, color.White, false)
	vector.StrokeLine(screen, float32(PlayerPaddle.X), float32(PlayerPaddle.Y)-float32(PlayerPaddle.Size/2), float32(PlayerPaddle.X), float32(PlayerPaddle.Y)+float32(PlayerPaddle.Size/2), 5, color.White, false)
	vector.StrokeLine(screen, float32(EnemyPaddle.X), float32(EnemyPaddle.Y)-float32(EnemyPaddle.Size/2), float32(EnemyPaddle.X), float32(EnemyPaddle.Y)+float32(EnemyPaddle.Size/2), 5, color.White, false)
	op := &text.DrawOptions{}
	op.GeoM.Scale(0.25, 0.25)
	op.GeoM.Translate(60, 10)
	text.Draw(screen, fmt.Sprintf("%s %d", "Score:", score), face, op)
}

func Simulatemenu() {
	uiengine.AddUi(&ui, 40, 180, 50, 30, 3, "Play", "PLY")
	uiengine.AddUi(&ui, 190, 180, 50, 30, 3, "Quit", "EXT")
}

func Simulategame() {
	uiengine.AddUi(&ui, 5, 20, 50, 30, 3, "End", "MNU")
	PlayerPaddle.Y = mouse.Y
	Ball.X = Ball.X + Ball.Xvel
	Ball.Y = Ball.Y + Ball.Yvel
	if Ball.Y < 5 || Ball.Y > 235 {
		Ball.Yvel = 0 - Ball.Yvel
	}

	EnemyPaddle.Y = int(Ball.Y) + enimofset

	if int(Ball.X) > EnemyPaddle.X-10 && int(Ball.X) < EnemyPaddle.X+10 && Ball.Xvel > 0 {
		if int(Ball.Y) > EnemyPaddle.Y-(EnemyPaddle.Size/2) && int(Ball.Y) < EnemyPaddle.Y+(EnemyPaddle.Size/2) {
			Ball.Xvel = 0 - math.Abs(Ball.Xvel*1.1)
		}
	}

	if int(Ball.X) > PlayerPaddle.X-10 && int(Ball.X) < PlayerPaddle.X+10 && Ball.Xvel < 0 {
		if int(Ball.Y) > PlayerPaddle.Y-(PlayerPaddle.Size/2) && int(Ball.Y) < PlayerPaddle.Y+(PlayerPaddle.Size/2) {
			Ball.Xvel = math.Abs(Ball.Xvel * 1.1)
			score++
			enimofset = rand.Intn(80-(score*3)) - (80-(score*3))/2
			if score > topscore {
				topscore = score
			}
		}
	}

	if Ball.X < 0 {
		state = false
	}
	if Ball.X > 320 {
		Ball = Balltype{
			50,
			50,
			2,
			3,
		}
		score = score + 2
		enimofset = 0
	}
}
