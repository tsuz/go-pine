package pine

import (
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestEMA(t *testing.T) {
	opts := SeriesOpts{
		Interval: 300,
		Max:      100,
	}
	now := time.Now()
	five := now.Add(5 * time.Minute)
	ten := now.Add(10 * time.Minute)
	fifteen := now.Add(15 * time.Minute)
	twenty := now.Add(20 * time.Minute)
	hl2 := NewOHLCProp(OHLCPropClose)
	ema := NewEMA(hl2, 4)
	data := []OHLCV{
		OHLCV{
			C: 3.1,
			S: now,
		},
		OHLCV{
			C: 2.8,
			S: five,
		},
		OHLCV{
			C: 3.4,
			S: ten,
		},
		OHLCV{
			C: 2.1,
			S: fifteen,
		},
		OHLCV{
			C: 9,
			S: twenty,
		},
	}
	ema1 := 2.85
	ema2 := 5.31
	io := []struct {
		time   time.Time
		output *float64
	}{
		{
			time:   now,
			output: nil,
		},
		{
			time:   five,
			output: nil,
		},
		{
			time:   ten,
			output: nil,
		},
		{
			time:   fifteen,
			output: &ema1,
		},
		{
			time:   twenty,
			output: &ema2,
		},
	}
	s, err := NewSeries(data, opts)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error init series"))
	}
	name := "ema"
	s.AddIndicator(name, ema)

	for i, o := range io {
		v := s.GetValueForInterval(o.time)
		if v.Indicators[name] == nil && o.output == nil {
			// ok
			continue
		}
		if v.Indicators[name] == nil || o.output == nil {
			t.Fatalf("expected both to be non nil but got %+v vs %+v at idx: %d", v.Indicators[name], o.output, i)
		}
		if !isWithin(*(v.Indicators[name]), *(o.output), 0.001) {
			t.Errorf("expected: %+v but got %+v for idx: %d", *(o.output), *(v.Indicators[name]), i)
		}
	}

}
