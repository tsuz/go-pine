package pine

import (
	"time"

	"github.com/pkg/errors"
)

type chg struct {
	lastUpdate OHLCV
	lookback   int
	opts       *SeriesOpts
	src        Indicator
}

// NewChange creates a new Change indicator
func NewChange(i Indicator, lookback int) Indicator {
	return &chg{
		src:      i,
		lookback: lookback,
	}
}

func (i *chg) GetValueForInterval(t time.Time) *Interval {
	v1 := i.src.GetValueForInterval(t)
	if v1 == nil {
		return nil
	}
	v2 := i.src.GetValueForInterval(t.Add(-1 * time.Duration(i.lookback*i.opts.Interval) * time.Second))
	if v2 == nil {
		// handle empty case values
		return nil
	}
	return &Interval{
		StartTime: t,
		Value:     v1.Value - v2.Value,
	}
}

func (i *chg) Update(v OHLCV) error {
	if err := i.src.Update(v); err != nil {
		return errors.Wrap(err, "error received from src in Change")
	}
	return nil
}

func (i *chg) ApplyOpts(opts SeriesOpts) error {
	if opts.Max < i.lookback {
		return errors.New("SeriesOpts max cannot be less than Change lookback value")
	}
	i.opts = &opts
	return nil
}
