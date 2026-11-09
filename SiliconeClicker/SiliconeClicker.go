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
	count int
}

type Upgrade struct {
	Name    string
	Power   float32
	Image   *ebiten.Image
	Machine string
	Cost    int
	Got     bool
}

var state string = "Main"

var mouse mousetype

var Silicone float32

var SiliconeRate float32

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
	SiliconeRate = 0
	for i := 0; i < len(Machines); i++ {
		SiliconeRate = SiliconeRate + float32(Machines[i].count)*Machines[i].Speed
	}
	Silicone = Silicone + (SiliconeRate / 60)
	mouse.JSTIn = ebiten.IsMouseButtonPressed(ebiten.MouseButton0) && !mouse.In
	mouse.In = ebiten.IsMouseButtonPressed(ebiten.MouseButton0)
	mouse.X, mouse.Y = ebiten.CursorPosition()
	switch state {
	case "Main":
		{
			var err = SimulateMain(g)
			if err != nil {
				return err
			}
		}
	case "Upgrades":
		{
			var err = SimulateUpgrade(g)
			if err != nil {
				return err
			}
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
	case "Upgrades":
		{
			DrawUpgrade(screen)
		}
	}
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	return 320, 240
}

var Machines []Machine

var Upgrades []Upgrade

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("SiliconeClicker")

	Upgrades = []Upgrade{
		{
			"Glass Matrix",
			1.5,
			loadImage("Assets/Upgrades/Slice1.png"),
			"Slicer",
			50,
			false,
		},
		{
			"Better Toner",
			1.5,
			loadImage("Assets/Upgrades/Print1.png"),
			"Printer",
			300,
			false,
		},
		{
			"Better LEDs",
			1.5,
			loadImage("Assets/Upgrades/Expo1.png"),
			"Exposer",
			800,
			false,
		},
		{
			"Citrus Juice",
			2,
			loadImage("Assets/Upgrades/Etch1.png"),
			"Etcher",
			1500,
			false,
		},
		{
			"New Servos",
			1.5,
			loadImage("Assets/Upgrades/Assemble1.png"),
			"Assembler",
			5000,
			false,
		},
	}

	Machines = []Machine{
		{
			"Slicer",
			1,
			loadImage("Assets/Machines/Slicer.png"),
			15,
			0,
		},
		{
			"Printer",
			5,
			loadImage("Assets/Machines/Printer.png"),
			50,
			0,
		},
		{
			"Exposer",
			20,
			loadImage("Assets/Machines/Exposer.png"),
			150,
			0,
		},
		{
			"Etcher",
			60,
			loadImage("Assets/Machines/Etcher.png"),
			1000,
			0,
		},
		{
			"Assembler",
			120,
			loadImage("Assets/Machines/Assembler.png"),
			3000,
			0,
		},
	}

	ebiten.RunGame(&Game{})
}

func SimulateMain(g *Game) error {
	if mouse.JSTIn {
		if mouse.X > 25 && mouse.X < 75 && mouse.Y > 80 && mouse.Y < 130 {
			Silicone++
		}
		if mouse.X < 110 && mouse.Y > 200 {
			return ebiten.Termination
		}
		if mouse.X > 110 && mouse.Y > 200 && mouse.X < 230 {
			state = "Upgrades"
		}
		for i := 0; i < len(Machines); i++ {
			if mouse.X > 245 && mouse.Y > i*45 && mouse.Y < (i*45)+45 {
				if Silicone >= float32(Machines[i].Cost) {
					Silicone = Silicone - float32(Machines[i].Cost)
					Machines[i].count++
					Machines[i].Cost = int(float32(Machines[i].Cost) * 1.2)
				}
			}
		}
	}
	return nil
}

