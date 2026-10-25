package main

import (
	"bytes"
	_ "fmt"
	"image/color"
	_ "image/color"
	"image/png"
	_ "math"
	_ "math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	_ "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	_ "github.com/hajimehoshi/ebiten/v2/vector"
)

var data, _ = os.ReadFile("assets/font.ttf")
var r = bytes.NewReader(data)
var src, _ = text.NewGoTextFaceSource(r)

var face = &text.GoTextFace{
	Source:    src,
	Size:      32,
	Direction: text.DirectionLeftToRight,
}

type mousetype struct {
	X  int
	Y  int
	In bool
}

var state string = "Main"

var mouse mousetype

type Game struct {
	face text.Face
}

var SilaconeImg, _ = loadImage("Assets/Silicone.png")

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
	switch state {
	case "Main":
		{
			SimulateMain()
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch state {
	case "Main":
		{
			DrawMain(screen)
		}
	}
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("SiliconeClicker")
	ebiten.RunGame(&Game{})
}

func SimulateMain() {

}

func DrawMain(screen *ebiten.Image) {
	screen.Fill(color.RGBA{97, 87, 160, 1})
	vector.StrokeLine(screen, 100, 0, 100, 220, 5, color.RGBA{30, 30, 20, 1}, false)
	vector.StrokeLine(screen, 0, 200, 320, 200, 5, color.RGBA{30, 30, 20, 1}, false)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(25, 80)
	op.GeoM.Scale(1, 1)
	screen.DrawImage(SilaconeImg, op)
}
