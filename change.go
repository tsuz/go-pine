package pine

import (
	"time"

	"github.com/pkg/errors"
)

type chg struct {
	lastUpdate OHLCV
	lookback   int
	opts       *SeriesOpts
	chgopts    *ChangeOpts
	src        Indicator
}

type ChangeOpts struct {
	DiffType ChangeDiffType
}

type ChangeDiffType int

const (
	ChangeDiffTypeDiff ChangeDiffType = iota
	ChangeDiffTypeRatio
)

// NewChange creates a new Change indicator
func NewChange(i Indicator, lookback int, opts *ChangeOpts) Indicator {
	return &chg{
		chgopts:  opts,
		lookback: lookback,
		src:      i,
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
	var computed float64
	if i.chgopts != nil && i.chgopts.DiffType == ChangeDiffTypeRatio {
		computed = v1.Value / v2.Value
	} else {
		computed = v1.Value - v2.Value
	}
	return &Interval{
		StartTime: t,
		Value:     computed,
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
