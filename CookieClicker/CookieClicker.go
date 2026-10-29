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
	Size:      64,
	Direction: text.DirectionLeftToRight,
}

type mousetype struct {
	X  int
	Y  int
	In bool
}

var mouse mousetype
var State string = "Main"
var CookieImg, _ = loadImage("Assets/Cookie.png")

type Game struct {
	face text.Face
}

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
	switch State {
	case "Main":
		{
			SimulateMain()
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch State {
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
	ebiten.SetWindowTitle("Cookie Clicker")
	ebiten.RunGame(&Game{})
}

func DrawMain(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 100, 0})
	for i := 0; i < 16; i++ {
		vector.StrokeLine(screen, float32(i*20), 0, float32(i*20), 240, 10, color.RGBA{0, 0, 100, 0}, false)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(50, 90)
	screen.DrawImage(CookieImg, op)
}

func SimulateMain() {

}
