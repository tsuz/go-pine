package pine

import (
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"
)

// TestSeriesCrossunder tests no data scenario
func TestSeriesCrossunder(t *testing.T) {

	data := OHLCVStaticTestData()

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	c := OHLCVAttr(series, OHLCPropClose)
	o := OHLCVAttr(series, OHLCPropOpen)
	co, err := Crossunder(c, o)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error Crossunder"))
	}
	if co == nil {
		t.Error("Expected co to be non nil but got nil")
	}
}

// TestSeriesCrossunderNoIteration tests this sceneario where there's no iteration yet
func TestSeriesCrossunderNoIteration(t *testing.T) {

	data := OHLCVStaticTestData()
	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	c := OHLCVAttr(series, OHLCPropClose)
	o := OHLCVAttr(series, OHLCPropOpen)
	co, err := Crossunder(c, o)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error Crossunder"))
	}
	if co == nil {
		t.Error("Expected co to be non nil but got nil")
	}
}

// TestSeriesCrossunderIteration tests the output against TradingView's expected values
func TestSeriesCrossunderIteration(t *testing.T) {
	data := OHLCVStaticTestData()
	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}
	// array in order of Middle, Upper, Lower
	tests := []float64{
		0.0,
		0.0,
		0.0,
		1.0,
		0.0,
		1.0,
		0.0,
		0.0,
		1.0,
		0.0,
	}

	for i, v := range tests {
		series.Next()

		c := OHLCVAttr(series, OHLCPropClose)
		o := OHLCVAttr(series, OHLCPropOpen)
		co, err := Crossunder(c, o)
		if err != nil {
			t.Fatal(errors.Wrap(err, "error Crossunder"))
		}

		// Lower line
		if *co.Val() != v {
			t.Errorf("Expected lower to be %+v but got %+v for iteration: %d", v, *co.Val(), i)
		}
	}
}

func BenchmarkCrossunder(b *testing.B) {
	// run the Fib function b.N times
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)

	for n := 0; n < b.N; n++ {
		series.Next()
		c := OHLCVAttr(series, OHLCPropClose)
		o := OHLCVAttr(series, OHLCPropOpen)
		Crossunder(c, o)
	}
}

func ExampleCrossunder() {
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	c := OHLCVAttr(series, OHLCPropClose)
	o := OHLCVAttr(series, OHLCPropOpen)
	co, err := Crossunder(c, o)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error Crossunder"))
	}
	log.Printf("Did Crossunder? = %t", *co.Val() == 1.0)
}
