package pine

import (
	"fmt"
)

// EMA generates a ValueSeries of exponential moving average.
func EMA(p ValueSeries, l int64) ValueSeries {
	key := fmt.Sprintf("ema:%s:%d", p.ID(), l)
	ema := getCache(key)
	if ema == nil {
		ema = NewValueSeries()
	}

	if p == nil || p.GetCurrent() == nil {
		return ema
	}

	// current available value
	stop := p.GetCurrent()

	ema = getEMA(stop, p, ema, l)

	setCache(key, ema)

	ema.SetCurrent(stop.t)

	return ema
}

func getEMA(stop *Value, vs ValueSeries, ema ValueSeries, l int64) ValueSeries {

	var mul float64 = 2.0 / float64(l+1.0)
	firstVal := ema.GetLast()

	if firstVal == nil {
		firstVal = vs.GetFirst()
	}

	if firstVal == nil {
		// if nothing is available, then nothing can be done
		return ema
	}

	itervt := firstVal.t

	var fseek int64
	var ftot float64

	for {
		v := vs.Get(itervt)
		if v == nil {
			break
		}
		e := ema.Get(itervt)
		if e != nil && v.next == nil {
			break
		}
		if e != nil {
			itervt = v.next.t
			continue
		}

		// get previous ema
		if v.prev != nil {
			prevv := vs.Get(v.prev.t)
			preve := ema.Get(prevv.t)
			// previous ema exists, just do multiplication to that
			if preve != nil {
				nextEMA := (v.v-preve.v)*mul + preve.v
				ema.Set(v.t, nextEMA)
				continue
			}
		}

		// previous value does not exist. just keep adding until multplication is required
		fseek++
		ftot = ftot + v.v

		if fseek == l {
			avg := ftot / float64(fseek)
			ema.Set(v.t, avg)
		}

		if v.next == nil {
			break
		}
		if v.t.Equal(stop.t) {
			break
		}
		itervt = v.next.t
	}

	return ema
}
