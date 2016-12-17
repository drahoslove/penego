package main

import (
	"fmt"
	"time"
	"penego/net"
	"penego/gui"
)


func main() {

	var (
		network net.Net
		err error
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
		g := &net.Place{Tokens:1} // generator
		e := &net.Place{Description: "exit"}
		t := &net.Transition{
			Origins: net.Arcs{{1,g}},
			Targets: net.Arcs{{1,g},{2,e}},
			TimeFunc: net.GetExponentialTimeFunc(30*time.Second),
		}
		network = net.New(net.Places{g, e}, net.Transitions{t})
	} else {
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

	sim := net.NewSimulation(0, time.Millisecond, network)
	sim.DoEveryTime = func () {
		fmt.Println(sim.GetNow(), network.Places())
	}

	for i := 0; i < 10; i++ {
		net.TrueRandomSeed()
		sim.Run()
	}

	////////////////////////////////
	gui.OnRedraw(func() {
		for i := 0; i < 10; i++ {
			gui.DrawPlace(60*(i%5), (60*(i-i%5))/5, i);
		}

		gui.DrawTransition(-150, 75)
		gui.DrawTransition(-50, 60)
	})
	gui.Run() // blocks

}

