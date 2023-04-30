package pine

import (
	"fmt"
	"math"
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
func Variance(p ValueSeries, l int64) ValueSeries {
	key := fmt.Sprintf("variance:%s:%d", p.ID(), l)
	vari := getCache(key)
	if vari == nil {
		vari = NewValueSeries()
	}

	// current available value
	stop := p.GetCurrent()
	if stop == nil {
		return vari
	}

	if vari.Get(stop.t) != nil {
		return vari
	}

	sma := SMA(p, l)

	meanv := sma.Get(stop.t)
	if meanv == nil {
		return vari
	}
	diff := SubConstNoCache(p, meanv.v)
	sqrt := Pow(diff, 2)
	sum := SumNoCache(sqrt, int(l))
	denom := math.Max(float64(l-1), 1)
	vari = DivConstNoCache(sum, denom)

	vari.SetCurrent(stop.t)

	setCache(key, vari)

	return vari
}
