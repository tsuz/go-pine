package pine

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
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
