package main

import (
	"fmt"
	"math/rand"
	"penego/gui"
	"penego/net"
	"time"
	"os"
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

	// this petri net:

	/**
	 *
	 *   (1)<-----
	 *    |       |
	 *    |       |    2    exit
	 *     ----->[ ]------->( )
	 *         exp(30s)
	 */

	if true {
		// can be done likek this:
		g := &net.Place{Tokens: 1} // generator
		e := &net.Place{Description: "exit"}
		t := &net.Transition{
			Origins:  net.Arcs{{1, g}},
			Targets:  net.Arcs{{1, g}, {2, e}},
			TimeFunc: net.GetExponentialTimeFunc(30 * time.Second),
			Description: "gen",
		}
		network = net.New(net.Places{g, e}, net.Transitions{t})
	} else {
		// or like this:
		network, err = net.Parse(`
			g (1)
			e ( ) "exit"
			----
			g -> [exp(30us)] "gen" -> g, 2*e
		`)
		if err != nil {
			panic(err)
		}
	}

	////////////////////////////////

	fmt.Println(network)

	sim := net.NewSimulation(0, time.Millisecond, network)
	sim.DoEveryTime = func() {
		fmt.Println(sim.GetNow(), network.Places())
	}

	for i := 0; i < 10; i++ {
		net.TrueRandomSeed()
		sim.Run()
	}

	////////////////////////////////

	gui.Run(func() { // runs this anon func in goroutine
		// time.Sleep(time.Second * 2) // show splash for 2 seconds

		noise := func() int {
			return rand.Int() % 4
		}
		gui.OnRedraw(func() {
			noise()
			places := network.Places()
			transitions := network.Transitions()

			for i, p := range places {
				x := i * 90 - len(places)/2 * 90
				y := -90
				gui.DrawPlace(x, y, p.Tokens, p.Description)
			}

			for i, t := range transitions {
				x := i * 90 - len(transitions)/2 * 90 - 90/2
				y := +60
				gui.DrawTransition(x, y, t.TimeFunc.String(), t.Description)
			}

		})

		// shake for 10 seconds
		tick := time.After(time.Second * 10)
		shake: for {
			select {
			case <-tick:
				break shake
			default:
				gui.ForceRedraw() // synced with frameloop, waits for frame start
				time.Sleep(time.Second/20)
			}
		}
		// chill forever
		for true {
			time.Sleep(time.Second/1)
		}
	}) // returns when func returns

}
