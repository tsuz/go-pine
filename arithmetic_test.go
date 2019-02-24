package pine

import (
	"math"
	"testing"
	"time"

	"github.com/pkg/errors"
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
		outputs []float64
	}{
		{
			name:    "add",
			t:       ArithmeticAddition,
			outputs: []float64{14 + cval, 13 + cval},
		},
		{
			name:    "sub",
			t:       ArithmeticSubtraction,
			outputs: []float64{14 - cval, 13 - cval},
		},
		{
			name:    "mul",
			t:       ArithmeticMultiplication,
			outputs: []float64{14 * cval, 13 * cval},
		},
		{
			name:    "div",
			t:       ArithmeticDivision,
			outputs: []float64{14 / cval, 13 / cval},
		},
		{
			name:    "abs",
			t:       ArithmeticAbsDiff,
			outputs: []float64{math.Abs(14 - cval), math.Abs(13 - cval)},
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
		if *(nowv.Indicators[o.name]) != o.outputs[0] {
			t.Errorf("expected: %+v but got %+v at idx %d", o.outputs[0], *(nowv.Indicators[o.name]), i)
		}
		fivv := s.GetValueForInterval(fivemin)
		if *(fivv.Indicators[o.name]) != o.outputs[1] {
			t.Errorf("expected: %+v but got %+v at idx %d", o.outputs[1], *(fivv.Indicators[o.name]), i)
		}
	}
}
