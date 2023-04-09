package pine

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
)

// Stdev generates a ValueSeries of one standard deviation
//
// Simplified formula is s = sqrt(1 / (N-1) * sum(xi - x)^2)
//
// Parameters
// p - ValueSeries: source data
// l - lookback: lookback periods [1, ∞)
//
// TradingView's PineScript has an option to use an unbiased estimator, however; this function currently supports biased estimator.
// Any effort to add a bias correction factor is welcome.
//
// t=time.Time        | 1   |  2  | 3   | 4    | 5      |
// p=ValueSeries      | 13  | 15  | 11  | 19   | 21     |
// sma(p, 3)	      | nil | nil | 13  | 15   | 17     |
// p - sma(p, 3)(t=1) | nil | nil | nil | nil  | nil    |
// p - sma(p, 3)(t=2) | nil | nil | nil | nil  | nil    |
// p - sma(p, 3)(t=3) | 0   |  2  | -2  | 6    | 5      |
// p - sma(p, 3)(t=4) | -2  |  0  | -4  | 4    | 6      |
// p - sma(p, 3)(t=5) | -4  | -2  | -6  | 2    | 4      |
// Stdev(p, 3)		  | nil | nil |  2  | 4    | 5.2915 |
func Stdev(p ValueSeries, l int64) (ValueSeries, error) {
	key := fmt.Sprintf("stdev:%s:%d", p.ID(), l)
	stdev := getCache(key)
	if stdev == nil {
		stdev = NewValueSeries()
	}

	sma, err := SMA(p, l)
	if err != nil {
		return stdev, errors.Wrap(err, "error getting sma")
	}

	// current available value
	stop := p.GetCurrent()
	if stop == nil {
		return stdev, nil
	}

	meanv := sma.Get(stop.t)
	if meanv == nil {
		return stdev, nil
	}

	diff := p.SubConst(meanv.v)
	sqrt, err := Pow(diff, 2)
	if err != nil {
		return stdev, errors.Wrap(err, "error pow(2)")
	}
	sum, err := Sum(sqrt, int(l))
	if err != nil {
		return stdev, errors.Wrap(err, "error sum")
	}
	denom := math.Max(float64(l-1), 1)
	vari := sum.DivConst(denom)
	stdev, err = Pow(vari, 0.5)
	if err != nil {
		return stdev, errors.Wrap(err, "error pow(0.5)")
	}
	setCache(key, stdev)

	stdev.SetCurrent(stop.t)

	return stdev, nil
}
