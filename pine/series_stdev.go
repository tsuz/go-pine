package pine

import (
	"fmt"
)

// Stdev generates a ValueSeries of one standard deviation
//
// Simplified formula is
//   - s = sqrt(1 / (N-1) * sum(xi - x)^2)
//
// Parameters
//   - p - ValueSeries: source data
//   - l - int64: lookback periods [1, ∞)
//
// TradingView's PineScript has an option to use an unbiased estimator, however; this function currently supports biased estimator.
// Any effort to add a bias correction factor is welcome.
func Stdev(p ValueSeries, l int64) ValueSeries {
	key := fmt.Sprintf("stdev:%s:%d", p.ID(), l)
	stdev := getCache(key)
	if stdev == nil {
		stdev = NewValueSeries()
	}

	// current available value
	stop := p.GetCurrent()
	if stop == nil {
		return stdev
	}

	if stdev.Get(stop.t) != nil {
		return stdev
	}

	vari := Variance(p, l)

	stdev = Pow(vari, 0.5)

	setCache(key, stdev)

	stdev.SetCurrent(stop.t)

	return stdev
}
