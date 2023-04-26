package pine

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"
)

// TestSeriesWilderSMA tests no data scenario
//
// t=time.Time (no iteration) | |
// p=ValueSeries              | |
// wilderma=ValueSeries            | |
func TestSeriesWilderSMA(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}
	close := series.GetSeries(OHLCPropClose)

	mfi, err := WilderSMA(close, 3)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error WilderSMA"))
	}
	if mfi == nil {
		t.Error("Expected mfi to be non nil but got nil")
	}
}

// TestSeriesWilderSMANoIteration tests this sceneario where there's no iteration yet
//
// t=time.Time (no iteration) | 1  |  2   | 3  | 4  |
// p=ValueSeries              | 14 |  15  | 17 | 18 |
// mfi=ValueSeries            |    |      |    |    |
func TestSeriesWilderSMANoIteration(t *testing.T) {

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

	mfi, err := WilderSMA(series.GetSeries(OHLCPropClose), 3)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error WilderSMA"))
	}
	if mfi == nil {
		t.Error("Expected mfi to be non nil but got nil")
	}
}

// TestSeriesWilderSMAIteration tests the output against TradingView's expected values
func TestSeriesWilderSMAIteration(t *testing.T) {
	data := OHLCVStaticTestData()
	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	n := 4.0
	wsmafn := func(prevwsma, src float64) float64 {
		return (prevwsma*(n-1) + src) / n
	}
	wsma1 := (data[0].C + data[1].C + data[2].C + data[3].C) / 4
	wsma2 := wsmafn(wsma1, data[4].C)
	wsma3 := wsmafn(wsma2, data[5].C)
	wsma4 := wsmafn(wsma3, data[6].C)
	wsma5 := wsmafn(wsma4, data[7].C)
	wsma6 := wsmafn(wsma5, data[8].C)
	wsma7 := wsmafn(wsma6, data[9].C)

	tests := []*float64{
		nil,
		nil,
		nil,
		NewFloat64(wsma1),
		NewFloat64(wsma2),
		NewFloat64(wsma3),
		NewFloat64(wsma4),
		NewFloat64(wsma5),
		NewFloat64(wsma6),
		NewFloat64(wsma7),
	}

	for i, v := range tests {
		series.Next()
		c := series.GetSeries(OHLCPropClose)
		wsma, err := WilderSMA(c, int(n))
		if err != nil {
			t.Fatal(errors.Wrap(err, "error mfi"))
		}

		// mfi line
		if (wsma.Val() == nil) != (v == nil) {
			if wsma.Val() != nil {
				t.Errorf("Expected wsma to be nil: %t but got %+v for iteration: %d", v == nil, *wsma.Val(), i)
			} else {
				t.Errorf("Expected wsma to be: %+v but got %+v for iteration: %d", *v, wsma.Val(), i)
			}
			continue
		}
		if v != nil && fmt.Sprintf("%.03f", *v) != fmt.Sprintf("%.03f", *wsma.Val()) {
			t.Errorf("Expected wsma to be %+v but got %+v for iteration: %d", *v, *wsma.Val(), i)
		}
	}
}

func BenchmarkWilderSMA(b *testing.B) {
	// run the Fib function b.N times
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)

	for n := 0; n < b.N; n++ {
		series.Next()
		WilderSMA(series.GetSeries(OHLCPropClose), 12)
	}
}

func ExampleWilderSMA() {
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	mfi, err := WilderSMA(series.GetSeries(OHLCPropClose), 12)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error WilderSMA"))
	}
	log.Printf("WilderSMA line: %+v", mfi.Val())
}
