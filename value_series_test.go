package pine

import (
	"testing"
	"time"
)

func TestValueSeriesAddConst(t *testing.T) {
	a := NewValueSeries()
	now := time.Now()
	a.Set(now, 1)
	a.Set(now.Add(time.Duration(1000*1e6)), 2)

	b := a.AddConst(3)
	f := b.GetFirst()
	if f == nil {
		t.Fatalf("expected to be non nil but got nil")
	}
	if f.v != 4 {
		t.Errorf("expected %+v but got %+v", 4, f.v)
	}
	if f.next.v != 5 {
		t.Errorf("expected %+v but got %+v", 5, f.v)
	}
}

func TestValueSeriesDiv(t *testing.T) {
	a := NewValueSeries()
	now := time.Now()
	a.Set(now, 1)
	a.Set(now.Add(time.Duration(1000*1e6)), 2)

	b := NewValueSeries()
	b.Set(now, 4)
	b.Set(now.Add(time.Duration(1000*1e6)), 4)

	c := a.Div(b)
	c.SetCurrent(now)
	f := c.GetCurrent()
	if f == nil {
		t.Fatalf("expected to be non nil but got nil")
	}
	if f.v != 0.25 {
		t.Errorf("expected %+v but got %+v", 0.25, f.v)
	}
	if f.next.v != 0.5 {
		t.Errorf("expected %+v but got %+v", 0.5, f.v)
	}
}

func TestValueSeriesGetFirst(t *testing.T) {

	s := NewValueSeries()
	now := time.Now()
	s.Set(now, 1)
	s.Set(now.Add(time.Duration(1000*1e6)), 2)
	s.SetCurrent(now)
	f := s.GetFirst()
	if f == nil {
		t.Errorf("expected to be non nil but got nil")
	}
	if f.next == nil {
		t.Errorf("expected next to be non nil but got nil")
	}
	if f.next.v != 2 {
		t.Errorf("expected next value to be 2 but got  %+v", f.next.v)
	}
}
