package pine

// ATR generates a ValueSeries of average true range
//
// Function ATR returns the RMA of true range ValueSeries.
// True range is already generated by OHLCVSeries.
// True range is
//   - max(high - low, abs(high - close[1]), abs(low - close[1])).
//
// The arguments are:
//   - tr: ValueSeries - true range value
//   - length: int - lookback length to generate ATR. 1 is same as the current value.
//
// The return values are:
//   - ATR: ValueSeries - ATR
//   - err: error
func ATR(tr ValueSeries, l int64) (ValueSeries, error) {
	return RMA(tr, l)
}
