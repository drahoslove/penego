package net

import (
	"fmt"
	"time"
	"sort"
	"math"
	"strings"
	"strconv"
)


/******* types *******/

/* Net */

type Net struct {
	places Places
	transitions Transitions
}

func New(places Places, transitions Transitions) Net {
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
	Tokens int
	Description string
	id string
	initTokens int
}

func (p Place) String () string {
	return fmt.Sprintf("%s(%d)%s", p.id, p.Tokens, p.Description)
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


/* Arc */

type Arc struct {
	Weight int
	Place *Place
}

func (arc Arc) String() string {
	if arc.Weight > 1 {
		return fmt.Sprintf("%d*%s", arc.Weight, arc.Place.id)
	} else {
		return arc.Place.id
	}
}

/* Arcs */

type Arcs []Arc

func (arcs Arcs) String() string {
	arcsstr := make([]string, 0, len(arcs))
	for _, arc := range arcs {
		arcsstr = append(arcsstr, arc.String())
	}
	return strings.Join(arcsstr, ", ")
}

func (arcs *Arcs) Push(w int, place *Place) {
	*arcs = append(*arcs, Arc{w, place})
}

/* Transtition */

type Transition struct {
	Origins Arcs
	Targets Arcs
	Priority int
	TimeFunc *TimeFunc
	Description string
}

func (t Transition) String () string {
	prio := ""
	if t.Priority != 0 {
		prio = "p=" + strconv.Itoa(t.Priority)
	}
	return fmt.Sprintf("%s -> [%s%s]%s -> %s", t.Origins, t.TimeFunc, prio, t.Description, t.Targets)
}

/**
 * How many times can by transition fired with current marking on origins arcs
 */
func (t * Transition) getEnabilityMagnitude() int {
	enability := math.MaxInt64
	for _, arc := range t.Origins {
		arcEnability := arc.Place.Tokens / arc.Weight // posible fires for this arc
		if arcEnability < enability {
			enability = arcEnability
		}
	}
	return enability
}

func (t * Transition) isEnabled() bool {
	for _, arc := range t.Origins {
		if arc.Place.Tokens < arc.Weight {
			return false
		}
	}
	return true
}

func (t * Transition) doIn() {
	for _, arc := range t.Origins {
		arc.Place.Tokens -= arc.Weight
		if arc.Place.Tokens < 0 {
			panic("impossible transition done")
		}
	}
}

func (t * Transition) doOut() {
	for _, arc := range t.Targets {
		arc.Place.Tokens += arc.Weight
	}
}


/* Transitions */

type Transitions []*Transition

func (trans *Transitions) Push(tran Transition) {
	*trans = append(*trans, &tran)
}

func (trans *Transitions) Remove(i int) {
	*trans = append((*trans)[:i], (*trans)[i+1:]...)
}

func (trans Transitions) Len() int {
	return len(trans)
}

func (trans Transitions) Swap(i, j int) {
	trans[i], trans[j] = trans[j], trans[i]
}

func (trans Transitions) Less(i, j int) bool {
	if trans[i].TimeFunc == nil && trans[j].TimeFunc != nil {
		return true
	}
	if trans[j].TimeFunc == nil && trans[i].TimeFunc != nil {
		return false
	}
	return trans[i].Priority > trans[j].Priority
}


/* Event */

type Event struct {
	t time.Duration
	transition *Transition
}


/* Calendar */

type Calendar []Event

func (c Calendar) String() string {
	str := "c: "
	for _, event := range c {
		str += fmt.Sprintf("T=%s,%s | ", event.t, event.transition.Description)
	}
	return str
}

func (c *Calendar) Insert(event Event, i int) {
	*c = append((*c)[:i], append([]Event{event}, (*c)[i:]...)...)
}

func (c *Calendar) Remove(i int) {
	*c = append((*c)[:i], (*c)[i+1:]...)
}

func (c *Calendar) isEmpty() bool {
	return len(*c) == 0
}

func (c *Calendar) shift() (time.Duration, *Transition) {
	defer func() {
		*c = (*c)[1:]
	}()
	return (*c)[0].t, (*c)[0].transition
}

func (c *Calendar) insert(newTime time.Duration, tran *Transition) {
	if c.isEmpty() {
		c.Insert(Event{newTime, tran}, 0)
		return
	}
	i, event := 0, Event{};
	for i, event = range *c {
		if newTime < event.t {
			c.Insert(Event{newTime, tran}, i)
			return
		}
	}
	// not found, new is biggest
	c.Insert(Event{newTime, tran}, i+1)
}


/* Simulation */

type Simulation struct {
	startTime time.Duration
	endTime time.Duration
	now time.Duration
	net Net
	calendar Calendar
	stateChange func(time.Duration, time.Duration)
	paused bool
}

func (sim *Simulation) GetNow() time.Duration {
	return sim.now
}

/**
 * Check how much is enabled and how many times is already scheduled
 * and return difference
 * positive number means how many event should be scheduled
 * negative number means how many scheduled event should be canceled
 */
func (sim *Simulation) diffEnabilityVsScheduled(transition *Transition) int {
	enability := transition.getEnabilityMagnitude()
	eventCount := 0
	for _, event := range sim.calendar {
		if event.transition == transition {
			eventCount++
		}
	}
	return enability - eventCount
}

func (sim *Simulation) scheduleEnabledTimed() {
	for _, tran := range sim.net.transitions {
		if tran.TimeFunc != nil {
			max := sim.diffEnabilityVsScheduled(tran) // how many times schedule
			for i := 0; i < max; i++ {
				sim.calendar.insert(sim.now + (*tran.TimeFunc)(), tran)
			}
		}
	}
}

func (sim *Simulation) cancelUnenabledTimed() {
	subtractions := map[*Transition]int{}
	for _, tran := range sim.net.transitions {
		subtractions[tran] = sim.diffEnabilityVsScheduled(tran)
	}

	// remove excess
	for i := len(sim.calendar)-1; i >= 0; i-- { // TODO go through in random order each time
		tran := sim.calendar[i].transition
		if sub, ok := subtractions[tran]; ok && sub < 0 {
			sim.calendar.Remove(i)
			subtractions[tran]++
		}
	}
}

func (sim *Simulation) DoEveryStateChange(fun func(time.Duration, time.Duration)) {
	sim.stateChange = func(now, then time.Duration) {
		if fun != nil {
			fun(now, then)
		}
	}
}

// Run starts running simulation or continue in running from previously paused state
func (sim *Simulation) Run() {

	if !sim.paused {
		restartSeed()
		sim.now = sim.startTime
		sim.calendar = Calendar{}
		sim.net.restoreState()
	} // else use previous

	sim.paused = false


	sortedTransitions := make(Transitions, len(sim.net.transitions))
	copy(sortedTransitions, sim.net.transitions)
	sort.Sort(sortedTransitions)

	fireEvent := func(scheduledTran *Transition) {

		scheduledTran.doIn()
		sim.cancelUnenabledTimed()
		scheduledTran.doOut()

		countOfPasses := 0
		stabilize: // whenever transition is completed, start checking again from bigest priority
		countOfPasses++
		if countOfPasses > 1E6 {
			panic("too many transitions done in same time, possible loop")
		}
		for _, tran := range sortedTransitions { // TODO cycle transitions with same priority in random order
			if tran.TimeFunc != nil {
				break // no need to go further, rest are timed due to sort
			}
			if tran.isEnabled() {
				sim.stateChange(sim.now, sim.now)
				tran.doIn()
				sim.cancelUnenabledTimed()
				tran.doOut()
				goto stabilize
			}
		}

		sim.scheduleEnabledTimed() // might create new event in current time
	}

	fireEvent(&Transition{})

	for !sim.calendar.isEmpty() {
		if sim.paused {
			return
		}
		eventTime, tranToFireNow := sim.calendar.shift()
		sim.stateChange(sim.now, eventTime) // current time and next time
		if eventTime > sim.endTime {
			break
		}
		sim.now = eventTime
		fireEvent(tranToFireNow)

	}

}

// Pause pauses current simulation Run
// Intended to be called simultaneously with simulation.Run
func (sim *Simulation) Pause() {
	sim.paused = true
}


/******* exported functions *******/

func NewSimulation(startTime, endTime time.Duration, net Net) Simulation {
	net.saveState()
	return Simulation{startTime, endTime, 0, net, Calendar{}, nil, false}
}
