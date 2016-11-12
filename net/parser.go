package net

import (
	"fmt"
	"regexp"
	"strings"
	"strconv"
	"errors"
)


var (
	placeRE *regexp.Regexp
	transitionRE *regexp.Regexp
	emptyLineRE *regexp.Regexp
)

func compileRegExps () {

	const (
		SP = `[ \t]*`
		ID = `[a-zA-Z][a-zA-Z0-9_]*`
		NUM = `(0|([1-9][0-9]*))`
		STR = `"([^"]*")`
		CMNT = `//.*`
		IDS = SP+ID+SP+`(,`+SP+ID+SP+`)*`
		PRIO = `p=`+NUM
		TIME = NUM+`[smhd]?`
		EXP = `exp\(`+TIME+`\)`
		UNIF = TIME+`-`+TIME
		ATTR = `(`+PRIO+`)|(`+TIME+`)|(`+UNIF+`)|(`+EXP+`)`
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

	// IDS -> [] STR? -> IDS
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
	emptyLineRE = regexp.MustCompile(SP+`(`+CMNT+`)?`)

}



func Parse(input string) (transitions Transitions, places Places, err error) {

	if emptyLineRE == nil || placeRE == nil || transitionRE == nil {
		compileRegExps()
	}


	lines := strings.Split(input, "\n")


	/* ----------- parse places ----------- */

	namedPlaces := make(map[string]*Place)

	for _, line := range lines {
		if matches := placeRE.FindStringSubmatch(line); matches != nil {

			id := getSubmatchString(placeRE, line, "id")
			num, _ := strconv.Atoi(getSubmatchString(placeRE, line, "num"))
			desc := getSubmatchString(placeRE, line, "desc")

			if _, ok := namedPlaces[id]; ok {
				err = errors.New("id already used")
				return
			}

			namedPlaces[id] = &Place{
				Tokens: num,
				Description: desc,
			}
			fmt.Println(id, num, desc)
		} else {
			if emptyLineRE.FindAllString(line, -1) == nil {
				err = errors.New("syntax error")
				return
			}
		}
	}


	/* ----------- parse transitions ----------- */

	transitions = Transitions{}

	for _, line := range lines {
		if matches := transitionRE.FindStringSubmatch(line); matches != nil {

			listin := getSubmatchString(transitionRE, line, "in")
			listout := getSubmatchString(transitionRE, line, "out")
			attr := getSubmatchString(transitionRE, line, "attr")
			desc := getSubmatchString(transitionRE, line, "desc")

			fmt.Printf("%s->[%s]%s->%s\n", listin, attr, desc, listout)

		} else {
			if emptyLineRE.FindAllString(line, -1) == nil {
				err = errors.New("syntax error")
				return
			}
		}

	}


	return transitions, places, err
}


func getSubmatchString(re *regexp.Regexp, input string, name string) string {
	return re.ReplaceAllString(input, "${"+name+"}")
}
