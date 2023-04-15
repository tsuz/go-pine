package pine

import (
	"fmt"

	"github.com/pkg/errors"
)

// ROC calculates the percentage of change (rate of change) between the current value of `source` and its value `length` bars ago.
// It is calculated by the formula: 100 * change(src, length) / src[length].
//
// arguments are
//   - src: ValueSeries - Source data
//   - length: int - number of bars to lookback. 1 is the previous bar
func ROC(src ValueSeries, l int) (ValueSeries, error) {
	var err error
	if l < 1 {
		return nil, errors.New("length must be 1 or greater")
	}
	key := fmt.Sprintf("roc:%s:%s:%d", src.ID(), src.ID(), l)
	rocs := getCache(key)
	if rocs == nil {
		rocs = NewValueSeries()
	}

	// current available value
	stop := src.GetCurrent()

	if stop == nil {
		return rocs, nil
	}

	chg, err := Change(src, l)
	if err != nil {
		return rocs, errors.Wrapf(err, "error getting change")
	}

	rocs, err = roc(*stop, src, rocs, chg, l)
	if err != nil {
		return rocs, err
	}

	setCache(key, rocs)

	rocs.SetCurrent(stop.t)

	return rocs, err
}

func roc(stop Value, src, roc, chg ValueSeries, l int) (ValueSeries, error) {

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
		return roc, errors.New("expected continuous but found fragmented data and cannot continue")
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

	return roc, nil
}
