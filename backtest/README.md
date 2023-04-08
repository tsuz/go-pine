# Backtesting

Backtesting refers to the process of evaluating a trading strategy using historical data to see how it would have performed in the past.

## Example



```go


type mystrat struct{}

func (m *mystrat) OnNextOHLCV(strategy backtest.Strategy, s pine.OHLCVSeries, state map[string]interface{}) error {

	var short int64 = 2
	var long int64 = 20
	var span float64 = 10

	close := s.GetSeries(pine.OHLCPropClose)
	open := s.GetSeries(pine.OHLCPropOpen)

	basis, _ := pine.SMA(close, short)
	basis2, _ := pine.SMA(open, long)
	rsi, _ := pine.RSI(close, short)
	avg, _ := pine.SMA(rsi, long)

	basis3 := basis.Add(basis2)
	upperBB := basis3.AddConst(span)

	log.Printf("t: %+v, close: %+v, rsi: %+v, avg: %+v, upperBB: %+v", s.Current().S, close.Val(), rsi.Val(), avg.Val(), upperBB.Val())

	if rsi.Val() != nil {
		if *rsi.Val() < 30 {
			log.Printf("Entry: %+v", *rsi.Val())
			entry1 := backtest.EntryOpts{
				Side: backtest.Long,
			}
			strategy.Entry("Buy1", entry1)
		}
		if *rsi.Val() > 70 {
			log.Printf("Exit %+v", *rsi.Val())
			strategy.Exit("Buy1")
		}
	}

	return nil
}

func main() {
	b := &mystrat{}
	data := pine.OHLCVTestData(time.Now(), 25, 5*60*1000)
	series, _ := pine.NewOHLCVSeries(data)

	res, _ := backtest.RunBacktest(series, b)
	log.Printf("TotalClosedTrades %d, PercentProfitable: %.03f, NetProfit: %.03f", res.TotalClosedTrades, res.PercentProfitable, res.NetProfit)
}

```

## List of Features

General Strategy

| Item | Supported | 
|--|--|
| Initial Capital |  | 
| Pyramiding |  | 
| Commission |  | 
| Verify Price for limit orders|  | 
| Slippage |  | 
| Margin for long positions |  | 
| Margin for short positions |  | 
| Recalculate after order is filled |  | 
| Recalculate on every tick |  | 
| Recalculate on bar close |  | 
| Use bar magnifier  |  | 

Order Type

| Item |Supported | 
|--|--|
| Market Order| ✅ | 
| Limit Order|  | 
| Stop Order|  | 
| Trail Order|  | 


Entry Options

| Item | Supported |  
|--|--|
| id | ✅ | 
| direction | ✅ | 
| qty |  | 
| limit |  | 
| stop |  | 
| oca_name |  | 
| oca_type |  | 
| comment |  | 

Exit Options

| Item | Supported |  |
|--|--|--|
| id |  | 
| from_entry | ✅ | 
| qty |  | 
| qty_percent |  | 
| profit |  | 
| limit |  | 
| loss |  | 
| stop |  | 
| trail_price |  | 
| trail_points |  | 
| trail_offset |  | 
| oca_name |  | 
| oca_type |  | 
| comment |  | 
| comment_profit |  | 
| comment_loss |  | 
| comment_trailing |  | 

