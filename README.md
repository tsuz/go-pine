# go-pine
Pinescript to Golang

[![Build Status](https://travis-ci.org/tsuz/go-pine.svg?branch=dev)](https://travis-ci.org/tsuz/go-pine) 
[![codecov](https://codecov.io/gh/tsuz/go-pine/branch/master/graph/badge.svg)](https://codecov.io/gh/tsuz/go-pine)
[![Go Report Card](https://goreportcard.com/badge/tsuz/go-pine)](https://goreportcard.com/report/tsuz/go-pine) 


## Example

*PineScript*
```

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

```

*Golang*
```
// initiate
initialData := make([]pine.OHLCV)
opts := pine.SeriesOpts{
  Interval: 300,
  Max: 30,
  EmptyInst: pine.EmptyInstUseLastClose,
}
s, _ := pine.NewSeries(initialData, opts)

// load indicators
short := 5
long := 20
span := 10
source := pine.NewOHLCProp(pine.OHLCPropClose)
basis := pine.NewSMA(source, short)
basis2 := pine.NewSMA(source, long)
multi := pine.NewArithmetic(pine.ArithmeticMultiplication, basis, basis2)
upperBB := pine.NewArithmetic(pine.ArithmeticAddition, basis, span)
lowerBB := pine.NewArithmetic(pine.ArithmeticSubtraction, basis, span)

s.AddIndicator("upperBB", upperBB)
s.AddIndicator("lowerBB", lowerBB)
s.AddIndicator("multi", multi)

// then add OHLCV or exec and play
t := time.Now()
s.AddOHLCV(pine.OHLCV{O: 14, L: 10, H: 19, C: 14, V: 432, S: t })
s.AddOHLCV(pine.OHLCV{O: 15, L: 8, H: 18, C: 15, V: 192, S: t.Add(time.Minute * 5) })
s.AddOHLCV(pine.OHLCV{O: 16, L: 9, H: 16, C: 13, V: 325, S: t.Add(time.Minute * 10) })
s.AddOHLCV(pine.OHLCV{O: 17, L: 10, H: 19, C: 11, V: 82, S: t.Add(time.Minute * 15) })
...

// or, if you're relying on exec information
s.AddExec(pine.TPQ{Qty: 13, Px: 12.3, Timestamp: t.Add(time.Second * 4) })
s.AddExec(pine.TPQ{Qty: 18, Px: 12.5, Timestamp: t.Add(time.Second * 8) })
s.AddExec(pine.TPQ{Qty: 12, Px: 12.6, Timestamp: t.Add(time.Second * 8) })


v := s.GetValueForInterval(t)

log.Printf("OHLCV: %+v", v.OHLCV)
log.Printf("Indicator values: %+v", v.Indicators)

```


## Limitations

- Assumes initial data is sequential in time ascending order


## Features
TBD


