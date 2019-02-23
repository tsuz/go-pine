package pine

import (
	"log"
	"time"
)

type hl2 struct {
	last    OHLCV
	opts    *SeriesOpts
	timemap map[time.Time]float64
}

// NewHL2 returns midpoint of high and low: (h + l) / 2
func NewHL2() Indicator {
	return &hl2{
		timemap: make(map[time.Time]float64),
	}
}

func (i *hl2) GetValueForInterval(t time.Time) *Interval {
	v, ok := i.timemap[t]
	log.Printf("Getting timemap %+v %t", t, ok)
	if !ok {
		return nil
	}
	return &Interval{
		StartTime: t,
		Value:     v,
	}
}

func (i *hl2) Update(v OHLCV) error {
	if i.last.H == v.H && i.last.L == v.L {
		return nil
	}
	i.timemap[v.S] = (v.H + v.L) / 2
	log.Printf("Adding to timemap %+v => %+v", v.S, i.timemap[v.S])
	i.last = v
	return nil
}

func (i *hl2) ApplyOpts(opts SeriesOpts) error {
	i.opts = &opts
	// validate if needed
	return nil
}
