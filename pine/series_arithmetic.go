package pine

import (
	"fmt"
)

func Add(a, b ValueSeries) ValueSeries {
	return operation(a, b, "add", func(av, bv float64) float64 {
		return av + bv
	}, true)
}

func AddConst(a ValueSeries, c float64) ValueSeries {
	key := fmt.Sprintf("addconst:%+v", c)
	return operationConst(a, key, func(av float64) float64 {
		return av + c
	}, true)
}

func AddConstNoCache(a ValueSeries, c float64) ValueSeries {
	key := fmt.Sprintf("addconst:%+v", c)
	return operationConst(a, key, func(av float64) float64 {
		return av + c
	}, false)
}

func Copy(a ValueSeries) ValueSeries {
	return operation(a, a, "copy", func(av, _ float64) float64 {
		return av
	}, true)
}

func Div(a, b ValueSeries) ValueSeries {
	return operation(a, b, "div", func(av, bv float64) float64 {
		return av / bv
	}, true)
}

func DivNoCache(a, b ValueSeries) ValueSeries {
	return operation(a, b, "div", func(av, bv float64) float64 {
		return av / bv
	}, false)
}

func DivConst(a ValueSeries, c float64) ValueSeries {
	key := fmt.Sprintf("divconst:%+v", c)
	return operationConst(a, key, func(av float64) float64 {
		return av / c
	}, true)
}

func DivConstNoCache(a ValueSeries, c float64) ValueSeries {
	key := fmt.Sprintf("divconst:%+v", c)
	return operationConst(a, key, func(av float64) float64 {
		return av / c
	}, false)
}

func Mul(a, b ValueSeries) ValueSeries {
	return operation(a, b, "mul", func(av, bv float64) float64 {
		return av * bv
	}, true)
}

func MulConst(a ValueSeries, c float64) ValueSeries {
	key := fmt.Sprintf("mulconst:%+v", c)
	return operationConst(a, key, func(av float64) float64 {
		return av * c
	}, true)
}

func MulConstNoCache(a ValueSeries, c float64) ValueSeries {
	key := fmt.Sprintf("mulconst:%+v", c)
	return operationConst(a, key, func(av float64) float64 {
		return av * c
	}, false)
}

func ReplaceAll(a ValueSeries, c float64) ValueSeries {
	key := fmt.Sprintf("replace:%+v", c)
	return operation(a, a, key, func(av, bv float64) float64 {
		return c
	}, true)
}

func Sub(a, b ValueSeries) ValueSeries {
	return operation(a, b, "sub", func(av, bv float64) float64 {
		return av - bv
	}, true)
}

func SubConst(a ValueSeries, c float64) ValueSeries {
	key := fmt.Sprintf("subconst:%+v", c)
	return operationConst(a, key, func(av float64) float64 {
		return av - c
	}, true)
}

func SubConstNoCache(a ValueSeries, c float64) ValueSeries {
	key := fmt.Sprintf("subconst:%+v", c)
	return operationConst(a, key, func(av float64) float64 {
		return av - c
	}, false)
}
