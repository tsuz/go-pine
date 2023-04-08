package pine

import "time"

type OHLCV struct {
	O float64
	H float64
	L float64
	C float64
	V float64
	S time.Time

	prev *OHLCV
	next *OHLCV
}

func NewOHLCVWithSamePx(px, qty float64, t time.Time) OHLCV {
	return OHLCV{
		O: px,
		H: px,
		L: px,
		C: px,
		V: qty,
		S: t,
	}
}

func (o *OHLCV) Get(p OHLCProp) float64 {
	switch p {
	case OHLCPropOpen:
		return o.O
	case OHLCPropClose:
		return o.C
	}
	return o.C
}
