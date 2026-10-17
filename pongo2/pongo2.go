package main

import (
	"bytes"
	_ "fmt"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	_ "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	_ "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	_ "github.com/hajimehoshi/ebiten/v2/vector"
)

var data, _ = os.ReadFile("font.ttf")
var r = bytes.NewReader(data)
var src, _ = text.NewGoTextFaceSource(r)

var face = &text.GoTextFace{
	Source:    src,
	Size:      24,
	Direction: text.DirectionLeftToRight,
}

type uielement struct {
	X      float32
	Y      float32
	Xsize  float32
	Ysize  float32
	Border float32
	Text   string
	ID     string
}

type mousetype struct {
	X  int
	Y  int
	In bool
}

var state bool
var ui []uielement

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
	SimulateUi()
	return nil
}

func DrawUi(screen *ebiten.Image) {
	for i := 0; i < len(ui); i++ {
		vector.StrokeLine(screen, ui[i].X, ui[i].Y, ui[i].X+ui[i].Xsize, ui[i].Y, ui[i].Ysize, color.White, false)
		vector.StrokeLine(screen, ui[i].X+ui[i].Border, ui[i].Y, ui[i].X+ui[i].Xsize-ui[i].Border, ui[i].Y, ui[i].Ysize-(ui[i].Border*2), color.Black, false)
		op := &text.DrawOptions{}
		op.GeoM.Scale(0.6, 0.6)
		op.GeoM.Translate(float64(ui[i].X+5), float64(ui[i].Y-(ui[i].Ysize*0.5)+5))
		text.Draw(screen, ui[i].Text, face, op)
	}
}

func SimulateUi() (Id string) {
	if mouse.In {
		for i := 0; i < len(ui); i++ {
			if mouse.X > int(ui[i].X) && mouse.X < int(ui[i].X+ui[i].Xsize) {
				if mouse.Y > int(ui[i].Y) && mouse.Y < int(ui[i].Y+ui[i].Ysize) {
					return ui[i].ID
				}
			}
		}
	}
	return ""
}

func (g *Game) Draw(screen *ebiten.Image) {
	if state {
		Drawgame(screen)
	} else {
		Drawmenu(screen)
	}
	DrawUi(screen)
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
	op.GeoM.Translate(50, 20)
	op.GeoM.Scale(1.5, 1.5)
	op.ColorScale.SetB(0)
	op.ColorScale.SetG(0)
	text.Draw(screen, "Pongo", face, op)
}

func Drawgame(screen *ebiten.Image) {

}

func Simulatemenu() {
	ui = append(ui, uielement{30, 180, 80, 30, 3, "Play", "PLY"})
	ui = append(ui, uielement{180, 180, 80, 30, 3, "Quit", "EXT"})
}

func Simulategame() {

}
