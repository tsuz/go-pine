package pine

import (
	"testing"
	"time"
)

func TestValueSeriesGetFirst(t *testing.T) {

	s := NewValueSeries()
	now := time.Now()
	s.Push(now, 1)
	s.Push(now.Add(time.Duration(1000*1e6)), 2)
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
