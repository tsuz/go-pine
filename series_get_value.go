package pine

import (
	"log"
	"time"
)

func (s *series) GetValueForInterval(t time.Time) *Interval {
	log.Printf("GetValueForInterval %+v %+v", t, s.lastOHLC.S)
	if s.lastOHLC == nil {
		return nil
	}
	if !s.lastOHLC.S.Equal(t) {
		// if time is within interval, adjust it
		t = s.getLastIntervalFromTime(t)
		log.Printf("getLastIntervalFromTime %+v", t)
	}
	inds := make(map[string]*float64)
	for k, v := range s.items {
		val := v.GetValueForInterval(t)
		if val != nil {
			inds[k] = &val.Value
		}
	}
	v, ok := s.timemap[t]
	if !ok {
		log.Printf("Not ok for t %+v", t)
		return nil
	} else {
		log.Printf("ok")
	}
	return &Interval{
		StartTime:  v.S,
		OHLCV:      v,
		Indicators: inds,
	}
}
