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
	// DivConst(float64) ValueSeries
	// Mul(ValueSeries) ValueSeries
	// MulConst(float64) ValueSeries
	Sub(ValueSeries) ValueSeries
	// SubConst(float64) ValueSeries

	// Get gets the item by time in value series
	Get(time.Time) *Value
	// GetLast gets the last item in value series
	GetLast() *Value
	// GetFirst gets the first item in value series
	GetFirst() *Value

	Copy() ValueSeries

	// replace existing
	Set(time.Time, float64)
	SetAll(val float64)

	Val() *float64
	SetCurrent(time.Time) bool
	GetCurrent() *Value
}

type valueSeries struct {
	id  string
	cur *Value
	ord []int64
	sync.Mutex
	timemap map[int64]*Value
}

type Value struct {
	t    time.Time
	v    float64
	prev *Value
	next *Value
}

func NewValueSeries() ValueSeries {
	u := uuid.NewV4()
	v := &valueSeries{
		id:      u.String(),
		ord:     make([]int64, 0),
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

func (s *valueSeries) Add(v ValueSeries) ValueSeries {
	copied := s.Copy()
	f := s.GetFirst()
	for {
		if f == nil {
			break
		}
		newv := v.Get(f.t)
		if newv != nil {
			copied.Set(f.t, f.v+newv.v)
		}
		f = f.next
	}
	return copied
}

func (s *valueSeries) AddConst(v float64) ValueSeries {
	copied := s.Copy()
	f := s.GetFirst()
	for {
		if f == nil {
			break
		}
		copied.Set(f.t, f.v+v)
		f = f.next
	}
	return copied
}

func (s *valueSeries) Div(v ValueSeries) ValueSeries {
	copied := s.Copy()
	f := s.GetFirst()
	for {
		if f == nil {
			break
		}
		newv := v.Get(f.t)
		if newv != nil {
			copied.Set(f.t, f.v/newv.v)
		}
		f = f.next
	}
	return copied
}

func (s *valueSeries) Sub(v ValueSeries) ValueSeries {
	copied := s.Copy()
	f := s.GetFirst()
	for {
		if f == nil {
			break
		}
		newv := v.Get(f.t)
		if newv != nil {
			copied.Set(f.t, f.v-newv.v)
		}
		f = f.next
	}
	return copied
}

func (s *valueSeries) ID() string {
	return s.id
}

func (s *valueSeries) SetCurrent(t time.Time) bool {
	v, ok := s.timemap[t.Unix()]
	if !ok {
		return false
	}
	s.cur = v
	return true
}

func (s *valueSeries) GetCurrent() *Value {
	return s.cur
}

func (s *valueSeries) GetFirst() *Value {
	if len(s.ord) == 0 {
		return nil
	}

	val := s.getValue(s.ord[0])

	return val
}

func (s *valueSeries) GetLast() *Value {
	if len(s.ord) == 0 {
		return nil
	}
	return s.getValue(s.ord[len(s.ord)-1])
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

func (s *valueSeries) appendValue(v *Value) {
	s.ord = append(s.ord, v.t.Unix())
	s.setValue(v.t.Unix(), v)
}

func (s *valueSeries) setValue(t int64, v *Value) {
	s.timemap[t] = v
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

// Push will append at the end of the list. Replaces value if exists
func (s *valueSeries) Set(t time.Time, val float64) {
	curval := s.getValue(t.Unix())
	if curval != nil {
		// just replace the map
		v2 := &Value{
			next: s.timemap[t.Unix()].next,
			prev: s.timemap[t.Unix()].prev,
			t:    t,
			v:    val,
		}
		if curval.prev != nil {
			curval.prev.next = v2
		}
		if curval.next != nil {
			curval.next.prev = v2
		}
		if s.cur != nil && s.cur.t.Equal(t) {
			s.cur = v2
		}
		s.setValue(t.Unix(), v2)
		return
	}

	v := &Value{
		t: t,
		v: val,
	}

	// no existing values so no previous or next pointers
	if len(s.ord) == 0 {
		s.appendValue(v)
		return
	}

	prevt := s.ord[len(s.ord)-1]
	prev := s.timemap[prevt]
	v.prev = prev

	s.appendValue(v)
	curt := s.getValue(v.t.Unix())

	s.timemap[prevt].next = curt
}
