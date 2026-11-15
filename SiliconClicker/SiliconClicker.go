package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	_ "image/color"
	"image/png"
	_ "math"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	_ "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var data, _ = os.ReadFile("Assets/Font.ttf")
var r = bytes.NewReader(data)
var src, _ = text.NewGoTextFaceSource(r)

var audioContext = audio.NewContext(48000)

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
	Count int
	name2 string
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
	Text    string
}

type VanType struct {
	X  int
	Y  int
	Is bool
}

type SaveData struct {
	Silicon  float32
	Upgrades []Upgrade
	Machines []Machine
}

type options struct {
	X        float64
	Y        float64
	Xsize    float64
	Ysize    float64
	Text     string
	Image    *ebiten.Image
	TextSize float64
	Tint     bool
}

var Particles []Particle

var state string = "Main"

var mouse mousetype

//var Silicon float32 = 9999999999999

var Silicon float32

var SiliconRate float32

var SiliconClicked int

var MousePower float32

var ScrollState int

type Game struct {
	face text.Face
}

var SiliconImg = loadImage("Assets/Silicon.png")
var StatBoxImg = loadImage("Assets/StatBox.png")
var VanImg = loadImage("Assets/Van.png")

var SiliconWav *audio.Player
var VanWav *audio.Player
var ClickWav *audio.Player

func loadWav(path string) (*audio.Player, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	d, err := wav.Decode(audioContext, f)
	if err != nil {
		return nil, err
	}
	return audio.NewPlayer(audioContext, d)
}

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
	var _, Scroll = ebiten.Wheel()
	ScrollState += int(Scroll * 10)
	ScrollState = min(ScrollState, 0)

	SiliconRate = 0
	for i := 0; i < len(Machines); i++ {
		SiliconRate = SiliconRate + float32(Machines[i].Count)*Machines[i].Speed
	}
	Silicon = Silicon + (SiliconRate / 60)
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

func ButtonOptions() (op options) {
	return options{
		Xsize:    100,
		Ysize:    100,
		TextSize: 1,
		Tint:     false,
	}
}

func DrawButton(op options, screen *ebiten.Image) {
	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Scale(op.Xsize/float64(op.Image.Bounds().Size().X), op.Ysize/float64(op.Image.Bounds().Size().Y))
	op2.GeoM.Translate(op.X, op.Y)
	if op.Tint {
		op2.ColorScale.Set(0.5, 0.5, 0.5, 1)
	}
	screen.DrawImage(op.Image, op2)

	op3 := &text.DrawOptions{}
	op3.GeoM.Scale(op.TextSize, op.TextSize)
	op3.GeoM.Translate(op.X, op.Y)
	op3.GeoM.Translate(10, 5)
	text.Draw(screen, op.Text, face, op3)
}

