package pine

import "time"

type DataSource interface {
	// Populate is called to fetch more data
	// 	t (time.Time) - the last start time of the last existing OHLCV
	Populate(t time.Time) ([]OHLCV, error)
}
