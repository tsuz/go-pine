package pine

import "fmt"

func AddConst(a ValueSeries, c float64) ValueSeries {
	return operation(a, a, "addconst", func(av, bv float64) float64 {
		return av + c
	})
}

func Copy(a ValueSeries) ValueSeries {
	return operation(a, a, "copy", func(av, _ float64) float64 {
		return av
	})
}

func Div(a, b ValueSeries) ValueSeries {
	return operation(a, b, "div", func(av, bv float64) float64 {
		return av / bv
	})
}

func Sub(a, b ValueSeries) ValueSeries {
	return operation(a, b, "sub", func(av, bv float64) float64 {
		return av - bv
	})
}

func ReplaceAll(a ValueSeries, c float64) ValueSeries {
	key := fmt.Sprintf("replace:%s:%+v", a.ID(), c)
	return operation(a, a, key, func(av, bv float64) float64 {
		return c
	})
}
