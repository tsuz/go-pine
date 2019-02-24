package pine_test

import (
	pine "go-pine"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestArithmetic(t *testing.T) {
	opts := pine.SeriesOpts{
		Interval: 300,
		Max:      100,
	}
	now := time.Now()
	fivemin := now.Add(5 * time.Minute)
	hl2 := pine.NewOHLCProp(pine.OHLCPropOpen)
	cval := 3.21
	c := pine.NewConstant(cval)

	io := []struct {
		name    string
		t       pine.ArithmeticType
		outputs []float64
	}{
		{
			name:    "add",
			t:       pine.ArithmeticAddition,
			outputs: []float64{14 + cval, 13 + cval},
		},
		{
			name:    "sub",
			t:       pine.ArithmeticSubtraction,
			outputs: []float64{14 - cval, 13 - cval},
		},
		{
			name:    "mul",
			t:       pine.ArithmeticMultiplication,
			outputs: []float64{14 * cval, 13 * cval},
		},
		{
			name:    "div",
			t:       pine.ArithmeticDivision,
			outputs: []float64{14 / cval, 13 / cval},
		},
	}

	data := []pine.OHLCV{
		pine.OHLCV{
			O: 14,
			H: 15,
			L: 13,
			C: 14,
			V: 131,
			S: now,
		},
		pine.OHLCV{
			O: 13,
			H: 18,
			L: 10,
			C: 15,
			V: 12,
			S: fivemin,
		},
	}

	s, err := pine.NewSeries(data, opts)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error init series"))
	}
	for _, o := range io {
		ar := pine.NewArithmetic(o.t, hl2, c)
		s.AddIndicator(o.name, ar)
	}

	for i, o := range io {
		nowv := s.GetValueForInterval(now)
		if *(nowv.Indicators[o.name]) != o.outputs[0] {
			t.Errorf("expected: %+v but got %+v at idx %d", o.outputs[0], *(nowv.Indicators[o.name]), i)
		}
		fivv := s.GetValueForInterval(fivemin)
		if *(fivv.Indicators[o.name]) != o.outputs[1] {
			t.Errorf("expected: %+v but got %+v at idx %d", o.outputs[1], *(fivv.Indicators[o.name]), i)
		}
	}
}
