package backtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/tsuz/go-pine/pine"
)

type testShortLimitMystrat struct{}

func (m *testShortLimitMystrat) OnNextOHLCV(strategy Strategy, s pine.OHLCVSeries, state map[string]interface{}) error {

	close := pine.OHLCVAttr(s, pine.OHLCPropClose)

	entry1 := EntryOpts{
		Side:  Short,
		Limit: Px(16.4),
	}

	strategy.Entry("Short1", entry1)

	if *close.Val() < 15 {
		strategy.Exit("Short1")
	}

	return nil
}

// TestRunBacktestShortLimitOrder tests when backtest orders are executed using Limit orders
func TestRunBacktestShortLimitOrder(t *testing.T) {
	b := &testShortLimitMystrat{}
	data := pine.OHLCVTestData(time.Now(), 4, 5*60*1000)
	data[0].C = 15
	data[0].H = 15.2
	data[1].C = 16
	data[1].H = 16.9
	data[2].C = 14
	data[2].H = 15.9
	data[3].O = 15.8
	data[3].C = 16
	data[3].H = 16
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
	if fmt.Sprintf("%.03f", res.NetProfit) != "1.038" {
		t.Errorf("Expected NetProfit to be 1.038 but got %+v", res.NetProfit)
	}
}

// TestRunBacktestShortLimitNotExecuted tests limit orders are not executed
func TestRunBacktestShortLimitNotExecuted(t *testing.T) {
	b := &testShortLimitMystrat{}
	data := pine.OHLCVTestData(time.Now(), 4, 5*60*1000)
	data[0].C = 15
	data[0].H = 15
	data[1].C = 16
	data[1].H = 16
	data[2].C = 16
	data[2].H = 16
	data[3].C = 14
	data[3].H = 14

	series, _ := pine.NewOHLCVSeries(data)

	res, err := RunBacktest(series, b)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error runbacktest"))
	}

	if res.TotalClosedTrades != 0 {
		t.Errorf("Expected total trades to be 0 but got %d", res.TotalClosedTrades)
	}
}
