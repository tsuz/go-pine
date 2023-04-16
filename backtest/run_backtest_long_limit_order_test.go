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

	close := s.GetSeries(pine.OHLCPropClose)

	entry1 := EntryOpts{
		Side:  Long,
		Limit: Px(14.1),
	}

	strategy.Entry("Buy1", entry1)

	if *close.Val() > 16 {
		strategy.Exit("Buy1")
	}

	return nil
}

// TestRunBacktestLongLimitOrder tests when backtest orders are executed using Limit orders
func TestRunBacktestLongLimitOrder(t *testing.T) {
	b := &testLongLimitMystrat{}
	data := pine.OHLCVTestData(time.Now(), 4, 5*60*1000)
	data[0].C = 15
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
