package backtest

func (s *strategy) Cancel(ordID string) error {
	for _, v := range s.ordEntry {
		if v.OrdID == ordID {
			s.deleteEntryOrder(ordID)
		}
	}

	return nil
}

func (s *strategy) CancelAll() error {
	for _, v := range s.ordEntry {
		s.deleteEntryOrder(v.OrdID)
	}

	return nil
}
