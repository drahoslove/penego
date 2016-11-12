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

func compileRegExps () {

	const (
		SP = `[ \t]*`
		ID = `[a-zA-Z][a-zA-Z0-9_]*`
		NUM = `(0|([1-9][0-9]*))`
		STR = `"([^"]*")`
		CMNT = `//.*`
		IDS = SP+ID+SP+`(,`+SP+ID+SP+`)*`
		PRIO = `(p=(?P<prio>`+NUM+`))`
		TIME = `(`+NUM+`)([smhd]|(ms)|(us))?`
		FIX = `(?P<fix>`+TIME+`)`
		UNIF = `(?P<unf0>`+TIME+`)-(?P<unf1>`+TIME+`)`
		EXP = `exp\((?P<exp>`+TIME+`\))`
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
		`((?P<in>`+IDS+`)->)?`,
		`\[`,	// [
		`(?P<attr>`+ATTR+`)?`,
		`\]`,	// ]
		`(?P<desc>`+STR+`)?`,
		`(->(?P<out>`+IDS+`))?`,
		`(`+CMNT+`)?`,
		`$`,
	}, SP)

	/** compile regexps **/

	placeRE = regexp.MustCompile(placeREstr)
	transitionRE = regexp.MustCompile(transitionREstr)
	emptyLineRE = regexp.MustCompile(`^`+SP+`(`+CMNT+`)?$`)
	timeRE = regexp.MustCompile(TIME)

}



func Parse(input string) (transitions Transitions, places Places, err error) {

	if emptyLineRE == nil || placeRE == nil || transitionRE == nil {
		compileRegExps()
	}

	places = Places{}
	transitions = Transitions{}

	lines := strings.Split(input, "\n")

	namedPlaces := make(map[string]*Place)


	/* ----------- parse places ----------- */

	for i, line := range lines {
		if isPlaceDefinition(line) {

			id := getSubmatchString(placeRE, line, "id")
			num, _ := strconv.Atoi(getSubmatchString(placeRE, line, "num"))
			desc := getSubmatchString(placeRE, line, "desc")

			if _, ok := namedPlaces[id]; ok {
				err = errors.New("place with id >"+id+"< already defined")
				return
			}
			place := &Place{
				Tokens: num,
				Description: desc,
			}
			namedPlaces[id] = place
			places.Push(place)

		} else {
			if !isEmptyLine(line) && !isTransitionDefinition(line) {
				err = errors.New("syntax error at lin " + strconv.Itoa(i))
				return
			}
		}
	}


	/* ----------- parse transitions ----------- */

	for i, line := range lines {
		if isTransitionDefinition(line) {

			listin := strings.Split(getSubmatchString(transitionRE, line, "in"),",")
			listout := strings.Split(getSubmatchString(transitionRE, line, "out"), ",")
			attr := getSubmatchString(transitionRE, line, "attr") // TODO
			desc := getSubmatchString(transitionRE, line, "desc")


			getPlacesByList := func(list []string) Places {
				places := Places{}
				for _, id := range list {
					id := strings.TrimSpace(id)
					if id == "" {
						continue
					}
					if place, ok := namedPlaces[id]; !ok {
						err = errors.New("undefined place id >"+id+"< used in transition")
						return places
					} else {
						places.Push(place)
					}
				}
				return places
			}

			origins := getPlacesByList(listin)
			targets := getPlacesByList(listout)

			priority := 0
			var timeFunc TimeFunc = nil

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
					case unf0 != "":
						timeFunc = GetUniformTimeFunc(parseTime(unf0), parseTime(unf1))
					case exp != "":
						timeFunc = GetExponentialTimeFunc(parseTime(exp))
				}
			}


			transitions.Push(Transition{
				Origins: origins,
				Targets: targets,
				Priority: priority,
				TimeFunc: timeFunc,
				Description: desc,
			})

		} else {
			if !isEmptyLine(line) && !isPlaceDefinition(line) {
				err = errors.New("syntax error at line " + strconv.Itoa(i))
				return
			}
		}

	}


	return transitions, places, err
}

func isPlaceDefinition(line string) bool {
	return placeRE.FindStringSubmatch(line) != nil
}

func isTransitionDefinition(line string) bool {
	return transitionRE.FindStringSubmatch(line) != nil
}

func isEmptyLine(line string) bool {
	return len(emptyLineRE.FindAllString(line, 1)) != 0
}

func getSubmatchString(re *regexp.Regexp, input string, name string) string {
	return re.ReplaceAllString(input, "${"+name+"}")
}

func parseTime(tstr string) time.Duration {
	match := timeRE.FindStringSubmatch(tstr)
	timeInt, _ := strconv.Atoi(match[1])
	timeUnit := match[2]

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