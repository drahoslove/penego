package main

import (
	"fmt"
	"time"
	"./net"
	_"github.com/llgcode/draw2d"
)


func main() {

	// definice míst

	// g := &net.Place{Tokens:1} // generator studentů
	// f := &net.Place{Tokens:0, Description:"fronta"}
	// k := &net.Place{Tokens:4, Description:"kuchařky"}
	// v := &net.Place{Description:"výdej"}
	// s := &net.Place{Description:"stravování"}
	// o := &net.Place{} // odchod

	// places := net.Places{g, f, k, v, s, o}

	// // definice přechodů

	// transitions := net.Transitions{}

	// transitions.Push(net.Transition{
	// 	Origins: net.Places{g},
	// 	Targets: net.Places{g,f},
	// 	TimeFunc: net.GetExponentialTimeFunc(30*time.Second),
	// 	Description: "příchody studentů",
	// })
	// transitions.Push(net.Transition{
	// 	Origins: net.Places{f,k},
	// 	Targets: net.Places{v},
	// })
	// transitions.Push(net.Transition{
	// 	Origins: net.Places{v},
	// 	Targets: net.Places{s,k},
	// 	TimeFunc: net.GetExponentialTimeFunc(1*time.Minute),
	// })
	// transitions.Push(net.Transition{
	// 	Origins: net.Places{s},
	// 	Targets: net.Places{o},
	// 	TimeFunc: net.GetUniformTimeFunc(10*time.Minute, 15*time.Minute),
	// })
	// transitions.Push(Transition{
	// 	Origins: []*Place{o},
	// })

	// fmt.Println("Places:", places)
	// fmt.Println("Transitioms:", transitions)


	transitions, places, err := net.Parse(FORMAT_EXAMPLE)
	if err != nil {
		fmt.Println(err)
		return
	}

	_ = places
	for _, place := range places {
		fmt.Println(place)
	}
	for _, tran := range transitions {
		fmt.Println(tran)
	}

	sim := net.NewSimulation(0, 3*time.Hour, transitions)
	sim.DoEveryTime = func () {
		fmt.Println(sim.GetNow(), places)
	}
	sim.Run()

}



const FORMAT_EXAMPLE = `
// definice míst:

g (1) // generování studentů
f (0) "fronta"
k (5) "kuchařky"
v ( ) "výdej"
s ( ) "stravování"
o ( ) // odchod


// definice přechodů:

g	-> [exp(3m)] "příchod studentů" -> g,f
f,k	-> [] -> v
v	-> [exp(1m)] -> s,k
s 	-> [10m-15m] -> o
// o	-> [] "odchod"

z ( ) // ohlášena žloutenka
c ( ) // vyprazdňovací cyklus
i ( ) "karanténa"

[exp(100d)]  -> z
z,g -> [p=1] -> c
c,f	-> [p=3] -> c,o
c,v	-> [p=2] -> c,o,k
c,s	-> [p=1] -> c,o
c	-> [p=0] -> i
i	-> [10d] -> g

`