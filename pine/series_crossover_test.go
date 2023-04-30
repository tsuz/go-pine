package pine

import (
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"
)

// TestSeriesCrossover tests no data scenario
func TestSeriesCrossover(t *testing.T) {

	data := OHLCVStaticTestData()

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	c := OHLCVAttr(series, OHLCPropClose)
	o := OHLCVAttr(series, OHLCPropOpen)
	co := Crossover(c, o)
	if co == nil {
		t.Error("Expected co to be non nil but got nil")
	}
}

// TestSeriesCrossoverNoIteration tests this sceneario where there's no iteration yet
func TestSeriesCrossoverNoIteration(t *testing.T) {

	data := OHLCVStaticTestData()
	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	c := OHLCVAttr(series, OHLCPropClose)
	o := OHLCVAttr(series, OHLCPropOpen)
	co := Crossover(c, o)
	if co == nil {
		t.Error("Expected co to be non nil but got nil")
	}
}

// TestSeriesCrossoverIteration tests the output against TradingView's expected values
func TestSeriesCrossoverIteration(t *testing.T) {
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
		0.0,
		1.0,
		0.0,
		0.0,
		1.0,
		0.0,
		0.0,
	}

	for i, v := range tests {
		series.Next()

		c := OHLCVAttr(series, OHLCPropClose)
		o := OHLCVAttr(series, OHLCPropOpen)
		co := Crossover(c, o)
		if err != nil {
			t.Fatal(errors.Wrap(err, "error Crossover"))
		}

		// Lower line
		if *co.Val() != v {
			t.Errorf("Expected lower to be %+v but got %+v for iteration: %d", v, *co.Val(), i)
		}
	}
}

func TestMemoryLeakCrossover(t *testing.T) {
	testMemoryLeak(t, func(o OHLCVSeries) error {
		c := OHLCVAttr(o, OHLCPropClose)
		op := OHLCVAttr(o, OHLCPropOpen)
		Crossover(c, op)
		return nil
	})
}

func BenchmarkCrossover(b *testing.B) {
	// run the Fib function b.N times
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)

	for n := 0; n < b.N; n++ {
		series.Next()
		c := OHLCVAttr(series, OHLCPropClose)
		o := OHLCVAttr(series, OHLCPropOpen)
		Crossover(c, o)
	}
}

func ExampleCrossover() {
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	c := OHLCVAttr(series, OHLCPropClose)
	o := OHLCVAttr(series, OHLCPropOpen)
	co := Crossover(c, o)
	log.Printf("Did Crossover? = %t", *co.Val() == 1.0)
}
