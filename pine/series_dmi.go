package pine

import (
	"github.com/pkg/errors"
)

// DMI generates a ValueSeries of directional movement index.
func DMI(ohlcv OHLCVSeries, len, smoo int) (adx, dmip, dmim ValueSeries, err error) {

	h := ohlcv.GetSeries(OHLCPropHigh)
	stop := h.GetCurrent()
	if stop == nil {
		return
	}

	l := ohlcv.GetSeries(OHLCPropLow)
	tr := ohlcv.GetSeries(OHLCPropTR)

	up, err := Change(h, 1)
	if err != nil {
		return adx, dmip, dmim, errors.Wrap(err, "error change")
	}
	down, err := Change(l, 1)
	if err != nil {
		return adx, dmip, dmim, errors.Wrap(err, "error change")
	}
	plusdm := up.Operate(down, func(uv, dv float64) float64 {
		dv = dv * -1
		if uv > dv && uv > 0 {
			return uv
		}
		return 0
	})
	minusdm := down.Operate(up, func(dv, uv float64) float64 {
		dv = dv * -1
		if dv > uv && dv > 0 {
			return dv
		}
		return 0
	})
	trurange, err := RMA(tr, int64(len))
	if err != nil {
		return adx, dmip, dmim, errors.Wrap(err, "error trurange")
	}
	plusdmrma, err := RMA(plusdm, int64(len))
	if err != nil {
		return adx, dmip, dmim, errors.Wrap(err, "error RMA")
	}
	minusdmrma, err := RMA(minusdm, int64(len))
	if err != nil {
		return adx, dmip, dmim, errors.Wrap(err, "error RMA")
	}
	plus := plusdmrma.Div(trurange).MulConst(100)
	minus := minusdmrma.Div(trurange).MulConst(100)

	sum := plus.Add(minus)
	denom := sum.Operate(sum, func(a, b float64) float64 {
		if a == 0 {
			return 1
		}
		return a
	})

	adxrma, err := RMA(DiffAbs(plus, minus).Div(denom), 3)
	if err != nil {
		return adx, dmip, dmim, errors.Wrap(err, "error RMA for adx")
	}
	adx = adxrma.MulConst(100)

	return adx, plus, minus, nil
}
