package pine

import (
	"math"

	"github.com/pkg/errors"
	"github.com/twinj/uuid"
)

// OHLCVBaseSeries represents a series of OHLCV type (i.e. open, high, low, close, volume)
type OHLCVBaseSeries interface {
	ID() string

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
	u := uuid.NewV4()
	s := &ohlcvBaseSeries{
		id:   u.String(),
		max:  1000, // default maximum items
		vals: make(map[int64]OHLCV),
	}
	return s
}

type ohlcvBaseSeries struct {
	ds DataSource

	// current ohlcv
	cur *OHLCV

	id string

	first *OHLCV

	last *OHLCV

	// max number of candles. 0 means no limit. Defaults to 1000
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
	s.resize()
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

func (s *ohlcvBaseSeries) ID() string {
	return s.id
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
		var propVal *float64
		switch p {
		case OHLCPropClose:
			propVal = NewFloat64(v.C)
		case OHLCPropOpen:
			propVal = NewFloat64(v.O)
		case OHLCPropHigh:
			propVal = NewFloat64(v.H)
		case OHLCPropLow:
			propVal = NewFloat64(v.L)
		case OHLCPropVolume:
			propVal = NewFloat64(v.V)
		case OHLCPropTR, OHLCPropTRHL:
			if v.prev != nil {
				p := v.prev
				v1 := math.Abs(v.H - v.L)
				v2 := math.Abs(v.H - p.C)
				v3 := math.Abs(v.L - p.C)
				v := math.Max(v1, math.Max(v2, v3))
				propVal = NewFloat64(v)
			}
			if p == OHLCPropTRHL && v.prev == nil {
				d := v.H - v.L
				propVal = &d
			}
		case OHLCPropHLC3:
			propVal = NewFloat64((v.H + v.L + v.C) / 3)
		default:
			continue
		}
		if propVal != nil {
			vs.Set(v.S, *propVal)
		}
		v = v.next
	}

	if s.cur != nil {
		vs.SetCurrent(s.cur.S)
	}
	return vs
}

func (s *ohlcvBaseSeries) SetMax(m int64) {

	s.max = m

	s.resize()
}

func (s *ohlcvBaseSeries) resize() {
	m := s.max
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
