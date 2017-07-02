package main

import (
	"git.yo2.cz/drahoslav/penego/gui"
	"git.yo2.cz/drahoslav/penego/net"
)

func getDrawNet(network net.Net) gui.RedrawFunc {
	places := network.Places()
	transitions := network.Transitions()

	const BASE = 90.0

	posOfPlace := func(i int) gui.Pos {
		pos := gui.Pos{
			X: float64(i)*BASE - (float64(len(places))/2-0.5)*BASE,
			Y: 0,
		}
		if len(transitions) <= 1 {
			pos.Y += BASE
		}
		return pos
	}

	posOfTransition := func(i int) gui.Pos {
		pos := gui.Pos{
			X: float64(i)*BASE - (float64(len(transitions))/2)*BASE + BASE/2,
			Y: 4*BASE*float64(i%2) - 2*BASE,
		}
		if len(transitions) <= 1 {
			pos.Y += BASE
		}
		return pos
	}

	return func(screen *gui.Screen) {
		for i, p := range places {
			screen.DrawPlace(posOfPlace(i), p.Tokens, p.Description)
		}

		for ti, t := range transitions {
			screen.DrawTransition(posOfTransition(ti), t.TimeFunc.String(), t.Description)
			// arcs:
			for pi, p := range places {
				for _, arc := range t.Origins {
					if arc.Place == p {
						screen.DrawInArc(posOfPlace(pi), posOfTransition(ti), arc.Weight)
					}
				}
				for _, arc := range t.Targets {
					if arc.Place == p {
						screen.DrawOutArc(posOfTransition(ti), posOfPlace(pi), arc.Weight)
					}
				}
			}
		}
	}
}
