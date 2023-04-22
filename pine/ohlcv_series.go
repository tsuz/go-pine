/*
Pine represents core indicators written in the PineScript manual V5.

While this API looks similar to PineScript, keep in mind these design choices while integrating.

 1. Every indicator is derived from OHLCVSeries. OHLCVSeries contains information about the candle (i.e. OHLCV, true range, mid point etc) and indicators can use these data as its source.
 2. OHLCVSeries does not sort the order of the OHLCV values. The developer is responsible for providing the correct order.
 3. OHLCVSeries does not make assumptions about the time interval. The developer is responsible for specifying OHLCV's time as well as performing data manipulations before hand such as filling in empty intervals. One advantage of this is that each interval can be as small as an execution tick with a varying interval between them.
 4. OHLCV and indicators are in a series, meaning it will attempt to generate all values up to the specified high watermark. It is specified using either SetCurrent(time.Time) or calling Next() in the OHLCVSeries.
 5. OHLCVSeries differentiates OHLCV items by its start time (i.e. time.Time). Ensure all OHLCV have unique time.
*/
package pine

// OHLCVSeries represents a series of OHLCV type (i.e. open, high, low, close, volume)
type OHLCVSeries interface {
	OHLCVBaseSeries
}

// NewDynamicOHLCVSeries generates a dynamic OHLCV series
func NewDynamicOHLCVSeries(ohlcv []OHLCV, ds DataSource) (OHLCVSeries, error) {
	s := NewOHLCVBaseSeries()

	for _, v := range ohlcv {
		s.Push(v)
	}

	s.RegisterDataSource(ds)

	return s, nil
}

func NewOHLCVSeries(ohlcv []OHLCV) (OHLCVSeries, error) {
	s := NewOHLCVBaseSeries()

	for _, v := range ohlcv {
		s.Push(v)
	}

	return s, nil
}
