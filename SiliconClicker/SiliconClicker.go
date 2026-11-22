package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
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

type Prestiege struct {
	Brought   bool
	Cost      int
	X         float32
	Y         float32
	Connected []*Prestiege
	Image     *ebiten.Image
	Name      string
}

type SaveData struct {
	Silicon         float32
	Upgrades        []Upgrade
	Machines        []Machine
	Prestiege       []Prestiege
	PrestiegeCost   int
	PrestiegeTokens int
}

type options struct {
	X           float64
	Y           float64
	Xsize       float64
	Ysize       float64
	Text        string
	Image       *ebiten.Image
	TextSize    float64
	Tint        bool
	LeftMargin  float64
	TopMargin   float64
	BorderWidth int
	SideSize    int
}

type DisectedButton struct {
	Side   *ebiten.Image
	Corner *ebiten.Image
	Center *ebiten.Image
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

var PrestiegeCost int

var PrestiegeTokens int

type Game struct {
	face text.Face
}

var SiliconImg = loadImage("Assets/Images/Silicon.png")
var StatBoxImg = loadImage("Assets/Images/StatBox.png")
var VanImg = loadImage("Assets/Images/Van.png")
var XbuttonImg = loadImage("Assets/Images/Xbutton.png")

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
	return audioContext.NewPlayer(d)
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
	// Mouse
	var _, Scroll = ebiten.Wheel()
	ScrollState += int(Scroll * 30)
	ScrollState = min(ScrollState, 0)

	// Dev
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) &&
		ebiten.IsKeyPressed(ebiten.KeyArrowUp) &&
		ebiten.IsKeyPressed(ebiten.KeyArrowLeft) &&
		ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		Silicon = Silicon + float32(Machines[len(Machines)-1].Cost)
	}

	// Silicon/Sec calculation
	SiliconRate = 0
	for i := 0; i < len(Machines); i++ {
		SiliconRate = SiliconRate + float32(Machines[i].Count)*Machines[i].Speed
	}
	Silicon = Silicon + (SiliconRate / 60)
	mouse.JSTIn = ebiten.IsMouseButtonPressed(ebiten.MouseButton0) && !mouse.In
	mouse.In = ebiten.IsMouseButtonPressed(ebiten.MouseButton0)
	mouse.X, mouse.Y = ebiten.CursorPosition()

	// SimSwitch
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
	case "PrestigeMenu":
		{
			var err = SimulatePrestigeMenu(g)
			if err != nil {
				return err
			}

		}
	case "PrestigeTree":
		{
			var err = SimulatePrestigeTree(g)
			if err != nil {
				return err
			}

		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// DrawSwitch
	switch state {
	case "Main":
		{
			DrawMain(screen)
		}
	case "Upgrades":
		{
			DrawUpgrade(screen)
		}
	case "PrestigeMenu":
		{
			DrawPrestigeMenu(screen)
		}
	case "PrestigeTree":
		{
			DrawPrestigeTree(screen)
		}
	}
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	return 640, 480
}

var Machines []Machine

var Upgrades []Upgrade

var Prestieges []Prestiege

var Bases []*ebiten.Image

var Van VanType

func ButtonOptions() (op options) {
	return options{
		Xsize:       100,
		Ysize:       100,
		TextSize:    1,
		Tint:        false,
		LeftMargin:  10,
		TopMargin:   5,
		BorderWidth: 0,
		SideSize:    2,
	}
}

