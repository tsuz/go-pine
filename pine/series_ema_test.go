package pine

import (
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"
)

// TestSeriesEMANoData tests no data scenario
//
// t=time.Time (no iteration) | |
// p=ValueSeries              | |
// ema=ValueSeries            | |
func TestSeriesEMANoData(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	prop := series.GetSeries(OHLCPropClose)
	ema, err := EMA(prop, 2)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error EMA"))
	}
	if ema == nil {
		t.Error("Expected to be non nil but got nil")
	}
}

// TestSeriesEMANoIteration tests this sceneario where there's no iteration yet
//
// t=time.Time (no iteration) | 1  |  2   | 3  | 4  |
// p=ValueSeries              | 14 |  15  | 17 | 18 |
// ema=ValueSeries            |    |      |    |    |
func TestSeriesEMANoIteration(t *testing.T) {

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
	ema, err := EMA(prop, 2)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error EMA"))
	}
	if ema == nil {
		t.Error("Expected to be non-nil but got nil")
	}
}

// TestSeriesEMAIteration4 tests this scneario when the iterator is at t=4 is not at the end
//
// t=time.Time     | 1   |  2  | 3           | 4 (time here) |
// p=ValueSeries   | 13  | 15  | 17          | 18            |
// ema(close, 1)   | 13  | 15  | 17          | 18            |
// ema(close, 2)   | nil | 14  | 16          | 17.3333       |
// ema(close, 3)   | nil | nil | 15          | 16.5          |
// ema(close, 4)   | nil | nil | nil         | 15.75         |
// ema(close, 5)   | nil | nil | nil         | nil           |
func TestSeriesEMAIteration4(t *testing.T) {

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
			exp:      17.333333333333332,
		},
		{
			lookback: 3,
			exp:      16.5,
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
		ema, err := EMA(prop, int64(v.lookback))
		if err != nil {
			t.Fatal(errors.Wrap(err, "error EMA"))
		}

		if ema == nil {
			t.Errorf("Expected to be non nil but got nil at idx: %d", i)
		}
		if v.isNil && ema.Val() != nil {
			t.Error("expected to be nil but got non nil")
		}
		if !v.isNil && *ema.Val() != v.exp {
			t.Errorf("Expected to get %+v but got %+v for lookback %+v", v.exp, *ema.Val(), v.lookback)
		}
	}
}

// TestSeriesEMAIteration3 tests this scneario when the iterator is at t=4 is not at the end
//
// t=time.Time    | 1   |  2  | 3 (time here)  | 4             |
// p=ValueSeries  | 13  | 15  | 17             | 18            |
// ema(close, 1)  | 13  | 14  | 17             | 18            |
// ema(close, 2)  | nil | 14  | 16             | 17.3333       |
// ema(close, 3)  | nil | nil | 15             | 16.5          |
// ema(close, 4)  | nil | nil | nil            | 15.75         |
// ema(close, 5)  | nil | nil | nil            | nil           |
func TestSeriesEMAIteration3(t *testing.T) {

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
			exp:      16,
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
		ema, err := EMA(prop, int64(v.lookback))
		if err != nil {
			t.Fatal(errors.Wrap(err, "error EMA"))
		}

		if ema == nil {
			t.Errorf("Expected to be non nil but got nil at idx: %d", i)
		}
		if v.isNil && ema.Val() != nil {
			t.Error("expected to be nil but got non nil")
		}
		if !v.isNil && *ema.Val() != v.exp {
			t.Errorf("Expected to get %+v but got %+v for lookback %+v", v.exp, *ema.Val(), v.lookback)
		}
	}
}

// TestSeriesEMANested tests nested EMA
//

// TestSeriesEMAIteration3 tests this scneario when the iterator is at t=4 is not at the end
//
// t=time.Time 		 	  | 1   |  2  | 3              | 4 (time here) |
// p=ValueSeries          | 13  | 15  | 17             | 18            |
// ema(close, 1)          | 13  | 14  | 17             | 18            |
// ema(close, 2)          | nil | 14  | 16             | 17.3333       |
// ema(ema(close, 2), 2)  | nil | nil | 15             | 16.5555       |
func TestSeriesEMANested(t *testing.T) {

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

	testTable := []float64{15, 16.555555555555554}

	for _, v := range testTable {
		prop := series.GetSeries(OHLCPropClose)

		ema, err := EMA(prop, 2)
		if err != nil {
			t.Fatal(errors.Wrap(err, "error EMA"))
		}
		ema2, err := EMA(ema, 2)
		if err != nil {
			t.Fatal(errors.Wrap(err, "error EMA2"))
		}
		if *ema2.Val() != v {
			t.Errorf("expectd %+v but got %+v", v, *ema2.Val())
		}
		series.Next()
	}
}

// TestSeriesEMANotEnoughData tests this scneario when the lookback is more than the number of data available
//
// t=time.Time    | 1  |  2   | 3  | 4 (here)  |
// p=ValueSeries  | 14 |  15  | 17 | 18        |
// ema(close, 5)  | nil| nil  | nil| nil       |
func TestSeriesEMANotEnoughData(t *testing.T) {

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

		ema, err := EMA(prop, int64(v.lookback))
		if err != nil {
			t.Fatal(errors.Wrap(err, "error EMA"))
		}
		if ema == nil {
			t.Errorf("Expected to be non nil but got nil at idx: %d", i)
		}
		if ema.Val() != v.exp {
			t.Errorf("Expected to get %+v but got %+v for lookback %+v", v.exp, ema, v.lookback)
		}
	}
}

func ExampleEMA() {
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	for {
		if series.Next() == nil {
			break
		}

		close := series.GetSeries(OHLCPropClose)
		ema, err := EMA(close, 20)
		if err != nil {
			log.Fatal(errors.Wrap(err, "error EMA"))
		}
		log.Printf("EMA: %+v", ema.Val())
	}
}
