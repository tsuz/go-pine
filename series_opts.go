package pine

import "github.com/pkg/errors"

// SeriesOpts is options required for creating Series
type SeriesOpts struct {
	// interval in seconds
	Interval int
	// max number of OHLC bars to keep
	Max int
	// instruction when there are no execs during interval
	EmptyInst EmptyInst
}

// Validate validates series opts and returns error if not good
func (s *SeriesOpts) Validate() error {
	if s.Interval <= 0 {
		return errors.New("`Interval` must be positive")
	} else if s.Max <= 0 {
		return errors.New("`Max` must be positive")
	}
	return nil
}

// EmptyInst is instruction when no values are set for the interval
type EmptyInst int

const (
	// EmptyInstUseLastClose uses the last close value for open, high, low, close but zero for volume
	EmptyInstUseLastClose EmptyInst = iota
	// EmptyInstIgnore ignores intervals if empty
	EmptyInstIgnore
	// EmptyInstUseZeros uses zeros for open, high, low, close, and volume
	EmptyInstUseZeros
)
