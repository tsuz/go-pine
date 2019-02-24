# go-pine
Pinescript to Golang

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
source := pine.NewOHLCProp(pine.OHLCPropClose)
basis := pine.NewSMA(source, short)
basis2 := pine.NewSMA(source, long)
multi := pine.NewArithmetic(pine.ArithmeticMultiplication, basis, basis2)
upperBB := pine.NewArithmetic(pine.ArithmeticAddition, span)
lowerBB := pine.NewArithmetic(pine.ArithmeticSubtraction, span)

s.AddIndicator("upperBB", upperBB)
s.AddIndicator("lowerBB", lowerBB)
s.AddIndicator("multi", multi)

// then add OHLCV or exec and play
t := time.Now()
s.AddOHLCV(pine.OHLCV{O: 14, L: 10, H: 19, C: 15, V: 432, S: t })

v := s.GetValueForInterval(t)

log.Printf("OHLCV: %+v", v.OHLCV)
log.Printf("Indicator values: %+v", v.Indicators)

```


## Limitations

- Assumes initial data is sequential in time ascending order


## Features
TBD


