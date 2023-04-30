package pine

import (
	"fmt"
	"math"
)

// RSI generates a ValueSeries of relative strength index
//
// The formula for RSI is
//   - u = Count the number of p(t+1) - p(t) > 0 as gains
//   - d = Count the number of p(t+1) - p(t) < 0 as losses
//   - rs = ta.rma(u) / ta.rma(d)
//   - res = 100 - 100 / (1 + rs)
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
			prevv := v.prev
			// previous value exists
			if prevv != nil {
				prevr := rsiu.Get(prevv.t)

				// previous rsiu exists
				if prevr != nil {
					prevFirstVal := prevv
					removelb := 1
					for i := 1; i < int(l)+1; i++ {
						if prevFirstVal.prev == nil {
							break
						}
						removelb++
						prevFirstVal = prevFirstVal.prev
					}

					// was able to find previous value
					if int64(removelb) == l+1 {
						toAdd := math.Max(v.v-v.prev.v, 0)
						remval := math.Max(prevFirstVal.next.v-prevFirstVal.v, 0)
						newrsiu := prevr.v - remval + toAdd
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
			prevv := v.prev
			// previous value exists
			if prevv != nil {
				prevr := rsid.Get(prevv.t)

				// previous rsiu exists
				if prevr != nil {
					prevFirstVal := prevv
					removelb := 1
					for i := 1; i < int(l)+1; i++ {
						if prevFirstVal.prev == nil {
							break
						}
						removelb++
						prevFirstVal = prevFirstVal.prev
					}

					// was able to find previous value
					if int64(removelb) == l+1 {
						toAdd := math.Max(v.prev.v-v.v, 0)
						remval := math.Max(prevFirstVal.v-prevFirstVal.next.v, 0)
						newrsiu := prevr.v - remval + toAdd
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

	setCache(rsiukey, rsiu)
	setCache(rsidkey, rsid)

	rs := Div(rsiu, rsid)
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

	rmau := RMA(rsiu, l)
	rmad := RMA(rsid, l)

	rmadiv := Div(rmau, rmad)

	if rmadiv.GetCurrent() == nil {
		return rsi
	}

	hundred := ReplaceAll(vs, 100)

	res1 := Div(hundred, AddConst(rmadiv, 1.0))
	res2 := Sub(hundred, res1)

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
