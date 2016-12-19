package main

import (
	"fmt"
	"penego/gui"
	"penego/net"
	"time"
	"os"
	"io/ioutil"
	"github.com/pkg/profile"
)

func main() {
	if os.Getenv("PROFILE") != "" {
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	}

	var (
		network net.Net
	)

	// parse from file if given filename
	if len(os.Args) == 2 {
		filename := os.Args[1]
		filecontent, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		network, err = net.Parse(string(filecontent))
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("No pn file specified")
		return
	}


	////////////////////////////////

	fmt.Println(network)


	////////////////////////////////

	gui.Run(func(screen *gui.Screen) { // runs this anon func in goroutine
		// time.Sleep(time.Second * 2) // show splash for 2 seconds

		screen.SetRedrawFunc(func() {
			places := network.Places()
			transitions := network.Transitions()

			posOfPlace := func(i int) gui.Pos {
				return gui.Pos{
					X: i * 90 - len(places)/2 * 90,
					Y: 0,
				}
			}

			posOfTransition := func(i int) gui.Pos {
				return gui.Pos{
					X: i * 90 - len(transitions)/2 * 90 - 90/2,
					Y: 300 * (i % 2) - 150,
				}
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

			// screen.DrawInArc(posOfPlace(0), posOfTransition(0), 0)
			// screen.DrawOutArc(posOfTransition(2), posOfPlace(3), 0)


		})

		screen.ForceRedraw()

		////////////////

		sim := net.NewSimulation(0, time.Hour*24*10, network)

		sim.DoEveryStateChange(func(now, then time.Duration) {
			fmt.Println(sim.GetNow(), network.Places())
			screen.SetTitle(now.String())
			// time.Sleep((then-now)/1000) // render 1000× faster than reality

			// time.Sleep((then-now)) // render as fast as reality

			// time.Sleep(time.Second/5) // render to be comprehentable
			screen.ForceRedraw()
		})

		// for i:=0; i<10; i++ {
		// 	screen.ForceRedraw()
		// 	net.TrueRandomSeed()
		// 	sim.Run()
		// 	fmt.Println("----")
		// 	time.Sleep(time.Second*5)
		// }

		screen.SetTitle("done")
		for true {
			time.Sleep(time.Second/1)
		}
	}) // returns when func returns

}
