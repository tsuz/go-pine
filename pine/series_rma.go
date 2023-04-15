package pine

import (
	"fmt"
)

// RMA generates a ValueSeries of the exponentially weighted moving average with alpha = 1 / length.
func RMA(p ValueSeries, l int64) (ValueSeries, error) {
	key := fmt.Sprintf("rma:%s:%d", p.ID(), l)
	rma := getCache(key)
	if rma == nil {
		rma = NewValueSeries()
	}

	if p == nil || p.GetCurrent() == nil {
		return rma, nil
	}

	// current available value
	stop := p.GetCurrent()

	rma = getRMA(stop, p, rma, l)

	setCache(key, rma)

	rma.SetCurrent(stop.t)

	return rma, nil
}

func getRMA(stop *Value, vs ValueSeries, rma ValueSeries, l int64) ValueSeries {

	var mul float64 = 1.0 / float64(l)
	firstVal := rma.GetLast()

	if firstVal == nil {
		firstVal = vs.GetFirst()
	}

	if firstVal == nil {
		// if nothing is available, then nothing can be done
		return rma
	}

	itervt := firstVal.t

	var fseek int64
	var ftot float64

	for {
		v := vs.Get(itervt)
		if v == nil {
			break
		}
		e := rma.Get(itervt)
		if e != nil && v.next == nil {
			break
		}
		if e != nil {
			itervt = v.next.t
			continue
		}

		// get previous rma
		if v.prev != nil {
			prevv := vs.Get(v.prev.t)
			preve := rma.Get(prevv.t)
			// previous ema exists, just do multiplication to that
			if preve != nil {
				nextRMA := (preve.v)*(1-mul) + v.v*mul
				rma.Set(v.t, nextRMA)
				continue
			}
		}

		// previous value does not exist. just keep adding until multplication is required
		fseek++
		ftot = ftot + v.v

		if fseek == l {
			avg := ftot / float64(fseek)
			rma.Set(v.t, avg)
		}

		if v.next == nil {
			break
		}
		if v.t.Equal(stop.t) {
			break
		}
		itervt = v.next.t
	}

	return rma
}
