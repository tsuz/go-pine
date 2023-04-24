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


## Order Entry/Exit States

On each `OnNextOHLCV` method, `Entry()` and `Exit()` allows to enter and exit a position. 

### Entry()

- Calling `Entry()` without a limit price will execute on the next bar's open
- Calling `Entry()` with a limit price will execute at the limit price if the next bar's low is equal to or below the limit price. The scenario is the same for short orders using the bar's high value.

### Exit()

- Calling `Exit()` without a limit price will execute on the next bar's open.

### Entry() and Exit()

If an order entry and exit both exists for the same order ID, here is the expected behavior:

- If entry and exit are both exist, no trades will be executed.
