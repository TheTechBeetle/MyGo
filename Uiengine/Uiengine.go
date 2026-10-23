package uiengine

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type mousetype struct {
	X  int
	Y  int
	In bool
}

type UIElement struct {
	X      float32
	Y      float32
	Xsize  float32
	Ysize  float32
	Border float32
	Text   string
	ID     string
	Disp   bool
}

var mouse mousetype

func DrawUi(screen *ebiten.Image, ui []UIElement, face *text.GoTextFace) {
	mouse.In = ebiten.IsMouseButtonPressed(ebiten.MouseButton0)
	mouse.X, mouse.Y = ebiten.CursorPosition()
	for i := 0; i < len(ui); i++ {
		if ui[i].Disp {
			vector.StrokeLine(screen, ui[i].X, ui[i].Y, ui[i].X+ui[i].Xsize, ui[i].Y, ui[i].Ysize, color.White, false)
			vector.StrokeLine(screen, ui[i].X+ui[i].Border, ui[i].Y, ui[i].X+ui[i].Xsize-ui[i].Border, ui[i].Y, ui[i].Ysize-(ui[i].Border*2), color.Black, false)
			op := &text.DrawOptions{}
			op.GeoM.Scale(0.3, 0.3)
			op.GeoM.Translate(float64(ui[i].X+5), float64(ui[i].Y-(ui[i].Ysize*0.5)+5))
			text.Draw(screen, ui[i].Text, face, op)
		}
	}
}

func SimulateUi(ui []UIElement) (Id string) {
	mouse.In = ebiten.IsMouseButtonPressed(ebiten.MouseButton0)
	mouse.X, mouse.Y = ebiten.CursorPosition()
	if mouse.In {
		for i := 0; i < len(ui); i++ {
			if mouse.X > int(ui[i].X) && mouse.X < int(ui[i].X+ui[i].Xsize) {
				if mouse.Y > int(ui[i].Y-ui[i].Ysize) && mouse.Y < int(ui[i].Y+ui[i].Ysize) {
					return ui[i].ID
				}
			}
		}
	}
	return ""
}

func AddUi(
	UI *[]UIElement,
	X float32,
	Y float32,
	Xsize float32,
	Ysize float32,
	Border float32,
	Text string,
	ID string,
	Disp bool,
) {
	*UI = append(*UI, UIElement{X, Y, Xsize, Ysize, Border, Text, ID, Disp})
}
