package pine

import (
	"fmt"

	"github.com/pkg/errors"
)

// CCI generates a ValueSeries of exponential moving average.
func CCI(tp ValueSeries, l int64) (ValueSeries, error) {
	key := fmt.Sprintf("cci:%s:%d", tp.ID(), l)
	cci := getCache(key)
	if cci == nil {
		cci = NewValueSeries()
	}

	tpv := tp.GetCurrent()
	if tpv == nil {
		return cci, nil
	}

	ma, err := SMA(tp, l)
	if err != nil {
		return cci, errors.Wrap(err, "error sma")
	}

	// need moving average to perform this
	if ma.GetCurrent() == nil {
		return cci, nil
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

	return cci, nil
}
