package net

import(
	"time"
	"math"
	"math/big"
	"math/rand"
	truerand "crypto/rand"
	"strings"
)


/******* types *******/

/* TimeFunc */

type TimeFunc func() time.Duration

func (fn *TimeFunc) String() string {
	if repr, ok := timeFuncTextReprs[fn]; ok {
		return repr
	} else {
		return ""
	}
}

func (fn *TimeFunc) SetTextRepr(name string, args... time.Duration) {

	arguments := make([]string,0)

	for _, arg := range args {
		arguments = append(arguments, trimZeroUnits(arg.String()))
	}

	timeFuncTextReprs[fn] = func() string {
		switch name {
		case "const":
			return arguments[0]
		case "unif":
			return arguments[0] + ".." + arguments[1]
		case "exp":
			return "exp(" + arguments[0] + ")"
		default:
			return name + "(" + strings.Join(arguments, ",") + ")"
		}
	}()
}

/******* global vars *******/

var timeFuncTextReprs map[*TimeFunc] string
var startSeed int64 = 1


/******* exported functions *******/

/* timeFunc factories */

func GetConstantTimeFunc(duration time.Duration) *TimeFunc {
	fn := TimeFunc(func() time.Duration {
		return duration
	})
	fn.SetTextRepr("const", duration)
	return &fn
}

func GetUniformTimeFunc(from, to time.Duration) *TimeFunc {
	if from > to {
		from, to = to, from
	}
	fn := TimeFunc(func() time.Duration {
		return uniformTime(from, to)
	})
	fn.SetTextRepr("unif", from, to)
	return &fn
}

func GetExponentialTimeFunc(mean time.Duration) *TimeFunc {
	fn := TimeFunc(func() time.Duration {
		return exponentialTime(mean)
	})
	fn.SetTextRepr("exp", mean)
	return &fn
}


/**
 * Seed pseudo random generator with true random number.
 * This same seed is used at beginning of every simulation.Run()
 */
func TrueRandomSeed() {
	max := big.NewInt(math.MaxInt32)
	seed, _ := truerand.Int(truerand.Reader, max)
	startSeed = seed.Int64()
	rand.Seed(startSeed)
}


/******* unexported functions *******/

func init () {
	timeFuncTextReprs = make(map[*TimeFunc]string)
}

func restartSeed() {
	rand.Seed(startSeed)
}

func trimZeroUnits(input string) string {
	return strings.Replace(strings.Replace(input, "m0s", "m", 1), "h0m", "h", 1)
}

/* random functions*/

func uniformTime(from, to time.Duration) time.Duration {
	return from + time.Duration(rand.Int63n(int64(to-from)))
}

func exponentialTime(mean time.Duration) time.Duration {
	return time.Duration(rand.ExpFloat64() * float64(mean))
}