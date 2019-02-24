package pine

import "time"

type constant struct {
	val float64
}

func NewConstant(v float64) Indicator {
	return &constant{
		val: v,
	}
}

func (i *constant) GetValueForInterval(t time.Time) *Interval {
	return &Interval{
		StartTime: t,
		Value:     i.val,
	}
}

func (i *constant) Update(v OHLCV) error {
	return nil
}

func (i *constant) ApplyOpts(opts SeriesOpts) error {
	return nil
}
