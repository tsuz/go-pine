package pine

import (
	"fmt"
	"log"
	"testing"
	"time"
)

// TestSeriesPowNoData tests no data scenario
//
// t=time.Time (no iteration) | |
// p=ValueSeries              | |
// stdev=ValueSeries            | |
func TestSeriesPowNoData(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	prop := OHLCVAttr(series, OHLCPropClose)
	stdev := Pow(prop, 2.0)
	if stdev == nil {
		t.Error("Expected to be non nil but got nil")
	}
}

// TestSeriesPowNoIteration tests this sceneario where there's no iteration yet
//
// t=time.Time (no iteration) | 1  |  2   | 3  | 4  |
// p=ValueSeries              | 14 |  15  | 17 | 18 |
// pow=ValueSeries            |    |      |    |    |
func TestSeriesPowNoIteration(t *testing.T) {

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
	pow := Pow(prop, 2)
	if pow == nil {
		t.Error("Expected to be non-nil but got nil")
	}
}

// TestSeriesPowIteration tests this scneario
//
// t=time.Time       | 1     |  2    | 3     |
// p=ValueSeries     | 13    | 15    | 11    |
// pow(0.5)     	 | 3.606 | 3.873 | 3.317 |
// pow(2)       	 | 169   | 225   | 121   |
func TestSeriesPowIteration(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 3, 5*60*1000)
	data[0].C = 13
	data[1].C = 15
	data[2].C = 11

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	testTable := []struct {
		exp  float64
		vals []float64
	}{
		{
			exp:  0.5,
			vals: []float64{3.606, 3.873, 3.317},
		},
		{
			exp:  2,
			vals: []float64{169, 225, 121},
		},
	}

	for j := 0; j <= 2; j++ {
		series.Next()

		for i, v := range testTable {
			prop := OHLCVAttr(series, OHLCPropClose)
			pow := Pow(prop, v.exp)
			exp := v.vals[j]
			if exp == 0 {
				if pow.Val() != nil {
					t.Fatalf("expected nil but got non nil: %+v at vals item: %d, testtable item: %d", *pow.Val(), j, i)
				}
				// OK
			}
			if exp != 0 {
				if pow.Val() == nil {
					t.Fatalf("expected non nil: %+v but got nil at vals item: %d, testtable item: %d", exp, j, i)
				}
				if fmt.Sprintf("%.03f", exp) != fmt.Sprintf("%.03f", *pow.Val()) {
					t.Fatalf("expected %+v but got %+v at vals item: %d, testtable item: %d", exp, *pow.Val(), j, i)
				}
				// OK
			}
		}
	}
}

func TestMemoryLeakPow(t *testing.T) {
	testMemoryLeak(t, func(o OHLCVSeries) error {
		prop := OHLCVAttr(o, OHLCPropClose)
		Pow(prop, 2)
		return nil
	})
}

func ExamplePow() {
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	for {
		if v, _ := series.Next(); v == nil {
			break
		}

		close := OHLCVAttr(series, OHLCPropClose)
		added := AddConst(close, 3.0)
		pow := Pow(added, 2)
		log.Printf("Pow: %+v", pow.Val())
	}
}
