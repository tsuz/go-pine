package pine

import (
	"fmt"
)

// CCI generates a ValueSeries of exponential moving average.
func CCI(tp ValueSeries, l int64) ValueSeries {
	key := fmt.Sprintf("cci:%s:%d", tp.ID(), l)
	cci := getCache(key)
	if cci == nil {
		cci = NewValueSeries()
	}

	tpv := tp.GetCurrent()
	if tpv == nil {
		return cci
	}

	ma := SMA(tp, l)

	// need moving average to perform this
	if ma.GetCurrent() == nil {
		return cci
	}
	mav := ma.GetCurrent().v
	mdv := SubConstNoCache(tp, mav)

	// get absolute value
	mdvabs := OperateNoCache(mdv, mdv, "cci:absval", func(a, b float64) float64 {
		if a < 0 {
			return -1 * a
		}
		return a
	})
	mdvabssum := SumNoCache(mdvabs, int(l))
	md := DivConstNoCache(mdvabssum, float64(l))
	denom := MulConstNoCache(md, 0.015)
	cci = DivNoCache(mdv, denom)

	setCache(key, cci)

	cci.SetCurrent(tpv.t)

	return cci
}
