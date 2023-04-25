package pine

import (
	"fmt"

	"github.com/pkg/errors"
)

// MFI generates a ValueSeries of exponential moving average.
func MFI(o OHLCVSeries, l int64) (ValueSeries, error) {
	key := fmt.Sprintf("mfi:%s:%d", o.ID(), l)
	mfi := getCache(key)
	if mfi == nil {
		mfi = NewValueSeries()
	}

	hlc3 := o.GetSeries(OHLCPropHLC3)
	vol := o.GetSeries(OHLCPropVolume)

	hlc3c := hlc3.GetCurrent()
	if hlc3c == nil {
		return mfi, nil
	}

	chg, err := Change(hlc3, 1)
	if err != nil {
		return mfi, errors.Wrap(err, "error getting change")
	}

	u := hlc3.OperateWithNil(chg, func(a, b *float64) *float64 {
		var v float64
		// treat nil value as HLC3
		if b == nil {
			return a
		} else {
			v = *b
		}
		if v <= 0.0 {
			return NewFloat64(0.0)
		}
		return a
	})
	lo := hlc3.OperateWithNil(chg, func(a, b *float64) *float64 {
		var v float64
		// treat nil value as HLC3
		if b == nil {
			return a
		} else {
			v = *b
		}
		if v >= 0.0 {
			return NewFloat64(0.0)
		}
		return a
	})

	uv := vol.Mul(u)
	lv := vol.Mul(lo)

	upper, err := Sum(uv, int(l))
	if err != nil {
		return mfi, errors.Wrap(err, "error getting sum for higher")
	}

	lower, err := Sum(lv, int(l))
	if err != nil {
		return mfi, errors.Wrap(err, "error getting sum for lower")
	}

	hundo := hlc3.Copy()
	hundo.SetAll(100)

	mfi = hundo.Sub(hundo.Div(upper.Div(lower).AddConst(1)))

	setCache(key, mfi)

	mfi.SetCurrent(hlc3c.t)

	return mfi, nil
}
