package pine

import (
	"fmt"
	"log"
	"testing"
	"time"
)

// TestSeriesStdevNoData tests no data scenario
//
// t=time.Time (no iteration) | |
// p=ValueSeries              | |
// stdev=ValueSeries            | |
func TestSeriesStdevNoData(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	prop := OHLCVAttr(series, OHLCPropClose)
	stdev := Stdev(prop, 2)
	if stdev == nil {
		t.Error("Expected to be non nil but got nil")
	}
}

// TestSeriesStdevNoIteration tests this sceneario where there's no iteration yet
//
// t=time.Time (no iteration) | 1  |  2   | 3  | 4  |
// p=ValueSeries              | 14 |  15  | 17 | 18 |
// stdev=ValueSeries            |    |      |    |    |
func TestSeriesStdevNoIteration(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)
	data[0].C = 14
	data[1].C = 15
	data[2].C = 17
	data[3].C = 18

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	prop := OHLCVAttr(series, OHLCPropClose)
	stdev := RSI(prop, 2)
	if stdev == nil {
		t.Error("Expected to be non-nil but got nil")
	}
}

// TestSeriesStdevIteration tests this scneario
//
// t=time.Time        | 1   |  2  | 3   | 4    | 5      |
// p=ValueSeries      | 13  | 15  | 11  | 19   | 21     |
// sma(p, 3)	      | nil | nil | 13  | 15   | 17     |
// p - sma(p, 3)(t=1) | nil | nil | nil | nil  | nil    |
// p - sma(p, 3)(t=2) | nil | nil | nil | nil  | nil    |
// p - sma(p, 3)(t=3) | 0   |  2  | -2  | 6    | 5      |
// p - sma(p, 3)(t=4) | -2  |  0  | -4  | 4    | 6      |
// p - sma(p, 3)(t=5) | -4  | -2  | -6  | 2    | 4      |
// Stdev(p, 3)		  | nil | nil |  2  | 4    | 5.2915 |
func TestSeriesStdevIteration(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 5, 5*60*1000)
	data[0].C = 13
	data[1].C = 15
	data[2].C = 11
	data[3].C = 19
	data[4].C = 21

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	testTable := []float64{0, 0, 2, 4, 5.2915}

	for i, v := range testTable {
		series.Next()

		prop := OHLCVAttr(series, OHLCPropClose)
		stdev := Stdev(prop, 3)
		exp := v
		if exp == 0 {
			if stdev.Val() != nil {
				t.Fatalf("expected nil but got non nil: %+v  testtable item: %d", *stdev.Val(), i)
			}
			// OK
		}
		if exp != 0 {
			if stdev.Val() == nil {
				t.Fatalf("expected non nil: %+v but got nil  testtable item: %d", exp, i)
			}
			if fmt.Sprintf("%.04f", exp) != fmt.Sprintf("%.04f", *stdev.Val()) {
				t.Fatalf("expected %+v but got %+v  testtable item: %d", exp, *stdev.Val(), i)
			}
			// OK
		}
	}
}

// TestSeriesStdevNotEnoughData tests when the lookback is more than the number of data available
//
// t=time.Time     | 1  |  2   | 3  | 4 (here)  |
// p=ValueSeries   | 14 |  15  | 17 | 18        |
// stdev(close, 5) | nil| nil  | nil| nil       |
func TestSeriesStdevNotEnoughData(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)
	data[0].C = 13
	data[1].C = 15
	data[2].C = 11
	data[3].C = 18

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	series.Next()
	series.Next()
	series.Next()
	series.Next()

	testTable := []struct {
		lookback int
		exp      *float64
	}{
		{
			lookback: 5,
			exp:      nil,
		},
		{
			lookback: 6,
			exp:      nil,
		},
	}

	for i, v := range testTable {
		prop := OHLCVAttr(series, OHLCPropClose)

		stdev := Stdev(prop, int64(v.lookback))
		if stdev == nil {
			t.Errorf("Expected to be non nil but got nil at idx: %d", i)
		}
		if stdev.Val() != v.exp {
			t.Errorf("Expected to get %+v but got %+v for lookback %+v", v.exp, *stdev.Val(), v.lookback)
		}
	}
}

func TestMemoryLeakStdev(t *testing.T) {
	testMemoryLeak(t, func(o OHLCVSeries) error {
		prop := OHLCVAttr(o, OHLCPropClose)
		Stdev(prop, 12)
		return nil
	})
}

func ExampleStdev() {
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	for {
		if v, _ := series.Next(); v == nil {
			break
		}

		close := OHLCVAttr(series, OHLCPropClose)
		stdev := Stdev(close, 12)
		log.Printf("Stdev: %+v", stdev.Val())
	}
}