func DrawComplexButton(op options, screen *ebiten.Image) {
	W := op.BorderWidth
	Button := DisectedButton{
		op.Image.RecyclableSubImage(image.Rect(W, 0, int(op.Image.Bounds().Dx())-W, W)),
		op.Image.RecyclableSubImage(image.Rect(0, 0, W, W)),
		op.Image.RecyclableSubImage(image.Rect(W, W, int(op.Image.Bounds().Dx())-W, int(op.Image.Bounds().Dy())-W)),
	}

	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Scale((op.Xsize-float64(op.SideSize*W))/float64(Button.Center.Bounds().Size().X), (op.Ysize-float64(op.SideSize*W*2))/float64(Button.Center.Bounds().Size().Y))
	op2.GeoM.Translate(op.X, op.Y)
	op2.GeoM.Translate(float64(W*op.SideSize), float64(W*op.SideSize))
	screen.DrawImage(Button.Center, op2)

	for i := 0; i < 4; i++ {
		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Scale(float64(op.SideSize), float64(op.SideSize))
		op2.GeoM.Rotate(1.57079633 * float64(i))
		op2.GeoM.Translate(op.X, op.Y)
		if i == 1 || i == 2 {
			op2.GeoM.Translate(op.Xsize, 0)
		}
		if i == 2 || i == 3 {
			op2.GeoM.Translate(0, op.Ysize)
		}
		screen.DrawImage(Button.Corner, op2)

		op2 = &ebiten.DrawImageOptions{}
		if i == 1 || i == 3 {
			op2.GeoM.Scale((op.Ysize-float64(op.SideSize*W*2))/float64(Button.Side.Bounds().Size().X), float64(op.SideSize))
		} else {
			op2.GeoM.Scale((op.Xsize-float64(op.SideSize*W*2))/float64(Button.Side.Bounds().Size().X), float64(op.SideSize))
		}
		op2.GeoM.Translate(float64(W*op.SideSize), 0)
		op2.GeoM.Rotate(1.57079633 * float64(i))
		op2.GeoM.Translate(op.X, op.Y)
		if i == 1 || i == 2 {
			op2.GeoM.Translate(op.Xsize, 0)
		}
		if i == 2 || i == 3 {
			op2.GeoM.Translate(0, op.Ysize)
		}
		screen.DrawImage(Button.Side, op2)
	}
}

func DrawButton(op options, screen *ebiten.Image) {
	if op.BorderWidth == 0 {
		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Scale(op.Xsize/float64(op.Image.Bounds().Size().X), op.Ysize/float64(op.Image.Bounds().Size().Y))
		op2.GeoM.Translate(op.X, op.Y)
		if op.Tint {
			op2.ColorScale.Set(0.5, 0.5, 0.5, 1)
		}
		screen.DrawImage(op.Image, op2)
	} else {
		DrawComplexButton(op, screen)
	}

	op3 := &text.DrawOptions{}
	op3.GeoM.Scale(op.TextSize, op.TextSize)
	op3.GeoM.Translate(op.X, op.Y)
	op3.GeoM.Translate(op.LeftMargin, op.TopMargin)
	text.Draw(screen, op.Text, face, op3)
}

func TestButton(op options) (clicked bool) {
	return (mouse.X > int(op.X) &&
		mouse.X < int(op.X+(op.Xsize)) &&
		mouse.Y > int(op.Y) &&
		mouse.Y < int(op.Y+(op.Ysize)) &&
		(mouse.JSTIn || mouse.In && (ebiten.IsKeyPressed(ebiten.KeyShift))))
}

