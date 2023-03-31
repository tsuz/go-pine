package pine

import (
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"
)

// TestSeriesRMANoData tests no data scenario
//
// t=time.Time (no iteration) | |
// p=ValueSeries              | |
// rma=ValueSeries            | |
func TestSeriesRMANoData(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	prop := series.GetSeries(OHLCPropClose)
	rma, err := RMA(prop, 2)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error RMA"))
	}
	if rma == nil {
		t.Error("Expected to be non nil but got nil")
	}
}

// TestSeriesRMANoIteration tests this sceneario where there's no iteration yet
//
// t=time.Time (no iteration) | 1  |  2   | 3  | 4  |
// p=ValueSeries              | 14 |  15  | 17 | 18 |
// rma=ValueSeries            |    |      |    |    |
func TestSeriesRMANoIteration(t *testing.T) {

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

	prop := series.GetSeries(OHLCPropClose)
	rma, err := RMA(prop, 2)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error RMA"))
	}
	if rma == nil {
		t.Error("Expected to be non-nil but got nil")
	}
}

// TestSeriesRMAIteration4 tests this scneario when the iterator is at t=4 is not at the end
//
// t=time.Time     | 1   |  2  | 3           | 4 (time here) |
// p=ValueSeries   | 13  | 15  | 17          | 18            |
// rma(close, 1)   | 13  | 15  | 17          | 18            |
// rma(close, 2)   | nil | 14  | 15.5        | 16.75         |
// rma(close, 3)   | nil | nil | 15          | 16            |
// rma(close, 4)   | nil | nil | nil         | 15.75         |
// rma(close, 5)   | nil | nil | nil         | nil           |
func TestSeriesRMAIteration4(t *testing.T) {

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

	series.Next()
	series.Next()
	series.Next()
	series.Next()

	testTable := []struct {
		lookback int
		exp      float64
		isNil    bool
	}{
		{
			lookback: 1,
			exp:      18,
		},
		{
			lookback: 2,
			exp:      16.75,
		},
		{
			lookback: 3,
			exp:      16,
		},
		{
			lookback: 4,
			exp:      15.75,
		},
		{
			lookback: 5,
			exp:      0,
			isNil:    true,
		},
	}

	for i, v := range testTable {
		prop := series.GetSeries(OHLCPropClose)
		rma, err := RMA(prop, int64(v.lookback))
		if err != nil {
			t.Fatal(errors.Wrap(err, "error RMA"))
		}

		if rma == nil {
			t.Errorf("Expected to be non nil but got nil at idx: %d", i)
		}
		if v.isNil && rma.Val() != nil {
			t.Error("expected to be nil but got non nil")
		}
		log.Printf("Index %d", i)
		if !v.isNil && *rma.Val() != v.exp {
			t.Errorf("Expected to get %+v but got %+v for lookback %+v", v.exp, *rma.Val(), v.lookback)
		}
	}
}

// TestSeriesRMAIteration3 tests this scneario when the iterator is at t=4 is not at the end
//
// t=time.Time     | 1   |  2  | 3 (time here) | 4     |
// p=ValueSeries   | 13  | 15  | 17            | 18    |
// rma(close, 1)   | 13  | 15  | 17            | 18    |
// rma(close, 2)   | nil | 14  | 15.5          | 16.75 |
// rma(close, 3)   | nil | nil | 15            | 16    |
// rma(close, 4)   | nil | nil | nil           | 15.75 |
// rma(close, 5)   | nil | nil | nil           | nil   |
func TestSeriesRMAIteration3(t *testing.T) {

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

	series.Next()
	series.Next()
	series.Next()

	testTable := []struct {
		lookback int
		exp      float64
		isNil    bool
	}{
		{
			lookback: 1,
			exp:      17,
		},
		{
			lookback: 2,
			exp:      15.5,
		},
		{
			lookback: 3,
			exp:      15,
		},
		{
			lookback: 4,
			exp:      0,
			isNil:    true,
		},
		{
			lookback: 5,
			exp:      0,
			isNil:    true,
		},
	}

	for i, v := range testTable {
		prop := series.GetSeries(OHLCPropClose)
		rma, err := RMA(prop, int64(v.lookback))
		if err != nil {
			t.Fatal(errors.Wrap(err, "error RMA"))
		}

		if rma == nil {
			t.Errorf("Expected to be non nil but got nil at idx: %d", i)
		}
		if v.isNil && rma.Val() != nil {
			t.Error("expected to be nil but got non nil")
		}
		log.Printf("Index %d", i)
		if !v.isNil && *rma.Val() != v.exp {
			t.Errorf("Expected to get %+v but got %+v for lookback %+v", v.exp, *rma.Val(), v.lookback)
		}
	}
}

// TestSeriesRMANested tests nested RMA
//
// t=time.Time 	           | 1   |  2  | 3 (time here 1) | 4 (time here2)
// p=ValueSeries           | 13  | 15  | 17              | 18
// rma(close, 2)           | nil | 14  | 15.5            | 16.75
// rma(rma(close, 2), 2)   | nil | nil | 14.75           | 15.75
func TestSeriesRMANested(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 5, 5*60*1000)
	data[0].C = 13
	data[1].C = 15
	data[2].C = 17
	data[3].C = 18

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	series.Next()
	series.Next()
	series.Next()

	testTable := []float64{14.75, 15.75}

	for _, v := range testTable {
		prop := series.GetSeries(OHLCPropClose)

		rma, err := RMA(prop, 2)
		if err != nil {
			t.Fatal(errors.Wrap(err, "error RMA"))
		}
		rma2, err := RMA(rma, 2)
		if err != nil {
			t.Fatal(errors.Wrap(err, "error RMA2"))
		}
		if *rma2.Val() != v {
			t.Errorf("expectd %+v but got %+v", v, *rma2.Val())
		}
		series.Next()
	}
}

// TestSeriesRMANotEnoughData tests this scneario when the lookback is more than the number of data available
//
// t=time.Time    | 1  |  2   | 3  | 4 (here)  |
// p=ValueSeries  | 14 |  15  | 17 | 18        |
// rma(close, 5)  | nil| nil  | nil| nil       |
func TestSeriesRMANotEnoughData(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)
	data[0].C = 15
	data[1].C = 16
	data[2].C = 17
	data[3].C = 18

	log.Printf("Data[0].S, %+v, 3s: %+v", data[0].S, data[3].S)

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	series.Next()
	series.Next()
	series.Next()
	series.Next()

	testTable := []struct {
		lookback int
		exp      *float64
	}{
		{
			lookback: 5,
			exp:      nil,
		},
		{
			lookback: 6,
			exp:      nil,
		},
	}

	for i, v := range testTable {
		prop := series.GetSeries(OHLCPropClose)

		rma, err := RMA(prop, int64(v.lookback))
		if err != nil {
			t.Fatal(errors.Wrap(err, "error RMA"))
		}
		if rma == nil {
			t.Errorf("Expected to be non nil but got nil at idx: %d", i)
		}
		if rma.Val() != v.exp {
			t.Errorf("Expected to get %+v but got %+v for lookback %+v", v.exp, rma, v.lookback)
		}
	}
}
