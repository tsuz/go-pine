package pine_test

import (
	pine "go-pine"
	"math"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestOHLCProp(t *testing.T) {
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
	io := []struct {
		prop   pine.OHLCProp
		output []float64
	}{
		{
			prop:   pine.OHLCPropOpen,
			output: []float64{14, 13},
		},
		{
			prop:   pine.OHLCPropHigh,
			output: []float64{15, 18},
		},
		{
			prop:   pine.OHLCPropLow,
			output: []float64{13, 10},
		},
		{
			prop:   pine.OHLCPropClose,
			output: []float64{14, 15},
		},
		{
			prop:   pine.OHLCPropVolume,
			output: []float64{131, 12},
		},
		{
			prop:   pine.OHLCPropHL2,
			output: []float64{14, 14},
		},
		{
			prop:   pine.OHLCPropHLC3,
			output: []float64{14, 14.3333},
		},
	}
	for i, o := range io {
		s, err := pine.NewSeries(data, opts)
		if err != nil {
			t.Fatal(errors.Wrap(err, "error init series"))
		}
		p := pine.NewOHLCProp(o.prop)
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
