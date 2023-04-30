package pine

import (
	"fmt"
)

// MFI generates a ValueSeries of exponential moving average.
func MFI(o OHLCVSeries, l int64) ValueSeries {
	key := fmt.Sprintf("mfi:%s:%d", o.ID(), l)
	mfi := getCache(key)
	if mfi == nil {
		mfi = NewValueSeries()
	}

	hlc3 := OHLCVAttr(o, OHLCPropHLC3)
	vol := OHLCVAttr(o, OHLCPropVolume)

	hlc3c := hlc3.GetCurrent()
	if hlc3c == nil {
		return mfi
	}

	chg := Change(hlc3, 1)

	u := OperateWithNil(hlc3, chg, "mfiu", func(a, b *Value) *Value {
		var v float64
		// treat nil value as HLC3
		if b == nil {
			return a
		} else {
			v = b.v
		}
		if v <= 0.0 {
			return &Value{
				t: a.t,
				v: 0.0,
			}
			//  NewFloat64(0.0)
		}
		return a
	})
	lo := OperateWithNil(hlc3, chg, "mfil", func(a, b *Value) *Value {
		var v float64
		// treat nil value as HLC3
		if b == nil {
			return a
		} else {
			v = b.v
		}
		if v >= 0.0 {
			return &Value{
				t: a.t,
				v: 0.0,
			}
		}
		return a
	})

	uv := Mul(vol, u)
	lv := Mul(vol, lo)

	upper := Sum(uv, int(l))

	lower := Sum(lv, int(l))

	hundo := ReplaceAll(hlc3, 100)

	mfi = Sub(hundo, Div(hundo, AddConst(Div(upper, lower), 1)))

	mfi.SetCurrent(hlc3c.t)

	setCache(key, mfi)

	return mfi
}
