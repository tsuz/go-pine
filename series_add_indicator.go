package pine

import "github.com/pkg/errors"

func (s *series) AddIndicator(name string, i Indicator) error {
	// enforce series constraint
	i.ApplyOpts(s.opts)
	// update with current values downstream
	for _, v := range s.values {
		if err := i.Update(v); err != nil {
			return errors.Wrap(err, "error updating indicator")
		}
	}
	s.items[name] = i
	return nil
}
