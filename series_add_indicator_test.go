package pine

import (
	"testing"
	"time"
)

func TestSeriesAddIndicator(t *testing.T) {
	opts := SeriesOpts{
		Interval: 300,
		Max:      100,
	}
	now := time.Now()
	fivemin := now.Add(5 * time.Minute)
	data := []OHLCV{
		OHLCV{
			O: 14,
			H: 15,
			L: 13,
			C: 14,
			V: 131,
			S: now,
		},
		OHLCV{
			O: 13,
			H: 18,
			L: 9,
			C: 15,
			V: 12,
			S: fivemin,
		},
	}
	s, err := NewSeries(data, opts)
	if err != nil {
		t.Fatal(err)
	}
	hl2 := NewOHLCProp(OHLCPropHL2)
	if err := s.AddIndicator("hl2", hl2); err != nil {
		t.Fatal("error adding indicator")
	}

	v := s.GetValueForInterval(now)
	if v == nil || v.Indicators["hl2"] == nil {
		t.Fatalf("expected ohlcv to be non nil but got %+v", v)
	} else if *(v.Indicators["hl2"]) != 14 {
		t.Fatalf("expected HL2 (midpoint) to be 14 but got %+v", v.Value)
	}

	v = s.GetValueForInterval(fivemin)
	if v == nil || v.Indicators["hl2"] == nil {
		t.Fatalf("expected ohlcv to be non nil but got %+v", v)
	} else if *(v.Indicators["hl2"]) != 13.5 {
		t.Fatalf("expected HL2 (midpoint) to be 13.5 but got %+v", v.Value)
	}

}
