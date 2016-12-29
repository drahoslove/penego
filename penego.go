package main

import (
	"fmt"
	"time"
	"os"
	"io/ioutil"
	"flag"
	"github.com/pkg/profile"
	"penego/gui"
	"penego/net"
)

type TimeFlow int

const (
	// no waits, just jum to the end of simulation
	NoFlow TimeFlow = iota
	// render as fast as reality, or proportionally faster/slower
	ContinuousFlow
	// render continuously, with fixed waits between events, independent of simulation time
	NaturalFlow
)

func (flow TimeFlow) String() string {
	return map[TimeFlow]string{
		NoFlow: "no",
		ContinuousFlow: "continuous",
		NaturalFlow: "natural",
	}[flow]
}

func (flow *TimeFlow) Set(name string) error {
	val, ok := map[string]TimeFlow{
		"no": NoFlow,
		"continuous": ContinuousFlow,
		"natural": NaturalFlow,
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
		err error
	)

	// flags

	var (
		startTime = time.Duration(0)
		endTime = time.Hour * 24 * 1e5
		timeFlow = ContinuousFlow
		timeSpeed = uint(10)
		truerandom = false
		idle = true
		verbose = false
	)

	flag.DurationVar(&startTime, "start", startTime, "start time of simulation")
	flag.DurationVar(&endTime, "end", endTime, "end time of simulation")
	flag.Var(&timeFlow, "flow", "type of time flow\n\tno, continuous, or natural")
	flag.UintVar(&timeSpeed, "speed", timeSpeed, "time flow acceleration\n\tdifferent meaning for different -flow\n\t")
	flag.BoolVar(&truerandom, "truerandom", truerandom, "seed pseudorandom generator with true random seed on start")
	flag.BoolVar(&idle, "idle", idle, "preserve window after simulation ends")
	flag.BoolVar(&verbose, "v", verbose, "be more verbose")
	flag.Parse()


	////////////////////////////////

	// load network from file if given filename

	pnstring := `
		g (1)
		e ( ) "exit"
		----
		g -> [exp(1s)] -> g, 2*e
	`
	if flag.NArg() >= 1 {
		filename := flag.Arg(0)
		filecontent, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			return
		}
		pnstring = string(filecontent)
	} else {
		fmt.Println("No pn file specified, using example")
	}
	network, err = net.Parse(pnstring)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return
	}

	if verbose {
		fmt.Println(network)
	}


	////////////////////////////////

	gui.Run(func(screen *gui.Screen) { // runs this anon func in goroutine

		// show splash for 2 seconds
		time.Sleep(time.Second * 2)

		// how to draw
		screen.SetRedrawFunc(func() {
			places := network.Places()
			transitions := network.Transitions()

			const BASE = 90.0

			posOfPlace := func(i int) gui.Pos {
				pos := gui.Pos{
					X: float64(i) * BASE - (float64(len(places))/2 - 0.5) * BASE,
					Y: 0,
				}
				if len(transitions) <= 1 {
					pos.Y += BASE
				}
				return pos
			}

			posOfTransition := func(i int) gui.Pos {
				pos := gui.Pos{
					X: float64(i) * BASE - (float64(len(transitions))/2) * BASE + BASE/2,
					Y: 4 * BASE * float64(i % 2) - 2 * BASE,
				}
				if len(transitions) <= 1 {
					pos.Y += BASE
				}
				return pos
			}

			for i, p := range places {
				screen.DrawPlace(posOfPlace(i), p.Tokens, p.Description)
			}

			for ti, t := range transitions {
				screen.DrawTransition(posOfTransition(ti), t.TimeFunc.String(), t.Description)
				// arcs:
				for pi, p := range places {
					for _, arc := range t.Origins {
						if arc.Place == p {
							screen.DrawInArc(posOfPlace(pi), posOfTransition(ti), arc.Weight)
						}
					}
					for _, arc := range t.Targets {
						if arc.Place == p {
							screen.DrawOutArc(posOfTransition(ti), posOfPlace(pi), arc.Weight)
						}
					}
				}
			}

		})


		////////////////

		// draw initial state
		screen.ForceRedraw(true)

		sim := net.NewSimulation(startTime, endTime, network)

		sim.DoEveryStateChange(func(now, then time.Duration) {
			if verbose {
				fmt.Println(now, network.Places())
			}
			screen.SetTitle(now.String())
			screen.ForceRedraw(false) // donnt block

			switch timeFlow {
			case NoFlow:
			case NaturalFlow:
				time.Sleep((then-now) / time.Duration(timeSpeed))
			case ContinuousFlow:
				time.Sleep(time.Second / time.Duration(timeSpeed))
			}

		})

		// simulate

		if truerandom {
			net.TrueRandomSeed()
		}

		sim.Run() ////////////////// <--

		screen.ForceRedraw(true)
		if verbose {
			fmt.Println("----")
		}
		screen.SetTitle(sim.GetNow().String() + " done")

		// idle
		for idle {
			time.Sleep(time.Second)
		}

	}) // returns when func returns

}
