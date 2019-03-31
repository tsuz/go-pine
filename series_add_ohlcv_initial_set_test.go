package pine

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

func initDataset() []OHLCV {
	f, err := os.Open("./test_data.csv")
	if err != nil {
		panic(err)
	}
	r := csv.NewReader(f)

	data := make([]OHLCV, 0)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		o, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}
		h, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatal(err)
		}
		l, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Fatal(err)
		}
		c, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}
		v, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			log.Fatal(err)
		}
		s, err := time.Parse("2006-01-02 15:04:05", record[5])
		if err != nil {
			log.Fatal(err)
		}
		p := OHLCV{
			O: o,
			H: h,
			L: l,
			C: c,
			V: v,
			S: s.Add(24 * time.Hour),
		}
		data = append(data, p)
		fmt.Println(p)
	}
	return data
}

func TestSeriesAddOHLCVFromEmptyData(t *testing.T) {
	opts := SeriesOpts{
		Interval: 300,
		Max:      100,
	}
	data := initDataset()
	_, err := NewSeries(data, opts)
	if err != nil {
		t.Fatal(err)
	}
	s, err := NewSeries(data, opts)
	if err != nil {
		t.Fatal(err)
	}

	v := s.GetValueForInterval(data[len(data)-1].S)
	if v == nil {
		t.Fatal("expected v to be non nil but got nil")
	}
	h := v.OHLCV
	if h.O != 0.00000063 {
		t.Fatalf("expected new open to be 0.00000063 but got %+v", h.O)
	} else if h.H != 0.00000063 {
		t.Fatalf("expected new high to be 0.00000063 but got %+v", h.H)
	} else if h.L != 0.00000062 {
		t.Fatalf("expected new low to be 0.00000062 but got %+v", h.L)
	} else if h.V != 1971083 {
		t.Fatalf("expected vol to be 1971083 but got %+v", h.V)
	} else if h.C != 0.00000063 {
		t.Fatalf("expected close to be 0.00000063 but got %+v", h.C)
	} else if !h.S.Equal(data[len(data)-1].S) {
		t.Fatalf("expected time to be %+v but got %+v", h.S, data[len(data)-1].S)
	}
}
