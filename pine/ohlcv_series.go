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

import "math"

// OHLCVSeries represents a series of OHLCV type (i.e. open, high, low, close, volume)
type OHLCVSeries interface {

	// Current returns current ohlcv
	Current() *OHLCV

	// GetSeries returns series of values for a property
	GetSeries(OHLCProp) ValueSeries

	// GoToFirst sets the current value to first and returns that value
	GoToFirst() *OHLCV

	// Next moves the pointer to the next one
	Next() *OHLCV
}

func NewOHLCVSeries(ohlcv []OHLCV) (OHLCVSeries, error) {
	for idx := range ohlcv {
		if idx > 0 {
			ohlcv[idx].prev = &ohlcv[idx-1]
		}
		if idx < len(ohlcv)-1 {
			ohlcv[idx].next = &ohlcv[idx+1]
		}
	}
	s := &ohlcvSeries{
		ohlcv: ohlcv,
	}
	return s, nil
}

type ohlcvSeries struct {
	// current ohlcv
	cur *OHLCV

	// values
	ohlcv []OHLCV
}

func (s *ohlcvSeries) Current() *OHLCV {
	return s.cur
}

func (s *ohlcvSeries) GoToFirst() *OHLCV {
	if len(s.ohlcv) == 0 {
		return nil
	}
	s.cur = &s.ohlcv[0]
	return s.cur
}

func (s *ohlcvSeries) Next() *OHLCV {
	if s.cur == nil {
		// empty list, no movement
		if len(s.ohlcv) == 0 {
			return nil
		}
		s.cur = &s.ohlcv[0]
		return s.cur
	}
	// don't move the cursor just yet
	if s.cur.next == nil {
		return nil
	}
	s.cur = s.cur.next
	return s.cur
}

func (s *ohlcvSeries) GetSeries(p OHLCProp) ValueSeries {
	vs := NewValueSeries()
	for idx, v := range s.ohlcv {
		var propVal float64
		switch p {
		case OHLCPropClose:
			propVal = v.C
		case OHLCPropOpen:
			propVal = v.O
		case OHLCPropHigh:
			propVal = v.H
		case OHLCPropLow:
			propVal = v.L
		case OHLCPropVolume:
			propVal = v.V
		case OHLCPropTR:
			if idx > 0 {
				p := s.ohlcv[idx-1]
				propVal = math.Max(
					math.Abs(v.H-v.L),
					math.Max(
						math.Abs(v.H-p.C),
						math.Abs(v.L-p.C)))
			} else {
				propVal = math.Abs(v.H - v.L)
			}

		default:
			continue
		}
		vs.Set(v.S, propVal)
	}
	if s.cur != nil {
		vs.SetCurrent(s.cur.S)
	}
	return vs
}
