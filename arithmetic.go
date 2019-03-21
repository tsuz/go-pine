package pine

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/pkg/errors"
)

// ArithmeticType defines the arthmetic operation
type ArithmeticType int

const (
	// ArithmeticAddition adds values
	ArithmeticAddition ArithmeticType = iota
	// ArithmeticSubtraction subtracts values
	ArithmeticSubtraction
	// ArithmeticMultiplication multiplies values
	ArithmeticMultiplication
	// ArithmeticDivision divides values
	ArithmeticDivision
	// ArithmeticAbsDiff shows absolute difference math.Abs(a-b)
	ArithmeticAbsDiff
	// ArithmeticMax shows maximum of the two
	ArithmeticMax
	// ArithmeticMin shows minimum of the two
	ArithmeticMin
)

type arith struct {
	a Indicator
	b Indicator
	o ArithmeticOpts
	t ArithmeticType
}

// NewArithmetic generates arithmetic operation on the output of two indicators
func NewArithmetic(t ArithmeticType, a Indicator, b Indicator, o ArithmeticOpts) Indicator {
	return &arith{
		a: a,
		b: b,
		o: o,
		t: t,
	}
}

func (i *arith) GetValueForInterval(t time.Time) *Interval {
	// validate if needed
	a := i.a.GetValueForInterval(t)
	b := i.b.GetValueForInterval(t)
	v := i.generateValue(a, b)
	if v == nil {
		return nil
	}
	return &Interval{
		StartTime: t,
		Value:     *v,
	}
}

func (i *arith) generateValue(ai, bi *Interval) *float64 {
	if ai == nil || bi == nil {
		switch i.o.NilHandlInst {
		case NilValueReturnNil:
			return nil
		case NilValueReturnZero:
			val := 0.0
			return &val
		}
	}
	var val decimal.Decimal
	a := decimal.NewFromFloat(ai.Value)
	b := decimal.NewFromFloat(bi.Value)
	switch i.t {
	case ArithmeticAddition:
		val = a.Add(b)
	case ArithmeticSubtraction:
		val = a.Sub(b)
	case ArithmeticMultiplication:
		val = a.Mul(b)
	case ArithmeticDivision:
		val = a.Div(b)
	case ArithmeticAbsDiff:
		val = a.Sub(b).Abs()
	case ArithmeticMax:
		if a.GreaterThan(b) {
			val = a
		} else {
			val = b
		}
	case ArithmeticMin:
		if a.LessThan(b) {
			val = a
		} else {
			val = b
		}
	}
	f64, _ := val.Float64()
	return &f64
}

func (i *arith) Update(v OHLCV) error {
	// validate if needed
	if err := i.a.Update(v); err != nil {
		return errors.Wrap(err, "error updating in addition")
	}
	if err := i.b.Update(v); err != nil {
		return errors.Wrap(err, "error updating in addition")
	}
	return nil
}

func (i *arith) ApplyOpts(opts SeriesOpts) error {
	// validate if needed
	if err := i.a.ApplyOpts(opts); err != nil {
		return errors.Wrap(err, "error applying opts in addition")
	}
	if err := i.b.ApplyOpts(opts); err != nil {
		return errors.Wrap(err, "error applying opts in addition")
	}
	return nil
}
