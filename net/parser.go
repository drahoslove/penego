package net

import (
	"regexp"
	"strings"
	"strconv"
	"time"
	"errors"
)


var (
	placeRE *regexp.Regexp
	transitionRE *regexp.Regexp
	emptyLineRE *regexp.Regexp
	timeRE *regexp.Regexp
)

func init () {

	const (
		SP = `[ \t]*`
		ID = `[a-zA-Z][a-zA-Z0-9_]*`
		NUM = `(0|([1-9][0-9]*))`
		STR = `"[^"]*"`
		CMNT = `((//)|(--)).*`
		ARC = SP+`(`+NUM+SP+`\*`+SP+`)?`+ID+SP
		ARCS = ARC+`(,`+ARC+`)*`
		PRIO = `p=(?P<prio>`+NUM+`)`
		TIME = `(?P<t>`+NUM+`)(?P<u>[smhd]|(ms)|(us))?`
		FIX = `(?P<fix>`+TIME+`)`
		UNIF = `(?P<unf0>`+TIME+`)(-|(..))(?P<unf1>`+TIME+`)`
		EXP = `exp\((?P<exp>`+TIME+`)\)`
		ATTR = `(`+PRIO+`)|(`+FIX+`)|(`+UNIF+`)|(`+EXP+`)`
	)


	/** prepare regexps strings **/

	// ID ( NUM? ) STR?
	placeREstr := strings.Join([]string{
		`^`,
		`(?P<id>`+ID+`)`,
		`\(`,
		`(?P<num>`+NUM+`)?`,
		`\)`,
		`(?P<desc>`+STR+`)?`,
		`(`+CMNT+`)?`,
		`$`,
	}, SP)

	// IDS -> [ ATTR? ] STR? -> IDS
	transitionREstr := strings.Join([]string{
		`^`,
		`((?P<in>`+ARCS+`)->)?`,
		`\[`,	// [
		`(?P<attr>`+ATTR+`)?`,
		`\]`,	// ]
		`(?P<desc>`+STR+`)?`,
		`(->(?P<out>`+ARCS+`))?`,
		`(`+CMNT+`)?`,
		`$`,
	}, SP)

	/** compile regexps **/

	placeRE = regexp.MustCompile(placeREstr)
	transitionRE = regexp.MustCompile(transitionREstr)
	emptyLineRE = regexp.MustCompile(`^`+SP+`(`+CMNT+`)?$`)
	timeRE = regexp.MustCompile(TIME)

}


func Parse(input string) (net Net, err error) {

	net.places = Places{}
	net.transitions = Transitions{}

	lines := strings.Split(input, "\n")

	namedPlaces := make(map[string]*Place)


	/* ----------- parse places ----------- */

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if isPlaceDefinition(line) {

			id := getSubmatchString(placeRE, line, "id")
			num, _ := strconv.Atoi(getSubmatchString(placeRE, line, "num"))
			desc := getSubmatchString(placeRE, line, "desc")

			if _, exists := namedPlaces[id]; exists {
				err = errors.New("place with id `"+id+"` is already defined")
				return
			}
			place := &Place{
				Tokens: num,
				Description: unPack(desc), // strip first and last char
				id: id,
			}
			namedPlaces[id] = place
			net.places.Push(place)

		} else {
			if !isEmptyLine(line) && !isTransitionDefinition(line) {
				err = errors.New("syntax error at line " + strconv.Itoa(i))
				return
			}
		}
	}


	/* ----------- parse transitions ----------- */

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if isTransitionDefinition(line) {

			listin := getSubmatchString(transitionRE, line, "in")
			listout := getSubmatchString(transitionRE, line, "out")
			attr := getSubmatchString(transitionRE, line, "attr")
			desc := getSubmatchString(transitionRE, line, "desc")


			getArcsByList := func(list string) Arcs {
				arcs := Arcs{}
				list = strings.TrimSpace(list)
				if list == "" {
					return arcs
				}
				for _, pair := range strings.Split(list, ",") {
					pair := strings.Split(strings.TrimSpace(pair), "*")
					id := strings.TrimSpace(pair[len(pair)-1])
					w := 1
					if len(pair) == 2 {
						w, _ = strconv.Atoi(strings.TrimSpace(pair[0]))
					}
					if place, exists := namedPlaces[id]; !exists {
						err = errors.New("undefined place id `"+id+"` used in transition")
						return arcs
					} else {
						for _, arc := range arcs {
							if arc.Place == place {
								err = errors.New("place `"+place.id+"` used multiple times in one side of transition")
							}
						}
						arcs.Push(w, place)
					}
				}
				return arcs
			}

			origins := getArcsByList(listin)
			targets := getArcsByList(listout)

			// changes `[] -> n` to `S -> [] -> n,S`
			// where S is hidden place creating self loop
			if len(origins) == 0 {
				selfLoopPlace := &Place{Tokens: 1, id: "."}
				origins.Push(1, selfLoopPlace)
				targets.Push(1, selfLoopPlace)
			}

			priority := 0
			var timeFunc *TimeFunc

			if attr != "" {
				prio := getSubmatchString(transitionRE, line, "prio")
				fix := getSubmatchString(transitionRE, line, "fix")
				unf0 := getSubmatchString(transitionRE, line, "unf0")
				unf1 := getSubmatchString(transitionRE, line, "unf1")
				exp := getSubmatchString(transitionRE, line, "exp")
				switch {
					case prio != "":
						priority, _ = strconv.Atoi(prio)
					case fix != "":
						timeFunc = GetConstantTimeFunc(parseTime(fix))
					case unf0 != "" && unf1 != "":
						timeFunc = GetUniformTimeFunc(parseTime(unf0), parseTime(unf1))
					case exp != "":
						timeFunc = GetExponentialTimeFunc(parseTime(exp))
				}
			}


			net.transitions.Push(Transition{
				Origins: origins,
				Targets: targets,
				Priority: priority,
				TimeFunc: timeFunc,
				Description: unPack(desc),
			})

		} else {
			if !isEmptyLine(line) && !isPlaceDefinition(line) {
				err = errors.New("syntax error at line " + strconv.Itoa(i))
				return
			}
		}

	}

	return net, err
}

func isPlaceDefinition(line string) bool {
	return placeRE.MatchString(line)
}

func isTransitionDefinition(line string) bool {
	return transitionRE.MatchString(line)
}

func isEmptyLine(line string) bool {
	return emptyLineRE.MatchString(line)
}

func getSubmatchString(re *regexp.Regexp, input string, name string) string {
	return re.ReplaceAllString(input, "${"+name+"}")
}

func parseTime(tstr string) time.Duration {
	timeStr := getSubmatchString(timeRE, tstr, "t")
	timeInt, _ := strconv.Atoi(timeStr)
	timeUnit := getSubmatchString(timeRE, tstr, "u")

	switch timeUnit {
	case "":
		return time.Duration(timeInt)
	case "d":
		return time.Duration(timeInt) * time.Hour * 24
	default:
		t, _ := time.ParseDuration(tstr)
		return t
	}
}

func unPack(str string) string {
	if len(str) > 2 {
		return string(str[1:len(str)-1])
	}
	return str
}