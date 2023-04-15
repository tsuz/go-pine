package pine

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
	// OHLCPropHL2 is the midpoint value of OHLC
	OHLCPropHL2
	// OHLCPropHLC3 is (high + low + close) / 3 of OHLC
	OHLCPropHLC3
	// OHLCPropTR is true range i.e. max(high - low, abs(high - close[1]), abs(low - close[1])).
	OHLCPropTR
)
