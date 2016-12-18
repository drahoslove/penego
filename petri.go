package main

import (
	"fmt"
	"time"
	"math/rand"
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

	go func () {
		time.Sleep(time.Second*2)

		noise := func() int {
			return rand.Int()%2
		}
		gui.OnRedraw(func() {
			for i := 0; i < 10; i++ {
				gui.DrawPlace(60*(i%5)+noise(), (60*(i-i%5))/5 + noise(), i);
			}

			gui.DrawTransition(-150+noise(), 75+noise())
			gui.DrawTransition(-50+noise(), 60+noise())
		})
		for count := 1000; count > 0; count-- {

			gui.ForceRedraw()

			time.Sleep(time.Second/20)
		}
	}()

	gui.Run() // blocks

}

