package pine

import (
	"time"

	"github.com/pkg/errors"
)

type prev struct {
	lastUpdate OHLCV
	lookback   int
	opts       *SeriesOpts
	src        Indicator
}

// NewPrevious looks back at previous intervals values
func NewPrevious(i Indicator, lookback int) Indicator {
	return &prev{
		lookback: lookback,
		src:      i,
	}
}

func (i *prev) GetValueForInterval(t time.Time) *Interval {
	v2 := i.src.GetValueForInterval(t.Add(-1 * time.Duration(i.lookback*i.opts.Interval) * time.Second))
	if v2 == nil {
		return nil
	}
	return &Interval{
		StartTime: t,
		Value:     v2.Value,
	}
}

func (i *prev) Update(v OHLCV) error {
	if err := i.src.Update(v); err != nil {
		return errors.Wrap(err, "error received from src in Change")
	}
	return nil
}

func (i *prev) ApplyOpts(opts SeriesOpts) error {
	if opts.Max < i.lookback {
		return errors.New("SeriesOpts max cannot be less than Change lookback value")
	}
	i.opts = &opts
	return nil
}
