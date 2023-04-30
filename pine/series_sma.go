package pine

import (
	"fmt"
)

type smaCalcItem struct {
	valuetot float64
	total    int64
	seeked   int64
}

// SMA generates a ValueSeries of simple moving averages
func SMA(p ValueSeries, l int64) ValueSeries {
	key := fmt.Sprintf("sma:%s:%d", p.ID(), l)
	sma := getCache(key)
	if sma == nil {
		sma = NewValueSeries()
	}
	if p == nil || p.GetCurrent() == nil {
		return sma
	}

	// current available value
	stop := p.GetCurrent()
	// where we left off last time
	val := sma.GetLast()

	var f *Value
	// if we have not generated any SMAs yet
	if val == nil {
		f = p.GetFirst()
	} else {

		// time has not advanced. return cache
		if val.t.Equal(stop.t) {
			return sma
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

	return sma
}

var cache map[string]ValueSeries = make(map[string]ValueSeries)

func getCache(key string) ValueSeries {
	return cache[key]
}

func setCache(key string, v ValueSeries) {
	cache[key] = v
}
