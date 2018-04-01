package net

import (
	"fmt"
	"strconv"
	"strings"
)

const MaxInt = int(^uint(0) >> 1)

/******* types *******/

/* Net */

type Net struct {
	places      Places
	transitions Transitions
}

func New(places Places, transitions Transitions) Net { // TODO make this a pointer type?
	return Net{places, transitions}
}

func (net *Net) Places() Places {
	return net.places
}

func (net *Net) Transitions() Transitions {
	return net.transitions
}

func (net Net) String() (str string) {
	for _, pl := range net.places {
		str += pl.String() + "\n"
	}
	str += "----\n"
	for _, tr := range net.transitions {
		str += tr.String() + "\n"
	}
	return
}

func (net *Net) Equals(another *Net) (bool, error) {
	if a, b := len(net.places), len(another.places); a != b {
		return false, fmt.Errorf("num of places does not match (%d != %d)", a, b)
	}
	if a, b := len(net.transitions), len(another.transitions); a != b {
		return false, fmt.Errorf("num of transitions does not match (%d != %d)", a, b)
	}

	pairedPs := map[int]bool{}
pairPlace:
	for _, p := range net.places {
		for j, pp := range another.places {
			if pairedPs[j] { // j already paired with prev i
				continue
			}
			if p.Equals(pp) {
				pairedPs[j] = true
				continue pairPlace
			}
		}
		return false, fmt.Errorf("no matching place") // no j pair for i
	}

	pairedTs := map[int]bool{}
pairTran:
	for _, t := range net.transitions {
		for j, tt := range another.transitions {
			if pairedTs[j] {
				continue
			}
			if t.Equals(tt) {
				pairedTs[j] = true
				continue pairTran
			}
		}
		return false, fmt.Errorf("no matching transition")
	}

	return true, nil
}

func (net *Net) saveState() {
	for _, place := range net.places {
		place.initTokens = place.Tokens
	}
}

func (net *Net) restoreState() {
	for _, place := range net.places {
		place.Tokens = place.initTokens
	}
}

/* Place */

type Place struct {
	Tokens      int
	Description string
	Id          string
	initTokens  int
}

func (p Place) String() string {
	return fmt.Sprintf("%s(%d)%s", p.Id, p.Tokens, p.Description)
}

func (p *Place) Equals(pp *Place) bool {
	if p.Tokens != pp.Tokens {
		return false
	}
	// if p.Description != pp.Description {
	// 	return false
	// }
	return true
}

func (p *Place) Hidden() bool {
	return p.Id == "."
}

/* Places */

type Places []*Place

func (places *Places) Push(place *Place) {
	*places = append(*places, place)
}

func (places Places) String() string {
	placestrs := make([]string, 0, len(places))
	for _, place := range places {
		placestrs = append(placestrs, place.String())
	}
	return strings.Join(placestrs, ", ")
}

func (places Places) Find(id string) *Place {
	for _, place := range places {
		if place.Id == id {
			return place
		}
	}
	return nil
}

/* Arc */

type Arc struct {
	Weight int
	Place  *Place
}

func (arc Arc) String() string {
	if arc.Weight > 1 {
		return fmt.Sprintf("%d*%s", arc.Weight, arc.Place.Id)
	} else {
		return arc.Place.Id
	}
}

func (a *Arc) Equals(aa *Arc) bool {
	if a.Weight != aa.Weight {
		return false
	}
	if !a.Place.Equals(aa.Place) {
		return false
	}
	return true
}

/* Arcs */

type Arcs []*Arc

func (arcs Arcs) String() string {
	arcsstr := make([]string, 0, len(arcs))
	for _, arc := range arcs {
		arcsstr = append(arcsstr, arc.String())
	}
	return strings.Join(arcsstr, ", ")
}

func (arcs *Arcs) Push(w int, place *Place) {
	*arcs = append(*arcs, &Arc{w, place})
}

func (a *Arcs) Equals(aa *Arcs) bool {
	pairedAs := map[int]bool{}

pairing:
	for _, arc := range *a {
		for j, another := range *aa {
			if pairedAs[j] {
				continue
			}
			if arc.Equals(another) {
				pairedAs[j] = true
				continue pairing
			}
		}
		return false
	}
	return true
}

/* Transtition */

type Transition struct {
	Id          string
	Origins     Arcs
	Targets     Arcs
	Priority    int
	TimeFunc    *TimeFunc
	Description string
}

func (t Transition) String() string {
	prio := ""
	if t.Priority != 0 {
		prio = "p=" + strconv.Itoa(t.Priority)
	}
	return fmt.Sprintf("%s -> %s[%s%s]%s -> %s", t.Origins, t.Id, t.TimeFunc, prio, t.Description, t.Targets)
}

func (t *Transition) Equals(tt *Transition) bool {
	if !t.Origins.Equals(&tt.Origins) {
		return false
	}
	if !t.Targets.Equals(&tt.Targets) {
		return false
	}
	if t.Priority != tt.Priority {
		return false
	}
	if t.Description != tt.Description {
		return false
	}

	return true
}

/**
 * How many times can by transition fired with current marking on origins arcs
 */
func (t *Transition) getEnabilityMagnitude() int {
	enability := MaxInt
	for _, arc := range t.Origins {
		arcEnability := arc.Place.Tokens / arc.Weight // posible fires for this arc
		if arcEnability < enability {
			enability = arcEnability
		}
	}
	return enability
}

func (t *Transition) isEnabled() bool {
	for _, arc := range t.Origins {
		if arc.Place.Tokens < arc.Weight {
			return false
		}
	}
	return true
}

func (t *Transition) doIn() {
	for _, arc := range t.Origins {
		arc.Place.Tokens -= arc.Weight
		if arc.Place.Tokens < 0 {
			panic("impossible transition done")
		}
	}
}

func (t *Transition) doOut() {
	for _, arc := range t.Targets {
		arc.Place.Tokens += arc.Weight
	}
}

/* Transitions */

type Transitions []*Transition

func (trans *Transitions) Push(tran *Transition) {
	*trans = append(*trans, tran)
}

func (trans *Transitions) Remove(i int) {
	*trans = append((*trans)[:i], (*trans)[i+1:]...)
}

/* following 3 methods are implemented to satisfy sort.Interface */
func (trans Transitions) Len() int {
	return len(trans)
}

func (trans Transitions) Swap(i, j int) {
	trans[i], trans[j] = trans[j], trans[i]
}

func (trans Transitions) Less(i, j int) bool {
	// timed have always lower priority
	if trans[i].TimeFunc == nil && trans[j].TimeFunc != nil {
		return true
	}
	if trans[j].TimeFunc == nil && trans[i].TimeFunc != nil {
		return false
	}
	// if both or niether are timed copmpare priority
	return trans[i].Priority > trans[j].Priority
}
