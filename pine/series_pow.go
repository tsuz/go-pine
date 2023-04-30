package pine

import (
	"fmt"
	"math"
)

// Pow generates a ValueSeries of values from power function
//
// Parameters
//   - p - ValueSeries: source data
//   - exp - float64: exponent of the power function
func Pow(src ValueSeries, exp float64) ValueSeries {
	key := fmt.Sprintf("pow:%s:%.8f", src.ID(), exp)
	pow := getCache(key)
	if pow == nil {
		pow = NewValueSeries()
	}

	// current available value
	stop := src.GetCurrent()
	if stop == nil {
		return pow
	}

	pow = getPow(*stop, pow, src, exp)
	// disable this for now
	// setCache(key, pow)

	pow.SetCurrent(stop.t)

	return pow
}

func getPow(stop Value, pow ValueSeries, src ValueSeries, exp float64) ValueSeries {

	var startNew *Value

	lastAvail := pow.GetLast()

	if lastAvail == nil {
		startNew = src.GetFirst()
	} else {
		v := src.Get(lastAvail.t)
		startNew = v.next
	}

	if startNew == nil {
		// if nothing is to start with, then nothing can be done
		return pow
	}

	// first new time
	itervt := startNew.t

	for {
		v := src.Get(itervt)
		if v == nil {
			break
		}
		newpow := math.Pow(v.v, exp)
		pow.Set(itervt, newpow)

		if v.next == nil {
			break
		}
		if v.t.Equal(stop.t) {
			break
		}

		itervt = v.next.t
	}

	return pow
}
