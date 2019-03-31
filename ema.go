package pine

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/pkg/errors"
)

type ema struct {
	lastUpdate OHLCV
	lookback   int
	opts       *SeriesOpts
	genval     map[time.Time]*TimeValue
	srcval     map[time.Time]*TimeValue
	genvalues  []*TimeValue
	srcvalues  []*TimeValue
	src        Indicator
}

// NewEMA creates a new EMA indicator
func NewEMA(i Indicator, lookback int) Indicator {
	return &ema{
		src:       i,
		lookback:  lookback,
		genval:    make(map[time.Time]*TimeValue, lookback),
		srcval:    make(map[time.Time]*TimeValue, lookback*2),
		srcvalues: make([]*TimeValue, 0, lookback),
		genvalues: make([]*TimeValue, 0, lookback*2),
	}
}

func (i *ema) GetValueForInterval(t time.Time) *Interval {
	v, ok := i.genval[t]
	if !ok {
		return nil
	}
	return &Interval{
		StartTime: t,
		Value:     v.Value,
	}
}

func (i *ema) shouldUpdate(v OHLCV) bool {
	downval := i.src.GetValueForInterval(v.S)
	if downval == nil {
		// src value does not exist so cannot generate
		return false
	}
	val, ok := i.srcval[v.S]
	if !ok {
		return true
	}
	if val.Value != downval.Value {
		// if src value is updated and our src value cache is not updated
		return true
	}
	return false
}

func (i *ema) generateEma(t time.Time) error {
	total := len(i.srcvalues)
	firstidx := total - i.lookback
	// not enough data
	if firstidx < 0 {
		return nil
	}
	tv := NewTimeValue(t, 0)
	if firstidx == 0 {
		// get SMA for initial value
		val := decimal.NewFromFloat(0.0)
		for j := firstidx; j < total; j++ {
			val = val.Add(decimal.NewFromFloat(i.srcvalues[j].Value))
		}
		avg := val.Div(decimal.NewFromFloat(float64(i.lookback)))
		tv.Value, _ = avg.Float64()
	} else if firstidx > 0 {
		// get previous value
		lastgen := len(i.genvalues) - 1
		if lastgen < 0 {
			return errors.New("Last gen value should have been there")
		}
		last := i.genvalues[lastgen]
		k := decimal.NewFromFloat(2.0).Div(decimal.NewFromFloat(float64(i.lookback + 1.0)))
		srcval := i.src.GetValueForInterval(t)
		if srcval == nil {
			return errors.New("srcval cannot be obtained in EMA")
		}
		tv.Value, _ = decimal.NewFromFloat(srcval.Value).
			Sub(decimal.NewFromFloat(last.Value)).
			Mul(k).
			Add(decimal.NewFromFloat(last.Value)).
			Float64()
	}
	if len(i.genvalues) == cap(i.genvalues) {
		var old *TimeValue
		old, i.genvalues = i.genvalues[0], i.genvalues[1:]
		delete(i.genval, old.Time)
	}
	i.genval[t] = tv
	i.genvalues = append(i.genvalues, tv)
	return nil
}

func (i *ema) Update(v OHLCV) error {
	if err := i.src.Update(v); err != nil {
		return errors.Wrap(err, "error received from src in EMA")
	}
	if !i.shouldUpdate(v) {
		return nil
	}
	i.lastUpdate = v

	// update src value and generate gen value
	val := i.src.GetValueForInterval(v.S)
	if val == nil {
		return errors.New("expected src to provide data but none was given")
	}
	src, ok := i.srcval[v.S]
	if !ok {
		tv := NewTimeValue(v.S, val.Value)
		if len(i.srcvalues) >= i.lookback*2 {
			// remove first
			var old *TimeValue
			old, i.srcvalues = i.srcvalues[0], i.srcvalues[1:]
			if old == nil {
				return errors.New("unexpected nil in ema array first el")
			}
			delete(i.srcval, old.Time)
		}
		i.srcval[v.S] = tv
		i.srcvalues = append(i.srcvalues, tv)
	} else if src.Value != val.Value {
		// source value has changed
		i.srcval[v.S].Value = src.Value
	}
	i.generateEma(v.S)

	return nil
}

func (i *ema) ApplyOpts(opts SeriesOpts) error {
	if opts.Max < i.lookback {
		return errors.New("SeriesOpts max cannot be less than EMA lookback value")
	}
	if err := i.src.ApplyOpts(opts); err != nil {
		return errors.Wrap(err, "error applying opts in source")
	}
	return nil
}
