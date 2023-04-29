package pine

import (
	"fmt"
)

// operationConst operates on a and b ValueSeries using op function. use ns as a unique cache identifier
func operationConst(a ValueSeries, ns string, op func(a float64) float64) ValueSeries {
	key := fmt.Sprintf("operationconst:%s:%s", a.ID(), ns)
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

		dest.Set(f.t, op(f.v))

		f = f.next
	}

	propagateCurrent(a, dest)

	setCache(key, dest)

	return dest
}
