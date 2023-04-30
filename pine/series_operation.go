package pine

import (
	"fmt"
)

// Operate operates on two series. Enabling caching means it starts from where it was left off.
func Operate(a, b ValueSeries, ns string, op func(b, c float64) float64) ValueSeries {
	return operation(a, b, ns, op, true)
}

// OperateNoCache operates on two series without caching
func OperateNoCache(a, b ValueSeries, ns string, op func(b, c float64) float64) ValueSeries {
	return operation(a, b, ns, op, false)
}

// operation operates on a and b ValueSeries using op function. use ns as a unique cache identifier
func operation(a, b ValueSeries, ns string, op func(a, b float64) float64, cache bool) ValueSeries {
	key := fmt.Sprintf("operation:%s:%s:%s", a.ID(), b.ID(), ns)
	dest := getCache(key)
	if dest == nil {
		dest = NewValueSeries()
	}

	firstaVal := operationGetStart(a, dest)

	// nowhere to start
	if firstaVal == nil {

		// propagate current pointer if needed
		propagateCurrent(a, dest)

		return dest
	}

	f := firstaVal
	for {
		if f == nil {
			break
		}

		newv := b.Get(f.t)

		if newv != nil {
			dest.Set(f.t, op(f.v, newv.v))
		}

		f = f.next
	}

	propagateCurrent(a, dest)

	if cache {
		setCache(key, dest)
	}

	return dest
}

func propagateCurrent(a, dest ValueSeries) {
	if cur := a.GetCurrent(); cur != nil {
		dest.SetCurrent(cur.t)
	}
}

func operationGetStart(a, dest ValueSeries) *Value {
	var firstaVal *Value
	destlast := dest.GetLast()
	if destlast == nil {
		firstaVal = a.GetFirst()
	} else if destlast != nil {
		if v := a.Get(destlast.t); v != nil && v.next != nil {
			firstaVal = v.next
		}
	}
	return firstaVal
}
