package pine

import "time"

type Indicator interface {
	ApplyOpts(opts SeriesOpts) error
	GetValueForInterval(t time.Time) *Interval
	Update(v OHLCV) error
}
