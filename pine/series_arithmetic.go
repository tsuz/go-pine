package pine

<<<<<<< HEAD
import (
	"fmt"
)

func Add(a, b ValueSeries) ValueSeries {
	return operation(a, b, "add", func(av, bv float64) float64 {
		return av + bv
	})
}

func AddConst(a ValueSeries, c float64) ValueSeries {
	key := fmt.Sprintf("addconst:%+v", c)
	return operationConst(a, key, func(av float64) float64 {
=======
import "fmt"

func AddConst(a ValueSeries, c float64) ValueSeries {
	return operation(a, a, "addconst", func(av, bv float64) float64 {
>>>>>>> 41691d6 (fix rsi memory leak)
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

func DivConst(a ValueSeries, c float64) ValueSeries {
	key := fmt.Sprintf("divconst:%+v", c)
	return operationConst(a, key, func(av float64) float64 {
		return av / c
	})
}

func Mul(a, b ValueSeries) ValueSeries {
	return operation(a, b, "mul", func(av, bv float64) float64 {
		return av * bv
	})
}

func MulConst(a ValueSeries, c float64) ValueSeries {
	key := fmt.Sprintf("mulconst:%+v", c)
	return operationConst(a, key, func(av float64) float64 {
		return av * c
	})
}

func ReplaceAll(a ValueSeries, c float64) ValueSeries {
	key := fmt.Sprintf("replace:%+v", c)
	return operation(a, a, key, func(av, bv float64) float64 {
		return c
	})
}

func Sub(a, b ValueSeries) ValueSeries {
	return operation(a, b, "sub", func(av, bv float64) float64 {
		return av - bv
	})
}

func SubConst(a ValueSeries, c float64) ValueSeries {
	key := fmt.Sprintf("subconst:%+v", c)
	return operationConst(a, key, func(av float64) float64 {
		return av - c
	})
}
