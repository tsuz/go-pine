package pine

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"
)

// TestSeriesCCI tests no data scenario
//
// t=time.Time (no iteration) | |
// p=ValueSeries              | |
// cci=ValueSeries            | |
func TestSeriesCCI(t *testing.T) {

	data := OHLCVStaticTestData()

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}
	tp := OHLCVAttr(series, OHLCPropHLC3)

	cci, err := CCI(tp, 3)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error CCI"))
	}
	if cci == nil {
		t.Error("Expected cci to be non nil but got nil")
	}
}

// TestSeriesCCINoIteration tests this sceneario where there's no iteration yet

// t=time.Time (no iteration) | 1  |  2   | 3  | 4  |
// p=ValueSeries              | 14 |  15  | 17 | 18 |
// cci=ValueSeries            |    |      |    |    |
func TestSeriesCCINoIteration(t *testing.T) {

	data := OHLCVStaticTestData()

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}
	tp := OHLCVAttr(series, OHLCPropHLC3)

	cci, err := CCI(tp, 3)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error CCI"))
	}
	if cci == nil {
		t.Error("Expected cci to be non nil but got nil")
	}
}

// TestSeriesCCIIteration tests the output against TradingView's expected values
func TestSeriesCCIIteration(t *testing.T) {
	data := OHLCVStaticTestData()
	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	tests := []*float64{
		nil,
		nil,
		nil,
		NewFloat64(-133.3),
		NewFloat64(65.3),
		NewFloat64(17.5),
		NewFloat64(-2.7),
		NewFloat64(-133.3),
		NewFloat64(22.2),
		NewFloat64(-98.6),
	}

	for i, v := range tests {
		series.Next()
		tp := OHLCVAttr(series, OHLCPropHLC3)
		cci, err := CCI(tp, 4)
		if err != nil {
			t.Fatal(errors.Wrap(err, "error cci"))
		}

		// cci line
		if (cci.Val() == nil) != (v == nil) {
			if cci.Val() != nil {
				t.Errorf("Expected cci to be nil: %t but got %+v for iteration: %d", v == nil, *cci.Val(), i)
			} else {
				t.Errorf("Expected cci to be: %+v but got %+v for iteration: %d", *v, cci.Val(), i)
			}
			continue
		}
		if v != nil && fmt.Sprintf("%.01f", *v) != fmt.Sprintf("%.01f", *cci.Val()) {
			t.Errorf("Expected cci to be %+v but got %+v for iteration: %d", *v, *cci.Val(), i)
		}
	}
}

// func TestMemoryLeakCCI(t *testing.T) {
// 	testMemoryLeak(t, func(o OHLCVSeries) error {
// 		c := OHLCVAttr(o, OHLCPropClose)
// 		_, err := CCI(c, 7)
// 		return err
// 	})
// }

func BenchmarkCCI(b *testing.B) {
	// run the Fib function b.N times
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)

	for n := 0; n < b.N; n++ {
		series.Next()
		tp := OHLCVAttr(series, OHLCPropHLC3)
		CCI(tp, 12)
	}
}

func ExampleCCI() {
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	tp := OHLCVAttr(series, OHLCPropHLC3)
	cci, err := CCI(tp, 12)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error CCI"))
	}
	log.Printf("CCI line: %+v", cci.Val())
}
