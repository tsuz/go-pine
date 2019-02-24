# go-pine
Pinescript to Golang

[![Build Status](https://travis-ci.org/tsuz/go-pine.svg?branch=dev)](https://travis-ci.org/tsuz/go-pine) 
[![codecov](https://codecov.io/gh/tsuz/go-pine/branch/master/graph/badge.svg)](https://codecov.io/gh/tsuz/go-pine)
[![Go Report Card](https://goreportcard.com/badge/tsuz/go-pine)](https://goreportcard.com/report/tsuz/go-pine) 


Assumes initial data is sequential in time ascending order
Only sequential OHLCV, exes are supported. If you're adding OHLCV or exec for an interval previous, it will not work

We don't support odd interval. It's increment from UTC
