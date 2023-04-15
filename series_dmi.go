package pine

// DMI generates a ValueSeries of directional movement index.
// the variable ema=ValueSeries is the exponentially weighted moving average values of p=ValueSeries
// ema may be behind where they should be with regards to p.GetCurrent()
// while ema catches up to where p.GetCurrent() is, the series should also contain
// all available average values between the last and up to p.GetCurrent()
//
// The formula for EMA is EMA=(closing price − previous day’s EMA)× smoothing constant as a decimal + previous day’s EMA
// where smoothing constant is 2 ÷ (number of time periods + 1)
// if the previous day's EMA is nil then it's the SMA of the lookback time.
// Using the above formula, the below example illustrates what EMA values look like
//
// t=time.Time (no iteration) | 1   |  2  | 3   | 4       |
// p=ValueSeries              | 13  | 15  | 17  | 18      |
// dmi(low, high, close, 1, 2)              | 13  | 15  | 17  | 18      |
// dmi(low, high, close, 1, 2)              | nil | 14  | 16  | 17.3333 |
// func DMI(l ValueSeries, h ValueSeries, c ValueSeries, len, smoo int64) (adx, dmip, dmim ValueSeries, err error) {
// 	key := fmt.Sprintf("dmi:%s:%d", c.ID(), l)
// 	dmi := getCache(key)
// 	if dmi == nil {
// 		dmi = NewValueSeries()
// 	}

// 	if c == nil || c.GetCurrent() == nil {
// 		return dmi, nil
// 	}

// 	// current available value
// 	stop := c.GetCurrent()
// 	stop2 := c.GetCurrent()
// 	stop3 := c.GetCurrent()

// 	if !stop.t.Equal(stop2.t) || !stop2.t.Equal(stop3.t) {
// 		return dmi, errors.New("Stop values must be equal for all three values series")
// 	}

// 	dmi = getDMI(stop, l, h, c, len, smoo)

// 	setCache(key, dmi)

// 	dmi.SetCurrent(stop.t)

// 	return dmi, nil
// }

// func getDMI(stop *Value, l ValueSeries, h ValueSeries, c ValueSeries, len, smoo int64) ValueSeries {

// 	var mul float64 = 2.0 / float64(len+1.0)
// 	firstVal := ema.GetLast()

// 	if firstVal == nil {
// 		firstVal = vs.GetFirst()
// 	}

// 	if firstVal == nil {
// 		// if nothing is available, then nothing can be done
// 		return ema
// 	}

// 	itervt := firstVal.t

// 	var fseek int64
// 	var ftot float64

// 	for {
// 		v := vs.Get(itervt)
// 		if v == nil {
// 			break
// 		}
// 		e := ema.Get(itervt)
// 		if e != nil && v.next == nil {
// 			break
// 		}
// 		if e != nil {
// 			itervt = v.next.t
// 			continue
// 		}

// 		// get previous ema
// 		if v.prev != nil {
// 			prevv := vs.Get(v.prev.t)
// 			preve := ema.Get(prevv.t)
// 			// previous ema exists, just do multiplication to that
// 			if preve != nil {
// 				nextEMA := (v.v-preve.v)*mul + preve.v
// 				ema.Set(v.t, nextEMA)
// 				continue
// 			}
// 		}

// 		// previous value does not exist. just keep adding until multplication is required
// 		fseek++
// 		ftot = ftot + v.v

// 		if fseek == l {
// 			avg := ftot / float64(fseek)
// 			ema.Set(v.t, avg)
// 		}

// 		if v.next == nil {
// 			break
// 		}
// 		if v.t.Equal(stop.t) {
// 			break
// 		}
// 		itervt = v.next.t
// 	}

// 	return ema
// }
