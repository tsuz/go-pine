package pine

import (
	"math"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestLinReg(t *testing.T) {
	interval := 300
	max := 100
	opts := SeriesOpts{
		Interval: interval,
		Max:      max,
	}
	close := NewOHLCProp(OHLCPropClose)
	linregname := "linreg"
	datatable := []struct {
		linereg int
		io      []struct {
			input  float64
			output map[string]float64
		}
	}{
		{
			linereg: 3,
			io: []struct {
				input  float64
				output map[string]float64
			}{
				{
					input:  4119.05,
					output: nil,
				}, {
					input:  4119.18,
					output: nil,
				}, {
					input: 4118.02,
					output: map[string]float64{
						linregname: 4118.235,
					},
				}, {
					input: 4118.01,
					output: map[string]float64{
						linregname: 4117.818,
					},
				}, {
					input: 4118,
					output: map[string]float64{
						linregname: 4118,
					},
				}, {
					input: 4113.35,
					output: map[string]float64{
						linregname: 4114.123,
					},
				},
			},
		},
		{
			linereg: 10,
			io: []struct {
				input  float64
				output map[string]float64
			}{
				{
					input:  0.00000628, // 2019/03/31 05:50
					output: nil,
				}, {
					input:  0.00000631,
					output: nil,
				}, {
					input:  0.00000631, // 2019/03/31 06:00
					output: nil,
				}, {
					input:  0.00000631,
					output: nil,
				}, {
					input:  0.00000631, // 2019/03/31 06:10
					output: nil,
				}, {
					input:  0.00000632,
					output: nil,
				}, {
					input:  0.00000633, // 2019/03/31 06:20
					output: nil,
				}, {
					input:  0.00000632,
					output: nil,
				}, {
					input:  0.00000632, // 2019/03/31 06:30
					output: nil,
				}, {
					input: 0.00000633,
					output: map[string]float64{
						linregname: 0.0000063315,
					},
				}, {
					input: 0.00000632, // 2019/03/31 06:40
					output: map[string]float64{
						linregname: 0.0000063267,
					},
				}, {
					input: 0.00000631,
					output: map[string]float64{
						linregname: 0.0000063224,
					},
				}, {
					input: 0.00000632, // 2019/03/31 06:50
					output: map[string]float64{
						linregname: 0.0000063215,
					},
				}, {
					input: 0.00000629,
					output: map[string]float64{
						linregname: 0.0000063096,
					},
				}, {
					input: 0.00000630, // 2019/03/31 07:00
					output: map[string]float64{
						linregname: 0.0000063024,
					},
				},
			},
		},
	}

	prettybad := 0.00001
	name := linregname
	for _, o := range datatable {
		now := time.Now()
		ohlcv := make([]OHLCV, 0)
		for j, v := range o.io {
			ohlcv = append(ohlcv, OHLCV{C: v.input, S: now.Add(time.Duration(j*interval) * time.Second)})
		}
		s, err := NewSeries(ohlcv, opts)
		if err != nil {
			t.Fatal(errors.Wrap(err, "error init series"))
		}
		linreg := NewLinReg(close, o.linereg)
		err = s.AddIndicator(name, linreg)
		if err != nil {
			t.Fatal(errors.Wrap(err, "error adding indicator "))
		}
		for j, o := range o.io {
			if err != nil {
				t.Fatal(errors.Wrap(err, "error init series"))
			}
			spectime := now.Add(time.Duration(j*interval) * time.Second)
			v := s.GetValueForInterval(spectime)
			if v.Indicators[name] == nil && o.output == nil {
				// nil value
				continue
			}
			if v.Indicators[name] == nil && o.output == nil {
				// value is nil for that name
				continue
			}
			if v.Indicators[name] == nil || o.output == nil {
				t.Fatalf("expected both to be non nil but got %+v vs %+v at idx: %d for %s", v.Indicators[name], o.output, j, name)
			}
			if math.Abs(*(v.Indicators[name])/o.output[name]-1) >= prettybad {
				t.Errorf("expected: %+v but got %+v for idx: %d for %s", (o.output)[name], *(v.Indicators[name]), j, name)
			}
		}
	}
}
