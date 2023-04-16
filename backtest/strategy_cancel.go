package backtest

func (s *strategy) Cancel(ordID string) error {
	for _, v := range s.ordEntry {
		if v.OrdID == ordID {
			delete(s.ordEntry, ordID)
		}
	}

	return nil
}

func (s *strategy) CancelAll() error {
	for _, v := range s.ordEntry {
		delete(s.ordEntry, v.OrdID)
	}

	return nil
}
