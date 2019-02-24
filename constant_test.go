package pine

import (
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestConstantInit(t *testing.T) {
	opts := SeriesOpts{
		Interval: 300,
		Max:      100,
	}

	now := time.Now()
	fivemin := now.Add(5 * time.Minute)
	constant := NewConstant(5.0)

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
	s, err := NewSeries(data, opts)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error init series"))
	}
	if err := s.AddIndicator("constant", constant); err != nil {
		t.Fatal(errors.Wrap(err, "expected constant to not error but errored"))
	}
	v := s.GetValueForInterval(now)
	if *(v.Indicators["constant"]) != 5.0 {
		t.Errorf("expected 5.0 but got %+v", *(v.Indicators["constant"]))
	}
	v = s.GetValueForInterval(fivemin)
	if *(v.Indicators["constant"]) != 5.0 {
		t.Errorf("expected 5.0 but got %+v", *(v.Indicators["constant"]))
	}
}
