package backtest

import (
	"github.com/pkg/errors"
	"github.com/tsuz/go-pine/pine"
)

// Runbacktest starts a backtest
func RunBacktest(series pine.OHLCVSeries, b BackTestable) (*BacktestResult, error) {
	strategy := NewStrategy()
	states := map[string]interface{}{}

	series.GoToFirst()

	for {
		if err := b.OnNextOHLCV(strategy, series, states); err != nil {
			return nil, errors.Wrap(err, "error calling OnNextOHLCV")
		}
		next := series.Next()
		if next == nil {
			break
		}
		if err := strategy.Execute(*next); err != nil {
			return nil, errors.Wrapf(err, "error executing next: %+v", *next)
		}
	}
	result := strategy.Result()
	return &result, nil
}
