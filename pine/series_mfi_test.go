package pine

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"
)

// TestSeriesMFI tests no data scenario
//
// t=time.Time (no iteration) | |
// p=ValueSeries              | |
// mfi=ValueSeries            | |
func TestSeriesMFI(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	mfi, err := MFI(series, 3)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error MFI"))
	}
	if mfi == nil {
		t.Error("Expected mfi to be non nil but got nil")
	}
}

// TestSeriesMFINoIteration tests this sceneario where there's no iteration yet
//
// t=time.Time (no iteration) | 1  |  2   | 3  | 4  |
// p=ValueSeries              | 14 |  15  | 17 | 18 |
// mfi=ValueSeries            |    |      |    |    |
func TestSeriesMFINoIteration(t *testing.T) {

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

	mfi, err := MFI(series, 3)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error MFI"))
	}
	if mfi == nil {
		t.Error("Expected mfi to be non nil but got nil")
	}
}

// TestSeriesMFIIteration tests the output against TradingView's expected values
func TestSeriesMFIIteration(t *testing.T) {
	data := OHLCVStaticTestData()
	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	tests := []*float64{
		nil,
		nil,
		nil,
		NewFloat64(38.856),
		NewFloat64(52.679),
		NewFloat64(27.212),
		NewFloat64(26.905),
		NewFloat64(28.794),
		NewFloat64(27.858),
		NewFloat64(31.572),
	}

	for i, v := range tests {
		series.Next()
		mfi, err := MFI(series, 4)
		if err != nil {
			t.Fatal(errors.Wrap(err, "error mfi"))
		}

		// mfi line
		if (mfi.Val() == nil) != (v == nil) {
			if mfi.Val() != nil {
				t.Errorf("Expected mfi to be nil: %t but got %+v for iteration: %d", v == nil, *mfi.Val(), i)
			} else {
				t.Errorf("Expected mfi to be: %+v but got %+v for iteration: %d", *v, mfi.Val(), i)
			}
			continue
		}
		if v != nil && fmt.Sprintf("%.03f", *v) != fmt.Sprintf("%.03f", *mfi.Val()) {
			t.Errorf("Expected mfi to be %+v but got %+v for iteration: %d", *v, *mfi.Val(), i)
		}
	}
}

func TestMemoryLeakMFI(t *testing.T) {
	testMemoryLeak(t, func(o OHLCVSeries) error {
		_, err := MFI(o, 12)
		return err
	})
}

func BenchmarkMFI(b *testing.B) {
	// run the Fib function b.N times
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)

	for n := 0; n < b.N; n++ {
		series.Next()
		MFI(series, 12)
	}
}

func ExampleMFI() {
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	mfi, err := MFI(series, 12)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error MFI"))
	}
	log.Printf("MFI line: %+v", mfi.Val())
}
