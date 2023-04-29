package pine

import (
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"
)

// TestSeriesChangeNoData tests no data scenario
//
// t=time.Time (no iteration) | |
// p=ValueSeries              | |
// change=ValueSeries      | |
func TestSeriesChangeNoData(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 0, 5*60*1000)

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	src := OHLCVAttr(series, OHLCPropClose)

	rsi, err := Change(src, 2)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error Change"))
	}
	if rsi == nil {
		t.Error("Expected to be non nil but got nil")
	}
}

// TestSeriesChangeNoIteration tests this sceneario where there's no iteration yet
//
// t=time.Time (no iteration)  | 1   |  2  | 3   | 4
// src=ValueSeries             | 11  | 14  | 12  | 13
// change(src, 1)			   | nil |  3  | -2  | 1
// change(src, 2)			   | nil | nil | 1   | -1
// change(src, 3)			   | nil | nil | nil | 2
func TestSeriesChangeNoIteration(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)
	data[0].C = 11
	data[1].C = 14
	data[2].C = 12
	data[3].C = 13

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	src := OHLCVAttr(series, OHLCPropClose)
	rsi, err := Change(src, 1)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error Change"))
	}
	if rsi == nil {
		t.Error("Expected to be non-nil but got nil")
	}
}

// TestSeriesChangeSuccess tests this scneario when the iterator is at t=4 is not at the end
//
// t=time.Time      | 1   |  2  | 3   | 4
// src=ValueSeries  | 11  | 14  | 12  | 13
// change(src, 1)	| nil |  3  | -2  | 1
// change(src, 2)	| nil | nil | 1   | -1
// change(src, 3)	| nil | nil | nil | 2
func TestSeriesChangeSuccess(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)
	data[0].C = 11
	data[1].C = 14
	data[2].C = 12
	data[3].C = 13

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	testTable := []struct {
		lookback int
		vals     []float64
	}{
		{
			lookback: 1,
			vals:     []float64{0, 3, -2, 1},
		},
		{
			lookback: 2,
			vals:     []float64{0, 0, 1, -1},
		},
		{
			lookback: 3,
			vals:     []float64{0, 0, 0, 2},
		},
	}

	for j := 0; j <= 3; j++ {
		series.Next()

		for i, v := range testTable {
			src := OHLCVAttr(series, OHLCPropClose)
			vw, err := Change(src, v.lookback)
			if err != nil {
				t.Fatal(errors.Wrap(err, "error Change"))
			}
			exp := v.vals[j]
			if exp == 0 {
				if vw.Val() != nil {
					t.Fatalf("expected nil but got non nil: %+v at vals item: %d, testtable item: %d", *vw.Val(), j, i)
				}
				// OK
			}
			if exp != 0 {
				if vw.Val() == nil {
					t.Fatalf("expected non nil: %+v but got nil at vals item: %d, testtable item: %d", exp, j, i)
				}
				if exp != *vw.Val() {
					t.Fatalf("expected %+v but got %+v at vals item: %d, testtable item: %d", exp, *vw.Val(), j, i)
				}
				// OK
			}
		}
	}
}

// TestSeriesChangeNotEnoughData tests this scneario when the lookback is more than the number of data available
//
// t=time.Time      | 1   |  2  | 3   | 4
// src=ValueSeries  | 11  | 14  | 12  | 13
// change(src, 1)	| nil |  3  | -2  | 1
// change(src, 2)	| nil | nil | 1   | -1
// change(src, 3)	| nil | nil | nil | 2
func TestSeriesChangeNotEnoughData(t *testing.T) {

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

	src := OHLCVAttr(series, OHLCPropClose)

	vw, err := Change(src, 4)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error Change"))
	}
	if vw.Val() != nil {
		t.Errorf("Expected nil but got %+v", *vw.Val())
	}
}

func BenchmarkChange(b *testing.B) {
	// run the Fib function b.N times
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	vals := OHLCVAttr(series, OHLCPropClose)

	for n := 0; n < b.N; n++ {
		series.Next()
		Change(vals, 5)
	}
}

func ExampleChange() {
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	for {
		if v, _ := series.Next(); v == nil {
			break
		}

		close := OHLCVAttr(series, OHLCPropClose)
		chg, err := Change(close, 12)
		if err != nil {
			log.Fatal(errors.Wrap(err, "error change"))
		}
		log.Printf("Change line: %+v", chg.Val())
	}
}
