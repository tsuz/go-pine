package pine

import (
	"github.com/pkg/errors"
)

func (s *series) AddIndicator(name string, i Indicator) error {
	// enforce series constraint
	if err := i.ApplyOpts(s.opts); err != nil {
		return errors.Wrap(err, "error applying opts")
	}
	// update with current values downstream
	for _, v := range s.values {
		if err := i.Update(v); err != nil {
			return errors.Wrap(err, "error updating indicator")
		}
	}
	s.items[name] = i
	return nil
}
