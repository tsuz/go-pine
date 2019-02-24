package pine

import (
	"time"

	"github.com/pkg/errors"
)

type addition struct {
	a Indicator
	b Indicator
}

func NewAddition(a Indicator, b Indicator) Indicator {
	return &addition{
		a: a,
		b: b,
	}
}

func (i *addition) GetValueForInterval(t time.Time) *Interval {
	// validate if needed
	a := i.a.GetValueForInterval(t)
	if a == nil {
		return nil
	}
	b := i.b.GetValueForInterval(t)
	if b == nil {
		return nil
	}
	v := a.Value + b.Value
	return &Interval{
		StartTime: t,
		Value:     v,
	}
}

func (i *addition) Update(v OHLCV) error {
	// validate if needed
	if err := i.a.Update(v); err != nil {
		return errors.Wrap(err, "error updating in addition")
	}
	if err := i.b.Update(v); err != nil {
		return errors.Wrap(err, "error updating in addition")
	}
	return nil
}

func (i *addition) ApplyOpts(opts SeriesOpts) error {
	// validate if needed
	if err := i.a.ApplyOpts(opts); err != nil {
		return errors.Wrap(err, "error applying opts in addition")
	}
	if err := i.b.ApplyOpts(opts); err != nil {
		return errors.Wrap(err, "error applying opts in addition")
	}
	return nil
}
