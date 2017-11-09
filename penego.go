package main // import "git.yo2.cz/drahoslav/penego"

import (
	"flag"
	"fmt"
	"git.yo2.cz/drahoslav/penego/compose"
	"git.yo2.cz/drahoslav/penego/gui"
	"git.yo2.cz/drahoslav/penego/net"
	"github.com/pkg/profile"
	"github.com/sqweek/dialog"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const EXAMPLE = `
	g (1)
	e ( ) "exit"
	----
	g -> [exp(1s)] -> g, 2*e
`

type State int

const (
	Initial State = iota
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

	var (
		network net.Net
		err     error
	)

	// flags

	var (
		startTime  = time.Duration(0)
		endTime    = time.Hour * 24 * 1e5
		timeFlow   = ContinuousFlow
		timeSpeed  = uint(10)
		trueRandom = false
		noClose    = true
		verbose    = false
		autoStart  = false
	)

	flag.DurationVar(&startTime, "start", startTime, "start `time` of simulation")
	flag.DurationVar(&endTime, "end", endTime, "end `time` of simulation")
	flag.Var(&timeFlow, "flow", "type of time flow\n\tno, continuous, or natural")
	flag.UintVar(&timeSpeed, "speed", timeSpeed, "time flow acceleration\n\tdifferent meaning for different -flow\n\t")
	flag.BoolVar(&trueRandom, "truerandom", trueRandom, "seed pseudo random generator with true random seed on start")
	flag.BoolVar(&noClose, "noclose", noClose, "preserve window after simulation ends")
	flag.BoolVar(&verbose, "v", verbose, "be more verbose")
	flag.BoolVar(&autoStart, "autostart", autoStart, "automatic start")
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
	parse := func(pnString string) (network net.Net) {
		network, err = net.Parse(pnString)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			return
		}
		if verbose {
			fmt.Println(network)
		}
		return
	}

	filename := flag.Arg(0)

	if len(filename) > 0 {
		pnString = read(filename)
	} else {
		fmt.Println("No pn file specified, using example")
	}
	network = parse(pnString)

	////////////////////////////////

	gui.Run(func(screen *gui.Screen) { // runs this anon func in goroutine

		var state State = Splash

		// how to draw
		var composeNet = compose.GetSimple(network)

		var onStateChange = func(before, now time.Duration) {
			switch timeFlow {
			case NoFlow:
			case NaturalFlow:
				time.Sleep((now - before) / time.Duration(timeSpeed))
			case ContinuousFlow:
				time.Sleep(time.Second / time.Duration(timeSpeed))
			}
			if verbose {
				fmt.Println(now, network.Places())
			}
			screen.SetTitle(now.String())
			screen.ForceRedraw(false) // block
		}

		var sim net.Simulation

		foo := func() {}
		_ = foo

		reloader := makeFileWatcher(func(filename string) {
			pnString = read(filename)
			network = parse(pnString)
			composeNet = compose.GetSimple(network)
			sim.Stop()
			state = Initial
		})
		defer reloader.close()

		reloader.watch(filename)
		reloader.action()

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

		open := func() {
			go func() {
				filename, err := dialog.File().Filter("Penego notation", "pn").SetStartDir(".").Load()
				if verbose {
					fmt.Println(filename)
				}
				if err != nil {
					return
				}
				reloader.watch(filename)
				reloader.action()
			}()
		}

		// up bar commands
		screen.RegisterControl(0, "Q", gui.AlwaysIcon(gui.QuitIcon), "quit", quit, gui.True)
		screen.RegisterControl(0, "O", gui.AlwaysIcon(gui.FileIcon), "open", open, gui.True)
		screen.RegisterControl(0, "R", gui.AlwaysIcon(gui.ReloadIcon), "reload", reloader.action, reloader.isOn)

		// down bar commands (simulation related)
		screen.RegisterControl(1, "R", gui.AlwaysIcon(gui.PrevIcon), "reset", reset, gui.True)
		screen.RegisterControl(1, "space", func() gui.Icon {
			if state != Running {
				return gui.PlayIcon
			} else {
				return gui.PauseIcon
			}
		}, "play/pause", playPause, gui.True)

		for state != Exit {
			switch state {
			case Splash:
				// show splash for 2 seconds
				screen.SetRedrawFuncToSplash("Penego")
				time.Sleep(time.Second * 1)
				state = Initial
			case Initial:
				sim = net.NewSimulation(startTime, endTime, network)
				sim.DoEveryStateChange(onStateChange)
				if trueRandom {
					net.TrueRandomSeed()
				}
				screen.SetRedrawFunc(gui.RedrawFunc(composeNet))
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
