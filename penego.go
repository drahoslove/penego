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
	NoFlow TimeFlow = iota
	ContinuousFlow
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
		startTime = time.Duration(0)
		endTime = time.Duration(^uint(0) >> 1)
		timeFlow = ContinuousFlow
		timeSpeed = uint(10)
		verbose = false
		idle = true
	)

	flag.DurationVar(&startTime, "start", startTime, "start time of simulation")
	flag.DurationVar(&endTime, "end", endTime, "end time of simulation")
	flag.Var(&timeFlow, "flow", "type of time flow\n\tno, continuous, or natural")
	flag.UintVar(&timeSpeed, "speed", timeSpeed, "time flow acceleration\n\tdifferent meaning for different -flow\n\t")
	flag.BoolVar(&idle, "idle", idle, "preserve window after simulation ends")
	flag.BoolVar(&verbose, "v", verbose, "be more verbose")
	flag.Parse()

	// parse from file if given filename
	if flag.NArg() >= 1 {
		filename := flag.Arg(0)
		filecontent, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			return
		}
		network, err = net.Parse(string(filecontent))
	} else {
		fmt.Println("No pn file specified, using example")
		network, err = net.Parse(`
			g (1)
			e ( ) "exit"
			----
			g -> [exp(1s)] -> g, 2*e
		`)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return
	}


	////////////////////////////////

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

			for i, t := range transitions {
				screen.DrawTransition(posOfTransition(i), t.TimeFunc.String(), t.Description)
				// arcs:
				for j, p := range places {
					for _, arc := range t.Origins {
						if arc.Place == p {
							screen.DrawInArc(posOfPlace(j), posOfTransition(i), arc.Weight)
						}
					}
					for _, arc := range t.Targets {
						if arc.Place == p {
							screen.DrawOutArc(posOfTransition(i), posOfPlace(j), arc.Weight)
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

			switch timeFlow {

			case NoFlow:
				// nothing just jum to the end of simulation

			case NaturalFlow:
				// render as fast as reality, or proportionally faster/slower
				screen.ForceRedraw(false)
				time.Sleep((then-now) / time.Duration(timeSpeed))

			case ContinuousFlow:
				// render continuously, with fixed waits between events, independent of simulation time
				screen.ForceRedraw(false) // dont block
				time.Sleep(time.Second / time.Duration(timeSpeed))

			}

		})

		// simulate
		screen.ForceRedraw(true)
		net.TrueRandomSeed()
		sim.Run()
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
