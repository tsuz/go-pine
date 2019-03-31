package pine

import (
	"log"
	"math"
	"time"

	"github.com/shopspring/decimal"

	"github.com/pkg/errors"
)

// NewStdDev creates a new standard deviation indicator
func NewStdDev(i Indicator, lookback int) Indicator {
	return &stddev{
		src:       i,
		lookback:  lookback,
		genval:    make(map[time.Time]*TimeValue, lookback),
		srcval:    make(map[time.Time]*TimeValue, lookback*2),
		srcvalues: make([]*TimeValue, 0, lookback),
		genvalues: make([]*TimeValue, 0, lookback*2),
	}
}

type stddev struct {
	lastUpdate OHLCV
	lookback   int
	opts       *SeriesOpts
	genval     map[time.Time]*TimeValue
	srcval     map[time.Time]*TimeValue
	genvalues  []*TimeValue
	srcvalues  []*TimeValue
	src        Indicator
}

func (i *stddev) GetValueForInterval(t time.Time) *Interval {
	v, ok := i.genval[t]
	if !ok {
		return nil
	}
	return &Interval{
		StartTime: t,
		Value:     v.Value,
	}
}

func (i *stddev) shouldUpdate(v OHLCV) bool {
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

func (i *stddev) Update(v OHLCV) error {
	if err := i.src.Update(v); err != nil {
		return errors.Wrap(err, "error received from src in StdDev")
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
				return errors.New("unexpected nil in StdDev array first el")
			}
			delete(i.srcval, old.Time)
		}
		i.srcval[v.S] = tv
		i.srcvalues = append(i.srcvalues, tv)
	} else if src.Value != val.Value {
		// source value has changed
		i.srcval[v.S].Value = val.Value
	}
	i.generateStdDev(v.S)

	return nil
}

func (i *stddev) ApplyOpts(opts SeriesOpts) error {
	if opts.Max < i.lookback {
		return errors.New("SeriesOpts max cannot be less than StdDev lookback value")
	}
	if err := i.src.ApplyOpts(opts); err != nil {
		return errors.Wrap(err, "error applying opts in source")
	}
	if i.opts == nil || i.opts.Max != opts.Max {
		i.genval = make(map[time.Time]*TimeValue, opts.Max)
		i.srcval = make(map[time.Time]*TimeValue, i.lookback+opts.Max)
		i.srcvalues = make([]*TimeValue, 0, i.lookback+opts.Max)
		i.genvalues = make([]*TimeValue, 0, opts.Max)
	}
	i.opts = &opts
	return nil
}

func (i *stddev) generateStdDev(t time.Time) error {
	total := len(i.srcvalues)
	firstidx := total - i.lookback
	log.Printf("total %+v lookback %+v %t", total, i.lookback, firstidx < 0)
	if firstidx < 0 {
		return nil
	}

	avgtot := decimal.NewFromFloat(0.0)
	for j := firstidx; j < total; j++ {
		avgtot = avgtot.Add(decimal.NewFromFloat(i.srcvalues[j].Value))
	}
	avg := avgtot.Div(decimal.NewFromFloat(float64(i.lookback)))
	stddevtot := decimal.NewFromFloat(0.0)
	for j := firstidx; j < total; j++ {
		diff := decimal.NewFromFloat(i.srcvalues[j].Value).Sub(avg)
		t := diff.Pow(decimal.NewFromFloat(2.0))
		stddevtot = stddevtot.Add(t)
	}
	stddevtot = stddevtot.Div(decimal.NewFromFloat(float64(i.lookback)))
	tot, _ := stddevtot.Float64()
	stddev := math.Sqrt(tot)
	tv := NewTimeValue(t, stddev)
	_, ok := i.genval[t]
	if !ok {
		if len(i.genvalues) == cap(i.genvalues) {
			var old *TimeValue
			old, i.genvalues = i.genvalues[0], i.genvalues[1:]
			delete(i.genval, old.Time)
			log.Printf("removing %+v", old)
		}
		i.genval[t] = tv
		i.genvalues = append(i.genvalues, tv)
	} else {
		i.genval[t] = tv
	}
	log.Printf("end val %+v %+v %+v", stddev, t, tv)
	return nil
}

const SqrtMaxIter = 100000

// SqrtRound returns the square root of d, the result will have
// precision digits after the decimal point. The bool precise returns whether the precision was reached
func SqrtRound(d decimal.Decimal, precision int32) (decimal.Decimal, bool) {
	cutoff := decimal.New(1, -precision)
	lo := decimal.Zero
	hi := d
	var mid decimal.Decimal
	for i := 0; i < SqrtMaxIter; i++ {
		//mid = (lo+hi)/2;
		mid = lo.Add(hi).DivRound(decimal.New(2, 0), precision)
		if mid.Mul(mid).Sub(d).Abs().LessThan(cutoff) {
			return mid, true
		}
		if mid.Mul(mid).GreaterThan(d) {
			hi = mid
		} else {
			lo = mid
		}
	}
	return mid, false
}
