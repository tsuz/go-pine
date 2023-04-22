package pine

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"
)

// TestSeriesROCNoData tests no data scenario
//
// t=time.Time (no iteration) | |
// p=ValueSeries              | |
// change=ValueSeries      | |
func TestSeriesROCNoData(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 0, 5*60*1000)

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	src := series.GetSeries(OHLCPropClose)

	rsi, err := ROC(src, 2)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error ROC"))
	}
	if rsi == nil {
		t.Error("Expected to be non nil but got nil")
	}
}

// TestSeriesROCNoIteration tests this sceneario where there's no iteration yet
//
// t=time.Time (no iteration)  | 1   |  2      | 3   	  | 4
// src=ValueSeries             | 11  | 14      | 12       | 13
// roc(src, 1)	               | nil | 27.2727 | -14.286  | 8.3333
// roc(src, 2)	               | nil | nil     | 9.090909 | 7.1429
// roc(src, 3)	               | nil | nil     | nil      | 18.1818
func TestSeriesROCNoIteration(t *testing.T) {

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

	src := series.GetSeries(OHLCPropClose)
	rsi, err := ROC(src, 1)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error ROC"))
	}
	if rsi == nil {
		t.Error("Expected to be non-nil but got nil")
	}
}

// TestSeriesROCSuccess tests this scneario when the iterator is at t=4 is not at the end
//
// t=time.Time      | 1   |  2      | 3        | 4
// src=ValueSeries  | 11  | 14      | 12       | 13
// roc(src, 1)	    | nil | 27.2727 | -14.2857 | 8.3333
// roc(src, 2)	    | nil | nil     | 9.090909 | -7.1429
// roc(src, 3)	    | nil | nil     | nil      | 18.1818
func TestSeriesROCSuccess(t *testing.T) {

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
			vals:     []float64{0, 27.2727, -14.2857, 8.3333},
		},
		{
			lookback: 2,
			vals:     []float64{0, 0, 9.090909, -7.1429},
		},
		{
			lookback: 3,
			vals:     []float64{0, 0, 0, 18.1818},
		},
	}

	for j := 0; j <= 3; j++ {
		series.Next()

		for i, v := range testTable {
			src := series.GetSeries(OHLCPropClose)
			vw, err := ROC(src, v.lookback)
			if err != nil {
				t.Fatal(errors.Wrap(err, "error ROC"))
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
				if fmt.Sprintf("%.4f", exp) != fmt.Sprintf("%.4f", *vw.Val()) {
					t.Fatalf("expected %+v but got %+v at vals item: %d, testtable item: %d", exp, *vw.Val(), j, i)
				}
				// OK
			}
		}
	}
}

// TestSeriesROCNotEnoughData tests this scneario when the lookback is more than the number of data available
//
// t=time.Time      | 1   |  2      | 3        | 4
// src=ValueSeries  | 11  | 14      | 12       | 13
// roc(src, 1)	    | nil | 27.2727 | -14.2857 | 8.3333
// roc(src, 2)	    | nil | nil     | 9.090909 | -7.1429
// roc(src, 3)	    | nil | nil     | nil      | 18.1818
func TestSeriesROCNotEnoughData(t *testing.T) {

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

	src := series.GetSeries(OHLCPropClose)

	vw, err := ROC(src, 4)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error ROC"))
	}
	if vw.Val() != nil {
		t.Errorf("Expected nil but got %+v", *vw.Val())
	}
}

func BenchmarkROC(b *testing.B) {
	// run the Fib function b.N times
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	vals := series.GetSeries(OHLCPropClose)

	for n := 0; n < b.N; n++ {
		series.Next()
		ROC(vals, 5)
	}
}

func ExampleROC() {
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	for {
		if v, _ := series.Next(); v == nil {
			break
		}

		close := series.GetSeries(OHLCPropClose)
		roc, err := ROC(close, 4)
		if err != nil {
			log.Fatal(errors.Wrap(err, "error geting roc"))
		}
		log.Printf("ROC: %+v", roc.Val())
	}
}
