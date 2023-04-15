# Go Pine

Backtesting tool written in Golang inspired by PineScript from TradingView.

[![Build Status](https://dl.circleci.com/status-badge/img/gh/tsuz/go-pine/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/tsuz/go-pine/tree/main)
[![docs godocs](https://img.shields.io/badge/docs-godoc-brightgreen.svg?style=flat)](https://godoc.org/github.com/tsuz/go-pine)
[![codecov](https://codecov.io/gh/tsuz/go-pine/branch/main/graph/badge.svg?token=1EeuK2Ro6F)](https://codecov.io/gh/tsuz/go-pine)
[![Go Report Card](https://goreportcard.com/badge/tsuz/go-pine)](https://goreportcard.com/report/tsuz/go-pine) 
[![HitCount](http://hits.dwyl.io/tsuz/go-pine.svg)](http://hits.dwyl.io/tsuz/go-pine)
[![Maintainability](https://api.codeclimate.com/v1/badges/ba4f05de8cb12c615695/maintainability)](https://codeclimate.com/github/tsuz/go-pine/maintainability)

> Note: This library is under heavy development

## Requirements

- Golang v1.20 (recommended)

## Example

### Backtest

See [backtest example][2].


## Supported Features

The functions are listed in the [Pine Script reference manual V5][1]

Language Operators

| Pine Script | Go Pine |
|--|--|
| != | .NotEq() | 
| == | .Eq() | 
| < | |
| <= | |
| > | |
| >= | |
| + | ValueSeries.Add() | 
| - | ValueSeries.Sub() | 
| * | ValueSeries.Mul() |
| / | ValueSeries.Div() |
| % | |

Mathematical Operators

| Pine Script | Go Pine |
|--|--|
| math.abs | |
| math.acos | |
| math.asin | |
| math.atan | |
| math.avg | |
| math.ceil | |
| math.cos | |
| math.e | |
| math.exp | |
| math.floor | |
| math.log | |
| math.log10 | |
| math.max | |
| math.min | |
| math.phi | |
| math.pi | |
| math.pow | pine.Pow |
| math.random | |
| math.round | |
| math.round_to_mintick | |
| math.rphi | |
| math.sign | |
| math.sin | |
| math.sqrt | use pine.Pow(src, 0.5) |
| math.sum | pine.Sum |
| math.tan | |
| math.todegrees | |
| math.toradians | |

Technical Indicators

| Pine Script | Go Pine |
|--|--|
| ta.alma | | 
| ta.accdist | |
| ta.atr | pine.ATR | 
| ta.barssince | | 
| ta.bb | | 
| ta.bbw | | 
| ta.cci | | 
| ta.change | pine.Change() | 
| ta.cmo | | 
| ta.cog | | 
| ta.correlation | | 
| ta.cross | | 
| ta.crossover | | 
| ta.crossunder | | 
| ta.cum | | 
| ta.dev | | 
| ta.dmi | | 
| ta.ema |  pine.EMA() | 
| ta.falling | | 
| ta.highest | | 
| ta.highestbars | | 
| ta.hma | | 
| ta.iii | | 
| ta.kc | | 
| ta.kcw | | 
| ta.linreg | | 
| ta.lowest | | 
| ta.lowestbars | | 
| ta.macd | pine.MACD() | 
| ta.max | | 
| ta.median | | 
| ta.mfi | | 
| ta.min | | 
| ta.mode | | 
| ta.mom | | 
| ta.nvi | | 
| ta.obv | | 
| ta.percentile_linear_interpolation | | 
| ta.percentile_nearest_rank | | 
| ta.percentrank | | 
| ta.pivot_point_levels | | 
| ta.pivothigh | | 
| ta.pivotlow | | 
| ta.pvi | | 
| ta.pvt | | 
| ta.range | | 
| ta.rising | | 
| ta.rma | pine.RMA() | 
| ta.roc | pine.ROC() | 
| ta.rsi | pine.RSI() | 
| ta.sar | | 
| ta.sma | pine.SMA()  | 
| ta.stdev | pine.Stdev() | 
| ta.stoch | | 
| ta.supertrend | | 
| ta.swma | | 
| ta.tr | OHLCVSeries.getSeries(OHLCPropTR) | 
| ta.tsi | | 
| ta.valuewhen | pine.ValueWhen() | 
| ta.variance | pine.Variance() | 
| ta.vwap | | 
| ta.vwma | | 
| ta.wad | | 
| ta.wma | | 
| ta.wpr | | 
| ta.wvad | | 

## Data Integrity

This library does not make assumptions about the initial OHLCV data which means the developer is responsible for generating the OHLCV slice in an ascending order with correct intervals. The technical analysis indicators uses each candle as a period and so if there are missing time period (i.e. no executions), then it will skip that interval. 

`time.Time` is sometimes used as the unique identifier for `Value` struct so avoid having duplicate time.


[1]: https://www.tradingview.com/pine-script-reference/v5/


[2]: backtest/README.md
