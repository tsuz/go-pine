package pine

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"
)

// TestSeriesATRNoData tests no data scenario
//
// t=time.Time (no iteration) | |
// p=ValueSeries              | |
// atr=ValueSeries            | |
func TestSeriesATRNoData(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	prop := OHLCVAttr(series, OHLCPropClose)
	atr, err := ATR(prop, 2)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error ATR"))
	}
	if atr == nil {
		t.Error("Expected to be non nil but got nil")
	}
}

// TestSeriesATRNoIteration tests this sceneario where there's no iteration yet
//
// t=time.Time (no iteration) | 1  |  2   | 3  | 4  |
// p=ValueSeries              | 14 |  15  | 17 | 18 |
// atr=ValueSeries            |    |      |    |    |
func TestSeriesATRNoIteration(t *testing.T) {

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
	atr, err := ATR(prop, 2)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error RSI"))
	}
	if atr == nil {
		t.Error("Expected to be non-nil but got nil")
	}
}

// TestSeriesATRIteration tests this scneario
//
// t=time.Time        | 1     |  2    | 3     | 4
// close=ValueSeries  | 13    | 15    | 11    | 19
// high=ValueSeries   | 13.2  | 16    | 11.2  | 19
// low=ValueSeries    | 12.6  | 14.8  | 11    | 19
// TR=ValueSeries     | 0.6   | 3     | 4     | 8
// ATR=ValueSeries    | nil   | nil   | 2.5333| 4.3556
func TestSeriesATRIteration(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)
	data[0].C = 13
	data[0].H = 13.2
	data[0].L = 12.6
	data[1].C = 15
	data[1].H = 16
	data[1].L = 14.8
	data[2].C = 11
	data[2].H = 11.2
	data[2].L = 11
	data[3].C = 19
	data[3].H = 19
	data[3].L = 19

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	testTable := []float64{0, 0, 2.5333, 4.3556}

	for i, v := range testTable {
		series.Next()

		prop := OHLCVAttr(series, OHLCPropTRHL)
		atr, err := ATR(prop, 3)
		if err != nil {
			t.Fatal(errors.Wrap(err, "error ATR"))
		}
		exp := v
		if exp == 0 {
			if atr.Val() != nil {
				t.Fatalf("expected nil but got non nil: %+v  testtable item: %d", *atr.Val(), i)
			}
			// OK
		}
		if exp != 0 {
			if atr.Val() == nil {
				t.Fatalf("expected non nil: %+v but got nil  testtable item: %d", exp, i)
			}
			if fmt.Sprintf("%.04f", exp) != fmt.Sprintf("%.04f", *atr.Val()) {
				t.Fatalf("expected %+v but got %+v  testtable item: %d", exp, *atr.Val(), i)
			}
			// OK
		}
	}
}

// TestSeriesATRNotEnoughData tests when the lookback is more than the number of data available
//
// t=time.Time     | 1  |  2   | 3  | 4 (here)  |
// p=ValueSeries   | 14 |  15  | 17 | 18        |
// atr(close, 5) | nil| nil  | nil| nil       |
func TestSeriesATRNotEnoughData(t *testing.T) {

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

		atr, err := ATR(prop, int64(v.lookback))
		if err != nil {
			t.Fatal(errors.Wrap(err, "error RSI"))
		}
		if atr == nil {
			t.Errorf("Expected to be non nil but got nil at idx: %d", i)
		}
		if atr.Val() != v.exp {
			t.Errorf("Expected to get %+v but got %+v for lookback %+v", v.exp, *atr.Val(), v.lookback)
		}
	}
}

func ExampleATR() {
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)

	for {
		if v, _ := series.Next(); v == nil {
			break
		}
		tr := OHLCVAttr(series, OHLCPropTR)
		atr, _ := ATR(tr, 3)
		if atr.Val() != nil {
			log.Printf("ATR value: %+v", *atr.Val())
		}
	}
}
