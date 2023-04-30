package pine

import (
	"fmt"
	"log"
)

func OperateWithNil(a, b ValueSeries, ns string, op func(a, b *Value) *Value) ValueSeries {
	key := fmt.Sprintf("operationwnil:%s:%s:%s", a.ID(), b.ID(), ns)
	dest := getCache(key)
	if dest == nil {
		log.Printf("New operation series %+v", key)
		dest = NewValueSeries()
	}

	// copied := NewValueSeries()
	f := a.GetFirst()
	for {
		if f == nil {
			break
		}

		newv := b.Get(f.t)

		if val := op(f, newv); val != nil {
			dest.Set(val.t, val.v)
		}

		f = f.next
	}

	propagateCurrent(b, dest)

	setCache(key, dest)

	return dest
}
