package pine

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/pkg/errors"
)

type Series interface {
	AddIndicator(name string, i Indicator) error
	AddExec(v TPQ) error
	GetValueForInterval(t time.Time) *Interval
}

type Interval struct {
	StartTime  time.Time
	OHLCV      *OHLCV
	Value      float64
	Indicators map[string]*float64
}

type series struct {
	items    map[string]Indicator
	lastExec TPQ
	lastOHLC *OHLCV
	opts     SeriesOpts
	values   []OHLCV
	timemap  map[time.Time]*OHLCV
}

func (s *series) getLatestInterval() *OHLCV {
	return s.lastOHLC
}

func NewSeries(ohlcv []OHLCV, opts SeriesOpts) (Series, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "error validating seriesopts")
	}
	tm := make(map[time.Time]*OHLCV)
	s := &series{
		items:   make(map[string]Indicator),
		opts:    opts,
		timemap: tm,
		values:  make([]OHLCV, 0, opts.Max),
	}
	s.initValues(ohlcv)
	return s, nil
}

func (s *series) initValues(values []OHLCV) {
	for _, v := range values {
		s.insertInterval(v)
	}
	// last := len(s.values)
	// lastidx := last - 1
	// for i := 0; i < last; i++ {
	// 	v := s.values[i]
	// 	t := s.getLastIntervalFromTime(v.S)
	// 	v.S = t
	// 	s.timemap[t] = &v
	// 	if i == lastidx {
	// 		s.lastOHLC = &v
	// 	}
	// }
}

func (s *series) insertInterval(v OHLCV) {
	log.Printf("v.S %+v", v.S)
	t := s.getLastIntervalFromTime(v.S)
	v.S = t
	_, ok := s.timemap[t]
	log.Printf("insertInterval %+v %+v %v", t, v, ok)
	if !ok {
		s.values = append(s.values, v)
		s.timemap[t] = &v
		s.lastOHLC = &v
	}
}

func (s *series) AddIndicator(name string, i Indicator) error {
	// enforce series constraint
	i.ApplyOpts(s.opts)
	// update with current values downstream
	for _, v := range s.values {
		if err := i.Update(v); err != nil {
			return errors.Wrap(err, "error updating indicator")
		}
	}
	s.items[name] = i
	return nil
}

func (s *series) getLastIntervalFromTime(t time.Time) time.Time {
	year, month, day := t.UTC().Date()
	st := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	m := s.getMultiplierDiff(t, st)
	return st.Add(time.Duration(m*s.opts.Interval) * time.Second)
}

func (s *series) getMultiplierDiff(t time.Time, st time.Time) int {
	diff := t.Sub(st).Seconds()
	return int(diff / float64(s.opts.Interval))
}

func (s *series) AddExec(v TPQ) error {
	log.Printf("Add exec %+v", v)
	start := s.getLastIntervalFromTime(v.Timestamp)
	if s.lastOHLC == nil {
		// create first one
		ohlcv := NewOHLCVWithSamePx(v.Px, v.Qty, start)
		s.insertInterval(ohlcv)
	} else if s.lastOHLC.S.Equal(start) {
		// update this interval
		itvl := s.lastOHLC
		itvl.C = v.Px
		itvl.V += v.Qty
		if v.Px > itvl.H {
			itvl.H = v.Px
		} else if v.Px < itvl.L {
			itvl.L = v.Px
		}
	} else if start.Sub(s.lastOHLC.S).Seconds() > 0 {
		// calculate how many intervals are missing
		switch s.opts.EmptyInst {
		case EmptyInstUseLastClose:
			// figure out how many
			m := s.getMultiplierDiff(v.Timestamp, s.lastOHLC.S)
			usenew := int(math.Max(float64(m-1), float64(0)))
			for i := 0; i < m; i++ {
				var px, qty float64
				if i == usenew {
					px = v.Px
					qty = v.Qty
				} else {
					px = s.lastOHLC.C
					qty = 0
				}
				newt := s.lastOHLC.S.Add(time.Duration(s.opts.Interval) * time.Second)
				ohlcv := NewOHLCVWithSamePx(px, qty, newt)
				s.insertInterval(ohlcv)
			}
		case EmptyInstUseZeros:
			// figure out how many
			m := s.getMultiplierDiff(v.Timestamp, s.lastOHLC.S)
			usenew := int(math.Max(float64(m-1), float64(0)))
			for i := 0; i < m; i++ {
				var px, qty float64
				if i == usenew {
					px = v.Px
					qty = v.Qty
				} else {
					px = 0
					qty = 0
				}
				newt := s.lastOHLC.S.Add(time.Duration(s.opts.Interval) * time.Second)
				ohlcv := NewOHLCVWithSamePx(px, qty, newt)
				s.insertInterval(ohlcv)
			}
		case EmptyInstIgnore:
			v := NewOHLCVWithSamePx(v.Px, v.Qty, start)
			s.insertInterval(v)
		default:
			return fmt.Errorf("Unsupported interval: %+v", s.opts.EmptyInst)
		}
	}
	s.lastExec = v
	return nil
}

func (s *series) updatePoint(v TPQ, st time.Time) {
	if s.lastOHLC != nil && s.lastOHLC.S.Equal(st) {
		s.lastOHLC.S = st
		s.lastOHLC.V += v.Qty
		s.lastOHLC.C = v.Px
		if s.lastOHLC.H < v.Px {
			s.lastOHLC.H = v.Px
		} else if s.lastOHLC.L > v.Px {
			s.lastOHLC.L = v.Px
		}
	} else {
		newv := OHLCV{
			O: v.Px,
			L: v.Px,
			H: v.Px,
			C: v.Px,
			V: v.Qty,
			S: v.Timestamp,
		}
		var old OHLCV
		if len(s.values) == s.opts.Max {
			log.Printf("Deleting as %+v %+v", s.values, s.opts.Max)
			old, s.values = s.values[0], s.values[1:]
			delete(s.timemap, old.S)
		}
		s.values = append(s.values, newv)
		s.lastOHLC = &newv
		s.timemap[st] = &newv
	}
}

func (s *series) GetValueForInterval(t time.Time) *Interval {
	log.Printf("GetValueForInterval %+v %+v", t, s.lastOHLC.S)
	if s.lastOHLC == nil {
		return nil
	}
	if !s.lastOHLC.S.Equal(t) {
		// if time is within interval, adjust it
		t = s.getLastIntervalFromTime(t)
		log.Printf("getLastIntervalFromTime %+v", t)
	}
	inds := make(map[string]*float64)
	for k, v := range s.items {
		val := v.GetValueForInterval(t)
		inds[k] = &val.Value
	}
	v, ok := s.timemap[t]
	if !ok {
		log.Printf("Not ok for t %+v", t)
		return nil
	} else {
		log.Printf("ok")
	}
	return &Interval{
		StartTime:  v.S,
		OHLCV:      v,
		Indicators: inds,
	}
}

func (s *series) getOHLCV(t time.Time) *OHLCV {
	return s.timemap[t]
}
