package pine

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// RSI generates a ValueSeries of relative strength index
// the variable rsi=ValueSeries is the relative strength values of p=ValueSeries
// This ValueSeries guarantees to contain values up to p.GetCurrent()
//
// The formula for RSI is
// u = Count the number of p(t+1) - p(t) > 0 as gains
// d = Count the number of p(t+1) - p(t) < 0 as losses
// rs = ta.rma(u) / ta.rma(d)
// res = 100 - 100 / (1 + rs)
// Using the above formula, the below example illustrates what EMA values look like
//
// t=time.Time (no iteration) | 1   |  2  | 3    | 4          | 5  |
// p=ValueSeries              | 13  | 15  | 11   | 18         | 20 |
// u(close, 2)                | nil | nil |  2   | 7          | 9  |
// d(close, 2)                | nil | nil |  4   | 4          | 0  |
// rma(u(close,2), 2)		  | nil | nil |  1   | 4.5        | 8  |
// rma(d(close,2), 2)		  | nil | nil |  2 	 | 2          | 2  |
// rsi(close, 2)			  | nil | nil | 33.33| 69.2307692 | 20 |
func RSI(p ValueSeries, l int64) (ValueSeries, error) {
	key := fmt.Sprintf("rsi:%s:%d", p.ID(), l)
	rsi := getCache(key)
	if rsi == nil {
		rsi = NewValueSeries()
	}

	if p == nil || p.GetCurrent() == nil {
		return rsi, nil
	}

	// current available value
	stop := p.GetCurrent()

	rsi = getRSI(stop, p, rsi, l)

	setCache(key, rsi)

	rsi.SetCurrent(stop.t)

	return rsi, nil
}

// getRSIU generates sum gains
func getRSIU(stop *Value, vs ValueSeries, rsiu ValueSeries, l int64) ValueSeries {
	firstVal := rsiu.GetLast()

	if firstVal == nil {
		firstVal = vs.GetFirst()
	}

	if firstVal == nil {
		// if nothing is available, then nothing can be done
		return rsiu
	}

	itervt := firstVal.t

	var fseek int64
	var ftot float64

	for {
		log.Printf("RSIU evaluating t:%+v", itervt)
		v := vs.Get(itervt)
		if v == nil {
			break
		}
		e := rsiu.Get(itervt)
		if e != nil && v.next == nil {
			break
		}
		if e != nil {
			itervt = v.next.t
			continue
		}

		// get previous value
		if v.prev != nil {
			log.Printf("RSIU evaluating t:%+v prevval exists", itervt)
			prevv := v.prev
			// previous value exists
			if prevv != nil {
				log.Printf("RSIU evaluating t:%+v prevval received", itervt)
				prevr := rsiu.Get(prevv.t)

				// previous rsiu exists
				if prevr != nil {
					log.Printf("RSIU evaluating t:%+v prevval rsiu exist", itervt)
					prevFirstVal := prevv
					removelb := 1
					for i := 1; i < int(l)+1; i++ {
						if prevFirstVal.prev == nil {
							break
						}
						removelb++
						prevFirstVal = prevFirstVal.prev
					}

					log.Printf("RSIU evaluating t:%+v removelb: %d, prevFirstVal: %+v", itervt, removelb, prevFirstVal.v)
					// was able to find previous value
					if int64(removelb) == l+1 {
						toAdd := math.Max(v.v-v.prev.v, 0)
						remval := math.Max(prevFirstVal.next.v-prevFirstVal.v, 0)
						newrsiu := prevr.v - remval + toAdd
						log.Printf("RSIU evaluating t:%+v set after %+v", v.t, newrsiu)
						rsiu.Set(v.t, newrsiu)
						continue
					}
				}
			}
		}

		// previous rsiu does not exist. just keep adding until multiplication is required
		fseek++
		if v.prev != nil {
			ftot = ftot + math.Max(v.v-v.prev.v, 0)
		}

		if fseek == l+1 {
			log.Printf("RSIU evaluating t:%+v set first %+v", itervt, ftot)
			rsiu.Set(v.t, ftot)
		}

		if v.next == nil {
			break
		}
		if v.t.Equal(stop.t) {
			break
		}
		itervt = v.next.t
	}

	return rsiu
}

