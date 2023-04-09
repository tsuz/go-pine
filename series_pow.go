package pine

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
)

// Sum generates a ValueSeries of summation of previous values
//
// Parameters
// p - ValueSeries: source data
// l - lookback: lookback periods [1, ∞)
//
// Example:
// t=time.Time       | 1     |  2    | 3     |
// p=ValueSeries     | 13    | 15    | 11    |
// pow(0.5)     	 | 3.606 | 3.873 | 3.317 |
// pow(2)       	 | 169   | 225   | 121   |
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

	setCache(key, pow)

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
