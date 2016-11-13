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

	// this ne net:

	/**
	 *
	 *   (1)<-----
	 *    |       |
	 *    |       |       exit
	 *     ----->[ ]----->( )
	 *         exp(30s)
	 */

	// can be done likek this:

	g := &net.Place{Tokens:1} // generator
	e := &net.Place{Description: "exit"}
	t := &net.Transition{
		Origins: net.Places{g},
		Targets: net.Places{g,e},
		TimeFunc: net.GetExponentialTimeFunc(30*time.Second),
	}
	network = net.New(net.Places{g, e}, net.Transitions{t})

	// or like this:

	network, err = net.Parse(`
		g (1) // geneartor
		e ( ) "exit"
		g -> [exp(30s)] -> g, e
	`)


	////////////////////////////////

	network, err = net.Parse(NOTATION_EXAMPLE)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(network)

	sim := net.NewSimulation(0, 3*time.Hour, network)
	sim.DoEveryTime = func () {
	}

	for i := 0; i < 10; i++ {
		net.TrueRandomSeed()
		sim.Run()
		fmt.Println(sim.GetNow(), network.Places())
	}

}


