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
		err     error
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
	}


	// this petri net:

	/**
	 *
	 *   (1)<-----
	 *    |       |
	 *    |       |    2    exit
	 *     ----->[ ]------->( )
	 *         exp(30s)
	 */

	if false {
		// can be created like this:
		g := &net.Place{Tokens: 1} // generator
		e := &net.Place{Description: "exit"}
		t := &net.Transition{
			Origins:  net.Arcs{{1, g}},
			Targets:  net.Arcs{{1, g}, {2, e}},
			TimeFunc: net.GetExponentialTimeFunc(30 * time.Microsecond),
			Description: "gen",
		}
		network = net.New(net.Places{g, e}, net.Transitions{t})
	}
	if false {
		// or like this:
		network, err = net.Parse(`
			g (1)
			e ( ) "exit"
			----
			g -> [exp(30us)] -> g, 2*e
		`)
		if err != nil {
			panic(err)
		}
	}

	////////////////////////////////

	fmt.Println(network)

	// sim := net.NewSimulation(0, time.Millisecond, network)
	// sim.DoEveryTime = func() {
	// 	fmt.Println(sim.GetNow(), network.Places())
	// }

	// for i := 0; i < 10; i++ {
	// 	net.TrueRandomSeed()
	// 	sim.Run()
	// }

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

		sim := net.NewSimulation(0, time.Millisecond, network)
		sim.DoEveryTimeChange = func() {
			fmt.Println(sim.GetNow(), network.Places())
			screen.ForceRedraw()
			time.Sleep(time.Second/5)
		}

		net.TrueRandomSeed()
		sim.Run()
		fmt.Println("----")

		for true {
			time.Sleep(time.Second/1)
		}
	}) // returns when func returns

}
