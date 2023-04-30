package pine

import "math"

func DiffAbs(a, b ValueSeries) ValueSeries {
	return Operate(a, b, "diffabs", func(av, bv float64) float64 {
		d := av - bv
		return math.Abs(d)
	})
}
