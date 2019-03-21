package pine

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func TestArithmetic(t *testing.T) {
	opts := SeriesOpts{
		Interval: 300,
		Max:      100,
	}
	now := time.Now()
	fivemin := now.Add(5 * time.Minute)
	hl2 := NewOHLCProp(OHLCPropOpen)
	cval := 15.21
	c := NewConstant(cval)

	io := []struct {
		name    string
		t       ArithmeticType
		outputs []decimal.Decimal
	}{
		{
			name: "add",
			t:    ArithmeticAddition,
			outputs: []decimal.Decimal{
				decimal.NewFromFloat(14).Add(decimal.NewFromFloat(float64(cval))),
				decimal.NewFromFloat(13).Add(decimal.NewFromFloat(float64(cval))),
			},
		},
		{
			name: "sub",
			t:    ArithmeticSubtraction,
			outputs: []decimal.Decimal{
				decimal.NewFromFloat(14).Sub(decimal.NewFromFloat(float64(cval))),
				decimal.NewFromFloat(13).Sub(decimal.NewFromFloat(float64(cval))),
			},
		},
		{
			name: "mul",
			t:    ArithmeticMultiplication,
			outputs: []decimal.Decimal{
				decimal.NewFromFloat(14).Mul(decimal.NewFromFloat(float64(cval))),
				decimal.NewFromFloat(13).Mul(decimal.NewFromFloat(float64(cval))),
			},
		},
		{
			name: "div",
			t:    ArithmeticDivision,
			outputs: []decimal.Decimal{
				decimal.NewFromFloat(14).Div(decimal.NewFromFloat(float64(cval))),
				decimal.NewFromFloat(13).Div(decimal.NewFromFloat(float64(cval))),
			},
		},
		{
			name: "abs",
			t:    ArithmeticAbsDiff,
			outputs: []decimal.Decimal{
				decimal.NewFromFloat(14).Sub(decimal.NewFromFloat(float64(cval))).Abs(),
				decimal.NewFromFloat(13).Sub(decimal.NewFromFloat(float64(cval))).Abs(),
			},
		},
		{
			name: "max",
			t:    ArithmeticMax,
			outputs: []decimal.Decimal{
				decimal.NewFromFloat(15.21),
				decimal.NewFromFloat(15.21),
			},
		},
		{
			name: "min",
			t:    ArithmeticMin,
			outputs: []decimal.Decimal{
				decimal.NewFromFloat(14),
				decimal.NewFromFloat(13),
			},
		},
	}

	data := []OHLCV{
		OHLCV{
			O: 14,
			H: 15,
			L: 13,
			C: 14,
			V: 131,
			S: now,
		},
		OHLCV{
			O: 13,
			H: 18,
			L: 10,
			C: 15,
			V: 12,
			S: fivemin,
		},
	}

	s, err := NewSeries(data, opts)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error init series"))
	}
	for _, o := range io {
		ar := NewArithmetic(o.t, hl2, c, ArithmeticOpts{})
		s.AddIndicator(o.name, ar)
	}

	for i, o := range io {
		nowv := s.GetValueForInterval(now)
		fst, _ := o.outputs[0].Float64()
		snd, _ := o.outputs[1].Float64()
		if *(nowv.Indicators[o.name]) != fst {
			t.Errorf("expected: %+v but got %+v at idx %d", fst, *(nowv.Indicators[o.name]), i)
		}
		fivv := s.GetValueForInterval(fivemin)
		if *(fivv.Indicators[o.name]) != snd {
			t.Errorf("expected: %+v but got %+v at idx %d", snd, *(fivv.Indicators[o.name]), i)
		}
	}
}
