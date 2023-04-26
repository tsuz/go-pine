package pine

import (
	"fmt"
	"testing"
	"time"
)

func TestNewOHLCVGetSeries(t *testing.T) {
	start := time.Now()
	data := OHLCVTestData(start, 3, 5*60*1000)

	s, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	testTable := []struct {
		prop []OHLCProp
		vals []float64
	}{
		{
			prop: []OHLCProp{OHLCPropClose, OHLCPropHigh, OHLCPropLow, OHLCPropOpen, OHLCPropHLC3},
			vals: []float64{data[0].C, data[0].H, data[0].L, data[0].O, (data[0].H + data[0].L + data[0].C) / 3},
		},
		{
			prop: []OHLCProp{OHLCPropClose, OHLCPropHigh, OHLCPropLow, OHLCPropOpen, OHLCPropHLC3},
			vals: []float64{data[1].C, data[1].H, data[1].L, data[1].O, (data[1].H + data[1].L + data[1].C) / 3},
		},
		{
			prop: []OHLCProp{OHLCPropClose, OHLCPropHigh, OHLCPropLow, OHLCPropOpen, OHLCPropHLC3},
			vals: []float64{data[2].C, data[2].H, data[2].L, data[2].O, (data[2].H + data[2].L + data[2].C) / 3},
		},
	}

	for i, v := range testTable {
		// move to next iteration
		s.Next()

		for j, p := range v.prop {
			vals := s.GetSeries(p)
			val := vals.Val()
			if *val != v.vals[j] {
				t.Errorf("Expected %+v to bs %+v but got %+v for i: %d, j: %d", p, v.vals[j], val, i, j)
			}
		}
	}

	// if this is last, return nil
	if v, _ := s.Next(); v != nil {
		t.Errorf("Expected to be nil but got %+v", v)
	}
}

func TestNewOHLCVGetTrueRange(t *testing.T) {
	data := OHLCVStaticTestData()

	s, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	expVals := []*float64{
		nil,
		NewFloat64(6.8),
		NewFloat64(8.5),
		NewFloat64(7.9),
		NewFloat64(8.3),
		NewFloat64(6.3),
		NewFloat64(6.6),
		NewFloat64(9.6),
		NewFloat64(8.0),
		NewFloat64(7.6),
	}

	for i, v := range expVals {
		// move to next iteration
		s.Next()
		vals := s.GetSeries(OHLCPropTR)

		if (vals.Val() == nil) != (v == nil) {
			t.Errorf("Expected %+v but got %+v for i: %d", v, vals.Val(), i)
			continue
		}
		if v == nil {
			continue
		}
		if fmt.Sprintf("%0.2f", *vals.Val()) != fmt.Sprintf("%0.2f", *v) {
			t.Errorf("expected %+v but got %+v", *v, *vals.Val())
		}
	}

	// if this is last, return nil
	if v, _ := s.Next(); v != nil {
		t.Errorf("Expected to be nil but got %+v", v)
	}
}

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
		close := s.GetSeries(OHLCPropClose)
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
