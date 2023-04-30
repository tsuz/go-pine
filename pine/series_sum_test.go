package pine

import (
	"log"
	"testing"
	"time"
)

// TestSeriesSumNoData tests no data scenario
//
// t=time.Time (no iteration) | |
// p=ValueSeries              | |
// stdev=ValueSeries            | |
func TestSeriesSumNoData(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	prop := OHLCVAttr(series, OHLCPropClose)
	stdev := Sum(prop, 2)
	if stdev == nil {
		t.Error("Expected to be non nil but got nil")
	}
}

// TestSeriesSumNoIteration tests this sceneario where there's no iteration yet
//
// t=time.Time (no iteration) | 1  |  2   | 3  | 4  |
// p=ValueSeries              | 14 |  15  | 17 | 18 |
// sum=ValueSeries            |    |      |    |    |
func TestSeriesSumNoIteration(t *testing.T) {

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
	sum := Sum(prop, 2)
	if sum == nil {
		t.Error("Expected to be non-nil but got nil")
	}
}

// TestSeriesSumIteration tests this scneario
//
// t=time.Time       | 1   |  2  | 3    | 4    | 5  |
// p=ValueSeries     | 13  | 15  | 11   | 19   | 21 |
// sum(p, 1)	     | 13  | 15  | 11   | 19   | 21 |
// sum(p, 2)	     | nil | 28  | 26   | 30   | 40 |
// sum(p, 3)	     | nil | nil | 39   | 45   | 51 |
func TestSeriesSumIteration(t *testing.T) {

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

	testTable := []struct {
		lookback int
		vals     []float64
	}{
		{
			lookback: 1,
			vals:     []float64{13, 15, 11, 19, 21},
		},
		{
			lookback: 2,
			vals:     []float64{0, 28, 26, 30, 40},
		},
		{
			lookback: 3,
			vals:     []float64{0, 0, 39, 45, 51},
		},
	}

	for j := 0; j <= 3; j++ {
		series.Next()

		for i, v := range testTable {
			prop := OHLCVAttr(series, OHLCPropClose)
			sum := Sum(prop, v.lookback)
			exp := v.vals[j]
			if exp == 0 {
				if sum.Val() != nil {
					t.Fatalf("expected nil but got non nil: %+v at vals item: %d, testtable item: %d", *sum.Val(), j, i)
				}
				// OK
			}
			if exp != 0 {
				if sum.Val() == nil {
					t.Fatalf("expected non nil: %+v but got nil at vals item: %d, testtable item: %d", exp, j, i)
				}
				if exp != *sum.Val() {
					t.Fatalf("expected %+v but got %+v at vals item: %d, testtable item: %d", exp, *sum.Val(), j, i)
				}
				// OK
			}
		}
	}
}

func TestMemoryLeakSum(t *testing.T) {
	testMemoryLeak(t, func(o OHLCVSeries) error {
		prop := OHLCVAttr(o, OHLCPropClose)
		Sum(prop, 10)
		return nil
	})
}

func ExampleSum() {
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	for {
		if v, _ := series.Next(); v == nil {
			break
		}

		close := OHLCVAttr(series, OHLCPropClose)
		// Get the sum of last 10 values
		sum := Sum(close, 10)
		log.Printf("Sum: %+v", sum.Val())
	}
}
