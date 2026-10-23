package main

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	_ "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/thetechbeetle/mygo/uiengine"
)

var data, _ = os.ReadFile("assets/font.otf")
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
	160,
	120,
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

var PlayImg, _ = loadImage("Assets/PlayImg.png")
var QuitImg, _ = loadImage("Assets/QuitImg.png")

func loadImage(path string) (*ebiten.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(img), nil
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
	//ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Pongo")
	ebiten.RunGame(&Game{})
}

func Drawmenu(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.GeoM.Scale(1.2, 1.2)
	op.GeoM.Translate(50, 40)
	op.ColorScale.SetB(0)
	op.ColorScale.SetG(0)
	text.Draw(screen, "Pongo", face, op)
	op = &text.DrawOptions{}
	op.GeoM.Translate(20, 5)
	op.GeoM.Scale(0.3, 0.3)
	text.Draw(screen, fmt.Sprintf("%s %d", "Score:", score), face, op)
	op = &text.DrawOptions{}
	op.GeoM.Translate(20, 70)
	op.GeoM.Scale(0.3, 0.3)
	text.Draw(screen, fmt.Sprintf("%s %d", "High Score:", topscore), face, op)
	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Scale(1.5, 1.5)
	op2.GeoM.Translate(32, 163)
	screen.DrawImage(PlayImg, op2)
	op2 = &ebiten.DrawImageOptions{}
	op2.GeoM.Scale(1.5, 1.5)
	op2.GeoM.Translate(170, 163)
	screen.DrawImage(QuitImg, op2)
}

func Drawgame(screen *ebiten.Image) {
	vector.FillCircle(screen, float32(Ball.X), float32(Ball.Y), 7, color.White, false)
	vector.StrokeLine(screen, float32(PlayerPaddle.X), float32(PlayerPaddle.Y)-float32(PlayerPaddle.Size/2), float32(PlayerPaddle.X), float32(PlayerPaddle.Y)+float32(PlayerPaddle.Size/2), 5, color.White, false)
	vector.StrokeLine(screen, float32(EnemyPaddle.X), float32(EnemyPaddle.Y)-float32(EnemyPaddle.Size/2), float32(EnemyPaddle.X), float32(EnemyPaddle.Y)+float32(EnemyPaddle.Size/2), 5, color.White, false)
	if !(Ball.X < 120 && Ball.Y < 50) {
		op := &text.DrawOptions{}
		op.GeoM.Scale(0.3, 0.3)
		op.GeoM.Translate(50, 3)
		text.Draw(screen, fmt.Sprintf("%s %d", "Score:", score), face, op)
		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Scale(0.7, 0.7)
		screen.DrawImage(QuitImg, op2)
	}
}

func Simulatemenu() {
	uiengine.AddUi(&ui, 40, 180, 75, 45, 3, "Play", "PLY", false)
	uiengine.AddUi(&ui, 190, 180, 75, 45, 3, "Quit", "EXT", false)
}

func Simulategame() {
	uiengine.AddUi(&ui, 5, 20, 50, 30, 3, "End", "MNU", false)

	Ball.X = Ball.X + Ball.Xvel
	Ball.Y = Ball.Y + Ball.Yvel
	if Ball.Y < 5 || Ball.Y > 235 {
		Ball.Yvel = 0 - Ball.Yvel
	}

	PlayerPaddle.Y = mouse.Y
	EnemyPaddle.Y = int(Ball.Y) + enimofset

	PlayerPaddle.Y = min(240-PlayerPaddle.Size/2, PlayerPaddle.Y)
	EnemyPaddle.Y = min(240-EnemyPaddle.Size/2, EnemyPaddle.Y)
	PlayerPaddle.Y = max(PlayerPaddle.Size/2, PlayerPaddle.Y)
	EnemyPaddle.Y = max(EnemyPaddle.Size/2, EnemyPaddle.Y)

	if int(Ball.X) > EnemyPaddle.X-10 && int(Ball.X) < EnemyPaddle.X+10 && Ball.Xvel > 0 {
		if int(Ball.Y) > EnemyPaddle.Y-(EnemyPaddle.Size/2) && int(Ball.Y) < EnemyPaddle.Y+(EnemyPaddle.Size/2) {
			Ball.Xvel = 0 - math.Abs(Ball.Xvel*1.1)
		}
	}

	if int(Ball.X) > PlayerPaddle.X-10 && int(Ball.X) < PlayerPaddle.X+10 && Ball.Xvel < 0 {
		if int(Ball.Y) > PlayerPaddle.Y-(PlayerPaddle.Size/2) && int(Ball.Y) < PlayerPaddle.Y+(PlayerPaddle.Size/2) {
			Ball.Xvel = math.Abs(Ball.Xvel * 1.1)
			score++
			if score > 20 {
				enimofset = 0
			} else {
				enimofset = rand.Intn(80-(score*3)) - (80-(score*3))/2
			}

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
			160,
			120,
			2,
			3,
		}
		score = score + 2
		enimofset = 0
	}
}
