package pine

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"
)

// TestSeriesVarianceNoData tests no data scenario
//
// t=time.Time (no iteration) | |
// p=ValueSeries              | |
// variance=ValueSeries            | |
func TestSeriesVarianceNoData(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	prop := OHLCVAttr(series, OHLCPropClose)
	variance, err := Variance(prop, 2)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error Variance"))
	}
	if variance == nil {
		t.Error("Expected to be non nil but got nil")
	}
}

// TestSeriesVarianceNoIteration tests this sceneario where there's no iteration yet
//
// t=time.Time (no iteration) | 1  |  2   | 3  | 4  |
// p=ValueSeries              | 14 |  15  | 17 | 18 |
// variance=ValueSeries            |    |      |    |    |
func TestSeriesVarianceNoIteration(t *testing.T) {

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
	variance, err := RSI(prop, 2)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error RSI"))
	}
	if variance == nil {
		t.Error("Expected to be non-nil but got nil")
	}
}

// TestSeriesVarianceIteration tests this scneario
//
// t=time.Time        | 1   |  2  | 3   | 4    | 5   |
// p=ValueSeries      | 13  | 15  | 11  | 19   | 21  |
// sma(p, 3)	      | nil | nil | 13  | 15   | 17  |
// p - sma(p, 3)(t=1) | nil | nil | nil | nil  | nil |
// p - sma(p, 3)(t=2) | nil | nil | nil | nil  | nil |
// p - sma(p, 3)(t=3) | 0   |  2  | -2  | 6    | 5   |
// p - sma(p, 3)(t=4) | -2  |  0  | -4  | 4    | 6   |
// p - sma(p, 3)(t=5) | -4  | -2  | -6  | 2    | 4   |
// Variance(p, 3)     | nil | nil |  4  | 16   | 28  |
func TestSeriesVarianceIteration(t *testing.T) {

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

	testTable := []float64{0, 0, 4, 16, 28}

	for i, v := range testTable {
		series.Next()

		prop := OHLCVAttr(series, OHLCPropClose)
		variance, err := Variance(prop, 3)
		if err != nil {
			t.Fatal(errors.Wrap(err, "error Variance"))
		}
		exp := v
		if exp == 0 {
			if variance.Val() != nil {
				t.Fatalf("expected nil but got non nil: %+v  testtable item: %d", *variance.Val(), i)
			}
			// OK
		}
		if exp != 0 {
			if variance.Val() == nil {
				t.Fatalf("expected non nil: %+v but got nil  testtable item: %d", exp, i)
			}
			if fmt.Sprintf("%.04f", exp) != fmt.Sprintf("%.04f", *variance.Val()) {
				t.Fatalf("expected %+v but got %+v  testtable item: %d", exp, *variance.Val(), i)
			}
			// OK
		}
	}
}

// TestSeriesVarianceNotEnoughData tests when the lookback is more than the number of data available
//
// t=time.Time         | 1  |  2   | 3  | 4 (here)  |
// p=ValueSeries       | 14 |  15  | 17 | 18        |
// variance(close, 5) | nil| nil  | nil| nil       |
func TestSeriesVarianceNotEnoughData(t *testing.T) {

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

		variance, err := Variance(prop, int64(v.lookback))
		if err != nil {
			t.Fatal(errors.Wrap(err, "error RSI"))
		}
		if variance == nil {
			t.Errorf("Expected to be non nil but got nil at idx: %d", i)
		}
		if variance.Val() != v.exp {
			t.Errorf("Expected to get %+v but got %+v for lookback %+v", v.exp, *variance.Val(), v.lookback)
		}
	}
}

// func TestMemoryLeakVariance(t *testing.T) {
// 	testMemoryLeak(t, func(o OHLCVSeries) error {
// 		prop := OHLCVAttr(o, OHLCPropClose)
// 		_, err := Variance(prop, 10)
// 		return err
// 	})
// }

func BenchmarkVariance(b *testing.B) {
	// run the Fib function b.N times
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	vals := OHLCVAttr(series, OHLCPropClose)

	for n := 0; n < b.N; n++ {
		series.Next()
		Variance(vals, 5)
	}
}

func ExampleVariance() {
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	for {
		if v, _ := series.Next(); v == nil {
			break
		}

		close := OHLCVAttr(series, OHLCPropClose)
		variance, err := Variance(close, 20)
		if err != nil {
			log.Fatal(errors.Wrap(err, "error geting variance"))
		}
		log.Printf("Variance: %+v", variance.Val())
	}
}
