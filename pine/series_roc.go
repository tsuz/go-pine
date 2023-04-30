package pine

import (
	"fmt"
)

// ROC calculates the percentage of change (rate of change) between the current value of `source` and its value `length` bars ago.
// It is calculated by the formula
//
//   - 100 * change(src, length) / src[length].
//
// arguments are
//   - src: ValueSeries - Source data
//   - length: int - number of bars to lookback. 1 is the previous bar
func ROC(src ValueSeries, l int) ValueSeries {

	key := fmt.Sprintf("roc:%s:%s:%d", src.ID(), src.ID(), l)
	rocs := getCache(key)
	if rocs == nil {
		rocs = NewValueSeries()
	}

	// current available value
	stop := src.GetCurrent()

	if stop == nil {
		return rocs
	}

	chg := Change(src, l)

	rocs = roc(*stop, src, rocs, chg, l)

	setCache(key, rocs)

	rocs.SetCurrent(stop.t)

	return rocs
}

func roc(stop Value, src, roc, chg ValueSeries, l int) ValueSeries {

	var val *Value

	lastvw := roc.GetCurrent()
	if lastvw != nil {
		val = src.Get(lastvw.t)
		if val != nil {
			val = val.next
		}
	} else {
		val = src.GetFirst()
	}

	if val == nil {
		return roc
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
			chgv := chg.Get(val.t)
			if chgv != nil {
				v := 100 * chgv.v / vwappend
				roc.Set(val.t, v)
			}
		}

		val = val.next
	}

	return roc
}
