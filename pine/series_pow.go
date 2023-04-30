package pine

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
)

// Pow generates a ValueSeries of values from power function
//
// Parameters
//   - p - ValueSeries: source data
//   - exp - float64: exponent of the power function
func Pow(src ValueSeries, exp float64) (ValueSeries, error) {
	var err error
	key := fmt.Sprintf("pow:%s:%.8f", src.ID(), exp)
	pow := getCache(key)
	if pow == nil {
		pow = NewValueSeries()
	}

	// current available value
	stop := src.GetCurrent()
	if stop == nil {
		return pow, nil
	}

	pow, err = getPow(*stop, pow, src, exp)
	if err != nil {
		return pow, errors.Wrap(err, "error getsum")
	}

	// disable this for now
	// setCache(key, pow)

	pow.SetCurrent(stop.t)

	return pow, nil
}

func getPow(stop Value, pow ValueSeries, src ValueSeries, exp float64) (ValueSeries, error) {

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
		return pow, nil
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

	return pow, nil
}
