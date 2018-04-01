package compose

import (
	"git.yo2.cz/drahoslav/penego/draw"
	"git.yo2.cz/drahoslav/penego/net"
)

type Composable interface {
}

type Composition struct {
	places      map[*net.Place]draw.Pos
	transitions map[*net.Transition]draw.Pos
	ghosts      map[Composable]draw.Pos
}

func NewComposition() Composition {
	return Composition{
		make(map[*net.Place]draw.Pos),
		make(map[*net.Transition]draw.Pos),
		make(map[Composable]draw.Pos),
	}
}

func (comp Composition) HitTest(x, y float64) Composable {
	for place, pos := range comp.places {
		if hitPlace(x, y, pos) {
			return place
		}
	}
	for transition, pos := range comp.transitions {
		if hitTransition(x, y, pos) {
			return transition
		}
	}
	return nil
}

func (comp Composition) Move(node Composable, x, y float64) {
	pos := snap(x, y, 15)
	switch node := node.(type) {
	case *net.Transition:
		comp.transitions[node] = pos
	case *net.Place:
		comp.places[node] = pos
	}
	delete(comp.ghosts, node)
}

func (comp Composition) GhostMove(node Composable, x, y float64) {
	comp.ghosts[node] = snap(x, y, 15)
}

func (comp Composition) DrawWith(drawer draw.Drawer) {
	setStyle := func(node Composable) func(Composable) {
		if _, isGhosted := comp.ghosts[node]; isGhosted {
			drawer.SetStyle(draw.FadedStyle)
			return func(Composable) {
				drawer.SetStyle(draw.FadedStyle)
			}
		} else {
			drawer.SetStyle(draw.DefaultStyle)
			return func(node Composable) {
				if _, isGhosted := comp.ghosts[node]; isGhosted {
					drawer.SetStyle(draw.FadedStyle)
				} else {
					drawer.SetStyle(draw.DefaultStyle)
				}
			}
		}
	}
	for place, pos := range comp.places {
		setStyle(place)
		drawer.DrawPlace(pos, place.Tokens, place.Description)
	}

	for tran, pos := range comp.transitions {
		orSetStyle := setStyle(tran)
		drawer.DrawTransition(pos, tran.TimeFunc.String(), tran.Description)
		for _, arc := range tran.Origins {
			if arc.Place.Hidden() {
				continue
			}
			from := comp.places[arc.Place]
			orSetStyle(arc.Place)
			drawer.DrawInArc(from, pos, arc.Weight)
		}
		for _, arc := range tran.Targets {
			if arc.Place.Hidden() {
				continue
			}
			to := comp.places[arc.Place]
			orSetStyle(arc.Place)
			drawer.DrawOutArc(pos, to, arc.Weight)
		}
	}
	drawer.SetStyle(draw.HighlightedStyle)
	for node, pos := range comp.ghosts {
		switch node := node.(type) {
		case *net.Place:
			place := node
			drawer.DrawPlace(pos, node.Tokens, node.Description)
			for tran, tranPos := range comp.transitions {
				for _, arc := range tran.Origins {
					if arc.Place == place {
						drawer.DrawInArc(pos, tranPos, arc.Weight)
					}
				}
				for _, arc := range tran.Targets {
					if arc.Place == place {
						drawer.DrawOutArc(tranPos, pos, arc.Weight)
					}
				}
			}
		case *net.Transition:
			tran := node
			drawer.DrawTransition(pos, tran.TimeFunc.String(), tran.Description)
			for _, arc := range tran.Origins {
				if arc.Place.Hidden() {
					continue
				}
				from := comp.places[arc.Place]
				drawer.DrawInArc(from, pos, arc.Weight)
			}
			for _, arc := range tran.Targets {
				if arc.Place.Hidden() {
					continue
				}
				to := comp.places[arc.Place]
				drawer.DrawOutArc(pos, to, arc.Weight)
			}
		}
	}
}

// basic "dumb" way to draw a net
func GetSimple(network net.Net) Composition {
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

	composition := NewComposition()

	// compute positions of places
	for i, place := range places {
		composition.places[place] = posOfPlace(i)
	}
	// compute positions of transitions
	for ti, tran := range transitions {
		composition.transitions[tran] = posOfTransition(ti)
	}

	return composition
}
