package pine_test

import (
	pine "go-pine"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestAddition(t *testing.T) {
	opts := pine.SeriesOpts{
		Interval: 300,
		Max:      100,
	}
	now := time.Now()
	fivemin := now.Add(5 * time.Minute)
	hl2 := pine.NewOHLCProp(pine.OHLCPropOpen)
	sma := pine.NewSMA(hl2, 1)
	c := pine.NewConstant(5.0)
	add := pine.NewAddition(sma, c)

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

	if *(nowv.Indicators[name]) != 19.0 {
		t.Errorf("expected: 19.0 but got %+v", *(nowv.Indicators[name]))
	}
	fivv := s.GetValueForInterval(fivemin)
	if *(fivv.Indicators[name]) != 18.0 {
		t.Errorf("expected: 18.0 but got %+v", *(fivv.Indicators[name]))
	}
}
