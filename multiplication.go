package pine

import (
	"time"

	"github.com/pkg/errors"
)

type mult struct {
	a Indicator
	b Indicator
}

// NewMultiplication multiplies the output of two indicators
func NewMultiplication(a Indicator, b Indicator) Indicator {
	return &mult{
		a: a,
		b: b,
	}
}

func (i *mult) GetValueForInterval(t time.Time) *Interval {
	// validate if needed
	a := i.a.GetValueForInterval(t)
	if a == nil {
		return nil
	}
	b := i.b.GetValueForInterval(t)
	if b == nil {
		return nil
	}
	v := a.Value * b.Value
	return &Interval{
		StartTime: t,
		Value:     v,
	}
}

func (i *mult) Update(v OHLCV) error {
	// validate if needed
	if err := i.a.Update(v); err != nil {
		return errors.Wrap(err, "error updating in mult")
	}
	if err := i.b.Update(v); err != nil {
		return errors.Wrap(err, "error updating in mult")
	}
	return nil
}

func (i *mult) ApplyOpts(opts SeriesOpts) error {
	// validate if needed
	if err := i.a.ApplyOpts(opts); err != nil {
		return errors.Wrap(err, "error applying opts in mult")
	}
	if err := i.b.ApplyOpts(opts); err != nil {
		return errors.Wrap(err, "error applying opts in mult")
	}
	return nil
}