// getRSID generates sum gains
func getRSID(stop *Value, vs ValueSeries, rsid ValueSeries, l int64) ValueSeries {
	firstVal := rsid.GetLast()

	if firstVal == nil {
		firstVal = vs.GetFirst()
	}

	if firstVal == nil {
		// if nothing is available, then nothing can be done
		return rsid
	}

	itervt := firstVal.t

	var fseek int64
	var ftot float64

	for {
		log.Printf("RSIU evaluating t:%+v", itervt)
		v := vs.Get(itervt)
		if v == nil {
			break
		}
		e := rsid.Get(itervt)
		if e != nil && v.next == nil {
			break
		}
		if e != nil {
			itervt = v.next.t
			continue
		}

		// get previous value
		if v.prev != nil {
			log.Printf("RSIU evaluating t:%+v prevval exists", itervt)
			prevv := v.prev
			// previous value exists
			if prevv != nil {
				log.Printf("RSIU evaluating t:%+v prevval received", itervt)
				prevr := rsid.Get(prevv.t)

				// previous rsiu exists
				if prevr != nil {
					log.Printf("RSIU evaluating t:%+v prevval rsiu exist", itervt)
					prevFirstVal := prevv
					removelb := 1
					for i := 1; i < int(l)+1; i++ {
						if prevFirstVal.prev == nil {
							break
						}
						removelb++
						prevFirstVal = prevFirstVal.prev
					}

					log.Printf("RSIU evaluating t:%+v removelb: %d, prevFirstVal: %+v", itervt, removelb, prevFirstVal.v)
					// was able to find previous value
					if int64(removelb) == l+1 {
						toAdd := math.Max(v.prev.v-v.v, 0)
						remval := math.Max(prevFirstVal.v-prevFirstVal.next.v, 0)
						newrsiu := prevr.v - remval + toAdd
						log.Printf("RSIU evaluating t:%+v set after %+v", v.t, newrsiu)
						rsid.Set(v.t, newrsiu)
						continue
					}
				}
			}
		}

		// previous rsiu does not exist. just keep adding until multiplication is required
		fseek++
		if v.prev != nil {
			ftot = ftot + math.Max(v.prev.v-v.v, 0)
		}

		if fseek == l+1 {
			log.Printf("RSIU evaluating t:%+v set first %+v", itervt, ftot)
			rsid.Set(v.t, ftot)
		}

		if v.next == nil {
			break
		}
		if v.t.Equal(stop.t) {
			break
		}
		itervt = v.next.t
	}

	return rsid
}

func getRSI(stop *Value, vs ValueSeries, rsi ValueSeries, l int64) ValueSeries {

	rsiukey := fmt.Sprintf("rsiu:%s:%d", vs.ID(), l)
	rsiu := getCache(rsiukey)
	if rsiu == nil {
		rsiu = NewValueSeries()
	}
	rsidkey := fmt.Sprintf("rsid:%s:%d", vs.ID(), l)
	rsid := getCache(rsidkey)
	if rsid == nil {
		rsid = NewValueSeries()
	}

	rsiu = getRSIU(stop, vs, rsiu, l)
	rsid = getRSID(stop, vs, rsid, l)

	rsiu.SetCurrent(stop.t)
	rsid.SetCurrent(stop.t)

	rs := rsiu.Div(rsid)
	rsn := rs.GetFirst()
	// set inifinity to 100
	for {
		if rsn == nil {
			break
		}
		if math.IsInf(rsn.v, 1) {
			rs.Set(rsn.t, 100) // set infinity to 100
		}
		rsn = rsn.next
	}

	rmau, err := RMA(rsiu, l)
	if err != nil {
		panic(errors.Wrap(err, "Error calling RMA"))
	}
	rmad, err := RMA(rsid, l)
	if err != nil {
		panic(errors.Wrap(err, "Error calling RMA"))
	}

	rmadiv := rmau.Div(rmad)

	if rmadiv.GetCurrent() == nil {
		return rsi
	}

	hundred1 := vs.Copy()
	hundred1.SetAll(100)
	hundred2 := vs.Copy()
	hundred2.SetAll(100)

	log.Printf("rmadiv  Val: %+v", rmadiv.GetFirst())
	res1 := hundred2.Div(rmadiv.AddConst(1.0))

	log.Printf("res1  Val: %+v", res1.GetCurrent())
	res2 := hundred1.Sub(res1)

	log.Printf("res2  Val: %+v", res2.GetCurrent())
	firstVal := rsi.GetLast()

	if firstVal == nil {
		firstVal = vs.GetFirst()
	}

	if firstVal == nil {
		// if nothing is available, then nothing can be done
		return rsi
	}

	itervt := firstVal.t

	for {
		v := vs.Get(itervt)
		v2 := res2.Get(itervt)
		if v == nil {
			break
		}
		e := rsi.Get(itervt)
		if e != nil && v.next == nil {
			break
		}
		if e != nil {
			itervt = v.next.t
			continue
		}
		if v2 != nil {
			log.Printf("RSI set t:%+v, v:%+v", v.t, v2.v)
			rsi.Set(v.t, v2.v)
		}

		if v.next == nil {
			break
		}
		if v.t.Equal(stop.t) {
			break
		}
		itervt = v.next.t
	}

	return rsi
}