func DrawMain(screen *ebiten.Image) {
	screen.Fill(color.RGBA{97, 87, 160, 1})
	vector.StrokeLine(screen, 110, 240, 110, 202.5, 5, color.RGBA{30, 30, 20, 1}, false)
	vector.StrokeLine(screen, 0, 200, 227.5, 200, 5, color.RGBA{30, 30, 20, 1}, false)
	vector.StrokeLine(screen, 230, 0, 230, 240, 5, color.RGBA{30, 30, 20, 1}, false)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(25, 80)
	op.GeoM.Scale(1, 1)
	screen.DrawImage(SiliconeImg, op)

	op2 := &text.DrawOptions{}
	op2.GeoM.Translate(10, 30)
	op2.GeoM.Scale(0.5, 0.5)
	text.Draw(screen, (fmt.Sprintf("%s %.2f", "Silicone: ", Silicone)), face, op2)
	op2.GeoM.Translate(0, 30)
	text.Draw(screen, (fmt.Sprintf("%s %.2f", "Silicone/Sec: ", SiliconeRate)), face, op2)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 0.7)
	op.GeoM.Translate(5, 205)

	screen.DrawImage(StatBoxImg, op)

	op2 = &text.DrawOptions{}
	op2.GeoM.Scale(1, 0.8)
	op2.GeoM.Translate(23, 207)
	text.Draw(screen, "Quit", face, op2)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 0.7)
	op.GeoM.Translate(120, 205)
	screen.DrawImage(StatBoxImg, op)

	op2 = &text.DrawOptions{}
	op2.GeoM.Scale(0.6, 0.5)
	op2.GeoM.Translate(128, 215)
	text.Draw(screen, "Upgrades", face, op2)

	for i := 0; i < len(Machines); i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(1.5, 1.5)
		op.GeoM.Translate(245, float64(i*45)+5)
		if int(Silicone) < Machines[i].Cost {
			op.ColorScale.Scale(0.5, 0.5, 0.5, 1)
		}
		screen.DrawImage(Machines[i].Image, op)

		if mouse.X > 245 && mouse.Y > i*50 && mouse.Y < (i*50+54) {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(1.5, 1.5)
			op.GeoM.Translate(153, 40)
			screen.DrawImage(StatBoxImg, op)

			op2 := &text.DrawOptions{}
			op2.GeoM.Scale(0.3, 0.3)
			op2.GeoM.Translate(163, 50)
			text.Draw(screen, Machines[i].Name, face, op2)
			op2.GeoM.Translate(0, 15)
			text.Draw(screen, (fmt.Sprintf("%s %d", "Cost: ", Machines[i].Cost)), face, op2)
			op2.GeoM.Translate(0, 15)
			text.Draw(screen, (fmt.Sprintf("%s %.f", "Speed: ", Machines[i].Speed)), face, op2)
			op2.GeoM.Translate(0, 15)
			text.Draw(screen, (fmt.Sprintf("%s %d", "Count: ", Machines[i].count)), face, op2)
		}
	}
}

func SimulateUpgrade(g *Game) error {
	if mouse.JSTIn {
		if mouse.X > 20 && mouse.X < 45 && mouse.Y > 15 && mouse.Y < 40 {
			state = "Main"
		}
		for i := 0; i < len(Upgrades); i++ {
			if mouse.X > 50 && mouse.Y > (i*30)+45 && mouse.X < 70 && mouse.Y < (i*30)+65 {
				if Silicone >= float32(Upgrades[i].Cost) {
					Silicone = Silicone - float32(Upgrades[i].Cost)
					Upgrades[i].Got = true
					for i2 := 0; i2 < len(Machines); i2++ {
						if Machines[i2].Name == Upgrades[i].Machine {
							Machines[i2].Speed = Machines[i2].Speed * Upgrades[i].Power
						}
					}
				}
			}
		}
	}
	return nil
}

func DrawUpgrade(screen *ebiten.Image) {
	screen.Fill(color.RGBA{97, 87, 160, 1})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(5.5, 4)
	op.GeoM.Translate(20, 20)
	screen.DrawImage(StatBoxImg, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.5, 0.5)
	op.GeoM.Translate(20, 15)
	screen.DrawImage(StatBoxImg, op)

	op2 := &text.DrawOptions{}
	op2.GeoM.Scale(0.8, 0.6)
	op2.GeoM.Translate(23, 18)
	text.Draw(screen, "X", face, op2)

	for i := 0; i < len(Upgrades); i++ {
		if !Upgrades[i].Got {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(2, 2)
			op.GeoM.Translate(50, float64(i*30)+45)
			if int(Silicone) < Upgrades[i].Cost {
				op.ColorScale.Scale(0.5, 0.5, 0.5, 1)
			}
			screen.DrawImage(Upgrades[i].Image, op)

			op2 = &text.DrawOptions{}
			op2.GeoM.Scale(0.5, 0.5)
			op2.GeoM.Translate(80, float64(i*30)+48)
			text.Draw(screen, fmt.Sprintf("%s %s %d", Upgrades[i].Name, ", Cost: ", Upgrades[i].Cost), face, op2)
		}
	}
}
