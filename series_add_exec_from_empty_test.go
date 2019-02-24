package pine

import (
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestSeriesAddExecFromEmptyData(t *testing.T) {
	opts := SeriesOpts{
		Interval: 300,
		Max:      100,
	}
	_, err := NewSeries(nil, opts)
	if err != nil {
		t.Fatal(err)
	}
	close := NewOHLCProp(OHLCPropClose)
	sma := NewSMA(close, 2)
	now := time.Now()
	fivemin := now.Add(5 * time.Minute)
	ten := now.Add(10 * time.Minute)
	fifteen := now.Add(15 * time.Minute)
	data := []OHLCV{}
	s, err := NewSeries(data, opts)
	if err != nil {
		t.Fatal(err)
	}
	s.AddIndicator("sma", sma)

	tpqs := []TPQ{
		TPQ{
			Qty:       5,
			Px:        14,
			Timestamp: now,
		},
		TPQ{
			Qty:       2,
			Px:        15,
			Timestamp: now,
		},
		TPQ{
			Qty:       3,
			Px:        13,
			Timestamp: now,
		},
		TPQ{
			Qty:       7,
			Px:        14,
			Timestamp: now,
		},
		TPQ{
			Qty:       10,
			Px:        13,
			Timestamp: fivemin,
		},
		TPQ{
			Qty:       1,
			Px:        18,
			Timestamp: fivemin,
		},
		TPQ{
			Qty:       3,
			Px:        10,
			Timestamp: fivemin,
		},
		TPQ{
			Qty:       10,
			Px:        15,
			Timestamp: fivemin,
		},
		TPQ{
			Qty:       13,
			Px:        14,
			Timestamp: fifteen,
		},
	}
	smafivemin := 14.5
	smatenmin := 15.0
	io := []struct {
		time  time.Time
		ohlcv OHLCV
		sma   *float64
	}{
		{
			time: now,
			ohlcv: OHLCV{
				O: 14,
				H: 15,
				L: 13,
				C: 14,
				V: 17,
			},
			sma: nil,
		},
		{
			time: fivemin,
			ohlcv: OHLCV{
				O: 13,
				H: 18,
				L: 10,
				C: 15,
				V: 24,
			},
			sma: &smafivemin,
		},
		{
			time: ten,
			ohlcv: OHLCV{
				O: 15,
				H: 15,
				L: 15,
				C: 15,
				V: 0,
			},
			sma: &smatenmin,
		},
	}

	for _, tpq := range tpqs {
		if err := s.AddExec(tpq); err != nil {
			t.Fatal(errors.Wrapf(err, "error adding exec: %+v", tpq))
		}
	}

	for i, o := range io {
		v := s.GetValueForInterval(o.time)
		if v == nil {
			t.Fatal("expected v to be non nil but got nil")
		}
		h := v.OHLCV
		if h.O != o.ohlcv.O {
			t.Fatalf("expected new open to be 14 but got %+v", o.ohlcv.O)
		} else if h.H != o.ohlcv.H {
			t.Fatalf("expected new high to be 15 but got %+v", o.ohlcv.H)
		} else if h.L != o.ohlcv.L {
			t.Fatalf("expected new high to be 13 but got %+v", o.ohlcv.H)
		} else if h.C != o.ohlcv.C {
			t.Fatalf("expected close to be 14 but got %+v", o.ohlcv.C)
		} else if h.V != o.ohlcv.V {
			t.Fatalf("expected vol to be 17 but got %+v", o.ohlcv.V)
		}
		sma := v.Indicators["sma"]
		if sma == nil && o.sma == nil {
			// ok
			continue
		}
		if sma == nil || o.sma == nil {
			t.Fatalf("expected to be %+v but got %+v at idx %d", o.sma, sma, i)
		}
		if *sma != *o.sma {
			t.Fatalf("expected value to be %+v but got %+v at idx %d", *o.sma, *sma, i)
		}
	}
}
