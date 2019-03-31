package pine

import (
	"time"

	"github.com/shopspring/decimal"

	"gonum.org/v1/gonum/stat"

	"github.com/pkg/errors"
)

type linreg struct {
	lastUpdate OHLCV
	lookback   int
	opts       *SeriesOpts
	genval     map[time.Time]*TimeValue
	srcval     map[time.Time]*TimeValue
	genvalues  []*TimeValue
	srcvalues  []*TimeValue
	src        Indicator
}

// NewLinReg creates a new LinReg indicator
func NewLinReg(i Indicator, lookback int) Indicator {
	return &linreg{
		src:       i,
		lookback:  lookback,
		genval:    make(map[time.Time]*TimeValue, lookback),
		srcval:    make(map[time.Time]*TimeValue, lookback*2),
		srcvalues: make([]*TimeValue, 0, lookback),
		genvalues: make([]*TimeValue, 0, lookback*2),
	}
}

func (i *linreg) GetValueForInterval(t time.Time) *Interval {
	v, ok := i.genval[t]
	if !ok {
		return nil
	}
	return &Interval{
		StartTime: t,
		Value:     v.Value,
	}
}

func (i *linreg) shouldUpdate(v OHLCV) bool {
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

func (i *linreg) generateAvg(t time.Time) error {
	total := len(i.srcvalues)
	if total < i.lookback {
		return nil
	}
	j := i.lookback*-1 + 1
	srcs := make([]float64, 0, total)
	x := make([]float64, 0, total)

	for _, v := range i.srcvalues {
		x = append(x, float64(j))
		srcs = append(srcs, v.Value)
		j++
		if j > 0 {
			break
		}
	}
	alpha, beta := stat.LinearRegression(x, srcs, nil, false)
	yhat := decimal.NewFromFloat(beta).Mul(
		decimal.NewFromFloat(0.0),
	).Add(
		decimal.NewFromFloat(alpha),
	)
	ypred, _ := yhat.Float64()
	tv := NewTimeValue(t, ypred)
	_, ok := i.genval[t]
	if !ok {
		if len(i.genvalues) == cap(i.genvalues) {
			var old *TimeValue
			old, i.genvalues = i.genvalues[0], i.genvalues[1:]
			delete(i.genval, old.Time)
		}
		i.genval[t] = tv
		i.genvalues = append(i.genvalues, tv)
	} else {
		i.genval[t] = tv
	}
	return nil
}

func (i *linreg) Update(v OHLCV) error {
	if err := i.src.Update(v); err != nil {
		return errors.Wrap(err, "error received from src in LinReg")
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
		if len(i.srcvalues) >= i.lookback {
			// remove first
			var old *TimeValue
			old, i.srcvalues = i.srcvalues[0], i.srcvalues[1:]
			if old == nil {
				return errors.New("unexpected nil in linreg array first el")
			}
			delete(i.srcval, old.Time)
		}
		i.srcval[v.S] = tv
		i.srcvalues = append(i.srcvalues, tv)
	} else if src.Value != val.Value {
		// source value has changed
		i.srcval[v.S].Value = val.Value
	}
	i.generateAvg(v.S)

	return nil
}

func (i *linreg) ApplyOpts(opts SeriesOpts) error {
	if opts.Max < i.lookback {
		return errors.New("SeriesOpts max cannot be less than LinReg lookback value")
	}
	if err := i.src.ApplyOpts(opts); err != nil {
		return errors.Wrap(err, "error applying opts in source")
	}
	return nil
}
