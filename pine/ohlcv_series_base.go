package pine

import (
	"math"

	"github.com/pkg/errors"
)

// OHLCVBaseSeries represents a series of OHLCV type (i.e. open, high, low, close, volume)
type OHLCVBaseSeries interface {
	Push(OHLCV)

	Shift() bool

	Len() int

	// Current returns current ohlcv
	Current() *OHLCV

	// GetSeries returns series of values for a property
	GetSeries(OHLCProp) ValueSeries

	// GoToFirst sets the current value to first and returns that value
	GoToFirst() *OHLCV

	// Next moves the pointer to the next one.
	// If there is no next item, nil is returned and the pointer does not advance.
	// If there is no next item and a data source registered, it will attempt to fetch and append items if there are any
	Next() (*OHLCV, error)

	// registers data source for dynamic updates
	RegisterDataSource(DataSource)

	// set the maximum number of OHLCV items. This helps prevent high memory usage.
	SetMax(int64)
}

func NewOHLCVBaseSeries() OHLCVBaseSeries {
	s := &ohlcvBaseSeries{
		vals: make(map[int64]OHLCV),
	}
	return s
}

type ohlcvBaseSeries struct {
	ds DataSource

	// current ohlcv
	cur *OHLCV

	first *OHLCV

	last *OHLCV

	// max number of candles. 0 means no limit. Defaults to 0
	max int64

	vals map[int64]OHLCV
}

func (s *ohlcvBaseSeries) Push(o OHLCV) {
	s.vals[o.S.Unix()] = o
	if s.last != nil {
		o.prev = s.last
		s.last.next = &o
	}
	s.last = &o
	if s.first == nil {
		s.first = &o
	}
	s.resize(s.max)
}

func (s *ohlcvBaseSeries) Shift() bool {
	if s.first == nil {
		return false
	}
	delete(s.vals, s.first.S.Unix())
	s.first = s.first.next
	return true
}

func (s *ohlcvBaseSeries) Len() int {
	return len(s.vals)
}

func (s *ohlcvBaseSeries) Current() *OHLCV {
	return s.cur
}

func (s *ohlcvBaseSeries) GoToFirst() *OHLCV {
	s.cur = s.first
	return s.cur
}

func (s *ohlcvBaseSeries) fetchAndAppend() (bool, error) {
	more, err := s.ds.Populate(s.cur.S)
	if err != nil {
		return false, errors.Wrap(err, "error populating")
	}
	for _, v := range more {
		s.Push(v)
	}
	return len(more) > 0, nil
}

func (s *ohlcvBaseSeries) Next() (*OHLCV, error) {
	if s.cur == nil {
		if len(s.vals) == 0 {
			return nil, nil
		}
		// set first one if nil
		s.cur = s.first
		return s.cur, nil
	}
	if s.cur.next == nil {
		if s.ds != nil {
			found, err := s.fetchAndAppend()
			if err != nil {
				return nil, errors.Wrap(err, "error populating")
			}
			if !found {
				return nil, nil
			}
			return s.Next()
		}
		return nil, nil
	}
	s.cur = s.cur.next
	return s.cur, nil
}

func (s *ohlcvBaseSeries) RegisterDataSource(ds DataSource) {
	s.ds = ds
}

func (s *ohlcvBaseSeries) GetSeries(p OHLCProp) ValueSeries {
	vs := NewValueSeries()
	v := s.first
	for {
		if v == nil {
			break
		}
		var propVal float64
		switch p {
		case OHLCPropClose:
			propVal = v.C
		case OHLCPropOpen:
			propVal = v.O
		case OHLCPropHigh:
			propVal = v.H
		case OHLCPropLow:
			propVal = v.L
		case OHLCPropVolume:
			propVal = v.V
		case OHLCPropTR:
			if v.prev != nil {
				p := v.prev
				propVal = math.Max(
					math.Abs(v.H-v.L),
					math.Max(
						math.Abs(v.H-p.C),
						math.Abs(v.L-p.C)))
			} else {
				propVal = math.Abs(v.H - v.L)
			}

		default:
			continue
		}
		vs.Set(v.S, propVal)
		v = v.next
	}

	if s.cur != nil {
		vs.SetCurrent(s.cur.S)
	}
	return vs
}

func (s *ohlcvBaseSeries) SetMax(m int64) {

	// set upon exit
	defer func() {
		s.max = m
	}()

	s.resize(m)
}

func (s *ohlcvBaseSeries) resize(m int64) {
	// set to unlimited, nothing to perform
	if m == 0 {
		return
	}
	for {
		if int64(s.Len()) <= m {
			break
		}
		s.Shift()
	}
}
