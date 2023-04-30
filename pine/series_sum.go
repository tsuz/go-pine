package pine

import (
	"fmt"
	"log"
)

// Sum generates a ValueSeries of summation of previous values
//
// Parameters
//   - p - ValueSeries: source data
//   - l - int: lookback periods [1, ∞)
func Sum(p ValueSeries, l int) (ValueSeries, error) {

	key := fmt.Sprintf("sum:%s:%d", p.ID(), l)
	sum := getCache(key)
	if sum == nil {
		log.Printf("New sum series: %+v", key)
		sum = NewValueSeries()
	}

	sum = generateSum(p, sum, l)

	setCache(key, sum)

	return sum, nil
}

// SumNoCache generates sum without caching
func SumNoCache(p ValueSeries, l int) ValueSeries {
	sum := NewValueSeries()
	return generateSum(p, sum, l)
}

func generateSum(p, sum ValueSeries, l int) ValueSeries {
	// current available value
	stop := p.GetCurrent()
	if stop == nil {
		return sum
	}
	sum, _ = getSum(*stop, sum, p, l)
	sum.SetCurrent(stop.t)
	return sum
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
