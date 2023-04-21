# Backtest Docs

## Order Types

If an order with the same ID is already pending, it is possible to modify the order. If there is no order with the specified ID, a new order is placed. To deactivate an entry order, the command `strategy.Cancel()` or `strategy.CancelAll()` should be used.

These order types are supported.

- Market Order
- Limit Order


### Market Order

Market order will be executed on the next OHLCV bar.

### Limit Order

Limit order will be placed and open until either it is executed or cancelled.


