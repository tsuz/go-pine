package pine

// Crossover generates ValueSeries of ketler channel's middle, upper and lower in that order.
func Crossover(a, b ValueSeries) (ValueSeries, error) {
	c := a.OperateWithNil(b, func(av, bv *Value) *Value {
		if av == nil || bv == nil {
			return nil
		}
		zero := &Value{
			t: av.t,
			v: 0,
		}
		if av.prev == nil || bv.prev == nil {
			return zero
		}
		if av.v > bv.v && av.prev.v < bv.prev.v {
			return &Value{
				t: av.t,
				v: 1.0,
			}
		}
		return zero
	})

	return c, nil
}
