package pine

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
)

// Variance generates a ValueSeries of variance.
// Variance is the expectation of the squared deviation of a series from its mean (ta.sma, and it informally measures how far a set of numbers are spread out from their mean.
//
// Simplified formula is
//   - v = (1 / (N-1) * sum(xi - x)^2)
//
// Parameters
//   - p - ValueSeries: source data
//   - l - lookback: lookback periods [1, ∞)
//
// TradingView's PineScript has an option to use an unbiased estimator, however; this function currently supports biased estimator.
// Any effort to add a bias correction factor is welcome.
func Variance(p ValueSeries, l int64) (ValueSeries, error) {
	key := fmt.Sprintf("variance:%s:%d", p.ID(), l)
	vari := getCache(key)
	if vari == nil {
		vari = NewValueSeries()
	}

	// current available value
	stop := p.GetCurrent()
	if stop == nil {
		return vari, nil
	}

	if vari.Get(stop.t) != nil {
		return vari, nil
	}

	sma, err := SMA(p, l)
	if err != nil {
		return vari, errors.Wrap(err, "error getting sma")
	}

	meanv := sma.Get(stop.t)
	if meanv == nil {
		return vari, nil
	}

	diff := p.SubConst(meanv.v)
	sqrt, err := Pow(diff, 2)
	if err != nil {
		return vari, errors.Wrap(err, "error pow(2)")
	}
	sum, err := Sum(sqrt, int(l))
	if err != nil {
		return vari, errors.Wrap(err, "error sum")
	}
	denom := math.Max(float64(l-1), 1)
	vari = sum.DivConst(denom)

	setCache(key, vari)

	vari.SetCurrent(stop.t)

	return vari, nil
}
