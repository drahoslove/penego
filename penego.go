package main

import (
	"fmt"
	"time"
	"os"
	"io/ioutil"
	"github.com/pkg/profile"
	"penego/gui"
	"penego/net"
)

func main() {
	if os.Getenv("PROFILE") != "" {
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	}

	var (
		network net.Net
		err error
		startTime = time.Duration(0)
		endTime = time.Duration(int(^uint(0) >> 1))
	)

	// TODO parse additinal flags
	_ = startTime
	_ = endTime

	// parse from file if given filename
	if len(os.Args) == 2 {
		filename := os.Args[1]
		filecontent, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
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
		panic(err)
	}


	////////////////////////////////

	fmt.Println(network)


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
		screen.ForceRedraw()

		sim := net.NewSimulation(startTime, endTime, network)

		sim.DoEveryStateChange(func(now, then time.Duration) {
			fmt.Println(now, network.Places())
			screen.SetTitle(now.String())
			screen.ForceRedraw()

			// time.Sleep((then-now)/1000) // render 1000× faster than reality
			// time.Sleep((then-now)) // render as fast as reality
			time.Sleep(time.Second/5) // render to be comprehentable
		})

		// simulate
		screen.ForceRedraw()
		net.TrueRandomSeed()
		sim.Run()
		fmt.Println("----")
		screen.SetTitle("done")

		// idle
		for true {
			time.Sleep(time.Second)
		}

	}) // returns when func returns

}
