package pine

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type emaCalcItem struct {
	valuetot float64
	total    int64
	seeked   int64
}

// EMA generates a ValueSeries of exponential moving average
// the variable ema=ValueSeries is the exponentially weighted moving average values of p=ValueSeries
// ema may be behind where they should be with regards to p.GetCurrent()
// while ema catches up to where p.GetCurrent() is, the series should also contain
// all available average values between the last and up to p.GetCurrent()
//
// The formula for EMA is EMA=(closing price − previous day’s EMA)× smoothing constant as a decimal + previous day’s EMA
// where smoothing constant is 2 ÷ (number of time periods + 1)
// if the previous day's EMA is nil then it's the SMA of the lookback time.
// Using the above formula, the below example illustrates what EMA values look like
//
// t=time.Time (no iteration) | 1   |  2  | 3   | 4       |
// p=ValueSeries              | 13  | 15  | 17  | 18      |
// ema(close, 1)              | 13  | 15  | 17  | 18      |
// ema(close, 2)              | nil | 14  | 16  | 17.3333 |
func EMA(p ValueSeries, l int64) (ValueSeries, error) {
	if p == nil || p.GetCurrent() == nil {
		log.Infof("p is nil %t, p.GetCurrent is nil %t", p == nil, p.GetCurrent() == nil)
		return nil, nil
	}
	key := fmt.Sprintf("ema:%s:%d", p.ID(), l)
	ema := getCache(key)
	if ema == nil {
		log.Printf("Creating new EMA value series")
		ema = NewValueSeries()
	}

	// current available value
	stop := p.GetCurrent()
	val := ema.Get(stop.t)

	// if ema exists for this time period, return what we have
	if val != nil {
		return ema, nil
	}

	ema = getEMA(stop, p, ema, l)

	setCache(key, ema)

	ema.SetCurrent(stop.t)

	return ema, nil
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
			log.Printf("time to stop (%t) or v is nil (%t), t=%+v", v == nil, v.t.Equal(stop.t), itervt)
			break
		}
		e := ema.Get(itervt)
		if e != nil && v.next == nil {
			log.Printf("next value doesnt exist. ending  %+v", v.t)
			break
		}
		if e != nil {
			itervt = v.next.t
			log.Printf("ema exists skpping at %+v", v.t)
			continue
		}

		// get previous ema
		if v.prev != nil {
			prevv := vs.Get(v.prev.t)
			preve := ema.Get(prevv.t)
			// previous ema exists, just do multiplication to that
			if preve != nil {
				nextEMA := (v.v-preve.v)*mul + preve.v
				log.Printf("pushed ema to existing for %+v: %+v", v.t, nextEMA)
				ema.Push(v.t, nextEMA)
				continue
			}
		}

		// previous value does not exist. just keep adding until multplication is required
		fseek++
		ftot = ftot + v.v

		if fseek == l {
			avg := ftot / float64(fseek)
			log.Printf("pushed new ema for %+v: %+v", v.t, avg)
			ema.Push(v.t, avg)
		}

		if v.next == nil {
			log.Printf("Next doesn't exist breaking at %+v", v.t)
			break
		}
		if v.t.Equal(stop.t) {
			log.Printf("time to stop %+v", stop.t)
			break
		}
		itervt = v.next.t
	}

	// emaval := ema.Get(v.t)
	// if emaval != nil {
	// 	log.Printf("EMA value already exists at %+v", v.v)
	// 	return ema
	// }

	// pushEMA := func(ema ValueSeries, v *Value, emaval *Value, mul float64) {
	// 	newval := (v.v-emaval.prev.v)*mul + emaval.v
	// 	log.Printf("Adding EMA at value: %+v, ema: %+v", v.v, newval)
	// 	ema.Push(v.t, newval)
	// }

	// log.Printf("Get ema for v: %+v, l: %+v, ema is not nil: %t, current t: %+v, hasprev: %t", v.v, l, emaval != nil, v.t, v.prev != nil)

	// if emaval != nil && emaval.prev != nil {
	// 	pushEMA(ema, v, emaval, mul)
	// 	return ema
	// }
	// if v.prev == nil {
	// 	if l > 1 {
	// 		return ema
	// 	}
	// 	log.Printf("Adding EMA at value: %+v, ema: %+v", v.v, v.v)
	// 	ema.Push(v.t, v.v)
	// 	return ema
	// }
	// // needs to fetch previous values before proceeding
	// if emaval == nil || emaval.prev == nil {
	// 	// recursively get previous values
	// 	ema = getEMA(v.prev, ema, l)
	// 	emaval := ema.Get(v.t)
	// 	log.Printf("Done so looking up: %+v emaval: %+v", v.t, emaval)
	// 	if emaval != nil && emaval.prev != nil {
	// 		pushEMA(ema, v, emaval, mul)
	// 	}
	// }

	return ema
}
