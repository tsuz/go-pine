package pine

import (
	"fmt"
)

// operation operates on a and b ValueSeries using op function. use ns as a unique cache identifier
func operation(a, b ValueSeries, ns string, op func(a, b float64) float64) ValueSeries {
	key := fmt.Sprintf("operation:%s:%s:%s", a.ID(), b.ID(), ns)
	dest := getCache(key)
	if dest == nil {
		dest = NewValueSeries()
	}

	var firstaVal *Value
	destlast := dest.GetLast()
	if destlast == nil {
		firstaVal = a.GetFirst()
	} else if destlast != nil {
		if v := a.Get(destlast.t); v != nil && v.next != nil {
			firstaVal = v.next
		}
	}

	// nowhere to start
	if firstaVal == nil {
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

	if cur := a.GetCurrent(); cur != nil {
		dest.SetCurrent(cur.t)
	}

	setCache(key, dest)

	return dest
}
