package pine_test

import (
	pine "go-pine"
	"testing"
	"time"
)

func TestSeriesInitWithValues(t *testing.T) {
	opts := pine.SeriesOpts{
		Interval: 300,
		Max:      100,
	}
	now := time.Now()
	fivemin := now.Add(5 * time.Minute)
	data := []pine.OHLCV{
		pine.OHLCV{
			O: 14,
			H: 15,
			L: 13,
			C: 14,
			V: 131,
			S: now,
		},
		pine.OHLCV{
			O: 13,
			H: 18,
			L: 10,
			C: 15,
			V: 12,
			S: fivemin,
		},
	}
	s, err := pine.NewSeries(data, opts)
	if err != nil {
		t.Fatal(err)
	}
	v := s.GetValueForInterval(now)
	if v == nil || v.OHLCV == nil {
		t.Fatalf("expected ohlcv to be non nil but got %+v", v)
	} else if v.OHLCV.O != 14 {
		t.Fatalf("expected open to be 14 but got %+v", v.OHLCV.O)
	} else if v.OHLCV.H != 15 {
		t.Fatalf("expected high to be 15 but got %+v", v.OHLCV.H)
	} else if v.OHLCV.L != 13 {
		t.Fatalf("expected low to be 13 but got %+v", v.OHLCV.L)
	} else if v.OHLCV.C != 14 {
		t.Fatalf("expected close to be 14 but got %+v", v.OHLCV.C)
	} else if v.OHLCV.V != 131 {
		t.Fatalf("expected volume to be 131 but got %+v", v.OHLCV.V)
	}

	v = s.GetValueForInterval(fivemin)
	if v == nil || v.OHLCV == nil {
		t.Fatalf("five min: expected ohlcv to be non nil but got %+v", v)
	} else if v.OHLCV.O != 13 {
		t.Fatalf("five min: expected open to be 13 but got %+v", v.OHLCV.O)
	} else if v.OHLCV.H != 18 {
		t.Fatalf("five min: expected high to be 18 but got %+v", v.OHLCV.H)
	} else if v.OHLCV.L != 10 {
		t.Fatalf("five min: expected low to be 10 but got %+v", v.OHLCV.L)
	} else if v.OHLCV.C != 15 {
		t.Fatalf("five min: expected close to be 15 but got %+v", v.OHLCV.C)
	} else if v.OHLCV.V != 12 {
		t.Fatalf("expected volume to be 12 but got %+v", v.OHLCV.V)
	}
}
