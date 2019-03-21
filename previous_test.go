package pine

import (
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestPrevious(t *testing.T) {
	opts := SeriesOpts{
		Interval: 300,
		Max:      100,
	}
	now := time.Now()
	five := now.Add(5 * time.Minute)
	ten := now.Add(10 * time.Minute)
	hl2 := NewOHLCProp(OHLCPropClose)
	prevname := "prev-diff"
	prev := NewPrevious(hl2, 1)
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
	}
	chgd1 := 3.1
	chgd2 := 2.8
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
				prevname: &chgd1,
			},
		},
		{
			time: ten,
			output: map[string]*float64{
				prevname: &chgd2,
			},
		},
	}
	s, err := NewSeries(data, opts)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error init series"))
	}
	s.AddIndicator(prevname, prev)

	for i, o := range io {
		v := s.GetValueForInterval(o.time)
		for _, name := range []string{prevname} {
			if v.Indicators[name] == nil && o.output == nil {
				// ok
				continue
			}
			if v.Indicators[name] == nil || o.output == nil {
				t.Fatalf("expected both to be non nil but got %+v vs %+v at idx: %d", v.Indicators[name], o.output, i)
			}
			if *(v.Indicators[name]) != *(o.output[name]) {
				t.Errorf("expected: %+v but got %+v for idx: %d", *(o.output[name]), *(v.Indicators[name]), i)
			}
		}
	}
}
