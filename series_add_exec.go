package pine

import (
	"fmt"
	"log"
	"math"
	"time"
)

func (s *series) AddExec(v TPQ) error {
	log.Printf("Add exec %+v", v)
	start := s.getLastIntervalFromTime(v.Timestamp)
	if s.lastOHLC == nil {
		// create first one
		ohlcv := NewOHLCVWithSamePx(v.Px, v.Qty, start)
		s.insertInterval(ohlcv)
		s.updateIndicators(ohlcv)
	} else if s.lastOHLC.S.Equal(start) {
		// update this interval
		itvl := s.lastOHLC
		itvl.C = v.Px
		itvl.V += v.Qty
		if v.Px > itvl.H {
			itvl.H = v.Px
		} else if v.Px < itvl.L {
			itvl.L = v.Px
		}
		log.Printf("Updating existing to %+v", s.lastOHLC)
		s.updateIndicators(*itvl)
	} else if start.Sub(s.lastOHLC.S).Seconds() > 0 {
		// calculate how many intervals are missing
		switch s.opts.EmptyInst {
		case EmptyInstUseLastClose:
			// figure out how many
			m := s.getMultiplierDiff(v.Timestamp, s.lastOHLC.S)
			usenew := int(math.Max(float64(m-1), float64(0)))
			for i := 0; i < m; i++ {
				var px, qty float64
				if i == usenew {
					px = v.Px
					qty = v.Qty
				} else {
					px = s.lastOHLC.C
					qty = 0
				}
				newt := s.lastOHLC.S.Add(time.Duration(s.opts.Interval) * time.Second)
				ohlcv := NewOHLCVWithSamePx(px, qty, newt)
				s.insertInterval(ohlcv)
				s.updateIndicators(ohlcv)
			}
			log.Printf("New updated %+v", s.lastOHLC)
		case EmptyInstUseZeros:
			// figure out how many
			m := s.getMultiplierDiff(v.Timestamp, s.lastOHLC.S)
			usenew := int(math.Max(float64(m-1), float64(0)))
			for i := 0; i < m; i++ {
				var px, qty float64
				if i == usenew {
					px = v.Px
					qty = v.Qty
				} else {
					px = 0
					qty = 0
				}
				newt := s.lastOHLC.S.Add(time.Duration(s.opts.Interval) * time.Second)
				ohlcv := NewOHLCVWithSamePx(px, qty, newt)
				s.insertInterval(ohlcv)
				s.updateIndicators(ohlcv)
			}
		case EmptyInstIgnore:
			v := NewOHLCVWithSamePx(v.Px, v.Qty, start)
			s.insertInterval(v)
			s.updateIndicators(v)
		default:
			return fmt.Errorf("Unsupported interval: %+v", s.opts.EmptyInst)
		}
	}
	s.lastExec = v
	return nil
}
