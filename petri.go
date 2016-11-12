package main

import (
	"fmt"
	"time"
	"./net"
	_"github.com/llgcode/draw2d"
)


func main() {

	// definice míst

	g := &net.Place{Tokens:1} // generator studentů
	f := &net.Place{Tokens:0, Description:"fronta"}
	k := &net.Place{Tokens:4, Description:"kuchařky"}
	v := &net.Place{Description:"výdej"}
	s := &net.Place{Description:"stravování"}
	o := &net.Place{} // odchod

	places := net.Places{g, f, k, v, s, o}

	// definice přechodů

	transitions := net.Transitions{}

	transitions.Push(net.Transition{
		Origins: net.Places{g},
		Targets: net.Places{g,f},
		TimeFunc: net.GetExponentialTimeFunc(30*time.Second),
		Description: "příchody studentů",
	})
	transitions.Push(net.Transition{
		Origins: net.Places{f,k},
		Targets: net.Places{v},
	})
	transitions.Push(net.Transition{
		Origins: net.Places{v},
		Targets: net.Places{s,k},
		TimeFunc: net.GetExponentialTimeFunc(1*time.Minute),
	})
	transitions.Push(net.Transition{
		Origins: net.Places{s},
		Targets: net.Places{o},
		TimeFunc: net.GetUniformTimeFunc(10*time.Minute, 15*time.Minute),
	})
	// transitions.Push(Transition{
	// 	Origins: []*Place{o},
	// })

	fmt.Println("Places:", places)
	fmt.Println("Transitioms:", transitions)

	sim := net.NewSimulation(0, 3*time.Hour, transitions)
	sim.DoEveryTime = func () {
		fmt.Println(sim.GetNow(), places)
	}
	sim.Run()


}


