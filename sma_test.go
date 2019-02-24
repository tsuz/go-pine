package pine_test

import (
	pine "go-pine"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestSMA(t *testing.T) {
	opts := pine.SeriesOpts{
		Interval: 300,
		Max:      100,
	}
	now := time.Now()
	fivemin := now.Add(5 * time.Minute)
	hl2 := pine.NewOHLCProp(pine.OHLCPropClose)
	sma1 := pine.NewSMA(hl2, 1)
	sma2 := pine.NewSMA(hl2, 2)
	sma3 := pine.NewSMA(hl2, 3)
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
	sma1res1 := 14.0
	sma1res2 := 15.0
	sma2res := 14.5
	io := []struct {
		name   string
		output []*float64
	}{
		{
			name:   "sma1",
			output: []*float64{&sma1res1, &sma1res2},
		},
		{
			name:   "sma2",
			output: []*float64{nil, &sma2res},
		},
		{
			name:   "sma3",
			output: []*float64{nil, nil},
		},
	}
	s, err := pine.NewSeries(data, opts)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error init series"))
	}
	s.AddIndicator("sma1", sma1)
	s.AddIndicator("sma2", sma2)
	s.AddIndicator("sma3", sma3)

	for i, o := range io {
		nowv := s.GetValueForInterval(now)
		if nowv.Indicators[o.name] == nil && o.output[0] == nil {
			// ok
			continue
		}
		if nowv.Indicators[o.name] == nil || o.output[0] == nil {
			t.Fatalf("expected both to be non nil but got %+v vs %+v at idx: %d", nowv.Indicators[o.name], o.output[0], i)
		}
		if *(nowv.Indicators[o.name]) != *(o.output[0]) {
			t.Errorf("expected: %+v but got %+v for idx: %d", *(o.output[0]), *(nowv.Indicators[o.name]), i)
		}

		fivv := s.GetValueForInterval(fivemin)
		if fivv.Indicators[o.name] == nil && o.output[1] == nil {
			// ok
			continue
		}
		if fivv.Indicators[o.name] == nil || o.output[1] == nil {
			t.Fatalf("expected both to be non nil but got %+v vs %+v at idx: %d", fivv.Indicators[o.name], o.output[1], i)
		}
		if *(fivv.Indicators[o.name]) != *(o.output[1]) {
			t.Errorf("expected: %+v but got %+v for idx: %d", *(o.output[1]), *(fivv.Indicators[o.name]), i)
		}
	}

}
