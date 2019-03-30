package pine

import (
	"math"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestLinReg(t *testing.T) {
	opts := SeriesOpts{
		Interval: 300,
		Max:      100,
	}
	close := NewOHLCProp(OHLCPropClose)
	linreg := NewLinReg(close, 3)
	data := []OHLCV{
		OHLCV{
			C: 4119.05,
		},
		OHLCV{
			C: 4119.18,
		},
		OHLCV{
			C: 4118.02,
		},
		OHLCV{
			C: 4118.01,
		},
		OHLCV{
			C: 4118,
		},
		OHLCV{
			C: 4113.35,
		},
	}
	now := time.Now()
	for i := range data {
		data[i].S = now.Add(time.Duration(i*5) * time.Minute)
	}
	io := []struct {
		output map[string]float64
	}{
		{
			output: nil,
		},
		{
			output: nil,
		},
		{
			output: map[string]float64{
				"linreg": 4118.235,
			},
		},
		{
			output: map[string]float64{
				"linreg": 4117.818,
			},
		},
		{
			output: map[string]float64{
				"linreg": 4118,
			},
		},
		{
			output: map[string]float64{
				"linreg": 4114.123,
			},
		},
	}

	prettybad := 0.005
	s, err := NewSeries(data, opts)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error init series"))
	}
	s.AddIndicator("linreg", linreg)

	name := "linreg"
	for i, o := range io {
		v := s.GetValueForInterval(data[i].S)
		if v.Indicators[name] == nil && o.output == nil {
			// nil value
			continue
		}
		if v.Indicators[name] == nil && o.output == nil {
			// value is nil for that name
			continue
		}
		if v.Indicators[name] == nil || o.output == nil {
			t.Fatalf("expected both to be non nil but got %+v vs %+v at idx: %d for %s", v.Indicators[name], o.output, i, name)
		}
		if math.Abs(*(v.Indicators[name])/o.output[name]-1) >= prettybad {
			t.Errorf("expected: %+v but got %+v for idx: %d for %s", (o.output)[name], *(v.Indicators[name]), i, name)
		}
	}
}
