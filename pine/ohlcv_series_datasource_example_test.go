package pine

import (
	"log"
	"time"
)

type mytestds struct {
	data2 []OHLCV
}

func MyNewTestDynamicDS(data2 []OHLCV) DataSource {
	return &testds{data2: data2}
}

// Populate is triggered when OHLCVSeries has reached the end and there are no next items
// Returning an empty OHLCV list if nothing else to add
func (t *mytestds) Populate(v time.Time) ([]OHLCV, error) {

	// Fetch data from API
	if t.data2[0].S.Sub(v) > 0 {
		return t.data2, nil
	}

	return []OHLCV{}, nil
}

func ExampleNewDynamicOHLCVSeries() {
	start := time.Now()
	data := OHLCVTestData(start, 3, 5*60*1000)
	data2 := OHLCVTestData(start.Add(3*5*time.Minute), 3, 5*60*1000)

	ds := MyNewTestDynamicDS(data2)
	s, _ := NewDynamicOHLCVSeries(data, ds)

	for {
		v, _ := s.Next()
		if v == nil {
			break
		}
		log.Printf("Close is %+v", v.C)
	}
}
