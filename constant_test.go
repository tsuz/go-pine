package pine_test

import (
	pine "go-pine"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestConstantInit(t *testing.T) {
	opts := pine.SeriesOpts{
		Interval: 300,
		Max:      100,
	}

	now := time.Now()
	fivemin := now.Add(5 * time.Minute)
	constant := pine.NewConstant(5.0)

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
