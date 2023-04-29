package pine

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"
)

// TestSeriesKC tests no data scenario
func TestSeriesKC(t *testing.T) {

	data := OHLCVStaticTestData()

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}
	close := OHLCVAttr(series, OHLCPropClose)

	m, u, l, err := KC(close, series, 3, 2.5, true)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error KC"))
	}
	if m == nil {
		t.Error("Expected kc to be non nil but got nil")
	}
	if u == nil {
		t.Error("Expected kc to be non nil but got nil")
	}
	if l == nil {
		t.Error("Expected kc to be non nil but got nil")
	}
}

// TestSeriesKCNoIteration tests this sceneario where there's no iteration yet
func TestSeriesKCNoIteration(t *testing.T) {

	data := OHLCVStaticTestData()
	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}
	close := OHLCVAttr(series, OHLCPropClose)

	m, u, l, err := KC(close, series, 3, 2.5, true)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error KC"))
	}
	if m == nil {
		t.Error("Expected kc to be non nil but got nil")
	}
	if u == nil {
		t.Error("Expected kc to be non nil but got nil")
	}
	if l == nil {
		t.Error("Expected kc to be non nil but got nil")
	}
}

// TestSeriesKCIteration tests the output against TradingView's expected values
func TestSeriesKCIteration(t *testing.T) {
	data := OHLCVStaticTestData()
	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}
	// array in order of Middle, Upper, Lower
	tests := [][]*float64{
		nil,
		nil,
		nil,
		{NewFloat64(16.33), NewFloat64(36.2), NewFloat64(-3.55)},
		{NewFloat64(17.52), NewFloat64(37.74), NewFloat64(-2.71)},
		{NewFloat64(16.19), NewFloat64(34.62), NewFloat64(-2.25)},
		{NewFloat64(15.47), NewFloat64(33.13), NewFloat64(-2.19)},
		{NewFloat64(13.68), NewFloat64(33.88), NewFloat64(-6.51)},
		{NewFloat64(14.09), NewFloat64(32.81), NewFloat64(-4.63)},
		{NewFloat64(12.57), NewFloat64(31.41), NewFloat64(-6.26)},
	}

	for i, v := range tests {
		series.Next()
		c := OHLCVAttr(series, OHLCPropClose)
		m, u, l, err := KC(c, series, 4, 2.5, false)
		if err != nil {
			t.Fatal(errors.Wrap(err, "error dmi"))
		}

		// list can be empty
		if v == nil {
			if m.Val() != nil || u.Val() != nil || l.Val() != nil {
				t.Errorf("Expected no values to be returned but got some at %d", i)
			}
			continue
		}

		// Middle line
		if v[0] != nil && fmt.Sprintf("%.01f", *v[0]) != fmt.Sprintf("%.01f", *m.Val()) {
			t.Errorf("Expected middle to be %+v but got %+v for iteration: %d", *v[0], *m.Val(), i)
		}

		// Upper line
		if v != nil && fmt.Sprintf("%.01f", *v[1]) != fmt.Sprintf("%.01f", *u.Val()) {
			t.Errorf("Expected upper to be %+v but got %+v for iteration: %d", *v[1], *u.Val(), i)
		}

		// Lower line
		if v != nil && fmt.Sprintf("%.01f", *v[2]) != fmt.Sprintf("%.01f", *l.Val()) {
			t.Errorf("Expected lower to be %+v but got %+v for iteration: %d", *v[2], *l.Val(), i)
		}
	}
}

func BenchmarkKC(b *testing.B) {
	// run the Fib function b.N times
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)

	for n := 0; n < b.N; n++ {
		series.Next()
		KC(OHLCVAttr(series, OHLCPropClose), series, 4, 2.5, false)
	}
}

func ExampleKC() {
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	m, u, l, err := KC(OHLCVAttr(series, OHLCPropClose), series, 4, 2.5, false)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error KC"))
	}
	log.Printf("KC middle line: %+v, upper: %+v, lower: %+v", m.Val(), u.Val(), l.Val())
}
