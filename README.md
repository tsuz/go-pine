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


Technical Indicators

| Pine Script | Go Pine |
|--|--|
| ta.alma | | 
| ta.atr | | 
| ta.barssince | | 
| ta.bb | | 
| ta.bbw | | 
| ta.cci | | 
| ta.change | | 
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
| ta.kc | | 
| ta.kcw | | 
| ta.linreg | | 
| ta.lowest | | 
| ta.lowestbars | | 
| ta.macd | | 
| ta.max | | 
| ta.median | | 
| ta.mfi | | 
| ta.min | | 
| ta.mode | | 
| ta.mom | | 
| ta.percentile_linear_interpolation | | 
| ta.percentile_nearest_rank | | 
| ta.percentrank | | 
| ta.pivot_point_levels | | 
| ta.pivothigh | | 
| ta.pivotlow | | 
| ta.range | | 
| ta.rising | | 
| ta.rma | pine.RMA() | 
| ta.roc | | 
| ta.rsi | pine.RSI() | 
| ta.sar | | 
| ta.sma | pine.SMA()  | 
| ta.stdev | | 
| ta.stoch | | 
| ta.supertrend | | 
| ta.swma | | 
| ta.tr | | 
| ta.tsi | | 
| ta.valuewhen | | 
| ta.variance | | 
| ta.vwap | | 
| ta.vwma | | 
| ta.wma | | 
| ta.wpr | | 
| ta.accdist | | 
| ta.iii | | 
| ta.nvi | | 
| ta.obv | | 
| ta.pvi | | 
| ta.pvt | | 
| ta.tr | | 
| ta.vwap | | 
| ta.wad | | 
| ta.wvad | | 

## Data Integrity

This library does not make assumptions about the initial OHLCV data which means the developer is responsible for generating the OHLCV slice in an ascending order with correct intervals. The technical analysis indicators uses each candle as a period and so if there are missing time period (i.e. no executions), then it will skip that interval. 

`time.Time` is sometimes used as the unique identifier for `Value` struct so avoid having duplicate time.


[1]: https://www.tradingview.com/pine-script-reference/v5/


[2]: backtest/README.md
