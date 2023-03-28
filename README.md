# Go Pine

Backtesting tool written in Golang inspired by PineScript from TradingView.

[![Build Status](https://dl.circleci.com/status-badge/img/gh/tsuz/go-pine/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/tsuz/go-pine/tree/main)
[![docs godocs](https://img.shields.io/badge/docs-godoc-brightgreen.svg?style=flat)](https://godoc.org/github.com/tsuz/go-pine)
[![codecov](https://codecov.io/gh/tsuz/go-pine/branch/main/graph/badge.svg?token=1EeuK2Ro6F)](https://codecov.io/gh/tsuz/go-pine)
[![Go Report Card](https://goreportcard.com/badge/tsuz/go-pine)](https://goreportcard.com/report/tsuz/go-pine) 
[![HitCount](http://hits.dwyl.io/tsuz/go-pine.svg)](http://hits.dwyl.io/tsuz/go-pine)
[![Maintainability](https://api.codeclimate.com/v1/badges/ba4f05de8cb12c615695/maintainability)](https://codeclimate.com/github/tsuz/go-pine/maintainability)

## Requirements

- Golang v1.20 (recommended)

## Example

### Backtest

**Pine Script**

```js

study("Pine test")

lengthshort = 5
lengthlong = 20
span = 10

source = close
basis = sma(source, lengthshort)
basis2 = sma(source, lengthlong)
multi = basis * basis2
upperBB = basis + span
lowerBB = basis - span

strategy.entry("Buy1", strategy.long, qty=1, limit=1234.56, stop=1211, comment="My Long Signal")
```

*Golang*

my_strategy.go

```go

type mystrat struct{
  ser: pine.Series
}

func NewMyStrat() (pine.BackTestable, error) {
  m := &mystrat{}
  return m, err
}

func (mystrat *m) OnNextOHLCV(strategy pine.Strategy, s pine.OHLCVSeries, states map[string]interface{}) error {

  short := 5
  long := 20
  span := 10
  source := pine.Close

  basis := s.GetSMA(source, short)
  basis2 := s.GetSMA(source, long)
  multi := basis.Add(basis2)
  upperBB := basis.AddConst(span)
  lowerBB := basis.AddConst(span)
  
  log.Printf("Get upper boundary", upperBB)
  log.Printf("Get lower boundary", lowerBB)

  entry1 := pine.EntryOpts {
    Comment: "My Long Signal",
    Limit: "1234.56",
    Stop: "1211",
    Qty: "1",
    Side: pine.Long,
  }
  strategy.Entry("Buy1", entry1)

  return nil
}
```

**main.go**

```go

s, _ := NewMyStrat()
opts := pine.OHLCVSeriesOpts{
  // OHLC interval in milliseconds. Below equates to a 5 minute interval.
  Interval: 300000,
  
  // The first OHLCV that will be fed into the backtest logic. This will also be used as the OHLCV's start offset
  StartTime: time.Date(2009, 1, 1, 12, 0, 0, 0, time.UTC),
  
  // How many look backs to cache. Defaults to 100.
  Max: 500,
}
ohlcv, _ := pine.NewOHLCVSeries(initialData, opts)
res, _ := pine.RunBacktest(ohlcv, s)

log.Printf("Results are %+v", res)
// NetProfit: 649%, Total Closed Trades: 436, Percent Profitable: 61.93%, Profit Factor: 1.622, Max Drawdown: -27.44%, Avg Trade: 14.89, Avg # Bars in Trade

```


## Supported Features

The functions are listed in the [Pine Script reference manual V5][1]

Language Operators

| Pine Script | Go Pine |
|--|--|
| != | .NotEq() | 
| == | .Eq() | 
| + | ArithmeticAddition | 
| - | AArithmeticSubtraction | 
| * | ArithmeticMultiplication |
| / | ArithmeticDivision |

Mathematical Operators

| Pine Script | Go Pine | 
|--|--|
| math.max | ArithmeticMax |
| math.min | ArithmeticMin |

Technical Indicators

| Pine Script | Go Pine |
|--|--|
| ta.median | .NewMedian()| 
| ta.sma | .NewSMA() | 
| ta.stdev | .NewStDev() | 



[1]: https://www.tradingview.com/pine-script-reference/v5/

