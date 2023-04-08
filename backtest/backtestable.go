package backtest

import (
	"time"

	"github.com/tsuz/go-pine"
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

type EntryOpts struct {
	Comment string
	Limit   string
	OrdID   string
	Side    Side
	Stop    string
	Qty     string
}

type Side string

const (
	Long  Side = "long"
	Short Side = "short"
)

type OrdMethod string

const (
	Limit  OrdMethod = "limit"
	Market OrdMethod = "market"
)

type Order struct {
	EntryPx float64
	OrdID   string
	Method  OrdMethod
}

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
