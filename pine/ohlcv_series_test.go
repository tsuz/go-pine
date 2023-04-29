package pine

import (
	"testing"
	"time"
)

func TestNewOHLCVSeriesPush(t *testing.T) {
	start := time.Now()
	data := OHLCVTestData(start, 3, 5*60*1000)
	empty := make([]OHLCV, 0)
	s, err := NewOHLCVSeries(empty)
	if err != nil {
		t.Fatal(err)
	}

	for i, v := range data {
		s.Push(v)

		if s.Len() != i+1 {
			t.Errorf("expected len of %d but got %d", i+1, s.Len())
		}
	}

	for i := 0; i < 3; i++ {
		s.Next()
		close := OHLCVAttr(s, OHLCPropClose)
		if *close.Val() != data[i].C {
			t.Errorf("expected %+v but got %+v", data[i].C, *close.Val())
		}
	}
}

func TestNewOHLCVSeriesShift(t *testing.T) {
	start := time.Now()
	data := OHLCVTestData(start, 3, 5*60*1000)

	s, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		s.Shift()
		if s.Len() != 3-(i+1) {
			t.Errorf("expected len of %d but got %d", 3-(i+1), s.Len())
		}
	}
}

func TestNewOHLCVSeriesMaxResize(t *testing.T) {
	start := time.Now()
	data := OHLCVTestData(start, 6, 5*60*1000)

	s, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}
	s.SetMax(3)

	for i := 0; i < 3; i++ {
		v, _ := s.Next()
		if v.C != data[i+3].C {
			t.Errorf("expected %+v but got %+v", v.C, data[i+3].C)
		}
	}
}

func TestNewOHLCVSeriesMaxCheckUponPush(t *testing.T) {
	start := time.Now()
	data := OHLCVTestData(start, 3, 5*60*1000)
	newv := OHLCVTestData(start.Add(3*5*time.Minute), 1, 5*60*1000)

	s, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}
	s.SetMax(3)

	s.Push(newv[0])

	for i := 0; i < 3; i++ {
		v, _ := s.Next()
		if i < 2 {
			if v.C != data[i+1].C {
				t.Errorf("expected %+v but got %+v for %d", data[i+1].C, v.C, i)
			}
		} else {
			if v.C != newv[0].C {
				t.Errorf("expected %+v but got %+v for %d", v.C, newv[0].C, i)
			}
		}
	}
}
