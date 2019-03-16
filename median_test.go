package pine

import (
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestMedian(t *testing.T) {
	opts := SeriesOpts{
		Interval: 300,
		Max:      100,
	}
	now := time.Now()
	five := now.Add(5 * time.Minute)
	ten := now.Add(10 * time.Minute)
	high := NewOHLCProp(OHLCPropHigh)
	low := NewOHLCProp(OHLCPropLow)
	sub := NewArithmetic(ArithmeticSubtraction, high, low, ArithmeticOpts{})
	even := "even"
	odd := "odd"
	indeven := NewMedian(sub, 2)
	indodd := NewMedian(sub, 3)
	data := []OHLCV{
		OHLCV{
			C: 3.1,
			H: 3.3,
			L: 3.0,
			S: now,
		},
		OHLCV{
			C: 2.8,
			H: 2.8,
			L: 2.4,
			S: five,
		},
		OHLCV{
			C: 3.4,
			H: 3.5,
			L: 2.8,
			S: ten,
		},
	}
	chg1even := 0.35
	chg2even := 0.55
	chg2odd := 0.4
	io := []struct {
		time   time.Time
		output map[string]*float64
	}{
		{
			time:   now,
			output: nil,
		},
		{
			time: five,
			output: map[string]*float64{
				"even": &chg1even,
				"odd":  nil,
			},
		},
		{
			time: ten,
			output: map[string]*float64{
				"even": &chg2even,
				"odd":  &chg2odd,
			},
		},
	}
	s, err := NewSeries(data, opts)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error init series"))
	}
	s.AddIndicator(even, indeven)
	s.AddIndicator(odd, indodd)

	for i, o := range io {
		v := s.GetValueForInterval(o.time)
		for _, name := range []string{even, odd} {
			if v.Indicators[name] == nil && o.output == nil {
				// nil value
				continue
			}
			if v.Indicators[name] == nil && (o.output)[name] == nil {
				// value is nil for that name
				continue
			}
			if v.Indicators[name] == nil || o.output == nil {
				t.Fatalf("expected both to be non nil but got %+v vs %+v at idx: %d for %s", v.Indicators[name], o.output, i, name)
			}
			if *(v.Indicators[name]) != *(o.output)[name] {
				t.Errorf("expected: %+v but got %+v for idx: %d for %s", *(o.output)[name], *(v.Indicators[name]), i, name)
			}
		}
	}
}
