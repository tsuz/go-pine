package pine

import (
	"time"

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
)

type arith struct {
	t ArithmeticType
	a Indicator
	b Indicator
}

// NewArithmetic generates arithmetic operation on the output of two indicators
func NewArithmetic(t ArithmeticType, a Indicator, b Indicator) Indicator {
	return &arith{
		a: a,
		b: b,
		t: t,
	}
}

func (i *arith) GetValueForInterval(t time.Time) *Interval {
	// validate if needed
	a := i.a.GetValueForInterval(t)
	if a == nil {
		return nil
	}
	b := i.b.GetValueForInterval(t)
	if b == nil {
		return nil
	}
	v := i.generateValue(a.Value, b.Value)
	return &Interval{
		StartTime: t,
		Value:     v,
	}
}

func (i *arith) generateValue(a, b float64) float64 {
	switch i.t {
	case ArithmeticAddition:
		return a + b
	case ArithmeticSubtraction:
		return a - b
	case ArithmeticMultiplication:
		return a * b
	case ArithmeticDivision:
		return a / b
	}
	return 0.0
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
