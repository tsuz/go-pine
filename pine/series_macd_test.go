package pine

import (
	"fmt"
	"log"
	"testing"
	"time"
)

// TestSeriesMACDNoData tests no data scenario
//
// t=time.Time (no iteration) | |
// p=ValueSeries              | |
// ema=ValueSeries            | |
func TestSeriesMACDNoData(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	prop := OHLCVAttr(series, OHLCPropClose)
	mline, sigline, histline := MACD(prop, 12, 26, 9)
	if mline == nil {
		t.Error("Expected macdline to be non nil but got nil")
	}
	if sigline == nil {
		t.Error("Expected sigline to be non nil but got nil")
	}
	if histline == nil {
		t.Error("Expected histline to be non nil but got nil")
	}
}

// TestSeriesMACDNoIteration tests this sceneario where there's no iteration yet
//
// t=time.Time (no iteration) | 1  |  2   | 3  | 4  |
// p=ValueSeries              | 14 |  15  | 17 | 18 |
// ema=ValueSeries            |    |      |    |    |
func TestSeriesMACDNoIteration(t *testing.T) {

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
	mline, sigline, histline := MACD(prop, 12, 26, 9)
	if mline == nil {
		t.Error("Expected macdline to be non nil but got nil")
	}
	if sigline == nil {
		t.Error("Expected sigline to be non nil but got nil")
	}
	if histline == nil {
		t.Error("Expected histline to be non nil but got nil")
	}
}

// TestSeriesMACDIteration tests this scneario when the iterator is at t=4 is not at the end
//
// t=time.Time                          				| 1   |  2  | 3   | 4       |
// p=ValueSeries                        				| 13  | 15  | 17  | 18      |
// ema(close, 1)                        				| 13  | 15  | 17  | 18      |
// ema(close, 2)                        				| nil | 14  | 16  | 17.3333 |
// MACD line = ema(close, 1) - ema(close,2)         	| nil |  1  |  1  |  0.6667 |
// Signal line = ema(ema(close, 1) - ema(close,2), 2) 	| nil | nil |  1  |  0.7778 |
// MACD Histogram = MACD line - Signal line     	    | nil | nil |  0  | -0.1111 |
func TestSeriesMACDIteration(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)
	data[0].C = 13
	data[1].C = 15
	data[2].C = 17
	data[3].C = 18

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	testTable := []struct {
		macd      *float64
		signal    *float64
		histogram *float64
	}{
		{
			macd:      nil,
			signal:    nil,
			histogram: nil,
		},
		{
			macd:      NewFloat64(1),
			signal:    nil,
			histogram: nil,
		},
		{
			macd:      NewFloat64(1),
			signal:    NewFloat64(1),
			histogram: NewFloat64(0),
		},
		{
			macd:      NewFloat64(0.6667),
			signal:    NewFloat64(0.7778),
			histogram: NewFloat64(-0.1111),
		},
	}

	for i, v := range testTable {
		series.Next()
		src := OHLCVAttr(series, OHLCPropClose)
		macd, signal, histogram := MACD(src, 1, 2, 2)

		// macd line
		if (macd.Val() == nil) != (v.macd == nil) {
			if macd.Val() != nil {
				t.Fatalf("Expected macd to be nil: %t but got %+v for iteration: %d", v.macd == nil, *macd.Val(), i)
			} else {
				t.Fatalf("Expected macd to be nil: %t but got %+v for iteration: %d", v.macd == nil, macd.Val(), i)
			}
		}
		if v.macd != nil && fmt.Sprintf("%.04f", *v.macd) != fmt.Sprintf("%.04f", *macd.Val()) {
			t.Errorf("Expected macd to be %+v but got %+v for iteration: %d", *v.macd, *macd.Val(), i)
		}

		// signal line
		if (signal.Val() == nil) != (v.signal == nil) {
			if signal.Val() != nil {
				t.Fatalf("Expected signal to be nil: %t but got %+v for iteration: %d", v.signal == nil, *signal.Val(), i)
			} else {
				t.Fatalf("Expected signal to be nil: %t but got %+v for iteration: %d", v.signal == nil, signal.Val(), i)
			}
		}
		if v.signal != nil && fmt.Sprintf("%.04f", *v.signal) != fmt.Sprintf("%.04f", *signal.Val()) {
			t.Errorf("Expected signal to be %+v but got %+v for iteration: %d", *v.signal, *signal.Val(), i)
		}

		// macd histogram
		if (histogram.Val() == nil) != (v.histogram == nil) {
			if histogram.Val() != nil {
				t.Fatalf("Expected histogram to be nil: %t but got %+v for iteration: %d", v.histogram == nil, *histogram.Val(), i)
			} else {
				t.Fatalf("Expected histogram to be nil: %t but got %+v for iteration: %d", v.histogram == nil, histogram.Val(), i)
			}
		}
		if v.histogram != nil && fmt.Sprintf("%.04f", *v.histogram) != fmt.Sprintf("%.04f", *histogram.Val()) {
			t.Errorf("Expected histogram to be %+v but got %+v for iteration: %d", *v.histogram, *histogram.Val(), i)
		}
	}
}

func TestMemoryLeakMACD(t *testing.T) {
	testMemoryLeak(t, func(o OHLCVSeries) error {
		MACD(OHLCVAttr(o, OHLCPropClose), 12, 26, 9)
		return nil
	})
}

func BenchmarkMACD(b *testing.B) {
	// run the Fib function b.N times
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	vals := OHLCVAttr(series, OHLCPropClose)

	for n := 0; n < b.N; n++ {
		series.Next()
		MACD(vals, 12, 26, 9)
	}
}

func ExampleMACD() {
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	close := OHLCVAttr(series, OHLCPropClose)
	mline, sigline, histline := MACD(close, 12, 26, 9)
	log.Printf("MACD line: %+v", mline.Val())
	log.Printf("Signal line: %+v", sigline.Val())
	log.Printf("Hist line: %+v", histline.Val())
}
