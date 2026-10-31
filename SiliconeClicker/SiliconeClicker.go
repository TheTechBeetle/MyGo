package main

import (
	"bytes"
	"fmt"
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
	X     int
	Y     int
	In    bool
	JSTIn bool
}

type Machine struct {
	Name  string
	Speed float32
	Image *ebiten.Image
	Cost  int
}

var state string = "Main"

var mouse mousetype

var Cookies int

type Game struct {
	face text.Face
}

var SiliconeImg = loadImage("Assets/Silicone.png")
var StatBoxImg = loadImage("Assets/StatBox.png")

func loadImage(path string) *ebiten.Image {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil
	}

	return ebiten.NewImageFromImage(img)
}

func (g *Game) Update() error {
	mouse.JSTIn = ebiten.IsMouseButtonPressed(ebiten.MouseButton0) && !mouse.In
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

var Machines [3]Machine

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("SiliconeClicker")

	Machines = [3]Machine{
		{
			"Slicer",
			1,
			loadImage("Assets/Upgrades/Slicer.png"),
			15,
		},
		{
			"Printer",
			5,
			loadImage("Assets/Upgrades/Printer.png"),
			50,
		},
		{
			"Exposer",
			20,
			loadImage("Assets/Upgrades/Exposer.png"),
			150,
		},
	}

	ebiten.RunGame(&Game{})
}

func SimulateMain() {
	if mouse.JSTIn {
		if mouse.X > 25 && mouse.X < 75 && mouse.Y > 80 && mouse.Y < 130 {
			Cookies++
		}
	}
	for i := 0; i < len(Machines); i++ {
		if mouse.X > 245 && mouse.Y > i*50 && mouse.Y < (i*50+54) {
			Cookies = Machines[i].Cost
		}
	}
}

func DrawMain(screen *ebiten.Image) {
	screen.Fill(color.RGBA{97, 87, 160, 1})
	vector.StrokeLine(screen, 100, 0, 100, 197.5, 5, color.RGBA{30, 30, 20, 1}, false)
	vector.StrokeLine(screen, 0, 200, 227.5, 200, 5, color.RGBA{30, 30, 20, 1}, false)
	vector.StrokeLine(screen, 230, 0, 230, 240, 5, color.RGBA{30, 30, 20, 1}, false)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(25, 80)
	op.GeoM.Scale(1, 1)
	screen.DrawImage(SiliconeImg, op)

	op2 := &text.DrawOptions{}
	op2.GeoM.Translate(10, 30)
	op2.GeoM.Scale(0.5, 0.5)
	text.Draw(screen, (fmt.Sprintf("%s %d", "Cookies: ", Cookies)), face, op2)

	for i := 0; i < len(Machines); i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(1.5, 1.5)
		op.GeoM.Translate(245, float64(i*50))
		screen.DrawImage(Machines[i].Image, op)

		if mouse.X > 245 && mouse.Y > i*50 && mouse.Y < (i*50+54) {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(1.5, 1.5)
			op.GeoM.Translate(160, 40)
			screen.DrawImage(StatBoxImg, op)

			op2 := &text.DrawOptions{}
			op2.GeoM.Scale(0.5, 0.5)
			op2.GeoM.Translate(1, 30)

			text.Draw(screen, (fmt.Sprintf("%s %d", "Cookies: ", Cookies)), face, op2)
		}
	}
}
