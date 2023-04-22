package pine

import (
	"log"
	"time"
)

func ExampleNewOHLCVSeries() {
	start := time.Now()
	// start = start time of the first OHLCV bar
	// 3 = 3 bars
	// 5*60*1000 = 5 minutes in milliseconds
	data := OHLCVTestData(start, 3, 5*60*1000)
	s, _ := NewOHLCVSeries(data)

	for {
		v, _ := s.Next()
		if v == nil {
			break
		}
		log.Printf("Close: %+v", v.C)
	}
}
