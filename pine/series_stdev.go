package pine

import (
	"fmt"

	"github.com/pkg/errors"
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
func Stdev(p ValueSeries, l int64) (ValueSeries, error) {
	key := fmt.Sprintf("stdev:%s:%d", p.ID(), l)
	stdev := getCache(key)
	if stdev == nil {
		stdev = NewValueSeries()
	}

	// current available value
	stop := p.GetCurrent()
	if stop == nil {
		return stdev, nil
	}

	if stdev.Get(stop.t) != nil {
		return stdev, nil
	}

	vari, err := Variance(p, l)
	if err != nil {
		return nil, errors.Wrap(err, "error getting variance")
	}

	stdev, err = Pow(vari, 0.5)
	if err != nil {
		return stdev, errors.Wrap(err, "error pow(0.5)")
	}
	setCache(key, stdev)

	stdev.SetCurrent(stop.t)

	return stdev, nil
}
