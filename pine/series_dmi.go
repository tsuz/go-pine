package pine

import (
	"fmt"
)

// DMI generates a ValueSeries of directional movement index.
func DMI(ohlcv OHLCVSeries, len, smoo int) (adx, plus, minus ValueSeries) {
	adxkey := fmt.Sprintf("adx:%s:%d:%d", ohlcv.ID(), len, smoo)
	adx = getCache(adxkey)
	if adx == nil {
		adx = NewValueSeries()
	}

	pluskey := fmt.Sprintf("plus:%s:%d:%d", ohlcv.ID(), len, smoo)
	plus = getCache(pluskey)
	if plus == nil {
		plus = NewValueSeries()
	}

	minuskey := fmt.Sprintf("minus:%s:%d:%d", ohlcv.ID(), len, smoo)
	minus = getCache(minuskey)
	if minus == nil {
		minus = NewValueSeries()
	}

	h := OHLCVAttr(ohlcv, OHLCPropHigh)
	stop := h.GetCurrent()
	if stop == nil {
		return
	}

	l := OHLCVAttr(ohlcv, OHLCPropLow)
	tr := OHLCVAttr(ohlcv, OHLCPropTRHL)

	up := Change(h, 1)
	down := Change(l, 1)
	plusdm := Operate(up, down, "dmi:uv", func(uv, dv float64) float64 {
		dv = dv * -1
		if uv > dv && uv > 0 {
			return uv
		}
		return 0
	})
	minusdm := Operate(down, up, "dmi:uv", func(dv, uv float64) float64 {
		dv = dv * -1
		if dv > uv && dv > 0 {
			return dv
		}
		return 0
	})
	trurange := RMA(tr, int64(len))
	plusdmrma := RMA(plusdm, int64(len))
	minusdmrma := RMA(minusdm, int64(len))
	plus = MulConst(Div(plusdmrma, trurange), 100)
	minus = MulConst(Div(minusdmrma, trurange), 100)

	sum := Add(plus, minus)
	denom := Operate(sum, sum, "dmi:denom", func(a, b float64) float64 {
		if a == 0 {
			return 1
		}
		return a
	})

	adxrma := RMA(Div(DiffAbs(plus, minus), denom), 3)
	adx = MulConst(adxrma, 100)

	setCache(adxkey, adx)
	setCache(pluskey, plus)
	setCache(minuskey, minus)

	return adx, plus, minus
}
