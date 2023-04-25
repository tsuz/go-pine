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
		ohlcv := generateOHLCV(s)
		ohlcv.S = s
		v = append(v, ohlcv)
		s = s.Add(time.Duration(intervalms * 1e6))
	}
	return v
}

func OHLCVStaticTestData() []OHLCV {
	start := time.Now()
	data := []OHLCV{
		OHLCV{O: 11.3, H: 19.7, L: 11.1, C: 16.5, V: 11.6},
		OHLCV{O: 12.9, H: 19.1, L: 12.3, C: 18.7, V: 13.0},
		OHLCV{O: 11.0, H: 18.8, L: 10.3, C: 18.2, V: 13.8},
		OHLCV{O: 19.2, H: 19.6, L: 11.7, C: 11.9, V: 15.9},
		OHLCV{O: 18.1, H: 19.5, L: 11.2, C: 19.3, V: 16.8},
		OHLCV{O: 19.4, H: 19.8, L: 13.5, C: 14.2, V: 19.1},
		OHLCV{O: 19.1, H: 19.5, L: 12.9, C: 14.4, V: 14.7},
		OHLCV{O: 10.6, H: 19.9, L: 10.3, C: 11.0, V: 11.7},
		OHLCV{O: 18.8, H: 19.0, L: 12.4, C: 14.7, V: 17.4},
		OHLCV{O: 17.1, H: 17.6, L: 10.0, C: 10.3, V: 15.0},
	}

	for i := range data {
		fivemin := 5 * time.Minute
		data[i].S = start.Add(time.Duration(i) * fivemin)
	}
	return data
}

func generateOHLCV(t time.Time) OHLCV {
	max := 20.0
	min := 10.0

	o := randVal(min, max)
	c := randVal(min, max)
	vol := randVal(min, max)
	v := OHLCV{
		O: o,
		C: c,
		V: vol,
	}

	h2 := math.Max(c, o)
	h := randVal(h2, max)
	v.H = h

	l2 := math.Min(c, o)
	l := randVal(min, l2)
	v.L = l

	return v
}

func randVal(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}
