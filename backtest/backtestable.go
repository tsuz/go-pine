package backtest

import (
	"time"

	"github.com/tsuz/go-pine/pine"
)

type BackTestable interface {
	OnNextOHLCV(Strategy, pine.OHLCVSeries, map[string]interface{}) error
}

type BacktestResult struct {
	ClosedOrd         []Position
	NetProfit         float64
	PercentProfitable float64
	ProfitableTrades  int64
	TotalClosedTrades int64
}

// EntryOpts is additional entry options
type EntryOpts struct {
	Comment string

	// Limit price is used if this value is non nil. If it's nil, market order is executed
	Limit *float64

	OrdID string
	Side  Side
	Stop  string
	Qty   string
}

// Px generates a non nil float64
func Px(v float64) *float64 {
	v2 := &v
	return v2
}

type Side string

const (
	Long  Side = "long"
	Short Side = "short"
)

type Position struct {
	EntryPx   float64
	ExitPx    float64
	EntryTime time.Time
	ExitTime  time.Time
	EntrySide Side
	OrdID     string
}

func (p Position) Profit() float64 {
	switch p.EntrySide {
	case Long:
		return p.ExitPx / p.EntryPx
	case Short:
		return p.EntryPx / p.ExitPx
	}
	return 0
}

func (b *BacktestResult) CalculateNetProfit() {
	start := 1.0
	for _, v := range b.ClosedOrd {
		p := v.Profit()
		start = start * p
	}
	b.NetProfit = start
}
