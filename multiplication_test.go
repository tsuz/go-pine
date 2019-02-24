package pine_test

import (
	pine "go-pine"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestMultiplication(t *testing.T) {
	opts := pine.SeriesOpts{
		Interval: 300,
		Max:      100,
	}
	now := time.Now()
	fivemin := now.Add(5 * time.Minute)
	hl2 := pine.NewOHLCProp(pine.OHLCPropOpen)
	sma := pine.NewSMA(hl2, 1)
	c := pine.NewConstant(182)
	add := pine.NewMultiplication(sma, c)

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
		t.Fatal(errors.Wrap(err, "error init series"))
	}
	name := "sma+const"
	s.AddIndicator(name, add)

	nowv := s.GetValueForInterval(now)

	if *(nowv.Indicators[name]) != 2548 {
		t.Errorf("expected: 2548 but got %+v", *(nowv.Indicators[name]))
	}
	fivv := s.GetValueForInterval(fivemin)
	if *(fivv.Indicators[name]) != 2366 {
		t.Errorf("expected: 2366 but got %+v", *(fivv.Indicators[name]))
	}
}
