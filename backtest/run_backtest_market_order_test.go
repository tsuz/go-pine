package backtest

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/tsuz/go-pine"
)

type testMystrat struct{}

func (m *testMystrat) OnNextOHLCV(strategy Strategy, s pine.OHLCVSeries, state map[string]interface{}) error {

	close := s.GetSeries(pine.OHLCPropClose)
	avg, _ := pine.SMA(close, 2)

	if avg.Val() != nil {
		log.Printf("*avg.Val() is %+v", *avg.Val())
		if *avg.Val() > 15.4 {
			entry1 := EntryOpts{
				Side: Long,
			}
			strategy.Entry("Buy1", entry1)
		}
		if *avg.Val() >= 16.0 {
			strategy.Exit("Buy1")
		}
	}

	return nil
}

// TestRunBacktestMarketOrder tests when backtest orders are executed using market orders
func TestRunBacktestMarketOrder(t *testing.T) {
	b := &testMystrat{}
	data := pine.OHLCVTestData(time.Now(), 4, 5*60*1000)
	data[0].C = 15
	data[1].C = 16
	data[2].O = 15.04
	data[2].C = 17
	data[3].O = 16.92
	data[3].C = 18
	series, _ := pine.NewOHLCVSeries(data)

	res, err := RunBacktest(series, b)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error runbacktest"))
	}

	if res.TotalClosedTrades != 1 {
		t.Errorf("Expected total trades to be 1 but got %d", res.TotalClosedTrades)
	}
	if res.PercentProfitable != 1 {
		t.Errorf("Expected pct profitable to be 1 but got %+v", res.PercentProfitable)
	}
	if res.ProfitableTrades != 1 {
		t.Errorf("Expected profitable trades to be 1 but got %+v", res.ProfitableTrades)
	}
	if fmt.Sprintf("%.03f", res.NetProfit) != "1.125" {
		t.Errorf("Expected NetProfit to be 1.125 but got %+v", res.NetProfit)
	}
}
