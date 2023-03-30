package pine

import (
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"
)

// TestSeriesSMANoIteration tests this sceneario where there's no iteration yet
//
// t=time.Time (no iteration) | 1  |  2   | 3  | 4  |
// p=ValueSeries              | 14 |  15  | 17 | 18 |
// sma=ValueSeries            |    |      |    |    |
func TestSeriesSMANoIteration(t *testing.T) {

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
	sma, err := SMA(prop, 2)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error SMA"))
	}
	if sma != nil {
		t.Errorf("Expected to be nil but got %+v", sma)
	}
}

// TestSeriesSMAIteration3 tests this scneario when the iterator is at t=3 is not at the end
//
// t=time.Time (no iteration) | 1  |  2   | 3  (here) | 4  |
// p=ValueSeries              | 13 |  15  | 17        | 18 |
// sma(close, 1)              |    |      | 17        |    |
// sma(close, 2)              |    |      | 16        |    |
// sma(close, 3)              |    |      | 15        |    |
// sma(close, 4)              |    |      | nil       |    |
func TestSeriesSMAIteration3(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)
	data[0].C = 13
	data[1].C = 15
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
	}

	for i, v := range testTable {
		prop := series.GetSeries(OHLCPropClose)

		sma, err := SMA(prop, int64(v.lookback))
		if err != nil {
			t.Fatal(errors.Wrap(err, "error SMA"))
		}

		if sma == nil {
			t.Errorf("Expected to be non nil but got nil at idx: %d", i)
		}
		if v.isNil && sma.Val() != nil {
			t.Error("expected to be nil but got non nil")
		}
		if !v.isNil && *sma.Val() != v.exp {
			t.Errorf("Expected to get %+v but got %+v for lookback %+v", v.exp, *sma.Val(), v.lookback)
		}
	}
}

// TestSeriesSMAIteration4 tests this scneario when the iterator is at t=4
//
// t=time.Time (no iteration) | 1  |  2   | 3  | 4 (here)  |
// p=ValueSeries              | 14 |  15  | 17 | 18        |
// sma(close, 1)              |    |      |    | 18        |
// sma(close, 2)              |    |      |    | 17.5      |
// sma(close, 3)              |    |      |    | 17        |
// sma(close, 4)              |    |      |    | 16.5      |
func TestSeriesSMAIteration4(t *testing.T) {

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
		exp      float64
	}{
		{
			lookback: 1,
			exp:      18,
		},
		{
			lookback: 2,
			exp:      17.5,
		},
		{
			lookback: 3,
			exp:      17,
		},
		{
			lookback: 4,
			exp:      16.5,
		},
	}

	for i, v := range testTable {
		prop := series.GetSeries(OHLCPropClose)

		sma, err := SMA(prop, int64(v.lookback))
		if err != nil {
			t.Fatal(errors.Wrap(err, "error SMA"))
		}
		if sma == nil {
			t.Errorf("Expected to be non nil but got nil at idx: %d", i)
		}
		if *sma.Val() != v.exp {
			t.Errorf("Expected to get %+v but got %+v for lookback %+v", v.exp, sma, v.lookback)
		}
	}
}

// TestSeriesSMANested tests nested SMA
//
// t=time.Time (no iteration) | 1  |  2    | 3     | 4 (here)  |
// p=ValueSeries              | 14 |  15   | 17    | 18        |
// sma(close, 2)              |    |  14.5 | 16    | 17.5      |
// sma(sma(close, 2), 2)      |    |       | 15.25 | 16.75     |
func TestSeriesSMANested(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)
	data[0].C = 14
	data[1].C = 15
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

	testTable := []float64{15.25, 16.75}

	for _, v := range testTable {
		prop := series.GetSeries(OHLCPropClose)

		sma, err := SMA(prop, 2)
		if err != nil {
			t.Fatal(errors.Wrap(err, "error SMA"))
		}
		sma2, err := SMA(sma, int64(2))
		if err != nil {
			t.Fatal(errors.Wrap(err, "error SMA2"))
		}
		if *sma2.Val() != v {
			t.Errorf("expectd %+v but got %+v", v, *sma2.Val())
		}
		series.Next()
	}
}

// TestSeriesSMANotEnoughData tests this scneario when the lookback is more than the number of data available
//
// t=time.Time    | 1  |  2   | 3  | 4 (here)  |
// p=ValueSeries  | 14 |  15  | 17 | 18        |
// sma(close, 5)  |    |      |    | nil       |
func TestSeriesSMANotEnoughData(t *testing.T) {

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

		sma, err := SMA(prop, int64(v.lookback))
		if err != nil {
			t.Fatal(errors.Wrap(err, "error SMA"))
		}
		if sma == nil {
			t.Errorf("Expected to be non nil but got nil at idx: %d", i)
		}
		if sma.Val() != v.exp {
			t.Errorf("Expected to get %+v but got %+v for lookback %+v", v.exp, sma, v.lookback)
		}
	}
}
