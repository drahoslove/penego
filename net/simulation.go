package net

import (
	"fmt"
	"sort"
	"time"
)

/* Event */

type Event struct {
	time       time.Duration
	transition *Transition
}

/* Calendar */

type Calendar []Event

func (c Calendar) String() string {
	str := "c: "
	for _, event := range c {
		str += fmt.Sprintf("T=%s,%s | ", event.time, event.transition.Description)
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
	return (*c)[0].time, (*c)[0].transition
}

func (c *Calendar) insertByTime(newTime time.Duration, tran *Transition) {
	if c.isEmpty() {
		c.Insert(Event{newTime, tran}, 0)
		return
	}
	i, event := 0, Event{}
	for i, event = range *c {
		if newTime < event.time {
			c.Insert(Event{newTime, tran}, i)
			return
		}
	}
	// not found, new is biggest
	c.Insert(Event{newTime, tran}, i+1)
}

/* Simulation */

type Simulation struct {
	startTime         time.Duration
	endTime           time.Duration
	now               time.Duration
	net               Net
	calendar          Calendar
	stateChange       func(time.Duration, time.Duration)
	paused            bool
	stopped           bool
	sortedTransitions Transitions
}

func NewSimulation(startTime, endTime time.Duration, net Net) Simulation {
	net.saveState()
	return Simulation{startTime, endTime, 0, net, Calendar{}, nil, false, false, nil}
}

/**
 * Get current time of simulation
 */
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
				sim.calendar.insertByTime(sim.now+(*tran.TimeFunc)(), tran)
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
	for i := len(sim.calendar) - 1; i >= 0; i-- { // TODO go through in random order each time
		tran := sim.calendar[i].transition
		if sub, ok := subtractions[tran]; ok && sub < 0 {
			sim.calendar.Remove(i)
			subtractions[tran]++
		}
	}
}

func (sim *Simulation) fireEvent(scheduledTran *Transition, before, now time.Duration) {
	scheduledTran.doIn()
	sim.cancelUnenabledTimed()
	scheduledTran.doOut()
	sim.stateChange(before, now)

	countOfPasses := 0
stabilize: // whenever transition is completed, start checking again from bigest priority
	countOfPasses++
	if countOfPasses > 1E3 {
		panic("too many transitions done in same time, possible loop")
	}
	for _, tran := range sim.sortedTransitions { // TODO cycle transitions with same priority in random order
		if tran.TimeFunc != nil {
			break // no need to go further, rest are timed due to sort
		}
		if sim.stopped {
			return
		}
		if tran.isEnabled() {
			if sim.paused {
				sim.calendar.Insert(Event{now, tran}, 0)
				return
			}
			tran.doIn()
			sim.cancelUnenabledTimed()
			tran.doOut()
			sim.stateChange(now, now)
			goto stabilize
		}
	}

	sim.scheduleEnabledTimed() // might create new event in current time
}

func (sim *Simulation) DoEveryStateChange(fun func(time.Duration, time.Duration)) {
	sim.stateChange = func(now, then time.Duration) {
		if fun != nil {
			fun(now, then)
		}
	}
}

func (sim *Simulation) Init() {
	restartSeed()
	sim.now = sim.startTime
	sim.calendar = Calendar{}
	sim.sortedTransitions = make(Transitions, len(sim.net.transitions))
	copy(sim.sortedTransitions, sim.net.transitions)
	sort.Sort(sim.sortedTransitions)

	// schedule empty tran
	sim.calendar.Insert(Event{sim.startTime, &Transition{}}, 0)
	// sim.Step() // this causes runtime error
}

func (sim *Simulation) Step() bool {
	if sim.calendar.isEmpty() {
		return false
	}
	eventTime, tranToFireNow := sim.calendar.shift()
	if eventTime > sim.endTime {
		return false
	}
	sim.stopped = false
	before := sim.now
	sim.now = eventTime
	sim.fireEvent(tranToFireNow, before, sim.now) // current time and time of event
	return true
}

// Run starts running simulation or continue in running from previously paused state
func (sim *Simulation) Run() {

	if !sim.paused { // start from beginning
		sim.Init()
	}

	sim.paused = false
	sim.stopped = false

	for {
		if sim.paused || sim.stopped {
			break
		}
		if !sim.Step() {
			break
		}
	}
}

// Pause pauses current simulation Run
// Intended to be called simultaneously with simulation.Run
func (sim *Simulation) Pause() {
	sim.paused = true
	sim.stopped = false
}

// Stops simulation and restore its initial state
func (sim *Simulation) Stop() {
	sim.net.restoreState()
	sim.paused = false
	sim.stopped = true
}
