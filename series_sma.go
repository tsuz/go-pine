package pine

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type smaCalcItem struct {
	valuetot float64
	total    int64
	seeked   int64
}

// SMA generates a ValueSeries of simple moving averages
// the variable sma=ValueSeries is the average values of p=ValueSeries
// sma may be behind where they should be with regards to p.GetCurrent()
// while sma catches up to where p.GetCurrent() is, the series should also contain
// all available average values between the last and up to p.GetCurrent()
//
// The below example illustrates sma needs to be generated for time=3,4
//
//	  t=time.Time  | 1 |  2  | 3 | 4                 | 5  |
//	p=ValueSeries  | 4 |  5  | 9 | 12 (p.GetCurrent) | 14 |
//
// sma=ValueSeries |   | 4.5 |   |                   |    |
func SMA(p ValueSeries, l int64) (ValueSeries, error) {
	if p == nil || p.GetCurrent() == nil {
		log.Infof("p is nil %t, p.GetCurrent is nil %t", p == nil, p.GetCurrent() == nil)
		return nil, nil
	}
	key := fmt.Sprintf("sma:%s:%d", p.ID(), l)
	sma := getCache(key)
	if sma == nil {
		sma = NewValueSeries()
	}

	// current available value
	stop := p.GetCurrent()
	// where we left off last time
	val := sma.GetLast()
	log.Printf("Sma lookback: %d, stop: %+v", l, stop.v)

	var f *Value
	// if we have not generated any SMAs yet
	if val == nil {
		f = p.GetFirst()
	} else {

		// time has not advanced. return cache
		if val.t.Equal(stop.t) {
			return sma, nil
		}

		// value exists, find where we need to start off
		v := p.Get(val.t)
		if v == nil {
			f = p.GetFirst()
		} else {
			for i := 0; i < int(l)-1; i++ {
				if v.prev == nil {
					break
				}
				v = v.prev
			}
			f = v
		}
	}

	// generate from the beginning
	calcs := make(map[int64]smaCalcItem)
	for {
		if f == nil {
			break
		}
		calcs[f.t.Unix()] = smaCalcItem{}
		toUpdate := make(map[int64]smaCalcItem)
		for k, v := range calcs {

			var nvt float64
			if f == nil {
				nvt = v.valuetot
			} else {
				nvt = v.valuetot + f.v
			}
			toUpdate[k] = smaCalcItem{
				valuetot: nvt,
				seeked:   v.seeked + 1,
				total:    v.total + 1,
			}
		}

		for k := range toUpdate {
			calcs[k] = toUpdate[k]
			// done seeking lookback times
			if toUpdate[k].seeked == l {
				if toUpdate[k].total > 0 {
					v := toUpdate[k].valuetot / float64(toUpdate[k].total)
					sma.Set(f.t, v)
				}
				delete(calcs, k)
			}
		}
		if f.t.Equal(stop.t) {
			break
		}
		f = f.next
	}

	setCache(key, sma)

	sma.SetCurrent(stop.t)

	return sma, nil
}

var cache map[string]ValueSeries = make(map[string]ValueSeries)

func getCache(key string) ValueSeries {
	return cache[key]
}

func setCache(key string, v ValueSeries) {
	cache[key] = v
}
