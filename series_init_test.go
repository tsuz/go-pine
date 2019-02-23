package pine_test

import (
	pine "go-pine"
	"testing"
)

func TestSeriesInit(t *testing.T) {
	opts := pine.SeriesOpts{
		Interval: 5,
		Max:      100,
	}
	_, err := pine.NewSeries(nil, opts)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSeriesInitWithError(t *testing.T) {
	badopts := []pine.SeriesOpts{
		pine.SeriesOpts{},
	}
	for i, opts := range badopts {
		_, err := pine.NewSeries(nil, opts)
		if err == nil {
			t.Fatalf("expected error but got none for index: %d", i)
		}
	}
}
