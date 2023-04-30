package pine

import (
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"
)

// TestSeriesCross tests no data scenario
func TestSeriesCross(t *testing.T) {

	data := OHLCVStaticTestData()

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	c := OHLCVAttr(series, OHLCPropClose)
	o := OHLCVAttr(series, OHLCPropOpen)
	co, err := Cross(c, o)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error Cross"))
	}
	if co == nil {
		t.Error("Expected co to be non nil but got nil")
	}
}

// TestSeriesCrossNoIteration tests this sceneario where there's no iteration yet
func TestSeriesCrossNoIteration(t *testing.T) {

	data := OHLCVStaticTestData()
	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	c := OHLCVAttr(series, OHLCPropClose)
	o := OHLCVAttr(series, OHLCPropOpen)
	co, err := Cross(c, o)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error Cross"))
	}
	if co == nil {
		t.Error("Expected co to be non nil but got nil")
	}
}

// TestSeriesCrossIteration tests the output against TradingView's expected values
func TestSeriesCrossIteration(t *testing.T) {
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
		1.0,
		1.0,
		0.0,
		1.0,
		1.0,
		0.0,
	}

	for i, v := range tests {
		series.Next()

		c := OHLCVAttr(series, OHLCPropClose)
		o := OHLCVAttr(series, OHLCPropOpen)
		co, err := Cross(c, o)
		if err != nil {
			t.Fatal(errors.Wrap(err, "error Cross"))
		}

		// Lower line
		if *co.Val() != v {
			t.Errorf("Expected lower to be %+v but got %+v for iteration: %d", v, *co.Val(), i)
		}
	}
}

func TestMemoryLeakCross(t *testing.T) {
	testMemoryLeak(t, func(o OHLCVSeries) error {
		c := OHLCVAttr(o, OHLCPropClose)
		op := OHLCVAttr(o, OHLCPropOpen)
		_, err := Cross(c, op)
		return err
	})
}

func BenchmarkCross(b *testing.B) {
	// run the Fib function b.N times
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)

	for n := 0; n < b.N; n++ {
		series.Next()
		c := OHLCVAttr(series, OHLCPropClose)
		o := OHLCVAttr(series, OHLCPropOpen)
		Cross(c, o)
	}
}

func ExampleCross() {
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	c := OHLCVAttr(series, OHLCPropClose)
	o := OHLCVAttr(series, OHLCPropOpen)
	co, err := Cross(c, o)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error Cross"))
	}
	log.Printf("Did Cross? = %t", *co.Val() == 1.0)
}
