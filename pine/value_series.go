package pine

import (
	"sync"
	"time"

	"github.com/twinj/uuid"
)

type ValueSeries interface {
	ID() string
	Add(ValueSeries) ValueSeries
	AddConst(float64) ValueSeries
	Div(ValueSeries) ValueSeries
	DivConst(float64) ValueSeries
	Mul(ValueSeries) ValueSeries
	MulConst(float64) ValueSeries
	Sub(ValueSeries) ValueSeries
	SubConst(float64) ValueSeries

	// Get gets the item by time in value series
	Get(time.Time) *Value
	// GetLast gets the last item in value series
	GetLast() *Value
	// GetFirst gets the first item in value series
	GetFirst() *Value

	// Gets size of the ValueSeries
	Len() int

	Copy() ValueSeries

	Set(time.Time, float64)
	SetAll(val float64)

	Shift() bool

	Val() *float64
	SetCurrent(time.Time) bool
	GetCurrent() *Value

	// set the maximum number of items.
	// This helps prevent allocating too much memory
	SetMax(int64)
}

type valueSeries struct {
	id    string
	cur   *Value
	first *Value
	last  *Value
	// max number of candles. 0 means no limit. Defaults to 1000
	max int64
	sync.Mutex
	timemap map[int64]*Value
}

type Value struct {
	t    time.Time
	v    float64
	prev *Value
	next *Value
}

// NewValueSeries creates an empty series that conforms to ValueSeries
func NewValueSeries() ValueSeries {
	u := uuid.NewV4()
	v := &valueSeries{
		id:      u.String(),
		max:     1000, // default maximum items
		timemap: make(map[int64]*Value),
	}
	return v
}

func (s *valueSeries) Copy() ValueSeries {
	newv := NewValueSeries()
	f := s.GetFirst()
	for {
		if f == nil {
			break
		}
		newv.Set(f.t, f.v)
		f = f.next
	}
	cur := s.GetCurrent()
	if cur != nil {
		newv.SetCurrent(cur.t)
	}
	return newv
}

func (s *valueSeries) operation(v ValueSeries, op func(a, b float64) float64) ValueSeries {
	copied := NewValueSeries()
	f := s.GetFirst()
	for {
		if f == nil {
			break
		}
		newv := v.Get(f.t)
		if newv != nil {
			copied.Set(f.t, op(f.v, newv.v))
		}
		f = f.next
	}
	cur := s.GetCurrent()
	if cur != nil {
		copied.SetCurrent(cur.t)
	}
	return copied
}

func (s *valueSeries) operationConst(op func(a float64) float64) ValueSeries {
	copied := NewValueSeries()
	f := s.GetFirst()
	for {
		if f == nil {
			break
		}
		copied.Set(f.t, op(f.v))
		f = f.next
	}
	cur := s.GetCurrent()
	if cur != nil {
		copied.SetCurrent(cur.t)
	}
	return copied
}

func (s *valueSeries) Add(v ValueSeries) ValueSeries {
	return s.operation(v, func(a, b float64) float64 {
		return a + b
	})
}

func (s *valueSeries) AddConst(c float64) ValueSeries {
	return s.operationConst(func(a float64) float64 {
		return a + c
	})
}

func (s *valueSeries) Div(v ValueSeries) ValueSeries {
	return s.operation(v, func(a, b float64) float64 {
		return a / b
	})
}

func (s *valueSeries) DivConst(v float64) ValueSeries {
	return s.operationConst(func(a float64) float64 {
		return a / v
	})
}

func (s *valueSeries) Len() int {
	return len(s.timemap)
}

func (s *valueSeries) Mul(v ValueSeries) ValueSeries {
	return s.operation(v, func(a, b float64) float64 {
		return a * b
	})
}

func (s *valueSeries) MulConst(v float64) ValueSeries {
	return s.operationConst(func(a float64) float64 {
		return a * v
	})
}

func (s *valueSeries) SetMax(m int64) {
	s.max = m
	s.resize()
}

func (s *valueSeries) Sub(v ValueSeries) ValueSeries {
	return s.operation(v, func(a, b float64) float64 {
		return a - b
	})
}

func (s *valueSeries) SubConst(v float64) ValueSeries {
	return s.operationConst(func(a float64) float64 {
		return a - v
	})
}

func (s *valueSeries) ID() string {
	return s.id
}

func (s *valueSeries) SetCurrent(t time.Time) bool {
	v, ok := s.timemap[t.Unix()]
	if !ok {
		s.cur = nil
		return false
	}
	s.cur = v
	return true
}

func (s *valueSeries) GetCurrent() *Value {
	return s.cur
}

func (s *valueSeries) GetFirst() *Value {
	return s.first
}

func (s *valueSeries) GetLast() *Value {
	return s.last
}

func (s *valueSeries) Val() *float64 {
	if s.cur == nil {
		return nil
	}
	return &s.cur.v
}

func (s *valueSeries) Get(t time.Time) *Value {
	return s.getValue(t.Unix())
}

func (s *valueSeries) getValue(t int64) *Value {
	return s.timemap[t]
}

func (s *valueSeries) setValue(t int64, v *Value) {
	s.timemap[t] = v
	s.resize()
}

// Push will append at the end of the list. Replaces value if exists
func (s *valueSeries) SetAll(val float64) {
	f := s.GetFirst()
	for {
		if f == nil {
			break
		}
		s.Set(f.t, val)
		f = f.next
	}
}

// Set appends to the end of the series. If same timestamp exists, its value will be replaced
func (s *valueSeries) Set(t time.Time, val float64) {
	curval := s.getValue(t.Unix())
	if curval != nil {
		// replace existing
		v2 := &Value{
			next: curval.next,
			prev: curval.prev,
			t:    t,
			v:    val,
		}
		if curval.prev != nil {
			curval.prev.next = v2
		}
		if curval.next != nil {
			curval.next.prev = v2
		}
		if s.cur == curval {
			s.cur = v2
		}
		if s.first == curval {
			s.first = v2
		}
		if s.last == curval {
			s.last = v2
		}
		s.setValue(t.Unix(), v2)
		return
	}

	v := &Value{
		t: t,
		v: val,
	}
	if s.last != nil {
		s.last.next = v
		v.prev = s.last
	}
	s.last = v
	if s.first == nil {
		s.first = v
	}
	s.setValue(t.Unix(), v)
}

func (s *valueSeries) resize() {
	// set to unlimited, nothing to perform
	if s.max == 0 {
		return
	}
	for {
		if int64(s.Len()) <= s.max {
			break
		}
		s.Shift()
	}
}

func (s *valueSeries) Shift() bool {
	if s.first == nil {
		return false
	}
	delete(s.timemap, s.first.t.Unix())
	s.first = s.first.next
	return true
}
