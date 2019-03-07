package pine

import (
	"fmt"
	"math"
	"time"

	"github.com/pkg/errors"
)

func (s *series) createNewOHLCV(v TPQ, start time.Time) error {
	// create first one
	ohlcv := NewOHLCVWithSamePx(v.Px, v.Qty, start)
	s.insertInterval(ohlcv)
	if err := s.updateIndicators(ohlcv); err != nil {
		return errors.Wrap(err, "error updating indicator")
	}
	return nil
}

func (s *series) updateLastOHLCV(v TPQ) error {
	itvl := s.lastOHLC
	itvl.C = v.Px
	itvl.V += v.Qty
	if v.Px > itvl.H {
		itvl.H = v.Px
	} else if v.Px < itvl.L {
		itvl.L = v.Px
	}
	if err := s.updateIndicators(*itvl); err != nil {
		return errors.Wrap(err, "error updating indicator")
	}
	return nil
}

func (s *series) updateAndFillGaps(v TPQ, start time.Time) error {
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
	return nil
}

func (s *series) AddExec(v TPQ) error {
	start := s.getLastIntervalFromTime(v.Timestamp)
	if s.lastOHLC == nil {
		if err := s.createNewOHLCV(v, start); err != nil {
			return errors.Wrap(err, "error creating new ohlcv")
		}
	} else if s.lastOHLC.S.Equal(start) {
		if err := s.updateLastOHLCV(v); err != nil {
			return errors.Wrap(err, "error creating new ohlcv")
		}
	} else if start.Sub(s.lastOHLC.S).Seconds() > 0 {
		// calculate how many intervals are missing
		if err := s.updateAndFillGaps(v, start); err != nil {
			return errors.Wrap(err, "error updating and filling gaps")
		}
	}
	s.lastExec = v
	return nil
}
