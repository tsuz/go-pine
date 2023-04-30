package pine

import (
	"fmt"
)

// MACD generates a ValueSeries of MACD (moving average convergence/divergence).
// It is supposed to reveal changes in the strength, direction, momentum, and duration of a trend in a stock's price.
//
// The formula for MACD is
//
//   - MACD Line: (12-day EMA - 26-day EMA)
//   - Signal Line: 9-day EMA of MACD Line
//   - MACD Histogram: MACD Line - Signal Line
//
// The arguments are:
//
//   - source: ValueSeries - source of data
//   - fastlen: int - fast len of MACD series
//   - slowlen: int - slow len of MACD series
//   - siglen: int - signal length of MACD series
//
// The return values are:
//   - macdLine: ValueSeries - MACD Line
//   - signalLine: ValueSeries - Signal Line
//   - histLine: ValueSeries - MACD Histogram
//   - err: error
func MACD(src ValueSeries, fastlen, slowlen, siglen int64) (ValueSeries, ValueSeries, ValueSeries) {
	macdlineKey := fmt.Sprintf("macdline:%s:%d:%d:%d", src.ID(), fastlen, slowlen, siglen)
	macdline := getCache(macdlineKey)
	if macdline == nil {
		macdline = NewValueSeries()
	}

	signalLineKey := fmt.Sprintf("macdsignal:%s:%d:%d:%d", src.ID(), fastlen, slowlen, siglen)
	signalLine := getCache(signalLineKey)
	if signalLine == nil {
		signalLine = NewValueSeries()
	}

	macdHistogramKey := fmt.Sprintf("macdhistogram:%s:%d:%d:%d", src.ID(), fastlen, slowlen, siglen)
	macdHistogram := getCache(macdHistogramKey)
	if macdHistogram == nil {
		macdHistogram = NewValueSeries()
	}

	// current available value
	stop := src.GetCurrent()

	if stop == nil {
		return macdline, signalLine, macdHistogram
	}

	// latest value exists
	if macdHistogram.Get(stop.t) != nil {
		return macdline, signalLine, macdHistogram
	}

	fast := EMA(src, fastlen)

	slow := EMA(src, slowlen)

	macdline = Sub(fast, slow)
	macdline.SetCurrent(stop.t)

	signalLine = EMA(macdline, siglen)
	signalLine.SetCurrent(stop.t)

	macdHistogram = Sub(macdline, signalLine)
	macdHistogram.SetCurrent(stop.t)

	setCache(macdlineKey, macdline)
	setCache(signalLineKey, signalLine)
	setCache(macdHistogramKey, macdHistogram)

	return macdline, signalLine, macdHistogram
}
