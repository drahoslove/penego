package main // import "git.yo2.cz/drahoslav/penego"

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"git.yo2.cz/drahoslav/penego/compose"
	"git.yo2.cz/drahoslav/penego/export"
	"git.yo2.cz/drahoslav/penego/gui"
	"git.yo2.cz/drahoslav/penego/net"
	"git.yo2.cz/drahoslav/penego/pnml"
	"git.yo2.cz/drahoslav/penego/storage"
	"github.com/pkg/profile"
	"github.com/skratchdot/open-golang/open"
)

const EXAMPLE = `
// This is an example new file, please edit it and save
# NET
g (1)
e ( ) "exit"
----
g -> p[exp(1s)] -> g, 2*e
`

type State int

const (
	New State = iota
	Initial
	Running
	Paused
	Stopped
	Splash
	Idle
	Exit
)

type TimeFlow int

const (
	NoFlow         TimeFlow = iota // no waits, just jum to the end of simulation
	ContinuousFlow                 // render as fast as reality, or proportionally faster/slower
	NaturalFlow                    // render continuously, with fixed waits between events, independent of simulation time
)

func (flow TimeFlow) String() string {
	return map[TimeFlow]string{
		NoFlow:         "no",
		ContinuousFlow: "continuous",
		NaturalFlow:    "natural",
	}[flow]
}

func (flow *TimeFlow) Set(name string) error {
	val, ok := map[string]TimeFlow{
		"no":         NoFlow,
		"continuous": ContinuousFlow,
		"natural":    NaturalFlow,
	}[name]
	if !ok {
		return fmt.Errorf("may be: no, continuous, natural")
	}
	*flow = val
	return nil
}

