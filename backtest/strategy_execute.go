package backtest

import "github.com/tsuz/go-pine/pine"

func (s *strategy) Execute(ohlcv pine.OHLCV) error {
	// convert open entry orders into open positions
	for _, v := range s.ordEntry {
		_, found := s.findPos(v.OrdID)
		if found {
			continue
		}

		entryPx := ohlcv.O

		// if limit order, see if it gets filled
		if v.Limit != nil {
			if v.Side == Long && *v.Limit < ohlcv.L {
				// long order not filled
				continue
			}
			if v.Side == Short && *v.Limit > ohlcv.H {
				// short order not filled
				continue
			}
			entryPx = *v.Limit
		}

		pos := Position{
			EntryPx:   entryPx,
			EntryTime: ohlcv.S,
			EntrySide: v.Side,
			OrdID:     v.OrdID,
		}
		s.setOpenPos(v.OrdID, pos)
	}
	s.ordEntry = make(map[string]EntryOpts)

	// convert positions into exit orders
	for id := range s.ordExit {
		p, found := s.findPos(id)
		if found {
			p.ExitPx = ohlcv.O
			p.ExitTime = ohlcv.S
			s.completePosition(p)
			s.deleteOpenPos(id)
		}
	}

	return nil
}
