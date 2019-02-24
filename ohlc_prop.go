package pine

import (
	"time"
)

// OHLCProp is a property of OHLC
type OHLCProp int

const (
	// OHLCPropClose is the close value of OHLC
	OHLCPropClose OHLCProp = iota
	// OHLCPropOpen is the open value of OHLC
	OHLCPropOpen
	// OHLCPropHigh is the high value of OHLC
	OHLCPropHigh
	// OHLCPropLow is the low value of OHLC
	OHLCPropLow
	// OHLCPropVolume is the volume value of OHLC
	OHLCPropVolume
)

type ohlcprop struct {
	prop    OHLCProp
	last    OHLCV
	timemap map[time.Time]float64
}

// NewOHLCProp returns property of OHLC in OHLCProps
func NewOHLCProp(p OHLCProp) Indicator {
	return &ohlcprop{
		prop:    p,
		timemap: make(map[time.Time]float64),
	}
}

func (i *ohlcprop) GetValueForInterval(t time.Time) *Interval {
	v, ok := i.timemap[t]
	if !ok {
		return nil
	}
	return &Interval{
		StartTime: t,
		Value:     v,
	}
}

func (i *ohlcprop) Update(v OHLCV) error {
	timesame := i.last.S == v.S
	var valssame bool
	var val float64
	switch i.prop {
	case OHLCPropClose:
		valssame = i.last.C == v.C
		val = v.C
	case OHLCPropHigh:
		valssame = i.last.H == v.H
		val = v.H
	case OHLCPropLow:
		valssame = i.last.L == v.L
		val = v.L
	case OHLCPropOpen:
		valssame = i.last.O == v.O
		val = v.O
	case OHLCPropVolume:
		valssame = i.last.V == v.V
		val = v.V
	}
	if valssame && timesame {
		return nil
	}
	i.timemap[v.S] = val
	i.last = v
	return nil
}

func (i *ohlcprop) ApplyOpts(opts SeriesOpts) error {
	// validate if needed
	return nil
}
