package pine

import (
	"fmt"

	"github.com/pkg/errors"
)

// Sum generates a ValueSeries of summation of previous values
//
// Parameters
// p - ValueSeries: source data
// l - lookback: lookback periods [1, ∞)
//
// Example:
// t=time.Time       | 1   |  2  | 3    | 4    | 5  |
// p=ValueSeries     | 13  | 15  | 11   | 19   | 21 |
// sum(p, 3)	     | nil | nil | 39   | 45   | 51 |
func Sum(p ValueSeries, l int) (ValueSeries, error) {
	var err error
	key := fmt.Sprintf("sum:%s:%d", p.ID(), l)
	sum := getCache(key)
	if sum == nil {
		sum = NewValueSeries()
	}

	// current available value
	stop := p.GetCurrent()
	if stop == nil {
		return sum, nil
	}

	sum, err = getSum(*stop, sum, p, l)
	if err != nil {
		return sum, errors.Wrap(err, "error getsum")
	}

	setCache(key, sum)

	sum.SetCurrent(stop.t)

	return sum, nil
}

func getSum(stop Value, sum ValueSeries, src ValueSeries, l int) (ValueSeries, error) {

	// keep track of the source values of sum, maximum of l+1 items
	sumSrc := make([]float64, 0)
	var startNew *Value

	lastAvail := sum.GetLast()

	if lastAvail == nil {
		startNew = src.GetFirst()
	} else {
		v := src.Get(lastAvail.t)
		startNew = v.next
	}

	if startNew == nil {
		// if nothing is to start with, then nothing can be done
		return sum, nil
	}

	// populate source values to be summed
	if lastAvail != nil {
		lastAvailv := src.Get(lastAvail.t)

		for {
			if lastAvailv == nil {
				break
			}

			srcv := src.Get(lastAvailv.t)
			// add at the beginning since we go backwards
			sumSrc = append([]float64{srcv.v}, sumSrc...)

			if len(sumSrc) == l {
				break
			}
			lastAvailv = lastAvailv.prev
		}
	}

	// first new time
	itervt := startNew.t

	for {
		v := src.Get(itervt)
		if v == nil {
			break
		}

		// append new source data
		sumSrc = append(sumSrc, v.v)

		var set bool

		// if previous exists, we just subtract from first value and add new value
		if v.prev != nil {
			e := sum.Get(v.prev.t)
			if e != nil && len(sumSrc) == l+1 {
				newsum := e.v - sumSrc[0] + v.v
				sum.Set(itervt, newsum)
				set = true
			}
		}

		if !set {
			if len(sumSrc) >= l {
				var ct int
				var tot float64
				for i := len(sumSrc) - 1; i >= 0; i-- {
					ct++
					tot = tot + sumSrc[i]
					if ct == l {
						break
					}
				}
				sum.Set(itervt, tot)
			}
		}

		if v.next == nil {
			break
		}
		if v.t.Equal(stop.t) {
			break
		}

		if len(sumSrc) > l+1 {
			sumSrc = sumSrc[1:]
		}
		itervt = v.next.t
	}

	return sum, nil
}
