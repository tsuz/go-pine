package pine

import "math"

func DiffAbs(a, b ValueSeries) ValueSeries {
	return a.Operate(b, func(av, bv float64) float64 {
		d := av - bv
		return math.Abs(d)
	})
}
