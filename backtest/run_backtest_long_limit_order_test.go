package backtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/tsuz/go-pine/pine"
)

type testLongLimitMystrat struct{}

func (m *testLongLimitMystrat) OnNextOHLCV(strategy Strategy, s pine.OHLCVSeries, state map[string]interface{}) error {

	close := pine.OHLCVAttr(s, pine.OHLCPropClose)

	if *close.Val() < 15 {

		entry1 := EntryOpts{
			Side:  Long,
			Limit: Px(14.1),
		}

		strategy.Entry("Buy1", entry1)
	}

	if *close.Val() > 16 {
		strategy.Exit("Buy1")
	}

	return nil
}

// TestRunBacktestLongLimitOrderImmediate tests when limit orders are executed immediately on the next candle
func TestRunBacktestLongLimitOrderImmediate(t *testing.T) {
	b := &testLongLimitMystrat{}
	data := pine.OHLCVTestData(time.Now(), 4, 5*60*1000)
	data[0].C = 14.9
	data[0].L = 14.6
	data[1].C = 16
	data[1].L = 14
	data[2].C = 16.1
	data[3].O = 16.1
	data[3].C = 17
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
	if fmt.Sprintf("%.03f", res.NetProfit) != "1.142" {
		t.Errorf("Expected NetProfit to be 1.142 but got %+v", res.NetProfit)
	}
}

// TestRunBacktestLongLimitOrderPersist that limit orders should persist for multiple candles
func TestRunBacktestLongLimitOrderPersist(t *testing.T) {
	b := &testLongLimitMystrat{}
	data := pine.OHLCVTestData(time.Now(), 5, 5*60*1000)
	data[0].C = 14.5 // limit order triggered
	data[0].L = 15
	data[1].O = 15
	data[1].C = 15
	data[1].L = 15
	data[2].C = 15
	data[2].L = 13 // limit order filled
	data[3].O = 16.1
	data[3].C = 17
	data[4].O = 16.1
	data[4].C = 17
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
	if fmt.Sprintf("%.03f", res.NetProfit) != "1.142" {
		t.Errorf("Expected NetProfit to be 1.142 but got %+v", res.NetProfit)
	}
}

// TestRunBacktestLongLimitOrderNotPersistAfterExit that limit orders should persist for multiple candles
func TestRunBacktestLongLimitOrderNotPersistAfterExit(t *testing.T) {
	b := &testLongLimitMystrat{}
	data := pine.OHLCVTestData(time.Now(), 4, 5*60*1000)
	data[0].C = 14.5 // <-- limit order triggered
	data[0].L = 15
	data[1].L = 13   // <-- limit order filled
	data[1].C = 17   // <-- exit triggered
	data[2].O = 15   // <-- exit filled
	data[2].L = 13   // <-- limit order should not be open anymore so no trigger
	data[2].C = 15   // <-- limit order should not be open anymore so no trigger
	data[3].O = 16.1 // <-- exit triggered
	data[3].C = 17   // <-- should not fill

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
	if fmt.Sprintf("%.03f", res.NetProfit) != "1.064" {
		t.Errorf("Expected NetProfit to be 1.154 but got %+v", res.NetProfit)
	}
}

// TestRunBacktestLongLimitNotExecuted tests limit orders are not executed
func TestRunBacktestLongLimitNotExecuted(t *testing.T) {
	b := &testLongLimitMystrat{}
	data := pine.OHLCVTestData(time.Now(), 4, 5*60*1000)
	data[0].C = 15
	data[0].L = 15
	data[1].C = 16
	data[1].L = 16
	data[2].C = 17
	data[2].L = 17
	data[3].C = 14.3
	data[3].L = 17

	series, _ := pine.NewOHLCVSeries(data)

	res, err := RunBacktest(series, b)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error runbacktest"))
	}

	if res.TotalClosedTrades != 0 {
		t.Errorf("Expected total trades to be 0 but got %d", res.TotalClosedTrades)
	}
}
