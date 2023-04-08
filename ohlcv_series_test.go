package pine

import (
	"testing"
	"time"
)

func TestNewOHLCVSeries(t *testing.T) {
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
			prop: []OHLCProp{OHLCPropClose, OHLCPropHigh, OHLCPropLow, OHLCPropOpen},
			vals: []float64{data[0].C, data[0].H, data[0].L, data[0].O},
		},
		{
			prop: []OHLCProp{OHLCPropClose, OHLCPropHigh, OHLCPropLow, OHLCPropOpen},
			vals: []float64{data[1].C, data[1].H, data[1].L, data[1].O},
		},
		{
			prop: []OHLCProp{OHLCPropClose, OHLCPropHigh, OHLCPropLow, OHLCPropOpen},
			vals: []float64{data[2].C, data[2].H, data[2].L, data[2].O},
		},
	}

	for _, v := range testTable {
		// move to next iteration
		s.Next()

		for j, p := range v.prop {
			vals := s.GetSeries(p)
			val := vals.Val()
			if *val != v.vals[j] {
				t.Errorf("Expected %+v to bs %+v but got %+v", p, v.vals[j], val)
			}
		}
	}

	// if this is last, return nil
	if v := s.Next(); v != nil {
		t.Errorf("Expected to be nil but got %+v", v)
	}
}
