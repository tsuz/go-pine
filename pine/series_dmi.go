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

	// smooth := func(a ValueSeries) (ValueSeries, error) {
	// 	return WilderSMA(a, smoo)
	// }

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

	// trs, err := smooth(tr)
	// if err != nil {
	// 	return adx, dmip, dmim, errors.Wrap(err, "error getting TR smooth")
	// }

	// plusdmraw := highMinusPrevHigh(h)
	// minusdmraw := prevLowMinusCurrentLow(l)
	// plusdm := plusdmraw.Operate(minusdmraw, func(plus, minus float64) float64 {
	// 	if plus > minus && plus > 0 {
	// 		return plus
	// 	}
	// 	return 0
	// })
	// minusdm := minusdmraw.Operate(minusdmraw, func(plus, minus float64) float64 {
	// 	if minus > plus && minus > 0 {
	// 		return minus
	// 	}
	// 	return 0
	// })

	// plusdms, err := smooth(plusdm)
	// if err != nil {
	// 	return adx, dmip, dmim, errors.Wrap(err, "error wildersma for plusdm")
	// }
	// dmip = plusdms.Div(trs).MulConst(100)

	// minusdms, err := smooth(minusdm)
	// if err != nil {
	// 	return adx, dmip, dmim, errors.Wrap(err, "error wildersma for minusdm")
	// }
	// dmim = minusdms.Div(trs).MulConst(100)

	// dmipsum, err := Sum(dmip, len)
	// if err != nil {
	// 	return adx, dmip, dmim, errors.Wrap(err, "error sum for dmip")
	// }
	// dmimsum, err := Sum(dmim, len)
	// if err != nil {
	// 	return adx, dmip, dmim, errors.Wrap(err, "error sum for dmim")
	// }
	// sumtot := dmipsum.Add(dmimsum)
	// adx = DiffAbs(dmip, dmim).Div(sumtot).MulConst(100)

	// adx.SetCurrent(stop.t)
	// dmip.SetCurrent(stop.t)
	// dmim.SetCurrent(stop.t)

	// return adx, dmip, dmim, nil
}

func highMinusPrevHigh(v ValueSeries) ValueSeries {
	copied := NewValueSeries()
	f := v.GetFirst()
	for {
		if f == nil {
			break
		}

		newv := v.Get(f.t)

		if newv.prev != nil {
			copied.Set(f.t, newv.v-newv.prev.v)
		}

		f = f.next
	}
	cur := v.GetCurrent()
	if cur != nil {
		copied.SetCurrent(cur.t)
	}
	return copied
}

func prevLowMinusCurrentLow(v ValueSeries) ValueSeries {
	copied := NewValueSeries()
	f := v.GetFirst()
	for {
		if f == nil {
			break
		}

		newv := v.Get(f.t)

		if newv.prev != nil {
			copied.Set(f.t, newv.prev.v-newv.v)
		}

		f = f.next
	}
	cur := v.GetCurrent()
	if cur != nil {
		copied.SetCurrent(cur.t)
	}
	return copied
}

// func getDMI(stop *Value, l ValueSeries, h ValueSeries, c ValueSeries, len, smoo int64) ValueSeries {

// 	var mul float64 = 2.0 / float64(len+1.0)
// 	firstVal := ema.GetLast()

// 	if firstVal == nil {
// 		firstVal = vs.GetFirst()
// 	}

// 	if firstVal == nil {
// 		// if nothing is available, then nothing can be done
// 		return ema
// 	}

// 	itervt := firstVal.t

// 	var fseek int64
// 	var ftot float64

// 	for {
// 		v := vs.Get(itervt)
// 		if v == nil {
// 			break
// 		}
// 		e := ema.Get(itervt)
// 		if e != nil && v.next == nil {
// 			break
// 		}
// 		if e != nil {
// 			itervt = v.next.t
// 			continue
// 		}

// 		// get previous ema
// 		if v.prev != nil {
// 			prevv := vs.Get(v.prev.t)
// 			preve := ema.Get(prevv.t)
// 			// previous ema exists, just do multiplication to that
// 			if preve != nil {
// 				nextEMA := (v.v-preve.v)*mul + preve.v
// 				ema.Set(v.t, nextEMA)
// 				continue
// 			}
// 		}

// 		// previous value does not exist. just keep adding until multplication is required
// 		fseek++
// 		ftot = ftot + v.v

// 		if fseek == l {
// 			avg := ftot / float64(fseek)
// 			ema.Set(v.t, avg)
// 		}

// 		if v.next == nil {
// 			break
// 		}
// 		if v.t.Equal(stop.t) {
// 			break
// 		}
// 		itervt = v.next.t
// 	}

// 	return ema
// }