func TestButton(op options) (clicked bool) {
	return (mouse.X > int(op.X) &&
		mouse.X < int(op.X+(op.Xsize)) &&
		mouse.Y > int(op.Y) &&
		mouse.Y < int(op.Y+(op.Ysize)) &&
		(mouse.JSTIn || mouse.In && (ebiten.IsKeyPressed(ebiten.KeyShift))))
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("SiliconClicker")

	SiliconWav, _ = loadWav("Assets/Silicon.wav")
	VanWav, _ = loadWav("Assets/Van.wav")
	ClickWav, _ = loadWav("Assets/Click.wav")

	Bases = []*ebiten.Image{
		loadImage("Assets/Upgrades/Tier1.png"),
		loadImage("Assets/Upgrades/Tier2.png"),
		loadImage("Assets/Upgrades/Tier3.png"),
		loadImage("Assets/Upgrades/Tier4.png"),
		loadImage("Assets/Upgrades/Tier5.png"),
		loadImage("Assets/Upgrades/MAX.png"),
	}

	Upgrades = []Upgrade{
		{
			"NPN Cursor",
			2,
			loadImage("Assets/Upgrades/Mouse.png"),
			"Mouse",
			150,
			0,
		},
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
		{
			"Thermo Crystals",
			2,
			loadImage("Assets/Upgrades/Wafer.png"),
			"Wafer Fab",
			20000,
			0,
		},
		{
			"Intagrated Casing",
			2,
			loadImage("Assets/Upgrades/ChipLine.png"),
			"Chip Packing",
			100000,
			0,
		},
		{
			"Nuclear Forge",
			2,
			loadImage("Assets/Upgrades/Foundry.png"),
			"Automated",
			2000000,
			0,
		},
		{
			"Entanglement",
			2,
			loadImage("Assets/Upgrades/Quantum.png"),
			"Quantum",
			10000000,
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
			"",
		},
		{
			"Printer",
			5,
			loadImage("Assets/Machines/Printer.png"),
			50,
			0,
			"",
		},
		{
			"Exposer",
			20,
			loadImage("Assets/Machines/Exposer.png"),
			150,
			0,
			"",
		},
		{
			"Etcher",
			60,
			loadImage("Assets/Machines/Etcher.png"),
			1000,
			0,
			"",
		},
		{
			"Assembler",
			120,
			loadImage("Assets/Machines/Assembler.png"),
			3000,
			0,
			"",
		},
		{
			"Wafer Fab",
			250,
			loadImage("Assets/Machines/Wafer.png"),
			12000,
			0,
			"",
		},
		{
			"Chip Packing",
			750,
			loadImage("Assets/Machines/ChipLine.png"),
			60000,
			0,
			"Line",
		},
		{
			"Automated",
			2250,
			loadImage("Assets/Machines/Foundry.png"),
			300000,
			0,
			"Foundry",
		},
		{
			"Quantum",
			6750,
			loadImage("Assets/Machines/Quantum.png"),
			1500000,
			0,
			"Fabricator",
		},
	}

	var Save, err = LoadGame()
	if err == nil {
		Silicon = Save.Silicon
		for i := range Machines {
			for _, savedMachine := range Save.Machines {
				if Machines[i].Name == savedMachine.Name {
					Machines[i].Count = savedMachine.Count
					Machines[i].Cost = savedMachine.Cost
					Machines[i].Speed = savedMachine.Speed
					break
				}
			}
		}
		for i := range Upgrades {
			for _, savedUpgrade := range Save.Upgrades {
				if Upgrades[i].Name == savedUpgrade.Name {
					Upgrades[i].Tier = savedUpgrade.Tier
					Upgrades[i].Cost = savedUpgrade.Cost
					break
				}
			}
		}
	}

	ebiten.RunGame(&Game{})
}

func SaveGame(state SaveData) error {
	b, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("save.json", b, 0644)
}

func LoadGame() (SaveData, error) {
	var s SaveData
	b, err := os.ReadFile("save.json")
	if err != nil {
		return s, err
	}
	if b == nil {
		return s, errors.New("NOFILE")
	}
	err = json.Unmarshal(b, &s)
	return s, err
}

func large(in float32) (out string) {
	if in > 1000000000 {
		out = fmt.Sprintf("%dB", int(in/1000000000))
	} else if in > 1000000 {
		out = fmt.Sprintf("%dM", int(in/1000000))
	} else if in > 1000 {
		out = fmt.Sprintf("%dK", int(in/1000))
	} else {
		out = fmt.Sprintf("%d", int(in))
	}
	return out
}

