package pine

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"
)

// TestSeriesDMI tests no data scenario
//
// t=time.Time (no iteration) | |
// p=ValueSeries              | |
// dmi=ValueSeries            | |
func TestSeriesDMI(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	adx, dip, dim, err := DMI(series, 15, 3)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error DMI"))
	}
	if dip == nil {
		t.Error("Expected dip to be non nil but got nil")
	}
	if dim == nil {
		t.Error("Expected dim to be non nil but got nil")
	}
	if adx == nil {
		t.Error("Expected adx to be non nil but got nil")
	}
}

// TestSeriesDMINoIteration tests this sceneario where there's no iteration yet
func TestSeriesDMINoIteration(t *testing.T) {

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

	adx, _, _, err := DMI(series, 3, 2)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error DMI"))
	}
	if adx == nil {
		t.Error("Expected dmi to be non nil but got nil")
	}
}

// TestSeriesDMIIteration tests the output against TradingView's expected values
func TestSeriesDMIIteration(t *testing.T) {
	data := OHLCVStaticTestData()
	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	// array in order of ADX, DI+, DI-
	tests := [][]*float64{
		nil,
		nil,
		nil,
		nil,
		{nil, NewFloat64(2.49), NewFloat64(7.78)},
		{nil, NewFloat64(2.96), NewFloat64(6.17)},
		{NewFloat64(45.43), NewFloat64(2.3), NewFloat64(6.82)},
		{NewFloat64(56.3), NewFloat64(1.6), NewFloat64(12.97)},
		{NewFloat64(63.55), NewFloat64(1.2), NewFloat64(9.7)},
		{NewFloat64(71.9), NewFloat64(0.9067), NewFloat64(14.99)},
	}

	for i, v := range tests {
		series.Next()
		adx, dmip, dmim, err := DMI(series, 4, 2)
		if err != nil {
			t.Fatal(errors.Wrap(err, "error dmi"))
		}

		// list can be empty
		if v == nil {
			// if adx.Val() != nil ||
			if dmip.Val() != nil || dmim.Val() != nil {
				t.Errorf("Expected no values to be returned but got some at %d", i)
			}
			continue
		}

		// ADX line
		if v[0] != nil && fmt.Sprintf("%.01f", *v[0]) != fmt.Sprintf("%.01f", *adx.Val()) {
			t.Errorf("Expected dmi to be %+v but got %+v for iteration: %d", *v[0], *adx.Val(), i)
		}

		// DMI+ line
		if v != nil && fmt.Sprintf("%.01f", *v[1]) != fmt.Sprintf("%.01f", *dmip.Val()) {
			t.Errorf("Expected dmi to be %+v but got %+v for iteration: %d", *v[1], *dmip.Val(), i)
		}

		// DMI- line
		if v != nil && fmt.Sprintf("%.01f", *v[2]) != fmt.Sprintf("%.01f", *dmim.Val()) {
			t.Errorf("Expected dmi to be %+v but got %+v for iteration: %d", *v[2], *dmim.Val(), i)
		}
	}
}

// func TestMemoryLeakDMI(t *testing.T) {
// 	testMemoryLeak(t, func(o OHLCVSeries) error {
// 		_, _, _, err := DMI(o, 4, 3)
// 		return err
// 	})
// }

func BenchmarkDMI(b *testing.B) {
	// run the Fib function b.N times
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)

	for n := 0; n < b.N; n++ {
		series.Next()
		DMI(series, 4, 3)
	}
}

func ExampleDMI() {
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	adx, dmip, dmim, err := DMI(series, 4, 3)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error DMI"))
	}
	log.Printf("ADX: %+v, DI+: %+v, DI-: %+v", adx.Val(), dmip.Val(), dmim.Val())
}
