package pine

import (
	"fmt"
	"math"
	"time"
)

func OHLCVAttr(o OHLCVSeries, p OHLCProp) ValueSeries {
	key := fmt.Sprintf("ohlcvattr:%s:%d", o.ID(), p)
	dest := getCache(key)
	if dest == nil {
		dest = NewValueSeries()
	}

	stop := o.Current()
	if stop == nil {
		return dest
	}
	dest = getOHLCVAttr(*stop, o, dest, p)

	dest.SetCurrent(stop.S)

	setCache(key, dest)

	return dest
}

func getOHLCVAttr(stop OHLCV, src OHLCVSeries, dest ValueSeries, p OHLCProp) ValueSeries {

	var startt time.Time

	firstVal := dest.GetLast()
	if firstVal != nil {
		if v := src.Get(firstVal.t); v != nil && v.next != nil {
			startt = v.next.S
		}
		// startt = firstVal.next.t
	} else if firstVal == nil {
		if v := src.GetFirst(); v != nil {
			startt = v.S
		}
	}

	if startt.IsZero() {
		return dest
	}

	ptr := startt

	for {
		v := src.Get(ptr)

		var propVal *float64
		switch p {
		case OHLCPropClose:
			propVal = NewFloat64(v.C)
		case OHLCPropOpen:
			propVal = NewFloat64(v.O)
		case OHLCPropHigh:
			propVal = NewFloat64(v.H)
		case OHLCPropLow:
			propVal = NewFloat64(v.L)
		case OHLCPropVolume:
			propVal = NewFloat64(v.V)
		case OHLCPropTR, OHLCPropTRHL:
			if v.prev != nil {
				p := v.prev
				v1 := math.Abs(v.H - v.L)
				v2 := math.Abs(v.H - p.C)
				v3 := math.Abs(v.L - p.C)
				v := math.Max(v1, math.Max(v2, v3))
				propVal = NewFloat64(v)
			}
			if p == OHLCPropTRHL && v.prev == nil {
				d := v.H - v.L
				propVal = &d
			}
		case OHLCPropHLC3:
			propVal = NewFloat64((v.H + v.L + v.C) / 3)
		default:
			continue
		}
		if propVal != nil {
			dest.Set(v.S, *propVal)
		}

		if v.next == nil {
			break
		}
		if v.S.Equal(stop.S) {
			break
		}

		ptr = v.next.S
	}

	return dest
}
