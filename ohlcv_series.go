package pine

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
	for _, v := range s.ohlcv {
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
