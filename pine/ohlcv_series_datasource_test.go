package pine

import (
	"testing"
	"time"

	"github.com/pkg/errors"
)

type testds struct {
	data2 []OHLCV
}

func NewTestDynamicDS(data2 []OHLCV) DataSource {
	return &testds{data2: data2}
}

func (t *testds) Populate(v time.Time) ([]OHLCV, error) {
	if t.data2[0].S.Sub(v) > 0 {
		return t.data2, nil
	}
	return []OHLCV{}, nil
}

func TestNewOHLCVSeriesFetchDataSource(t *testing.T) {
	start := time.Now()
	data := OHLCVTestData(start, 3, 5*60*1000)
	data2 := OHLCVTestData(start.Add(3*5*time.Minute), 3, 5*60*1000)

	ds := NewTestDynamicDS(data2)
	s, err := NewDynamicOHLCVSeries(data, ds)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < (len(data) + len(data2)); i++ {
		_, err := s.Next()
		if err != nil {
			t.Fatal(errors.Wrap(err, "error fetching next"))
		}
		close := s.GetSeries(OHLCPropClose)
		if i < len(data) {
			if *close.Val() != data[i].C {
				t.Errorf("expected %+v but got %+v", data[i].C, *close.Val())
			}
		} else {
			if *close.Val() != data2[i-len(data)].C {
				t.Errorf("expected %+v but got %+v", data2[i-len(data)].C, *close.Val())
			}
		}
	}
}
