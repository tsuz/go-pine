package pine

import (
	"log"
	"math"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestStdDev(t *testing.T) {
	itvl := 300
	opts := SeriesOpts{
		Interval: itvl,
		Max:      32,
	}
	name := "stddev"
	sdTests := struct {
		candles  []OHLCV
		expected []*Interval
	}{
		candles: []OHLCV{
			OHLCV{C: 52.22},
			OHLCV{C: 52.78},
			OHLCV{C: 53.02},
			OHLCV{C: 53.67},
			OHLCV{C: 53.67},
			OHLCV{C: 53.74},
			OHLCV{C: 53.45},
			OHLCV{C: 53.72},
			OHLCV{C: 53.39},
			OHLCV{C: 52.51},
			OHLCV{C: 52.32},
			OHLCV{C: 51.45},
			OHLCV{C: 51.60},
			OHLCV{C: 52.43},
			OHLCV{C: 52.47},
			OHLCV{C: 52.91},
			OHLCV{C: 52.07},
			OHLCV{C: 53.12},
			OHLCV{C: 52.77},
			OHLCV{C: 52.73},
			OHLCV{C: 52.09},
			OHLCV{C: 53.19},
			OHLCV{C: 53.73},
			OHLCV{C: 53.87},
			OHLCV{C: 53.85},
			OHLCV{C: 53.88},
			OHLCV{C: 54.08},
			OHLCV{C: 54.14},
			OHLCV{C: 54.50},
			OHLCV{C: 54.30},
			OHLCV{C: 54.40},
			OHLCV{C: 54.16},
		},
		expected: []*Interval{
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			&Interval{Value: 0.523018},
			&Interval{Value: 0.505411},
			&Interval{Value: 0.730122},
			&Interval{Value: 0.857364},
			&Interval{Value: 0.833642},
			&Interval{Value: 0.788707},
			&Interval{Value: 0.716251},
			&Interval{Value: 0.675498},
			&Interval{Value: 0.584679},
			&Interval{Value: 0.507870},
			&Interval{Value: 0.518353},
			&Interval{Value: 0.526061},
			&Interval{Value: 0.480964},
			&Interval{Value: 0.490176},
			&Interval{Value: 0.578439},
			&Interval{Value: 0.622905},
			&Interval{Value: 0.670093},
			&Interval{Value: 0.622025},
			&Interval{Value: 0.661064},
			&Interval{Value: 0.690358},
			&Interval{Value: 0.651152},
			&Interval{Value: 0.360466},
			&Interval{Value: 0.242959},
		},
	}

	prettybad := 0.005
	now := time.Now()
	for idx := range sdTests.candles {
		t := now.Add(time.Duration(idx*itvl) * time.Second)
		sdTests.candles[idx].S = t
	}

	s, err := NewSeries(sdTests.candles, opts)
	if err != nil {
		t.Fatal(err)
	}
	close := NewOHLCProp(OHLCPropClose)
	stddev := NewStdDev(close, 10)
	if err := s.AddIndicator(name, stddev); err != nil {
		t.Fatal(err)
	}
	for idx, exp := range sdTests.expected {
		log.Println(idx)
		tim := now.Add(time.Duration(idx*itvl) * time.Second)
		v := s.GetValueForInterval(tim)
		if v == nil {
			t.Fatal(errors.Wrap(err, "interval should not be nil"))
		}
		if exp == nil && v.Indicators[name] == nil {
			continue // ok
		}
		if exp == nil && v.Indicators[name] != nil {
			t.Errorf("expected v to be nil but got %+v at idx: %d", v, idx)
		}
		if v.Indicators[name] == nil {
			t.Errorf("expected indicator to have value but got none at idx %d", idx)
		} else if math.Abs(exp.Value-*v.Indicators[name])/exp.Value > prettybad {
			t.Errorf("expected %+v but got %+v for idx: %d", exp.Value, *v.Indicators[name], idx)
		}
	}
}
