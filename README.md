# go-pine
Pinescript to Golang

Assumes initial data is sequential in time ascending order
Only sequential OHLCV, exes are supported. If you're adding OHLCV or exec for an interval previous, it will not work

We don't support odd interval. It's increment from UTC