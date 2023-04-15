package backtest

import (
	"github.com/tsuz/go-pine/pine"
)

type Strategy interface {
	Execute(pine.OHLCV) error
	Entry(string, EntryOpts) error
	Exit(string) error
	Result() BacktestResult
}

type strategy struct {
	openPos  map[string]Position
	ordEntry map[string]EntryOpts
	ordExit  map[string]bool
	res      *BacktestResult
}

func NewStrategy() Strategy {
	s := strategy{
		res:      &BacktestResult{},
		openPos:  make(map[string]Position),
		ordEntry: make(map[string]EntryOpts),
		ordExit:  make(map[string]bool),
	}
	return &s
}

func (s *strategy) deleteEntryOrder(ordID string) {
	delete(s.ordEntry, ordID)
}

func (s *strategy) deleteOpenPos(ordID string) {
	delete(s.openPos, ordID)
}

func (s *strategy) deleteEntryExit(ordID string) {
	delete(s.ordExit, ordID)
}

func (s *strategy) findPos(ordID string) (Position, bool) {
	v, ok := s.openPos[ordID]
	return v, ok
}

func (s *strategy) findOrdEntry(ordID string) bool {
	_, ok := s.ordEntry[ordID]
	return ok
}

func (s *strategy) setEntryOrder(ordID string, v EntryOpts) {
	v.OrdID = ordID
	s.ordEntry[ordID] = v
}

func (s *strategy) setOpenPos(ordID string, v Position) {
	s.openPos[ordID] = v
}

func (s *strategy) setEntryExit(ordID string) {
	s.ordExit[ordID] = true
}

func (s *strategy) Entry(ordID string, opts EntryOpts) error {
	s.setEntryOrder(ordID, opts)
	return nil
}

func (s *strategy) Execute(ohlcv pine.OHLCV) error {
	// convert open entry orders into open positions
	for _, v := range s.ordEntry {
		_, found := s.findPos(v.OrdID)
		if found {
			continue
		}
		pos := Position{
			EntryPx:   ohlcv.O,
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

func (s *strategy) completePosition(p Position) {
	s.res.ClosedOrd = append(s.res.ClosedOrd, p)
	s.res.TotalClosedTrades++

	prof := p.Profit()
	if prof > 0 {
		s.res.ProfitableTrades++
		s.res.PercentProfitable = float64(s.res.ProfitableTrades) / float64(s.res.TotalClosedTrades)
	}
}

func (s *strategy) Exit(ordID string) error {
	s.setEntryExit(ordID)
	return nil
}

func (s *strategy) Result() BacktestResult {
	s.res.CalculateNetProfit()
	return *s.res
}
