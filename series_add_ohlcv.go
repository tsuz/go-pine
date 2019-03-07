package pine

import (
	"fmt"
	"math"
	"time"
)

func (s *series) AddOHLCV(v OHLCV) error {
	start := s.getLastIntervalFromTime(v.S)
	v.S = start
	if s.lastOHLC == nil {
		// create first one
		s.insertInterval(v)
	} else if s.lastOHLC.S.Equal(start) {
		// update this interval
		itvl := s.lastOHLC
		itvl.O = v.O
		itvl.H = v.H
		itvl.L = v.L
		itvl.C = v.C
		itvl.V = v.V
	} else if start.Sub(s.lastOHLC.S).Seconds() > 0 {
		// calculate how many intervals are missing
		switch s.opts.EmptyInst {
		case EmptyInstUseLastClose:
			// figure out how many
			m := s.getMultiplierDiff(v.S, s.lastOHLC.S)
			usenew := int(math.Max(float64(m-1), float64(0)))
			for i := 0; i < m; i++ {
				var px, qty float64
				var ohlcv OHLCV
				if i == usenew {
					ohlcv = v
				} else {
					px = s.lastOHLC.C
					qty = 0
					newt := s.lastOHLC.S.Add(time.Duration(s.opts.Interval) * time.Second)
					ohlcv = NewOHLCVWithSamePx(px, qty, newt)
				}
				s.insertInterval(ohlcv)
			}
		case EmptyInstUseZeros:
			// figure out how many
			m := s.getMultiplierDiff(v.S, s.lastOHLC.S)
			usenew := int(math.Max(float64(m-1), float64(0)))
			for i := 0; i < m; i++ {
				var px, qty float64
				var ohlcv OHLCV
				if i == usenew {
					ohlcv = v
				} else {
					px = 0
					qty = 0
					newt := s.lastOHLC.S.Add(time.Duration(s.opts.Interval) * time.Second)
					ohlcv = NewOHLCVWithSamePx(px, qty, newt)
				}
				s.insertInterval(ohlcv)
			}
		case EmptyInstIgnore:
			s.insertInterval(v)
		default:
			return fmt.Errorf("Unsupported interval: %+v", s.opts.EmptyInst)
		}
	}
	s.lastOHLC = &v
	return nil
}
