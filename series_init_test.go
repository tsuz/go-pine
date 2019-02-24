package pine

import (
	"testing"
)

func TestSeriesInit(t *testing.T) {
	opts := SeriesOpts{
		Interval: 5,
		Max:      100,
	}
	_, err := NewSeries(nil, opts)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSeriesInitWithError(t *testing.T) {
	badopts := []SeriesOpts{
		SeriesOpts{},
	}
	for i, opts := range badopts {
		_, err := NewSeries(nil, opts)
		if err == nil {
			t.Fatalf("expected error but got none for index: %d", i)
		}
	}
}