func main() {
	if os.Getenv("PROFILE") != "" {
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	}

	// TODO init elsewhere?
	pwd, _ := os.Getwd()
	storage.Of("export").
		Set("width", 1024).
		Set("height", 512).
		Set("background", false).
		Set("zoom", 0).
		Set("png.filename", pwd+string(filepath.Separator)+"image.png").
		Set("pdf.filename", pwd+string(filepath.Separator)+"image.pdf").
		Set("svg.filename", pwd+string(filepath.Separator)+"image.svg")
	storage.Of("settings").
		Set("linewidth", 2.0)
	storage.Of("gui.offset").
		Set("x", 0.0).
		Set("y", 0.0)

	var (
		network     net.Net
		composition compose.Composition
		err         error
	)

	// flags

	var (
		startTime  = time.Duration(0)
		endTime    = time.Hour * 24 * 1e5
		timeFlow   = ContinuousFlow
		timeSpeed  = uint(10)
		trueRandom = false
		noClose    = true
		autoStart  = false

		verbose = false
		input   = ""
		output  = ""
	)

	flag.DurationVar(&startTime, "start", startTime, "start `time` of simulation")
	flag.DurationVar(&endTime, "end", endTime, "end `time` of simulation")
	flag.Var(&timeFlow, "flow", "type of time flow\n\tno, continuous, or natural")
	flag.UintVar(&timeSpeed, "speed", timeSpeed, "time flow acceleration\n\tdifferent meaning for different -flow\n\t")
	flag.BoolVar(&trueRandom, "truerandom", trueRandom, "seed pseudo random generator with true random seed on start")
	flag.BoolVar(&noClose, "noclose", noClose, "preserve window after simulation ends")
	flag.BoolVar(&autoStart, "autostart", autoStart, "automatic start of simulation")
	flag.BoolVar(&verbose, "v", verbose, "be more verbose")

	flag.StringVar(&input, "i", input, "import file - *.(pnml|xml)")
	flag.StringVar(&output, "o", output, "export file - *.(png|svg|pdf)\n\t(this means no gui)")
	flag.Parse()

	////////////////////////////////

	// load network from file if given filename

	pnString := EXAMPLE

	read := func(filename string) string {
		fileContent, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
			return ""
		}
		return string(fileContent)
	}

	filename := flag.Arg(0)

	if len(filename) > 0 {
		pnString = read(filename)
	} else {
		log.Println("No penego file specified, using example")
	}
	network, composition = Parse(pnString)
	composition.Rotate()

	if input != "" {
		file, err := os.Open(input)
		if err != nil {
			log.Fatalln("cant open import file", err)
			return
		}
		defer file.Close()
		network, composition = pnml.Parse(file)
	}

	////////////////////////////////

	if output != "" { // headless mode
		err := export.ByName(output, composition.DrawWith)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	// else gui mode

	gui.Run(func(screen *gui.Screen) { // runs this anon func in goroutine

		storage.Of("settings").OnChange(func(st storage.Storage, key string) {
			screen.ForceRedraw(false)
		})
		storage.Of("export").OnChange(func(st storage.Storage, key string) {
			screen.ForceRedraw(false)
		})

		var state State = Splash

		var onStateChange = func(before, now time.Duration) {
			switch timeFlow {
			case NoFlow:
			case NaturalFlow:
				time.Sleep((now - before) / time.Duration(timeSpeed))
			case ContinuousFlow:
				time.Sleep(time.Second / time.Duration(timeSpeed))
			}
			if verbose {
				log.Println(now, network.Places())
			}
			screen.SetTitle(now.String())
			screen.ForceRedraw(false) // must not block
		}

		var sim net.Simulation

		foo := func() {}
		_ = foo

		reloader := makeFileWatcher(func(filename string) {
			sim.Stop()
			pnString = read(filename)
			network, composition = Parse(pnString)
			if verbose {
				log.Println(network)
			}
			state = New
		})
		defer reloader.close()

		reloader.watch(filename)

		// action functions:

		step := func() {
			switch state {
			case Paused:
				sim.Pause()
				sim.Step()
			}
		}

		playPause := func() {
			switch state {
			case Paused:
				state = Running
			case Running:
				state = Paused
				sim.Pause()
			}
		}
		reset := func() {
			switch state {
			case Running, Paused, Idle:
				sim.Stop()
				state = Initial
			}
		}
		quit := func() {
			screen.SetShouldClose(true)
		}

		create := func() {
			fpath := os.TempDir() + "/tempPenego.pn"
			f, err := os.Create(fpath)
			if err != nil {
				log.Fatal(err)
			}
			f.Write([]byte(EXAMPLE))
			f.Close()

			reloader.watch(fpath)
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "vim"
			}
			open.StartWith(fpath, editor)
		}

		open := func() {
			gui.LoadFile(func(filename string) {
				if verbose {
					fmt.Println(filename)
				}
				if err != nil {
					fmt.Fprintln(os.Stderr, "cant open file", err)
					return
				}
				screen.Reset()
				reloader.watch(filename)
			})
		}

		save := func() {
			gui.SaveFile(func(filename string) {
				file, err := os.Create(filename)
				if err != nil {
					fmt.Fprintln(os.Stderr, "cant save file", err)
					return
				}
				defer file.Close()
				str := Stringify(network, composition)
				file.WriteString(str)
				if verbose {
					fmt.Println(str)
				}
				reloader.watch(filename)
			})
		}

		doImport := func() {
			gui.LoadFile(func(filename string) {
				file, err := os.Open(filename)
				if err != nil {
					fmt.Fprintln(os.Stderr, "cant import file", err)
					return
				}
				defer file.Close()
				screen.Reset()
				network, composition = pnml.Parse(file)
				sim.Stop()
				state = New
				log.Println("net imported", filename)
				if verbose {
					fmt.Println(network)
				}
			})
		}

		doExport := func() {
			gui.ToggleExport(func(filename string) {
				err := export.ByName(filename, composition.DrawWith)
				if err != nil {
					log.Println(err)
				} else {
					log.Printf("image %s exported\n", filename)
				}
			})
		}

		settings := func() {
			gui.ToggleSettings()
		}

		center := func() {
			screen.Reset()
			composition.CenterTo(0, 0)
		}
		isCenter := func() bool {
			offset := storage.Of("gui.offset")
			ox, oy := offset.Float("x"), offset.Float("y")
			x, y := composition.FindCenter()
			return !(x == ox && y == oy)
		}
		rotate := func() {
			composition.Rotate()
		}

		// up bar commands
		screen.OnKey("Q", quit)
		// screen.RegisterControl(0, "Q", gui.AlwaysIcon(gui.QuitIcon), "quit", quit, gui.True)

		screen.RegisterControl(0, "N", gui.AlwaysIcon(gui.FileIcon), "new", create, gui.True)
		screen.RegisterControl(0, "O", gui.AlwaysIcon(gui.OpenIcon), "open", open, gui.True) // penego format
		screen.RegisterControl(0, "S", gui.AlwaysIcon(gui.SaveIcon), "save", save, gui.True) // penego format
		screen.RegisterControl(0, "R", gui.AlwaysIcon(gui.ReloadIcon), "reload", reloader.action, reloader.isOn)
		screen.RegisterControl(0, "I", gui.AlwaysIcon(gui.ImportIcon), "import net", doImport, gui.True)   // from pnml
		screen.RegisterControl(0, "E", gui.AlwaysIcon(gui.ExportIcon), "export image", doExport, gui.True) // to svg/pdf/png
		screen.RegisterControl(0, "P", gui.AlwaysIcon(gui.SettingsIcon), "settings", settings, gui.True)
		screen.RegisterControl(0, "C", gui.AlwaysIcon(gui.CenterOnIcon), "center net", center, isCenter)
		screen.RegisterControl(0, "T", gui.AlwaysIcon(gui.RotateIcon), "rotate net", rotate, gui.True)

		// down bar commands (simulation related)
		screen.RegisterControl(1, "home", gui.AlwaysIcon(gui.BeginIcon), "reset", reset, gui.True)
		screen.RegisterControl(1, "right", gui.AlwaysIcon(gui.NextStepIcon), "step", step, gui.True)
		screen.RegisterControl(1, "space", func() gui.Icon {
			if state != Running {
				return gui.PlayIcon
			} else {
				return gui.PauseIcon
			}
		}, "play/pause", playPause, gui.True)

		screen.OnMouseMove(true, func(x, y float64) bool {
			return composition.HitTest(x, y) != nil
		})

		screen.OnDrag(true, func(x, y, dx, dy, sx, sy float64, done bool) {
			node := composition.HitTest(sx, sy)
			if node != nil { // drag node
				if done {
					composition.Move(node, x, y)
				} else {
					composition.GhostMove(node, x, y)
				}
			} else { // pan view
				screen.Pan(dx, dy)
			}
			screen.ForceRedraw(false)
		})

		// main state machine

		for state != Exit {
			switch state {
			case Splash:
				// show splash for 2 seconds
				screen.SetRedrawFuncToSplash("Penego")
				time.Sleep(time.Second * 1)
				state = New
			case New:
				sim = net.NewSimulation(startTime, endTime, network)
				if trueRandom {
					net.TrueRandomSeed()
				}
				state = Initial
			case Initial:
				sim.Init()
				sim.DoEveryStateChange(onStateChange)
				screen.SetRedrawFunc(gui.RedrawFunc(composition.DrawWith))
				if autoStart {
					state = Running
				} else {
					state = Paused
				}
				screen.SetTitle(sim.GetNow().String() + " init")
			case Running:
				sim.Run()             ////////////////// <--
				if state != Running { // paused or stopped
					continue
				}
				screen.SetTitle(sim.GetNow().String() + " done")
				screen.ForceRedraw(true)
				if verbose {
					fmt.Println("----")
				}
				if noClose {
					state = Idle
				} else {
					state = Exit
				}
			case Paused:
				time.Sleep(time.Millisecond * 20)
				screen.SetTitle(sim.GetNow().String() + " paused")
			case Idle:
				time.Sleep(time.Millisecond * 20)
			}
		}

	}) // returns when func returns

}
