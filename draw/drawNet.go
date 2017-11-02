package draw

import (
	"git.yo2.cz/drahoslav/penego/net"
)

type Pos struct {
	X float64
	Y float64
}

type NetDrawer interface {
	DrawPlace(pos Pos, n int, description string)
	DrawTransition(pos Pos, attrs, description string)
	DrawInArc(from, to Pos, weight int)
	DrawOutArc(from, to Pos, weight int)
}

// basic "dumb" way to draw a net
func GetDrawNet(network net.Net) func(NetDrawer) {
	places := network.Places()
	transitions := network.Transitions()

	const BASE = 90.0

	posOfPlace := func(i int) Pos {
		pos := Pos{
			X: float64(i)*BASE - (float64(len(places))/2-0.5)*BASE,
			Y: 0,
		}
		if len(transitions) <= 1 {
			pos.Y += BASE
		}
		return pos
	}

	posOfTransition := func(i int) Pos {
		pos := Pos{
			X: float64(i)*BASE - (float64(len(transitions))/2)*BASE + BASE/2,
			Y: 4*BASE*float64(i%2) - 2*BASE,
		}
		if len(transitions) <= 1 {
			pos.Y += BASE
		}
		return pos
	}

	return func(drawer NetDrawer) {
		for i, p := range places {
			drawer.DrawPlace(posOfPlace(i), p.Tokens, p.Description)
		}

		for ti, t := range transitions {
			drawer.DrawTransition(posOfTransition(ti), t.TimeFunc.String(), t.Description)
			// arcs:
			for pi, p := range places {
				for _, arc := range t.Origins {
					if arc.Place == p {
						drawer.DrawInArc(posOfPlace(pi), posOfTransition(ti), arc.Weight)
					}
				}
				for _, arc := range t.Targets {
					if arc.Place == p {
						drawer.DrawOutArc(posOfTransition(ti), posOfPlace(pi), arc.Weight)
					}
				}
			}
		}
	}
}