func SimulateMain(g *Game) error {

	// hack
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) &&
		ebiten.IsKeyPressed(ebiten.KeyArrowUp) &&
		ebiten.IsKeyPressed(ebiten.KeyArrowLeft) &&
		ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		Silicon = Silicon + float32(Machines[len(Machines)-1].Cost)
	}

	// Silicon animation progresser
	SiliconClicked = SiliconClicked - 1

	// Spawn Van
	if int(rand.Float64()*600) == 1 {
		if int(rand.Float32()*3) == 1 {
			Van = VanType{0, int(rand.Float32() * 160), true}
		}
	}

	// Simulate Van
	if Van.Is {
		Van.X++
		if Van.X > 320 {
			Van.Is = false
		}
	}

	// Silicon
	op3 := ButtonOptions()
	op3.X = 30
	op3.Y = 80
	op3.Xsize = 50
	op3.Ysize = 50
	if TestButton(op3) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		Silicon += float32((Upgrades[0].Tier + 1) * (Upgrades[0].Tier + 1))
		SiliconWav.SetPosition(0)
		SiliconWav.Play()
		SiliconClicked = 5
		Particles = append(Particles, Particle{
			X:       rand.Float32() * 100,
			Y:       0,
			XVel:    0,
			YVel:    -0.4,
			Image:   SiliconImg,
			Gravity: 0.1,
			Timer:   240,
			Scale:   0.3,
		})
		Particles = append(Particles, Particle{
			X:       float32(mouse.X),
			Y:       float32(mouse.Y + -10),
			YVel:    -1.2,
			Gravity: 0.02,
			Timer:   60,
			Scale:   0.5,
			Text:    fmt.Sprintf("+%d", (Upgrades[0].Tier+1)*(Upgrades[0].Tier+1)),
		})
	}

	// Quit
	op3 = ButtonOptions()
	op3.X = 10
	op3.Y = 205
	op3.Xsize = 85
	op3.Ysize = 33

	if TestButton(op3) {
		SaveGame(SaveData{
			Silicon,
			Upgrades,
			Machines,
		})
		return ebiten.Termination
	}

	// Upgrades
	op3 = ButtonOptions()
	op3.X = 120
	op3.Y = 205
	op3.Xsize = 90
	op3.Ysize = 33

	if TestButton(op3) {
		state = "Upgrades"
		ScrollState = 0
		ClickWav.SetPosition(0)
		ClickWav.Play()
	}

	// Van collision
	if Van.Is {
		op3 = ButtonOptions()
		op3.X = float64(Van.X)
		op3.Y = float64(Van.Y)
		op3.Xsize = 54
		op3.Ysize = 36

		if TestButton(op3) {
			Van.Is = false
			VanWav.SetPosition(0)
			VanWav.Play()
			Particles = append(Particles, Particle{
				X:       float32(Van.X),
				Y:       float32(Van.Y),
				XVel:    0,
				YVel:    -0.4,
				Image:   VanImg,
				Gravity: 0.1,
				Timer:   240,
				Scale:   3,
			})
			Silicon = Silicon * 2
			for i := 0; i < 15; i++ {
				Particles = append(Particles, Particle{
					X:       rand.Float32() * 100,
					Y:       0,
					XVel:    0,
					YVel:    -0.4,
					Image:   SiliconImg,
					Gravity: 0.1,
					Timer:   240,
					Scale:   0.3,
				})
			}
		}
	}

	// Machines
	for i := 0; i < len(Machines); i++ {
		op3 = ButtonOptions()
		op3.X = 245
		op3.Y = float64((i*45)+ScrollState) + 5
		op3.Xsize = 70
		op3.Ysize = 45

		if TestButton(op3) {
			if Silicon >= float32(Machines[i].Cost) {
				if Machines[i].Name == "Assembler" && Machines[i].Count == 0 {
					VanWav.SetPosition(0)
					VanWav.Play()
				} else {
					ClickWav.SetPosition(0)
					ClickWav.Play()
				}
				Silicon = Silicon - float32(Machines[i].Cost)
				Machines[i].Count++
				Machines[i].Cost = int(float32(Machines[i].Cost) * 1.2)
			}
		}
	}

	// Particles
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

