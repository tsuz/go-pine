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
