package compose

import (
	"git.yo2.cz/drahoslav/penego/draw"
	"git.yo2.cz/drahoslav/penego/net"
)

type Composer func(draw.Drawer)

// basic "dumb" way to draw a net
func GetSimple(network net.Net) func(draw.Drawer) {
	places := network.Places()
	transitions := network.Transitions()

	const BASE = 90.0

	posOfPlace := func(i int) draw.Pos {
		pos := draw.Pos{
			X: float64(i)*BASE - (float64(len(places))/2-0.5)*BASE,
			Y: 0,
		}
		if len(transitions) <= 1 {
			pos.Y += BASE
		}
		return pos
	}

	posOfTransition := func(i int) draw.Pos {
		pos := draw.Pos{
			X: float64(i)*BASE - (float64(len(transitions))/2)*BASE + BASE/2,
			Y: 4*BASE*float64(i%2) - 2*BASE,
		}
		if len(transitions) <= 1 {
			pos.Y += BASE
		}
		return pos
	}

	return Composer(func(drawer draw.Drawer) {
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
	})
}
