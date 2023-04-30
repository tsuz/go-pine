package pine

import (
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"
)

// TestSeriesValueWhenNoData tests no data scenario
//
// t=time.Time (no iteration) | |
// p=ValueSeries              | |
// valueWhen=ValueSeries      | |
func TestSeriesValueWhenNoData(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 0, 5*60*1000)

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	prop := OHLCVAttr(series, OHLCPropClose)
	bs := NewValueSeries()

	rsi, err := ValueWhen(bs, prop, 2)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error ValueWhen"))
	}
	if rsi == nil {
		t.Error("Expected to be non nil but got nil")
	}
}

// TestSeriesValueWhenNoIteration tests this sceneario where there's no iteration yet
//
// t=time.Time (no iteration) | 1  |  2   | 3  | 4  |
// p=ValueSeries              | 14 |  15  | 17 | 18 |
// bs=ValueSeries             |1.0 | 0.0  |1.0 |0.0 |
// valuewhen(0)=ValueSeries   |    |      |    |    |
func TestSeriesValueWhenNoIteration(t *testing.T) {

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

	bs := NewValueSeries()
	bs.Set(data[0].S, 1)
	bs.Set(data[1].S, 0)
	bs.Set(data[2].S, 1)
	bs.Set(data[3].S, 0)

	prop := OHLCVAttr(series, OHLCPropClose)
	rsi, err := ValueWhen(bs, prop, 0)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error ValueWhen"))
	}
	if rsi == nil {
		t.Error("Expected to be non-nil but got nil")
	}
}

// TestSeriesValueWhenIteration5 tests this scneario when the iterator is at t=4 is not at the end
//
// t=time.Time (no iteration) | 1   |  2  | 3   | 4   | 5  | 6  |
// bs=ValueSeries             | 0   |  1  | 0   | 1   | 1  | 0  |
// src=ValueSeries            | 13  | 15  | 11  | 18  | 20 | 17 |
// valuewhen(0)=ValueSeries   | nil | 15  | 15  | 18  | 20 | 20 |
// valuewhen(1)=ValueSeries   | nil | nil | nil | 15  | 18 | 18 |
// valuewhen(2)=ValueSeries   | nil | nil | nil | nil | 15 | 15 |
func TestSeriesValueWhenIteration5(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 6, 5*60*1000)
	data[0].C = 13
	data[1].C = 15
	data[2].C = 11
	data[3].C = 18
	data[4].C = 20
	data[5].C = 17

	bs := NewValueSeries()
	bs.Set(data[0].S, 0)
	bs.Set(data[1].S, 1)
	bs.Set(data[2].S, 0)
	bs.Set(data[3].S, 1)
	bs.Set(data[4].S, 1)
	bs.Set(data[5].S, 0)

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	testTable := []struct {
		ocr  int
		vals []float64
	}{
		{
			ocr:  0,
			vals: []float64{0, 15, 15, 18, 20, 20},
		},
		{
			ocr:  1,
			vals: []float64{0, 0, 0, 15, 18, 18},
		},
		{
			ocr:  2,
			vals: []float64{0, 0, 0, 0, 15, 15},
		},
	}

	for j := 0; j <= 5; j++ {
		series.Next()

		for i, v := range testTable {
			prop := OHLCVAttr(series, OHLCPropClose)
			vw, err := ValueWhen(bs, prop, v.ocr)
			if err != nil {
				t.Fatal(errors.Wrap(err, "error ValueWhen"))
			}
			exp := v.vals[j]
			if exp == 0 {
				if vw.Val() != nil {
					t.Fatalf("expected nil but got non nil: %+v at vals item: %d, testtable item: %d", *vw.Val(), j, i)
				}
				// OK
			}
			if exp != 0 {
				if vw.Val() == nil {
					t.Fatalf("expected non nil: %+v but got nil at vals item: %d, testtable item: %d", exp, j, i)
				}
				if exp != *vw.Val() {
					t.Fatalf("expected %+v but got %+v at vals item: %d, testtable item: %d", exp, *vw.Val(), j, i)
				}
				// OK
			}
		}
	}
}

// TestSeriesValueWhenNotEnoughData tests this scneario when the lookback is more than the number of data available
//
// t=time.Time          | 1  |  2   | 3  | 4 (here)  |
// p=ValueSeries        | 14 |  15  | 17 | 18        |
// valuewhen(close, 5)  | nil| nil  | nil| nil       |
func TestSeriesValueWhenNotEnoughData(t *testing.T) {

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

	bs := NewValueSeries()
	bs.Set(data[0].S, 0)
	bs.Set(data[1].S, 1)
	bs.Set(data[2].S, 0)
	bs.Set(data[3].S, 1)

	prop := OHLCVAttr(series, OHLCPropClose)

	vw, err := ValueWhen(bs, prop, 5)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error ValueWhen"))
	}
	if vw.Val() != nil {
		t.Errorf("Expected nil but got %+v", *vw.Val())
	}
}

func TestMemoryLeakValueWhen(t *testing.T) {
	bs := NewValueSeries()
	testMemoryLeak(t, func(o OHLCVSeries) error {
		prop := OHLCVAttr(o, OHLCPropClose)
		if c := prop.GetCurrent(); c != nil {
			bs.Set(c.t, float64(int(c.v)%2))
		}
		_, err := ValueWhen(bs, prop, 10)
		return err
	})
}

func BenchmarkValueWhen(b *testing.B) {
	// run the Fib function b.N times
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	vals := OHLCVAttr(series, OHLCPropClose)

	bs := NewValueSeries()
	for _, v := range data {
		bs.Set(v.S, float64(int(v.C)%2))
	}

	for n := 0; n < b.N; n++ {
		series.Next()
		ValueWhen(bs, vals, 5)
	}
}

func ExampleValueWhen() {
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)

	// value series with 0.0 or 1.0 (true/false)
	bs := NewValueSeries()
	for _, v := range data {
		bs.Set(v.S, float64(int(v.C)%2))
	}

	for {
		if v, _ := series.Next(); v == nil {
			break
		}

		close := OHLCVAttr(series, OHLCPropClose)
		vw, err := ValueWhen(close, bs, 12)
		if err != nil {
			log.Fatal(errors.Wrap(err, "error geting ValueWhen"))
		}
		log.Printf("ValueWhen: %+v", vw.Val())
	}
}
