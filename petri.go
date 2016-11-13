package main

import (
	"fmt"
	"time"
	"./net"
	_"github.com/llgcode/draw2d"
)


const FORMAT_EXAMPLE = `
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
		places net.Places
		transitions net.Transitions
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
	places = net.Places{g, e}
	transitions = net.Transitions{t}

	// or like this:

	transitions, places, err = net.Parse(`
		g (1) // geneartor
		e ( ) "exit"
		g -> [exp(30s)] -> g, e
	`)


	////////////////////////////////

	transitions, places, err = net.Parse(FORMAT_EXAMPLE)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, place := range places {
		fmt.Println(place)
	}
	for _, tran := range transitions {
		fmt.Println(tran)
	}


	sim := net.NewSimulation(0, 3*time.Hour, transitions)
	sim.DoEveryTime = func () {
	}

	for i := 0; i < 10; i++ {
		net.TrueRandomSeed()
		sim.Run()
		fmt.Println(sim.GetNow(), places)
	}

}


