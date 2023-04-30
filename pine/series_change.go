package pine

import (
	"fmt"
)

// Change compares the current `source` value to its value `lookback` bars ago and returns the difference.
//
// arguments are
//   - src: ValueSeries - Source data to seek difference
//   - lookback: int - Lookback to compare the change
func Change(src ValueSeries, lookback int) ValueSeries {
	key := fmt.Sprintf("change:%s:%s:%d", src.ID(), src.ID(), lookback)
	chg := getCache(key)
	if chg == nil {
		chg = NewValueSeries()
	}

	// current available value
	stop := src.GetCurrent()

	if stop == nil {
		return chg
	}

	chg = change(*stop, src, chg, lookback)

	setCache(key, chg)

	chg.SetCurrent(stop.t)

	return chg
}

func change(stop Value, src, chg ValueSeries, l int) ValueSeries {

	var val *Value

	lastvw := chg.GetCurrent()
	if lastvw != nil {
		val = src.Get(lastvw.t)
		if val != nil {
			val = val.next
		}
	} else {
		val = src.GetFirst()
	}

	if val == nil {
		return chg
	}

	// populate src values
	condSrc := make([]float64, 0)

	prevVal := val
	for {
		prevVal = prevVal.prev
		if prevVal == nil {
			break
		}

		b := src.Get(prevVal.t)
		if b == nil {
			continue
		}

		srcv := src.Get(prevVal.t)
		// add at the beginning since we go backwards
		condSrc = append([]float64{srcv.v}, condSrc...)

		if len(condSrc) == (l + 1) {
			break
		}
	}

	// last available does not exist. start from first

	for {
		if val == nil {
			break
		}
		// update

		srcval := src.Get(val.t)
		if srcval != nil {
			condSrc = append(condSrc, srcval.v)
			if len(condSrc) > (l + 1) {
				condSrc = condSrc[1:]
			}
		}

		if len(condSrc) == (l + 1) {
			vwappend := condSrc[0]
			chg.Set(val.t, val.v-vwappend)
		}

		val = val.next
	}

	return chg
}

func NewFloat64(v float64) *float64 {
	v2 := v
	return &v2
}
