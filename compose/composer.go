// This package implements Composition
// which is a stuct defining 2-dimensional placeme of petri net elements
package compose

import (
	"fmt"
	"math"
	"strconv"
	"strings"

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

func New() Composition {
	return Composition{
		make(map[*net.Place]draw.Pos),
		make(map[*net.Transition]draw.Pos),
		make(map[Composable]draw.Pos),
	}
}

func (comp Composition) String() string {
	str := ""
	for place, pos := range comp.places {
		str += fmt.Sprintf("%s %v;%v\n", place.Id, pos.X, pos.Y)
	}
	str += fmt.Sprintf("----\n")
	for transition, pos := range comp.transitions {
		str += fmt.Sprintf("%s %v;%v\n", transition.Id, pos.X, pos.Y)
	}
	return str
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

func (comp Composition) FindCenter() (float64, float64) {
	left, right := math.Inf(+1), math.Inf(-1)
	top, bottom := math.Inf(+1), math.Inf(-1)
	/*
	                -
	               top
	              ____
	   - left -> |    | <- right +
	             |____|
	             bottom
	               +
	*/
	enhanceEdges := func(pos draw.Pos) {
		x, y := pos.X, pos.Y
		if x < left {
			left = x
		}
		if x > right {
			right = x
		}
		if y < top {
			top = y
		}
		if y > bottom {
			bottom = y
		}
	}

	for _, pos := range comp.places {
		enhanceEdges(pos)
	}

	for _, pos := range comp.transitions {
		enhanceEdges(pos)
	}

	return center(left, top, right, bottom)
}

// move whole composition so its center is at x, y
// note that is uses Movew which snaps to multiples of 15, so it might not end up exactly on those positions
func (comp Composition) CenterTo(x, y float64) {
	centerX, centerY := comp.FindCenter()

	deltaX := x - centerX
	deltaY := y - centerY

	for node, pos := range comp.places {
		comp.Move(node, pos.X+deltaX, pos.Y+deltaY)
	}

	for node, pos := range comp.transitions {
		comp.Move(node, pos.X+deltaX, pos.Y+deltaY)
	}
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

	// draw arcs
	for tran, pos := range comp.transitions {
		orSetStyle := setStyle(tran)
		for _, arc := range tran.Origins {
			if arc.Place.Hidden() {
				continue
			}
			from := comp.places[arc.Place]
			orSetStyle(arc.Place)
			if arc.Type == net.InhibitorArc {
				drawer.DrawInhibitorArc(from, pos)
			} else {
				drawer.DrawInArc(from, pos, arc.Weight)
			}
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

	// draw all places
	for place, pos := range comp.places {
		setStyle(place)
		drawer.DrawPlace(pos, place.Tokens, place.Description)
	}

	// draw all transtitions
	for tran, pos := range comp.transitions {
		setStyle(tran)
		drawer.DrawTransition(pos, tran.TimeFunc.String(), tran.Description)
	}

	// draw moving items last
	drawer.SetStyle(draw.HighlightedStyle)
	for node, pos := range comp.ghosts {
		switch node := node.(type) {
		case *net.Place:
			place := node
			drawer.DrawPlace(pos, node.Tokens, node.Description)
			for tran, tranPos := range comp.transitions {
				for _, arc := range tran.Origins {
					if arc.Place == place {
						if arc.Type == net.InhibitorArc {
							drawer.DrawInhibitorArc(pos, tranPos)
						} else {
							drawer.DrawInArc(pos, tranPos, arc.Weight)
						}
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
				if arc.Type == net.InhibitorArc {
					drawer.DrawInhibitorArc(from, pos)
				} else {
					drawer.DrawInArc(from, pos, arc.Weight)
				}
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

	composition := New()

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

func Parse(str string, network net.Net) Composition {
	composition := New()
	lines := strings.Split(str, "\n")
	for _, line := range lines {
		parts := strings.Split(line, " ")
		id := parts[0]
		if len(parts) > 1 {
			poss := strings.Split(parts[1], ";")
			x, _ := strconv.ParseFloat(poss[0], 64)
			y, _ := strconv.ParseFloat(poss[1], 64)
			for _, tran := range network.Transitions() {
				if tran.Id == id {
					composition.transitions[tran] = draw.Pos{x, y}
				}
			}
			for _, place := range network.Places() {
				if place.Id == id {
					composition.places[place] = draw.Pos{x, y}
				}
			}
		}
	}
	return composition
}