func drawParticles(screen *ebiten.Image) {
	for i := 0; i < len(Particles); i++ {
		p := Particles[i]
		if p.Text != "" {
			textOp := &text.DrawOptions{}
			textOp.GeoM.Scale(float64(p.Scale), float64(p.Scale))
			textOp.GeoM.Translate(float64(p.X), float64(p.Y))
			text.Draw(screen, p.Text, face, textOp)
		} else if p.Image != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(float64(p.Scale), float64(p.Scale))
			op.GeoM.Translate(float64(p.X), float64(p.Y))
			screen.DrawImage(p.Image, op)
		}
	}
}

func DrawMain(screen *ebiten.Image) {
	screen.Fill(color.RGBA{97, 87, 160, 1})
	vector.StrokeLine(screen, 110, 240, 110, 202.5, 5, color.RGBA{30, 30, 20, 1}, false)
	vector.StrokeLine(screen, 0, 200, 227.5, 200, 5, color.RGBA{30, 30, 20, 1}, false)
	vector.StrokeLine(screen, 230, 0, 230, 240, 5, color.RGBA{30, 30, 20, 1}, false)

	// Silicon
	op3 := ButtonOptions()
	op3.X = 30
	op3.Y = 80
	op3.Xsize = 50
	op3.Ysize = 50
	op3.Image = SiliconImg
	if SiliconClicked > 0 {
		op3.X = 30 - (float64(SiliconClicked))
		op3.Y = 80 - (float64(SiliconClicked))
		op3.Xsize = 50 + (float64(SiliconClicked) * 2)
		op3.Ysize = 50 + (float64(SiliconClicked) * 2)
	}

	DrawButton(op3, screen)
	drawParticles(screen)

	// Silicon Stats
	op2 := &text.DrawOptions{}
	op2.GeoM.Translate(10, 30)
	op2.GeoM.Scale(0.5, 0.5)
	text.Draw(screen, (fmt.Sprintf("%s %s", "Silicon: ", large(Silicon))), face, op2)
	op2.GeoM.Translate(0, 30)
	text.Draw(screen, (fmt.Sprintf("%s %s", "Silicon/Sec: ", large(SiliconRate))), face, op2)

	// Quit
	op3 = ButtonOptions()
	op3.X = 10
	op3.Y = 205
	op3.Xsize = 85
	op3.Ysize = 33
	op3.Image = StatBoxImg
	op3.Text = "Quit"
	op3.TextSize = 0.7
	DrawButton(op3, screen)

	// Upgrades

	op3 = ButtonOptions()
	op3.X = 120
	op3.Y = 205
	op3.Xsize = 90
	op3.Ysize = 33
	op3.Image = StatBoxImg
	op3.Text = "Upgrades"
	op3.TextSize = 0.5
	DrawButton(op3, screen)

	// Machines
	for i := 0; i < len(Machines); i++ {

		op3 = ButtonOptions()
		op3.X = 245
		op3.Y = float64((i*45)+ScrollState) + 5
		op3.Xsize = 70
		op3.Ysize = 45
		op3.Image = Machines[i].Image
		op3.Tint = Machines[i].Cost > int(Silicon)
		DrawButton(op3, screen)

		if mouse.X > 245 && mouse.Y > (i*45)+5+ScrollState && mouse.Y < (i*45)+59+ScrollState {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(1.5, 2)
			op.GeoM.Translate(153, 40)
			screen.DrawImage(StatBoxImg, op)

			op2 := &text.DrawOptions{}
			op2.GeoM.Scale(0.3, 0.3)
			op2.GeoM.Translate(163, 50)

			if Machines[i].name2 != "" {
				text.Draw(screen, Machines[i].Name, face, op2)
				op2.GeoM.Translate(0, 15)
				text.Draw(screen, Machines[i].name2, face, op2)
			} else {
				text.Draw(screen, Machines[i].Name, face, op2)
			}
			op2.GeoM.Translate(0, 15)
			text.Draw(screen, (fmt.Sprintf("%s %s", "Cost: ", large(float32(Machines[i].Cost)))), face, op2)
			op2.GeoM.Translate(0, 15)
			text.Draw(screen, (fmt.Sprintf("%s %s", "Speed: ", large(Machines[i].Speed))), face, op2)
			op2.GeoM.Translate(0, 15)
			text.Draw(screen, (fmt.Sprintf("%s %d", "Count: ", Machines[i].Count)), face, op2)
		}
	}

	// Van
	if Van.Is {
		op3 = ButtonOptions()
		op3.X = float64(Van.X)
		op3.Y = float64(Van.Y)
		op3.Xsize = 54
		op3.Ysize = 36
		op3.Image = VanImg
		DrawButton(op3, screen)
	}
}

