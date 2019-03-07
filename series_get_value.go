package pine

import (
	"time"
)

func (s *series) GetValueForInterval(t time.Time) *Interval {
	if s.lastOHLC == nil {
		return nil
	}
	if !s.lastOHLC.S.Equal(t) {
		// if time is within interval, adjust it
		t = s.getLastIntervalFromTime(t)
	}
	t = t.UTC()
	inds := make(map[string]*float64)
	for k, v := range s.items {
		val := v.GetValueForInterval(t)
		if val != nil {
			inds[k] = &val.Value
		}
	}
	v, ok := s.timemap[t]
	if !ok {
		return nil
	}
	return &Interval{
		StartTime:  v.S,
		OHLCV:      v,
		Indicators: inds,
	}
}
