package main

import (
	"bytes"
	"encoding/json"
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
	Silicone float32
	Upgrades []Upgrade
	Machines []Machine
}

var Particles []Particle

var state string = "Main"

var mouse mousetype

//var Silicone float32 = 9999999999999

var Silicone float32

var SiliconeRate float32

var SiliconeClicked int

type Game struct {
	face text.Face
}

var SiliconeImg = loadImage("Assets/Silicone.png")
var StatBoxImg = loadImage("Assets/StatBox.png")
var VanImg = loadImage("Assets/Van.png")

var SiliconeWav *audio.Player
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
	SiliconeRate = 0
	for i := 0; i < len(Machines); i++ {
		SiliconeRate = SiliconeRate + float32(Machines[i].Count)*Machines[i].Speed
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

	SiliconeWav, _ = loadWav("Assets/Silicone.wav")
	VanWav, _ = loadWav("Assets/Van.wav")
	ClickWav, _ = loadWav("Assets/Click.wav")

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
	err = json.Unmarshal(b, &s)
	return s, err
}

func SimulateMain(g *Game) error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) &&
		ebiten.IsKeyPressed(ebiten.KeyArrowUp) &&
		ebiten.IsKeyPressed(ebiten.KeyArrowLeft) &&
		ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		Silicone = Silicone + float32(Machines[len(Machines)-1].Cost)
	}
	SiliconeClicked = SiliconeClicked - 1

	if int(rand.Float64()*600) == 1 {
		if int(rand.Float32()*3) == 1 {
			Van = VanType{0, int(rand.Float32() * 160), true}
		}
	}
	if Van.Is {
		Van.X++
		if Van.X > 320 {
			Van.Is = false
		}
	}
	if (mouse.X > 25 && mouse.X < 75 && mouse.Y > 80 && mouse.Y < 130 && mouse.JSTIn) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		Silicone++
		SiliconeWav.SetPosition(0)
		SiliconeWav.Play()
		SiliconeClicked = 5
		Particles = append(Particles, Particle{
			X:       rand.Float32() * 100,
			Y:       0,
			XVel:    0,
			YVel:    -0.4,
			Image:   SiliconeImg,
			Gravity: 0.1,
			Timer:   240,
			Scale:   0.3,
		})
		Particles = append(Particles, Particle{
			X:       40 + rand.Float32()*20,
			Y:       90,
			YVel:    0.5,
			Gravity: 0.02,
			Timer:   60,
			Scale:   0.5,
			Text:    "+1",
		})
	}

	if mouse.JSTIn || (mouse.In && ebiten.IsKeyPressed(ebiten.KeyShift)) {
		if mouse.X > 120 && mouse.X < 180 && mouse.Y < 60 {
			ClickWav.SetPosition(0)
			ClickWav.Play()
			var Save, err = LoadGame()
			if !(err == nil) {
				return err
			}
			Silicone = Save.Silicone
			for i := 0; i < len(Machines); i++ {
				Machines[i].Count = Save.Machines[i].Count
				Machines[i].Cost = Save.Machines[i].Cost
				Machines[i].Speed = Save.Machines[i].Speed
			}
			for i := 0; i < len(Upgrades); i++ {
				Upgrades[i].Tier = Save.Upgrades[i].Tier
				Upgrades[i].Cost = Save.Upgrades[i].Cost
			}
		}

		if mouse.X > Van.X && mouse.X < Van.X+53 && mouse.Y > Van.Y && mouse.Y < Van.Y+36 && Van.Is {
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
			Silicone = Silicone * 2
			for i := 0; i < 15; i++ {
				Particles = append(Particles, Particle{
					X:       rand.Float32() * 100,
					Y:       0,
					XVel:    0,
					YVel:    -0.4,
					Image:   SiliconeImg,
					Gravity: 0.1,
					Timer:   240,
					Scale:   0.3,
				})
			}
		}
		if mouse.X < 110 && mouse.Y > 200 {
			SaveGame(SaveData{
				Silicone,
				Upgrades,
				Machines,
			})
			return ebiten.Termination
		}
		if mouse.X > 110 && mouse.Y > 200 && mouse.X < 230 {
			state = "Upgrades"
			ClickWav.SetPosition(0)
			ClickWav.Play()
		}
		for i := 0; i < len(Machines); i++ {
			if mouse.X > 245 && mouse.Y > i*45 && mouse.Y < (i*45)+45 {
				if Silicone >= float32(Machines[i].Cost) {
					if i == len(Machines)-1 && Machines[i].Count == 0 {
						VanWav.SetPosition(0)
						VanWav.Play()
					} else {
						ClickWav.SetPosition(0)
						ClickWav.Play()
					}
					Silicone = Silicone - float32(Machines[i].Cost)
					Machines[i].Count++
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

	op := &ebiten.DrawImageOptions{}
	if SiliconeClicked > 0 {
		var multi float64 = float64(SiliconeClicked)
		multi = multi / 50
		multi++
		op.GeoM.Scale(multi, multi)
		multi = float64(SiliconeClicked)
		multi = multi / 5
		op.GeoM.Translate(-1.25*multi, -4*multi)
	}
	op.GeoM.Translate(25, 80)
	screen.DrawImage(SiliconeImg, op)
	drawParticles(screen)

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
			text.Draw(screen, (fmt.Sprintf("%s %d", "Count: ", Machines[i].Count)), face, op2)
		}
	}

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(1.2, 0.6)
	op.GeoM.Translate(120, 10)

	screen.DrawImage(StatBoxImg, op)

	op2 = &text.DrawOptions{}
	op2.GeoM.Scale(0.5, 0.5)
	op2.GeoM.Translate(130, 15)
	text.Draw(screen, "Load", face, op2)

	if Van.Is {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(3, 3)
		op.GeoM.Translate(float64(Van.X), float64(Van.Y))
		screen.DrawImage(VanImg, op)
	}
}

func SimulateUpgrade(g *Game) error {
	if mouse.JSTIn || ebiten.IsKeyPressed(ebiten.KeySpace) || (mouse.In && ebiten.IsKeyPressed(ebiten.KeyShift)) {
		if mouse.X > 20 && mouse.X < 45 && mouse.Y > 15 && mouse.Y < 40 {
			state = "Main"
			ClickWav.SetPosition(0)
			ClickWav.Play()
		}
		for i := 0; i < len(Upgrades); i++ {
			if Upgrades[i].Tier < 5 {
				if mouse.X > 50 && mouse.Y > (i*30)+45 && mouse.X < 70 && mouse.Y < (i*30)+65 {
					if Silicone >= float32(Upgrades[i].Cost) {
						Silicone = Silicone - float32(Upgrades[i].Cost)
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