func SimulateUpgrade(g *Game) error {

	// Exit upgrades
	op3 := ButtonOptions()
	op3.X = 20
	op3.Y = 15
	op3.Xsize = 30
	op3.Ysize = 30

	if TestButton(op3) {
		state = "Main"
		ScrollState = 0
		ClickWav.SetPosition(0)
		ClickWav.Play()
	}

	// Upgrades
	for i := 0; i < len(Upgrades); i++ {
		if Upgrades[i].Tier < 5 {
			op3 = ButtonOptions()
			op3.X = 50
			op3.Y = float64((i*30)+ScrollState) + 45
			op3.Xsize = 25
			op3.Ysize = 25

			if TestButton(op3) {
				if Silicon >= float32(Upgrades[i].Cost) {
					Silicon = Silicon - float32(Upgrades[i].Cost)
					if Upgrades[i].Tier == 4 {
						VanWav.SetPosition(0)
						VanWav.Play()
					} else {
						ClickWav.SetPosition(0)
						ClickWav.Play()
					}
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
	return nil
}

func DrawUpgrade(screen *ebiten.Image) {
	screen.Fill(color.RGBA{97, 87, 160, 1})

	// Background
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(6, 4)
	op.GeoM.Translate(20, 20)
	screen.DrawImage(StatBoxImg, op)

	// Silicon readouts
	op2 := &text.DrawOptions{}
	op2.GeoM.Translate(170, 10)
	op2.GeoM.Scale(0.5, 0.5)
	text.Draw(screen, (fmt.Sprintf("%s %.2f", "Silicon: ", Silicon)), face, op2)

	// Exit upgrades
	op3 := ButtonOptions()
	op3.X = 20
	op3.Y = 15
	op3.Xsize = 30
	op3.Ysize = 30
	op3.Text = "X"
	op3.TextSize = 0.6
	op3.Image = StatBoxImg
	DrawButton(op3, screen)

	// Upgrades
	for i := 0; i < len(Upgrades); i++ {
		if (i*30)+ScrollState < 140 && (i*30)+ScrollState > -20 {

			op3 = ButtonOptions()
			op3.X = 50
			op3.Y = float64((i*30)+ScrollState) + 45
			op3.Xsize = 25
			op3.Ysize = 25
			op3.Tint = Upgrades[i].Cost > int(Silicon)

			if Upgrades[i].Tier < 5 {
				op3.Image = Bases[Upgrades[i].Tier]
				DrawButton(op3, screen)
				op3.Image = Upgrades[i].Image
				DrawButton(op3, screen)
			} else {
				op3.Image = Bases[0]
				DrawButton(op3, screen)
				op3.Image = Bases[5]
				DrawButton(op3, screen)
			}

			op2 = &text.DrawOptions{}
			op2.GeoM.Scale(0.4, 0.4)
			op2.GeoM.Translate(80, float64((i*30)+ScrollState)+48)
			text.Draw(screen, fmt.Sprintf("%s %s %d %s %s", Upgrades[i].Name, " ", Upgrades[i].Tier+1, ", Cost: ", large(float32(Upgrades[i].Cost))), face, op2)
		}
	}
}
