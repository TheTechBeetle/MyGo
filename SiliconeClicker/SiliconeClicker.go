package main

import (
	"bytes"
	"fmt"
	"image/color"
	_ "image/color"
	"image/png"
	_ "math"
	"math/rand"
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
	Tier    int
}

type Particle struct {
	X       float32
	Y       float32
	XVel    float32
	YVel    float32
	Image   *ebiten.Image
	Gravity float32
	Timer   int
	Scale   float32
}

type VanType struct {
	X  int
	Y  int
	Is bool
}

var Particles []Particle

var state string = "Main"

var mouse mousetype

//var Silicone float32 = 9999999999999

var Silicone float32

var SiliconeRate float32

type Game struct {
	face text.Face
}

var SiliconeImg = loadImage("Assets/Silicone.png")
var StatBoxImg = loadImage("Assets/StatBox.png")
var VanImg = loadImage("Assets/Van.png")

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

var Bases []*ebiten.Image

var Van VanType

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("SiliconeClicker")

	Bases = []*ebiten.Image{
		loadImage("Assets/Upgrades/Tier1.png"),
		loadImage("Assets/Upgrades/Tier2.png"),
		loadImage("Assets/Upgrades/Tier3.png"),
		loadImage("Assets/Upgrades/Tier4.png"),
		loadImage("Assets/Upgrades/Tier5.png"),
	}

	Upgrades = []Upgrade{
		{
			"Glass Matrix",
			1.5,
			loadImage("Assets/Upgrades/Slicer.png"),
			"Slicer",
			50,
			0,
		},
		{
			"Better Toner",
			1.5,
			loadImage("Assets/Upgrades/Printer.png"),
			"Printer",
			300,
			0,
		},
		{
			"Better LEDs",
			1.5,
			loadImage("Assets/Upgrades/Exposer.png"),
			"Exposer",
			800,
			0,
		},
		{
			"RadioAcid",
			2,
			loadImage("Assets/Upgrades/Etcher.png"),
			"Etcher",
			1500,
			0,
		},
		{
			"New Servos",
			2,
			loadImage("Assets/Upgrades/Assembler.png"),
			"Assembler",
			5000,
			0,
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
	if int(rand.Float64()*600) == 1 {
		if int(rand.Float32()*3) == 1 {
			Van = VanType{0, int(rand.Float32() * 240), true}
		}
	}
	if Van.Is {
		Van.X++
		if Van.X > 320 {
			Van.Is = false
		}
	}
	if mouse.JSTIn {
		if mouse.X > Van.X && mouse.X < Van.X+53 && mouse.Y > Van.Y && mouse.Y < Van.Y+36 {
			Van.Is = false
			Particles = append(Particles, Particle{
				float32(Van.X), float32(Van.Y),
				0, -0.4,
				VanImg,
				0.1,
				240,
				3,
			})
			Silicone = Silicone * 2
			for i := 0; i < 15; i++ {
				Particles = append(Particles, Particle{
					rand.Float32() * 100, 0,
					0, -0.4,
					SiliconeImg,
					0.1,
					240,
					0.3,
				})
			}
		}
		if mouse.X > 25 && mouse.X < 75 && mouse.Y > 80 && mouse.Y < 130 {
			Silicone++
			Particles = append(Particles, Particle{
				rand.Float32() * 100, 0,
				0, -0.4,
				SiliconeImg,
				0.1,
				240,
				0.3,
			})
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
	for i := 0; i < len(Particles); i++ {
		var Thisparticle = &Particles[i]

		Thisparticle.X = Thisparticle.X + Thisparticle.XVel
		Thisparticle.Y = Thisparticle.Y + Thisparticle.YVel

		Thisparticle.Timer = Thisparticle.Timer - 1
		Thisparticle.YVel = Thisparticle.YVel + Thisparticle.Gravity
		if Thisparticle.Timer <= 0 || Thisparticle.Y > 300 {
			Particles = append(Particles[:i], Particles[i+1:]...)
		}
	}

	return nil
}

func DrawMain(screen *ebiten.Image) {
	screen.Fill(color.RGBA{97, 87, 160, 1})
	for i := 0; i < len(Particles); i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(Particles[i].Scale), float64(Particles[i].Scale))
		op.GeoM.Translate(float64(Particles[i].X), float64(Particles[i].Y))
		screen.DrawImage(Particles[i].Image, op)
	}
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

	if Van.Is {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(3, 3)
		op.GeoM.Translate(float64(Van.X), float64(Van.Y))
		screen.DrawImage(VanImg, op)
	}
}

func SimulateUpgrade(g *Game) error {
	if mouse.JSTIn {
		if mouse.X > 20 && mouse.X < 45 && mouse.Y > 15 && mouse.Y < 40 {
			state = "Main"
		}
		for i := 0; i < len(Upgrades); i++ {
			if Upgrades[i].Tier < 5 {
				if mouse.X > 50 && mouse.Y > (i*30)+45 && mouse.X < 70 && mouse.Y < (i*30)+65 {
					if Silicone >= float32(Upgrades[i].Cost) {
						Silicone = Silicone - float32(Upgrades[i].Cost)
						Upgrades[i].Tier++
						Upgrades[i].Cost = Upgrades[i].Cost * 3
						for i2 := 0; i2 < len(Machines); i2++ {
							if Machines[i2].Name == Upgrades[i].Machine {
								Machines[i2].Speed = Machines[i2].Speed * Upgrades[i].Power
							}
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
	op.GeoM.Scale(6, 4)
	op.GeoM.Translate(20, 20)
	screen.DrawImage(StatBoxImg, op)

	op2 := &text.DrawOptions{}
	op2.GeoM.Translate(170, 10)
	op2.GeoM.Scale(0.5, 0.5)
	text.Draw(screen, (fmt.Sprintf("%s %.2f", "Silicone: ", Silicone)), face, op2)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.5, 0.5)
	op.GeoM.Translate(20, 15)
	screen.DrawImage(StatBoxImg, op)

	op2 = &text.DrawOptions{}
	op2.GeoM.Scale(0.8, 0.6)
	op2.GeoM.Translate(23, 18)
	text.Draw(screen, "X", face, op2)

	for i := 0; i < len(Upgrades); i++ {
		if Upgrades[i].Tier < 5 {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(2, 2)
			op.GeoM.Translate(50, float64(i*30)+45)
			if int(Silicone) < Upgrades[i].Cost {
				op.ColorScale.Scale(0.5, 0.5, 0.5, 1)
			}
			screen.DrawImage(Bases[Upgrades[i].Tier], op)
			screen.DrawImage(Upgrades[i].Image, op)

			op2 = &text.DrawOptions{}
			op2.GeoM.Scale(0.5, 0.5)
			op2.GeoM.Translate(80, float64(i*30)+48)
			text.Draw(screen, fmt.Sprintf("%s %s %d %s %d", Upgrades[i].Name, " ", Upgrades[i].Tier+1, ", Cost: ", Upgrades[i].Cost), face, op2)
		}
	}
}
