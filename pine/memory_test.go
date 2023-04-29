package pine

import (
	_ "net/http/pprof"
	"runtime"
	"testing"

	"time"

	"github.com/jinzhu/now"
	"github.com/pkg/errors"
)

type ds struct{}

func (d *ds) Populate(t time.Time) ([]OHLCV, error) {
	itvl := time.Minute
	next := t.Add(itvl)
	nextend := next.Add(itvl * 100)

	data := OHLCVTestData(nextend, 100, 60*1000)
	return data, nil
}

type testIndicator = func(o OHLCVSeries) error

func testMemoryLeak(t *testing.T, fn testIndicator) {

	tn := time.Now()
	fromTime := now.With(tn).BeginningOfMonth()

	d := &ds{}
	ohlcv, err := d.Populate(fromTime)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error populating"))
	}

	s, err := NewDynamicOHLCVSeries(ohlcv, d)
	if err != nil {
		panic(errors.Wrap(err, "error creating ohlcvseries"))
	}

	first := true
	v, err := s.Next()
	if err != nil {
		panic(errors.Wrap(err, "error next"))
	}

	i := 0
	var last uint64
	for {
		if v == nil && !first {
			break
		}
		if first {
			first = false
		}

		c := s.Current()
		if c == nil {
			break
		}

		if ierr := fn(s); ierr != nil {
			t.Error(errors.Wrap(err, "error getting"))
		}

		v, err = s.Next()
		if err != nil {
			t.Error(errors.Wrap(err, "error next"))
		}
		i++

		// get last one
		if i == 10000 {
			last = getMalloc()
			break
		}
	}

	// error if it's more than 25% allocated - an arbitrary value
	if last > 15000 {
		t.Errorf("Memory Leak. Memory allocation ended at %d", last)
	}
}

func getMalloc() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return m.Alloc / 1024
}
