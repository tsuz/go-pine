package pine

import (
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"
)

// TestSeriesSMANoIteration tests when no next() has been called
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

// TestSeriesSMAIteration tests when no next() has been called
func TestSeriesSMAIteration(t *testing.T) {

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
		if sma.Val() != v.exp {
			t.Errorf("Expected to get %+v but got %+v for lookback %+v", v.exp, sma, v.lookback)
		}
	}
}

func TestSeriesSMANotEnoughData(t *testing.T) {
}
