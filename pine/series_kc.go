package pine

import "github.com/pkg/errors"

// KC generates ValueSeries of ketler channel's middle, upper and lower in that order.
func KC(src ValueSeries, o OHLCVSeries, l int64, mult float64, usetr bool) (middle, upper, lower ValueSeries, err error) {

	lower = NewValueSeries()
	upper = NewValueSeries()
	middle = NewValueSeries()
	start := src.GetCurrent()

	if start == nil {
		return middle, upper, lower, nil
	}

	var span ValueSeries
	basis, err := EMA(src, l)
	if err != nil {
		return middle, upper, lower, errors.Wrap(err, "error EMA")
	}

	if usetr {
		span = OHLCVAttr(o, OHLCPropTR)
	} else {
		h := OHLCVAttr(o, OHLCPropHigh)
		l := OHLCVAttr(o, OHLCPropLow)
		span = h.Sub(l)
	}

	rangeEma, err := EMA(span, l)
	if err != nil {
		return middle, upper, lower, errors.Wrap(err, "error EMA")
	}

	middle = basis
	rangeEmaMul := rangeEma.MulConst(mult)
	upper = basis.Add(rangeEmaMul)
	lower = basis.Sub(rangeEmaMul)

	middle.SetCurrent(start.t)
	upper.SetCurrent(start.t)
	lower.SetCurrent(start.t)

	return middle, upper, lower, nil
}
