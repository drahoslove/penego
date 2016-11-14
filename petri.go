package main

import (
	"fmt"
	"time"
	"./net"
	_"github.com/llgcode/draw2d"
)


const NOTATION_EXAMPLE = `
// definice míst:

g (1) // generování studentů
f (0) "fronta"
k (5) "kuchařky"
v ( ) "výdej"
s ( ) "stravování"
o ( ) // odchod
z ( ) // ohlášena žloutenka
c ( ) // vyprazdňovací cyklus
i ( ) "karanténa"

// definice přechodů:

g	-> [exp(3m)] "příchod studentů" -> g,f
f,k	-> [] -> v
v	-> [exp(1m)] -> s,k
s 	-> [10m..15m] -> o
// o	-> [] "odchod"

[exp(100d)]  -> z
z,g -> [p=1] -> c
c,f	-> [p=3] -> c,o
c,v	-> [p=2] -> c,o,k
c,s	-> [p=1] -> c,o
c	-> [p=0] -> i
i	-> [10d] -> g

`

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
	 *    |       |   2   exit
	 *     ----->[ ]----->( )
	 *         exp(30s)
	 */

	// can be done likek this:

	g := &net.Place{Tokens:1} // generator
	e := &net.Place{Description: "exit"}
	t := &net.Transition{
		Origins: net.Arcs{{1,g}},
		Targets: net.Arcs{{1,g},{2,e}},
		TimeFunc: net.GetExponentialTimeFunc(30*time.Second),
	}
	network = net.New(net.Places{g, e}, net.Transitions{t})

	// or like this:
	network, err = net.Parse(`
		g (1)
		e ( ) "exit"
		----
		g -> [exp(3ms)] -> g, 2*e
	`)


	////////////////////////////////

	// network, err = net.Parse(NOTATION_EXAMPLE)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(network)

	sim := net.NewSimulation(0, time.Second, network)
	sim.DoEveryTime = func () {
		fmt.Println(sim.GetNow(), network.Places())
	}

	for i := 0; i < 10; i++ {
		net.TrueRandomSeed()
		sim.Run()
	}

}


