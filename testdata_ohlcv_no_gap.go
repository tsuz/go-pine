package pine

import (
	"math"
	"math/rand"
	"time"
)

// OHLCVTestData generates test data
func OHLCVTestData(start time.Time, period, intervalms int64) []OHLCV {
	s := start
	v := make([]OHLCV, 0)
	for i := 0; i < int(period); i++ {
		ohlcv := GenerateOHLCV(s)
		ohlcv.S = s
		v = append(v, ohlcv)
		s = s.Add(time.Duration(intervalms * 1e6))
	}
	return v
}

func GenerateOHLCV(t time.Time) OHLCV {
	max := 20.0
	min := 10.0

	o := Rand(min, max)
	c := Rand(min, max)
	v := OHLCV{
		O: o,
		C: c,
	}

	h2 := math.Max(c, o)
	h := Rand(h2, max)
	v.H = h

	l2 := math.Min(c, o)
	l := Rand(min, l2)
	v.L = l

	return v
}

func Rand(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}
