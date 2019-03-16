package pine

import (
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestChange(t *testing.T) {
	opts := SeriesOpts{
		Interval: 300,
		Max:      100,
	}
	now := time.Now()
	five := now.Add(5 * time.Minute)
	ten := now.Add(10 * time.Minute)
	hl2 := NewOHLCProp(OHLCPropClose)
	chgdiffname := "change-diff"
	chgdiff := NewChange(hl2, 1, nil)
	chgopts := &ChangeOpts{
		DiffType: ChangeDiffTypeRatio,
	}
	chgratio := NewChange(hl2, 1, chgopts)
	chgrationame := "change-ratio"
	data := []OHLCV{
		OHLCV{
			C: 3.1,
			S: now,
		},
		OHLCV{
			C: 2.8,
			S: five,
		},
		OHLCV{
			C: 3.4,
			S: ten,
		},
	}
	chgd1 := 2.8 - 3.1
	chgr1 := 2.8 / 3.1
	chgd2 := 3.4 - 2.8
	chgr2 := 3.4 / 2.8
	io := []struct {
		time   time.Time
		output map[string]*float64
	}{
		{
			time:   now,
			output: nil,
		},
		{
			time: five,
			output: map[string]*float64{
				chgdiffname:  &chgd1,
				chgrationame: &chgr1,
			},
		},
		{
			time: ten,
			output: map[string]*float64{
				chgdiffname:  &chgd2,
				chgrationame: &chgr2,
			},
		},
	}
	s, err := NewSeries(data, opts)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error init series"))
	}
	s.AddIndicator(chgrationame, chgratio)
	s.AddIndicator(chgdiffname, chgdiff)

	for i, o := range io {
		v := s.GetValueForInterval(o.time)
		for _, name := range []string{chgrationame, chgdiffname} {
			if v.Indicators[name] == nil && o.output == nil {
				// ok
				continue
			}
			if v.Indicators[name] == nil || o.output == nil {
				t.Fatalf("expected both to be non nil but got %+v vs %+v at idx: %d", v.Indicators[name], o.output, i)
			}
			if *(v.Indicators[name]) != *(o.output[name]) {
				t.Errorf("expected: %+v but got %+v for idx: %d", *(o.output[name]), *(v.Indicators[name]), i)
			}
		}
	}
}
