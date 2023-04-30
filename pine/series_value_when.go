package pine

import (
	"fmt"
)

// ValueWhen generates a ValueSeries of float
// arguments are
//   - bs: ValueSeries - Value Series where values are 0.0, 1.0 (boolean)
//   - src: ValueSeries - Value Series of the source
//   - ocr: int - The occurrence of the condition. The numbering starts from 0 and goes back in time, so '0' is the most recent occurrence of `condition`, '1' is the second most recent and so forth. Must be an integer >= 0.
func ValueWhen(bs, src ValueSeries, ocr int) ValueSeries {
	key := fmt.Sprintf("valuewhen:%s:%s:%d", bs.ID(), src.ID(), ocr)
	vw := getCache(key)
	if vw == nil {
		vw = NewValueSeries()
	}

	// current available value
	stop := src.GetCurrent()

	if stop == nil {
		return vw
	}

	vw = valueWhen(*stop, bs, src, vw, ocr)

	setCache(key, vw)

	vw.SetCurrent(stop.t)

	return vw
}

func valueWhen(stop Value, bs, src, vw ValueSeries, ocr int) ValueSeries {

	var val *Value

	lastvw := vw.GetCurrent()
	if lastvw != nil {
		val = bs.Get(lastvw.t)
		if val != nil {
			val = val.next
		}
	} else {
		val = bs.GetFirst()
	}

	if val == nil {
		return vw
	}

	// populate src values if condition=1.0
	condSrc := make([]float64, 0)

	prevVal := val
	for {
		prevVal = prevVal.prev
		if prevVal == nil {
			break
		}

		b := bs.Get(prevVal.t)
		if b == nil {
			continue
		}
		if b.v == 1 {
			srcv := src.Get(prevVal.t)
			// add at the beginning since we go backwards
			condSrc = append([]float64{srcv.v}, condSrc...)
		}

		if len(condSrc) == (ocr + 1) {
			break
		}
	}

	// last available does not exist. start from first

	for {
		if val == nil {
			break
		}
		// update
		if val.v == 1.0 {
			srcval := src.Get(val.t)
			if srcval != nil {
				condSrc = append(condSrc, srcval.v)
				if len(condSrc) > (ocr + 1) {
					condSrc = condSrc[1:]
				}
			}
		}

		if len(condSrc) == (ocr + 1) {
			vwappend := condSrc[0]
			vw.Set(val.t, vwappend)
		}

		val = val.next
	}

	return vw
}