func ResetMachines() {
	Upgrades = []Upgrade{
		{
			"NPN Cursor",
			2,
			loadImage("Assets/Upgrades/Upgrades/Mouse.png"),
			"Mouse",
			150,
			0,
		},
		{
			"Glass Matrix",
			1.5,
			loadImage("Assets/Upgrades/Upgrades/Slicer.png"),
			"Slicer",
			50,
			0,
		},
		{
			"Better Toner",
			1.5,
			loadImage("Assets/Upgrades/Upgrades/Printer.png"),
			"Printer",
			300,
			0,
		},
		{
			"Better LEDs",
			1.5,
			loadImage("Assets/Upgrades/Upgrades/Exposer.png"),
			"Exposer",
			800,
			0,
		},
		{
			"RadioAcid",
			2,
			loadImage("Assets/Upgrades/Upgrades/Etcher.png"),
			"Etcher",
			1500,
			0,
		},
		{
			"New Servos",
			2,
			loadImage("Assets/Upgrades/Upgrades/Assembler.png"),
			"Assembler",
			5000,
			0,
		},
		{
			"Thermo Crystals",
			2,
			loadImage("Assets/Upgrades/Upgrades/Wafer.png"),
			"Wafer Fab",
			20000,
			0,
		},
		{
			"Intagrated Casing",
			2,
			loadImage("Assets/Upgrades/Upgrades/ChipLine.png"),
			"Chip Packing",
			100000,
			0,
		},
		{
			"Nuclear Forge",
			2,
			loadImage("Assets/Upgrades/Upgrades/Foundry.png"),
			"Automated",
			2000000,
			0,
		},
		{
			"Entanglement",
			2,
			loadImage("Assets/Upgrades/Upgrades/Quantum.png"),
			"Quantum",
			10000000,
			0,
		},
		{
			"Atomic Precision",
			2,
			loadImage("Assets/Upgrades/Upgrades/NanoFab.png"),
			"NanoFab",
			400000000,
			0,
		},
		{
			"Neural Optimization",
			2.5,
			loadImage("Assets/Upgrades/Upgrades/AI.png"),
			"AI",
			2000000000,
			0,
		},
		{
			"Orbital Logistics",
			2.5,
			loadImage("Assets/Upgrades/Upgrades/Orbital.png"),
			"Orbital",
			15000000000,
			0,
		},
		{
			"Matter Reclamation",
			2.5,
			loadImage("Assets/Upgrades/Upgrades/Molecular.png"),
			"Molecular",
			800000000000,
			0,
		},
		{
			"Sturdy Circuits",
			3,
			loadImage("Assets/Upgrades/Upgrades/Plasma.png"),
			"Plasma Etcher",
			4600000000000,
			0,
		},
		{
			"Shift Cooling",
			3,
			loadImage("Assets/Upgrades/Upgrades/Dark.png"),
			"Dark Silicon",
			150000000000000,
			0,
		},
		{
			"Spacetime Routing",
			3,
			loadImage("Assets/Upgrades/Upgrades/Void.png"),
			"Void Core",
			10000000000000000,
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
		{
			"NanoFab",
			30000,
			loadImage("Assets/Machines/NanoFab.png"),
			10000000,
			0,
			"",
		},
		{
			"AI Compiler",
			100000,
			loadImage("Assets/Machines/AI.png"),
			100000000,
			0,
			"",
		},
		{
			"Orbital",
			3500000,
			loadImage("Assets/Machines/Orbital.png"),
			1800000000,
			0,
			"Foundry",
		},
		{
			"Molecular",
			80000000,
			loadImage("Assets/Machines/Molecular.png"),
			50000000000,
			0,
			"Assembler",
		},
		{
			"Plasma Etcher",
			500000000,
			loadImage("Assets/Machines/Plasma.png"),
			800000000000,
			0,
			"",
		},
		{
			"Dark Silicon",
			6000000000,
			loadImage("Assets/Machines/Dark.png"),
			20000000000000,
			0,
			"Reactor",
		},
		{
			"Void Core",
			40000000000,
			loadImage("Assets/Machines/Void.png"),
			1000000000000000,
			0,
			"Fabricator",
		},
	}
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("SiliconClicker")

	SiliconWav, _ = loadWav("Assets/Sounds/Silicon.wav")
	VanWav, _ = loadWav("Assets/Sounds/Van.wav")
	ClickWav, _ = loadWav("Assets/Sounds/Click.wav")

	Bases = []*ebiten.Image{
		loadImage("Assets/Upgrades/Bases/Tier1.png"),
		loadImage("Assets/Upgrades/Bases/Tier2.png"),
		loadImage("Assets/Upgrades/Bases/Tier3.png"),
		loadImage("Assets/Upgrades/Bases/Tier4.png"),
		loadImage("Assets/Upgrades/Bases/Tier5.png"),
	}

	ResetMachines()

	var Save, err = LoadGame()
	if err == nil {
		Silicon = Save.Silicon
		PrestiegeCost = Save.PrestiegeCost
		PrestiegeTokens = Save.PrestiegeTokens
		Prestieges = Save.Prestiege
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
	if in >= 1000000000000000000 {
		out = "IDUNU"
	} else if in >= 1000000000000000 {
		out = fmt.Sprintf("%.1fQ", in/1000000000000000)
	} else if in >= 1000000000000 {
		out = fmt.Sprintf("%.1fT", in/1000000000000)
	} else if in >= 1000000000 {
		out = fmt.Sprintf("%.1fB", in/1000000000)
	} else if in >= 1000000 {
		out = fmt.Sprintf("%.1fM", in/1000000)
	} else if in >= 1000 {
		out = fmt.Sprintf("%.1fK", in/1000)
	} else {
		out = fmt.Sprintf("%.1f", in)
	}
	return out
}

func SimulateMain(g *Game) error {

	// Scroll Limiter
	if ScrollState < -1000 {
		ScrollState = -1000
	}

	// Silicon animation progresser
	SiliconClicked = SiliconClicked - 1

	// Spawn Van
	if int(rand.Float64()*600) == 1 {
		if int(rand.Float32()*3) == 1 {
			Van = VanType{0, int(rand.Float32() * 300), true}
		}
	}

	// Simulate Van
	if Van.Is {
		Van.X++
		if Van.X > 640 {
			Van.Is = false
		}
	}

	// Silicon
	op3 := ButtonOptions()
	op3.X = 60
	op3.Y = 160
	op3.Xsize = 100
	op3.Ysize = 100
	if TestButton(op3) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		Silicon += float32((Upgrades[0].Tier + 1) * (Upgrades[0].Tier + 1))
		SiliconWav.SetPosition(0)
		SiliconWav.Play()
		SiliconClicked = 5
		Particles = append(Particles, Particle{
			X:       rand.Float32() * 250,
			Y:       0,
			XVel:    0,
			YVel:    -0.8,
			Image:   SiliconImg,
			Gravity: 0.2,
			Timer:   240,
			Scale:   0.6,
		})
		Particles = append(Particles, Particle{
			X:       float32(mouse.X),
			Y:       float32(mouse.Y + -10),
			YVel:    -2.4,
			Gravity: 0.04,
			Timer:   60,
			Scale:   1,
			Text:    fmt.Sprintf("+%d", (Upgrades[0].Tier+1)*(Upgrades[0].Tier+1)),
		})
	}

	// Quit
	op3 = ButtonOptions()
	op3.X = 20
	op3.Y = 410
	op3.Xsize = 85
	op3.Ysize = 66

	if TestButton(op3) {
		SaveGame(SaveData{
			Silicon,
			Upgrades,
			Machines,
			Prestieges,
			PrestiegeCost,
			PrestiegeTokens,
		})
		return ebiten.Termination
	}

	// Upgrades
	op3 = ButtonOptions()
	op3.X = 130
	op3.Y = 410
	op3.Xsize = 155
	op3.Ysize = 66

	if TestButton(op3) {
		state = "Upgrades"
		ScrollState = 0
		ClickWav.SetPosition(0)
		ClickWav.Play()
	}

	// Prestiege

	op3 = ButtonOptions()
	op3.X = 305
	op3.Y = 410
	op3.Xsize = 145
	op3.Ysize = 66
	if TestButton(op3) {
		state = "PrestigeMenu"
		ScrollState = 0
		ClickWav.SetPosition(0)
		ClickWav.Play()
	}

	// Van collision
	if Van.Is {
		op3 = ButtonOptions()
		op3.X = float64(Van.X)
		op3.Y = float64(Van.Y)
		op3.Xsize = 108
		op3.Ysize = 72

		if TestButton(op3) {
			Van.Is = false
			VanWav.SetPosition(0)
			VanWav.Play()
			Particles = append(Particles, Particle{
				X:       float32(Van.X),
				Y:       float32(Van.Y),
				XVel:    0,
				YVel:    -0.8,
				Image:   VanImg,
				Gravity: 0.2,
				Timer:   240,
				Scale:   6,
			})
			Silicon += (SiliconRate * 60)
			for i := 0; i < 15; i++ {
				Particles = append(Particles, Particle{
					X:       rand.Float32() * 250,
					Y:       0,
					XVel:    0,
					YVel:    -0.8,
					Image:   SiliconImg,
					Gravity: 0.2,
					Timer:   240,
					Scale:   0.6,
				})
			}
		}
	}

	// Machines
	for i := 0; i < len(Machines); i++ {
		op3 = ButtonOptions()
		op3.X = 490
		op3.Y = float64((i*90)+ScrollState) + 10
		op3.Xsize = 140
		op3.Ysize = 90

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
		if Thisparticle.Timer <= 0 || Thisparticle.Y > 640 {
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
	vector.StrokeLine(screen, 120, 480, 120, 405, 10, color.RGBA{30, 30, 20, 1}, false)
	vector.StrokeLine(screen, 295, 480, 295, 405, 10, color.RGBA{30, 30, 20, 1}, false)
	vector.StrokeLine(screen, 0, 400, 455, 400, 10, color.RGBA{30, 30, 20, 1}, false)
	vector.StrokeLine(screen, 460, 0, 460, 480, 10, color.RGBA{30, 30, 20, 1}, false)

	// Silicon
	op3 := ButtonOptions()
	op3.X = 60
	op3.Y = 160
	op3.Xsize = 100
	op3.Ysize = 100
	op3.Image = SiliconImg
	if SiliconClicked > 0 {
		op3.X = 60 - (float64(SiliconClicked))
		op3.Y = 160 - (float64(SiliconClicked))
		op3.Xsize = 100 + (float64(SiliconClicked) * 2)
		op3.Ysize = 100 + (float64(SiliconClicked) * 2)
	}

	DrawButton(op3, screen)
	drawParticles(screen)

	// Silicon Stats
	op2 := &text.DrawOptions{}
	op2.GeoM.Translate(20, 60)
	op2.GeoM.Scale(1, 1)
	text.Draw(screen, (fmt.Sprintf("%s %s", "Silicon: ", large(Silicon))), face, op2)
	op2.GeoM.Translate(0, 60)
	text.Draw(screen, (fmt.Sprintf("%s %s", "Silicon/Sec: ", large(SiliconRate))), face, op2)

	// Quit
	op3 = ButtonOptions()
	op3.X = 15
	op3.Y = 410
	op3.Xsize = 90
	op3.Ysize = 66
	op3.Image = StatBoxImg
	op3.Text = "Quit"
	op3.TextSize = 1
	op3.BorderWidth = 5
	DrawButton(op3, screen)

	// Upgrades

	op3 = ButtonOptions()
	op3.X = 135
	op3.Y = 410
	op3.Xsize = 150
	op3.Ysize = 66
	op3.Image = StatBoxImg
	op3.Text = "Upgrades"
	op3.TextSize = 0.9
	op3.BorderWidth = 5
	DrawButton(op3, screen)

	// Prestiege

	op3 = ButtonOptions()
	op3.X = 305
	op3.Y = 410
	op3.Xsize = 145
	op3.Ysize = 66
	op3.Image = StatBoxImg
	op3.Text = "Prestiege"
	op3.TextSize = 0.9
	op3.BorderWidth = 5
	DrawButton(op3, screen)

	// Machines
	for i := 0; i < len(Machines); i++ {

		op3 = ButtonOptions()
		op3.X = 490
		op3.Y = float64((i*90)+ScrollState) + 10
		op3.Xsize = 140
		op3.Ysize = 90
		op3.Image = Machines[i].Image
		op3.Tint = Machines[i].Cost > int(Silicon)
		DrawButton(op3, screen)

		if mouse.X > 490 && mouse.Y > (i*90)+10+ScrollState && mouse.Y < (i*90)+118+ScrollState {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(3.5, 4)
			op.GeoM.Translate(296, 80)
			screen.DrawImage(StatBoxImg, op)

			op2 := &text.DrawOptions{}
			op2.GeoM.Scale(0.6, 0.6)
			op2.GeoM.Translate(316, 100)

			if Machines[i].name2 != "" {
				text.Draw(screen, Machines[i].Name, face, op2)
				op2.GeoM.Translate(0, 30)
				text.Draw(screen, Machines[i].name2, face, op2)
			} else {
				text.Draw(screen, Machines[i].Name, face, op2)
			}
			op2.GeoM.Translate(0, 30)
			text.Draw(screen, (fmt.Sprintf("%s %s", "Cost: ", large(float32(Machines[i].Cost)))), face, op2)
			op2.GeoM.Translate(0, 30)
			text.Draw(screen, (fmt.Sprintf("%s %s", "Speed: ", large(Machines[i].Speed))), face, op2)
			op2.GeoM.Translate(0, 30)
			text.Draw(screen, (fmt.Sprintf("%s %d", "Count: ", Machines[i].Count)), face, op2)
		}
	}

	// Van
	if Van.Is {
		op3 = ButtonOptions()
		op3.X = float64(Van.X)
		op3.Y = float64(Van.Y)
		op3.Xsize = 108
		op3.Ysize = 72
		op3.Image = VanImg
		DrawButton(op3, screen)
	}
}

func SimulateUpgrade(g *Game) error {

	// Scroll Limiter
	if ScrollState < -800 {
		ScrollState = -800
	}

	// Exit upgrades
	op3 := ButtonOptions()
	op3.X = 40
	op3.Y = 30
	op3.Xsize = 60
	op3.Ysize = 60

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
			op3.X = 100
			op3.Y = float64((i*60)+ScrollState) + 90
			op3.Xsize = 50
			op3.Ysize = 50

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
	op3 := ButtonOptions()
	op3.X = 30
	op3.Y = 20
	op3.Xsize = 580
	op3.Ysize = 440
	op3.Image = StatBoxImg
	op3.BorderWidth = 10
	op3.SideSize = 10
	DrawButton(op3, screen)

	// Silicon readouts
	op2 := &text.DrawOptions{}
	op2.GeoM.Translate(150, 20)
	op2.GeoM.Scale(1, 1)
	text.Draw(screen, (fmt.Sprintf("%s %s", "Silicon: ", large(Silicon))), face, op2)

	// Exit upgrades
	op3 = ButtonOptions()
	op3.X = 40
	op3.Y = 30
	op3.Xsize = 60
	op3.Ysize = 60
	op3.Image = XbuttonImg
	DrawButton(op3, screen)

	// Upgrades
	for i := 0; i < len(Upgrades); i++ {
		if (i*60)+ScrollState < 280 && (i*60)+ScrollState > -40 {
			if Upgrades[i].Tier < 5 {
				op3 = ButtonOptions()
				op3.X = 100
				op3.Y = float64((i*60)+ScrollState) + 90
				op3.Xsize = 50
				op3.Ysize = 50
				op3.Tint = Upgrades[i].Cost > int(Silicon)
				op3.Image = Bases[Upgrades[i].Tier]
				DrawButton(op3, screen)
				op3.Image = Upgrades[i].Image
				DrawButton(op3, screen)
			} else {
				vector.StrokeLine(screen, 150,
					float32((i*60)+ScrollState)+110,
					480, float32((i*60)+ScrollState)+110,
					4, color.White, false)
			}

			op2 = &text.DrawOptions{}
			op2.GeoM.Scale(0.8, 0.8)
			op2.GeoM.Translate(160, float64((i*60)+ScrollState)+96)
			text.Draw(screen, fmt.Sprintf("%s %s %d %s %s", Upgrades[i].Name, " ", Upgrades[i].Tier+1, ", Cost: ", large(float32(Upgrades[i].Cost))), face, op2)
		}
	}
}

func SimulatePrestigeMenu(g *Game) (err error) {

	// Exit prestige
	op3 := ButtonOptions()
	op3.X = 40
	op3.Y = 30
	op3.Xsize = 90
	op3.Ysize = 90

	if TestButton(op3) {
		state = "Main"
		ScrollState = 0
		ClickWav.SetPosition(0)
		ClickWav.Play()
	}

	// Newgame Button
	op3 = ButtonOptions()
	op3.X = 150
	op3.Y = 30
	op3.Xsize = 250
	op3.Ysize = 90

	if TestButton(op3) && ebiten.IsKeyPressed(ebiten.KeyShift) {
		os.Remove("save.json")
		ResetMachines()
		Silicon = 0
	}

	return nil
}

func SimulatePrestigeTree(g *Game) (err error) {
	return nil
}

func DrawPrestigeMenu(screen *ebiten.Image) {
	screen.Fill(color.RGBA{97, 87, 160, 1})

	// Exit Button
	op3 := ButtonOptions()
	op3.X = 40
	op3.Y = 30
	op3.Xsize = 90
	op3.Ysize = 90
	op3.Image = XbuttonImg
	DrawButton(op3, screen)

	// Newgame Button
	op3 = ButtonOptions()
	op3.X = 150
	op3.Y = 30
	op3.Xsize = 250
	op3.Ysize = 90
	op3.Image = StatBoxImg
	op3.Text = "NewGame"
	op3.TextSize = 1.3
	op3.LeftMargin = 23
	op3.TopMargin = 15
	op3.BorderWidth = 5
	DrawButton(op3, screen)

	// Readouts
	op3 = ButtonOptions()
	op3.X = 30
	op3.Xsize = 350
	op3.Ysize = 60
	op3.Image = StatBoxImg
	op3.TextSize = 1
	op3.LeftMargin = 15
	op3.TopMargin = 5
	op3.BorderWidth = 5
	op3.Text = (fmt.Sprintf("%s %s", "Silicon: ", large(Silicon)))
	op3.Y = 140
	DrawButton(op3, screen)
	op3.Text = (fmt.Sprintf("%s %s", "Tokens: ", large(float32(PrestiegeTokens))))
	op3.Y = 210
	DrawButton(op3, screen)
	op3.Text = (fmt.Sprintf("%s %s", "Next Cost: ", large(float32(PrestiegeCost))))
	op3.Y = 280
	DrawButton(op3, screen)
}

func DrawPrestigeTree(screen *ebiten.Image) {

}
