package pine

import (
	"fmt"

	"github.com/pkg/errors"
)

// WilderSMA is a wilder's smoothed moving average. This is not exposed as a default in pinescript.
// This is not to be confused with ta.wma which stands for weighted moving average
func WilderSMA(a ValueSeries, len int) (ValueSeries, error) {
	key := fmt.Sprintf("wildersma:%s:%d", a.ID(), len)
	wsma := getCache(key)
	if wsma == nil {
		wsma = NewValueSeries()
	}

	stop := a.GetCurrent()
	if stop == nil {
		return wsma, nil
	}

	wsma, err := getWilderSMA(*stop, a, wsma, len)
	if err != nil {
		return wsma, errors.Wrap(err, "error wildersma")
	}

	setCache(key, wsma)

	if stop != nil {
		wsma.SetCurrent(stop.t)
	}

	return wsma, nil

	// return a.OperateWithNil(a, func(av *Value, _ *Value) *Value {
	// 	var prevw float64
	// 	var curv float64
	// 	if av == nil {
	// 		return nil
	// 	}

	// 	v := a.Get(av.t)
	// 	if v.prev != nil {
	// 		prevw = v.prev.v
	// 	}

	// 	curv = av.v

	// 	newv := (prevw + (curv - prevw)) / float64(len)

	// 	return &Value{
	// 		t: av.t,
	// 		v: newv,
	// 	}
	// })
}

func getWilderSMA(stop Value, src, wsma ValueSeries, len int) (ValueSeries, error) {
	firstVal := wsma.Get(stop.t)

	if firstVal == nil {
		firstVal = src.GetFirst()
	}

	if firstVal == nil {
		// if nothing is available, then nothing can be done
		return wsma, nil
	}

	itervt := firstVal.t

	var fseek int
	var ftot float64

	for {
		v := src.Get(itervt)
		if v == nil {
			break
		}
		e := wsma.Get(itervt)
		if e != nil && v.next == nil {
			break
		}
		if e != nil {
			itervt = v.next.t
			continue
		}

		// get previous ema
		if v.prev != nil {
			prevv := src.Get(v.prev.t)
			preve := wsma.Get(prevv.t)
			// previous wsma exists, just do multiplication to that
			if preve != nil {
				wsmav := (preve.v*float64(len) - preve.v + v.v) / float64(len)
				wsma.Set(v.t, wsmav)
				continue
			}
		}

		// previous value does not exist. just keep adding until multplication is required
		fseek++
		ftot = ftot + v.v

		if fseek == len {
			avg := ftot / float64(fseek)
			wsma.Set(v.t, avg)
		}

		if v.next == nil {
			break
		}
		if v.t.Equal(stop.t) {
			break
		}
		itervt = v.next.t
	}

	return wsma, nil
}
