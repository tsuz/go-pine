package pine

import (
	"math"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestOHLCProp(t *testing.T) {
	opts := SeriesOpts{
		Interval: 300,
		Max:      100,
	}
	now := time.Now()
	fivemin := now.Add(5 * time.Minute)
	data := []OHLCV{
		OHLCV{
			O: 14,
			H: 15,
			L: 13,
			C: 14,
			V: 131,
			S: now,
		},
		OHLCV{
			O: 13,
			H: 18,
			L: 10,
			C: 15,
			V: 12,
			S: fivemin,
		},
	}
	io := []struct {
		prop   OHLCProp
		output []float64
	}{
		{
			prop:   OHLCPropOpen,
			output: []float64{14, 13},
		},
		{
			prop:   OHLCPropHigh,
			output: []float64{15, 18},
		},
		{
			prop:   OHLCPropLow,
			output: []float64{13, 10},
		},
		{
			prop:   OHLCPropClose,
			output: []float64{14, 15},
		},
		{
			prop:   OHLCPropVolume,
			output: []float64{131, 12},
		},
		{
			prop:   OHLCPropHL2,
			output: []float64{14, 14},
		},
		{
			prop:   OHLCPropHLC3,
			output: []float64{14, 14.3333},
		},
	}
	for i, o := range io {
		s, err := NewSeries(data, opts)
		if err != nil {
			t.Fatal(errors.Wrap(err, "error init series"))
		}
		p := NewOHLCProp(o.prop)
		s.AddIndicator("val", p)
		nowv := s.GetValueForInterval(now)
		if !isWithin(*(nowv.Indicators["val"]), o.output[0], 0.001) {
			t.Errorf("expected: %+v but got %+v for idx: %d, first val", o.output[0], *(nowv.Indicators["val"]), i)
		}
		fivev := s.GetValueForInterval(fivemin)
		if !isWithin(*(fivev.Indicators["val"]), o.output[1], 0.001) {
			t.Errorf("expected: %+v but got %+v for idx: %d, seocnd val", o.output[1], *(fivev.Indicators["val"]), i)
		}
	}
}

func isWithin(a, b, e float64) bool {
	return math.Abs(a-b) < e
}
