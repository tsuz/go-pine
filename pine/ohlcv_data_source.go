package pine

import "time"

type DataSource interface {
	// Populate is called to fetch more data
	// 	t (time.Time) - the last start time of the last existing OHLCV
	//
	// Populate is triggered when OHLCVSeries has reached the end and there are no next items
	// Returning an empty OHLCV list if nothing else to add
	Populate(t time.Time) ([]OHLCV, error)
}
