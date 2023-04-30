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
	var start, last uint64
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

		if i == 100 {
			start = getMalloc()
		}

		// get last one
		if i == 3000 {
			last = getMalloc()
			break
		}
	}

	// error if allocated more than 15MB. This may not catch smaller increments of memory leak
	if last > start && last-start > 15000 {
		t.Errorf("Memory Leak. Memory allocation increased by %d, start: %d, end: %d", last-start, start, last)
	}
}

func getMalloc() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return m.Alloc / 1024
}
