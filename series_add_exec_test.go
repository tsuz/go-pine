package pine_test

import (
	pine "go-pine"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestSeriesAddExec(t *testing.T) {
	opts := pine.SeriesOpts{
		Interval: 300,
		Max:      100,
	}
	_, err := pine.NewSeries(nil, opts)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	fivemin := now.Add(5 * time.Minute)
	data := []pine.OHLCV{
		pine.OHLCV{
			O: 14,
			H: 15,
			L: 13,
			C: 14,
			V: 131,
			S: now,
		},
		pine.OHLCV{
			O: 13,
			H: 18,
			L: 10,
			C: 15,
			V: 12,
			S: fivemin,
		},
	}
	s, err := pine.NewSeries(data, opts)
	if err != nil {
		t.Fatal(err)
	}

	// This should update high, close, and volume
	tpqhigh := pine.TPQ{
		Timestamp: fivemin,
		Px:        20,
		Qty:       1,
	}
	if err := s.AddExec(tpqhigh); err != nil {
		t.Fatal(errors.Wrapf(err, "error adding exec: %+v", tpqhigh))
	}
	v := s.GetValueForInterval(fivemin)
	if v == nil {
		t.Fatal("expected v to be non nil but got nil")
	}
	h := v.OHLCV
	if h.O != 13 {
		t.Fatalf("expected new open to be 13 but got %+v", h.O)
	} else if h.H != 20 {
		t.Fatalf("expected new high to be 20 but got %+v", h.H)
	} else if h.V != 1+12 {
		t.Fatalf("expected vol to be 13 but got %+v", h.V)
	} else if h.C != 20 {
		t.Fatalf("expected close to be 20 but got %+v", h.C)
	}

	// This should update low, close, and volume
	tpqlow := pine.TPQ{
		Timestamp: fivemin,
		Px:        3,
		Qty:       4,
	}
	if err := s.AddExec(tpqlow); err != nil {
		t.Fatal(errors.Wrapf(err, "error adding exec: %+v", tpqlow))
	}
	v = s.GetValueForInterval(fivemin)
	if v == nil {
		t.Fatal("expected v to be non nil but got nil")
	}
	l := v.OHLCV
	if l.O != 13 {
		t.Fatalf("expected new open to be 13 but got %+v", h.O)
	} else if l.H != 20 {
		t.Fatalf("expected new high to be 20 but got %+v", h.H)
	} else if l.V != 1+12+4 {
		t.Fatalf("expected vol to be 13 but got %+v", h.V)
	} else if l.C != 3 {
		t.Fatalf("expected close to be 3 but got %+v", h.C)
	} else if l.L != 3 {
		t.Fatalf("expected close to be 3 but got %+v", h.L)
	}

	// This should create new interval
	tenmin := fivemin.Add(5 * time.Minute)
	tpqnew := pine.TPQ{
		Timestamp: tenmin,
		Px:        10,
		Qty:       9,
	}
	if err := s.AddExec(tpqnew); err != nil {
		t.Fatal(errors.Wrapf(err, "error adding exec: %+v", tpqnew))
	}
	v = s.GetValueForInterval(tenmin)
	if v == nil {
		t.Fatal("expected v to be non nil but got nil")
	}
	n := v.OHLCV
	if n.S.Sub(l.S).Seconds() != 300 {
		t.Fatalf("expected starting interval to have 300 second diff but got %+v", n.S.Sub(l.S).Seconds())
	} else if n.O != 10 {
		t.Fatalf("expected new open to be 10 but got %+v", n.O)
	} else if n.H != 10 {
		t.Fatalf("expected new high to be 10 but got %+v", n.H)
	} else if n.V != 9 {
		t.Fatalf("expected vol to be 9 but got %+v", n.V)
	} else if n.C != 10 {
		t.Fatalf("expected close to be 10 but got %+v", n.C)
	} else if n.L != 10 {
		t.Fatalf("expected close to be 10 but got %+v", n.L)
	}

	// This should create 2 intervals since this spans two intervals
	// refer to ExecInst
	twemin := tenmin.Add(10 * time.Minute)
	tpqtwe := pine.TPQ{
		Timestamp: twemin,
		Px:        14,
		Qty:       3,
	}
	if err := s.AddExec(tpqtwe); err != nil {
		t.Fatal(errors.Wrapf(err, "error adding exec: %+v", tpqtwe))
	}
	v = s.GetValueForInterval(twemin.Add(-5 * time.Minute))
	if v == nil {
		t.Fatal("expected v to be non nil but got nil")
	}
	n = v.OHLCV
	if n.S.Sub(l.S).Seconds() != 600 {
		t.Fatalf("expected starting interval to have 600 second diff but got %+v", n.S.Sub(l.S).Seconds())
	} else if n.O != 10 {
		t.Fatalf("expected new open to be 10 but got %+v", n.O)
	} else if n.H != 10 {
		t.Fatalf("expected new high to be 10 but got %+v", n.H)
	} else if n.V != 0 {
		t.Fatalf("expected vol to be 0 but got %+v", n.V)
	} else if n.C != 10 {
		t.Fatalf("expected close to be 10 but got %+v", n.C)
	} else if n.L != 10 {
		t.Fatalf("expected close to be 10 but got %+v", n.L)
	}

	v = s.GetValueForInterval(twemin)
	if v == nil {
		t.Fatal("expected v to be non nil but got nil")
	}
	n = v.OHLCV
	if n.S.Sub(l.S).Seconds() != 900 {
		t.Fatalf("expected starting interval to have 900 second diff but got %+v", n.S.Sub(l.S).Seconds())
	} else if n.O != 14 {
		t.Fatalf("expected new open to be 14 but got %+v", n.O)
	} else if n.H != 14 {
		t.Fatalf("expected new high to be 14 but got %+v", n.H)
	} else if n.V != 3 {
		t.Fatalf("expected vol to be 0 but got %+v", n.V)
	} else if n.C != 14 {
		t.Fatalf("expected close to be 14 but got %+v", n.C)
	} else if n.L != 14 {
		t.Fatalf("expected close to be 14 but got %+v", n.L)
	}
}
