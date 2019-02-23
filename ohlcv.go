package pine

import "time"

type OHLCV struct {
	O float64
	H float64
	L float64
	C float64
	V float64
	S time.Time
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
