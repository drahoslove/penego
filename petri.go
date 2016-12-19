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

			for i, p := range places {
				x := i * 90 - len(places)/2 * 90
				y := -90
				screen.DrawPlace(x, y, p.Tokens, p.Description)
			}

			for i, t := range transitions {
				x := i * 90 - len(transitions)/2 * 90 - 90/2
				y := +60
				screen.DrawTransition(x, y, t.TimeFunc.String(), t.Description)
			}
		})

		sim := net.NewSimulation(0, time.Hour*24*10, network)

		sim.DoEveryStateChange(func(now, then time.Duration) {
			fmt.Println(sim.GetNow(), network.Places())
			screen.SetTitle(now.String())
			// time.Sleep((then-now)/1000) // render 1000× faster than reality

			// time.Sleep((then-now)) // render as fast as reality

			// time.Sleep(time.Second/5) // render to be comprehentable
			screen.ForceRedraw()
		})

		for i:=0; i<10; i++ {
			screen.ForceRedraw()
			net.TrueRandomSeed()
			sim.Run()
			fmt.Println("----")
			time.Sleep(time.Second*5)
		}

		screen.SetTitle("done")
		for true {
			time.Sleep(time.Second/1)
		}
	}) // returns when func returns

}
