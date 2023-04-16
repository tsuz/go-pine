package backtest

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/tsuz/go-pine/pine"
)

type testCancelMystrat struct{}

func (m *testCancelMystrat) OnNextOHLCV(strategy Strategy, s pine.OHLCVSeries, state map[string]interface{}) error {

	close := s.GetSeries(pine.OHLCPropClose)

	if *close.Val() == 14 {
		entry1 := EntryOpts{
			Side:  Long,
			Limit: Px(12),
		}
		strategy.Entry("Buy1", entry1)
	}

	if *close.Val() == 13 {
		strategy.Cancel("Buy1")
	}

	strategy.Exit("Buy1")

	return nil
}

// TestRunBacktesttestCancelOrder tests canceling existing orders
func TestRunBacktesttestCancelOrder(t *testing.T) {
	b := &testCancelMystrat{}

	data := pine.OHLCVTestData(time.Now(), 3, 5*60*1000)
	data[0].C = 14
	data[1].L = 13
	data[1].C = 13
	data[2].L = 12
	data[2].C = 13
	series, _ := pine.NewOHLCVSeries(data)

	res, err := RunBacktest(series, b)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error runbacktest"))
	}

	if res.TotalClosedTrades != 0 {
		t.Errorf("Expected total trades to be 0 but got %d", res.TotalClosedTrades)
	}
}
