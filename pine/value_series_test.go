package pine

import (
	"math"
	"testing"
	"time"
)

func TestValueSeriesAdd(t *testing.T) {
	a := NewValueSeries()
	now := time.Now()
	a.Set(now, 1)
	a.Set(now.Add(time.Duration(1000*1e6)), 2)
	a.Set(now.Add(time.Duration(2000*1e6)), 4) // this doesn't exist in b

	b := NewValueSeries()
	b.Set(now, 4)
	b.Set(now.Add(time.Duration(1000*1e6)), 4)

	c := Add(a, b)
	c.SetCurrent(now)
	f := c.GetCurrent()
	if f == nil {
		t.Fatalf("expected to be non nil but got nil")
	}
	if f.v != 5 {
		t.Errorf("expected %+v but got %+v", 5, f.v)
	}
	if f.next.v != 6 {
		t.Errorf("expected %+v but got %+v", 6, f.v)
	}
	if f.next.next != nil {
		t.Errorf("expected nil but got %+v", f.next.next.v)
	}

	// current time is passed on
	a.SetCurrent(now.Add(time.Duration(1000 * 1e6)))
	d := Add(a, b)
	if *d.Val() != 6 {
		t.Errorf("expected 6 but got %+v", *d.Val())
	}
}

func TestValueSeriesAddConst(t *testing.T) {
	a := NewValueSeries()
	now := time.Now()
	t2 := now.Add(time.Duration(1000 * 1e6))
	a.Set(now, 1)
	a.Set(t2, 2)

	b := AddConst(a, 3)
	f := b.GetFirst()
	if f == nil {
		t.Fatalf("expected to be non nil but got nil")
	}
	if f.v != 4 {
		t.Errorf("expected %+v but got %+v", 4, f.v)
	}
	if f.next.v != 5 {
		t.Errorf("expected %+v but got %+v", 5, f.v)
	}

	// current time is passed on
	a.SetCurrent(t2)
	d := AddConst(a, 3)
	if *d.Val() != 5 {
		t.Errorf("expected 5 but got %+v", *d.Val())
	}
}

func TestValueSeriesOperator(t *testing.T) {
	a := NewValueSeries()
	now := time.Now()
	a.Set(now, 1)
	a.Set(now.Add(time.Duration(1000*1e6)), 2)
	a.Set(now.Add(time.Duration(2000*1e6)), 3)

	c := Operate(a, a, "testvalueseriesoperator", func(b, c float64) float64 {
		return math.Mod(b, 2)
	})

	f := c.GetFirst()
	if f == nil {
		t.Fatalf("expected to be non nil but got nil")
	}
	if f.v != 1 {
		t.Errorf("expected %+v but got %+v", 0, f.v)
	}
	if f.next.v != 0 {
		t.Errorf("expected %+v but got %+v", 0, f.next.v)
	}
	if f.next.next.v != 1 {
		t.Errorf("expected %+v but got %+v", 1, f.next.next.v)
	}
}

func TestValueSeriesOperatorWithNil(t *testing.T) {
	a := NewValueSeries()
	b := NewValueSeries()
	t1 := time.Now()
	t2 := t1.Add(time.Duration(1000 * 1e6))
	t3 := t2.Add(time.Duration(1000 * 1e6))
	a.Set(t1, 1)

	b.Set(t1, 1)
	b.Set(t2, 2)
	b.Set(t3, 3)

	c := OperateWithNil(b, a, "testoperatewithnil", func(bvalue, avalue *Value) *Value {
		if avalue == nil {
			return &Value{
				t: bvalue.t,
				v: 0,
			}
		}
		return &Value{
			t: avalue.t,
			v: avalue.v + bvalue.v,
		}
	})

	f := c.GetFirst()
	if f == nil {
		t.Fatalf("expected to be non nil but got nil")
	}

	if f.v != 2 {
		t.Errorf("expected %+v but got %+v", 2, f.v)
	}
	if f.next.v != 0 {
		t.Errorf("expected %+v but got %+v", 0, f.next.v)
	}
	if f.next.next.v != 0 {
		t.Errorf("expected %+v but got %+v", 0, f.next.next.v)
	}
}

func TestValueSeriesDiv(t *testing.T) {
	a := NewValueSeries()
	now := time.Now()
	a.Set(now, 1)
	a.Set(now.Add(time.Duration(1000*1e6)), 2)
	a.Set(now.Add(time.Duration(2000*1e6)), 3)

	b := NewValueSeries()
	b.Set(now, 4)
	b.Set(now.Add(time.Duration(1000*1e6)), 4)

	c := Div(a, b)
	c.SetCurrent(now)
	f := c.GetCurrent()
	if f == nil {
		t.Fatalf("expected to be non nil but got nil")
	}
	if f.v != 0.25 {
		t.Errorf("expected %+v but got %+v", 0.25, f.v)
	}
	if f.next.v != 0.5 {
		t.Errorf("expected %+v but got %+v", 0.5, f.v)
	}
	if f.next.next != nil {
		t.Errorf("expected nil but got %+v", f.next.next.v)
	}

	// current time is passed on
	a.SetCurrent(now.Add(time.Duration(1000 * 1e6)))
	d := Div(a, b)
	if *d.Val() != 0.5 {
		t.Errorf("expected .5 but got %+v", *d.Val())
	}
}

func TestValueSeriesDivConst(t *testing.T) {
	a := NewValueSeries()
	now := time.Now()
	a.Set(now, 1)
	a.Set(now.Add(time.Duration(1000*1e6)), 2)

	b := DivConst(a, 4)
	f := b.GetFirst()
	if f == nil {
		t.Fatalf("expected to be non nil but got nil")
	}
	if f.v != 0.25 {
		t.Errorf("expected %+v but got %+v", 0.25, f.v)
	}
	if f.next.v != 0.5 {
		t.Errorf("expected %+v but got %+v", 0.5, f.v)
	}

	// current time is passed on
	a.SetCurrent(now.Add(time.Duration(1000 * 1e6)))
	d := DivConst(a, 4)
	if *d.Val() != 0.5 {
		t.Errorf("expected .5 but got %+v", *d.Val())
	}
}

// TestValueSeriesSetMaxResize tests when set max is called after data is populated
func TestValueSeriesSetMaxResize(t *testing.T) {
	a := NewValueSeries()
	now := time.Now()
	t1 := time.Now()
	t2 := now.Add(time.Duration(1000 * 1e6))
	t3 := now.Add(time.Duration(2000 * 1e6))
	a.Set(t1, 1)
	a.Set(t2, 2)
	a.Set(t3, 4) // this doesn't exist in b
	a.SetMax(2)

	v1 := a.Get(t1)
	if v1 != nil {
		t.Errorf("expected to be nil but got %+v", v1.v)
	}
	v1 = a.Get(t2)
	if v1.v != 2 {
		t.Errorf("expected to be 2 but got %+v", v1.v)
	}
	v1 = a.Get(t3)
	if v1.v != 4 {
		t.Errorf("expected to be 4 but got %+v", v1.v)
	}
	if a.Len() != 2 {
		t.Errorf("expected to be 2 but got %+v", a.Len())
	}
}

// TestValueSeriesSetMaxPushResize tests when max is set and then data is populated
func TestValueSeriesSetMaxPushResize(t *testing.T) {
	a := NewValueSeries()
	a.SetMax(2)
	now := time.Now()
	t1 := time.Now()
	t2 := now.Add(time.Duration(1000 * 1e6))
	t3 := now.Add(time.Duration(2000 * 1e6))
	a.Set(t1, 1)
	a.Set(t2, 2)
	a.Set(t3, 4) // this doesn't exist in b

	v1 := a.Get(t1)
	if v1 != nil {
		t.Errorf("expected to be nil but got %+v", v1.v)
	}
	v1 = a.Get(t2)
	if v1.v != 2 {
		t.Errorf("expected to be 2 but got %+v", v1.v)
	}
	v1 = a.Get(t3)
	if v1.v != 4 {
		t.Errorf("expected to be 4 but got %+v", v1.v)
	}
	if a.Len() != 2 {
		t.Errorf("expected to be 2 but got %+v", a.Len())
	}
}

func TestValueSeriesMul(t *testing.T) {
	a := NewValueSeries()
	now := time.Now()
	a.Set(now, 1)
	a.Set(now.Add(time.Duration(1000*1e6)), 2)
	a.Set(now.Add(time.Duration(2000*1e6)), 3) // this doesn't exist in b

	b := NewValueSeries()
	b.Set(now, 4)
	b.Set(now.Add(time.Duration(1000*1e6)), 4)

	c := Mul(a, b)
	c.SetCurrent(now)
	f := c.GetCurrent()
	if f == nil {
		t.Fatalf("expected to be non nil but got nil")
	}
	if f.v != 4 {
		t.Errorf("expected %+v but got %+v", 4, f.v)
	}
	if f.next.v != 8 {
		t.Errorf("expected %+v but got %+v", 8, f.v)
	}
	if f.next.next != nil {
		t.Errorf("expected nil but got %+v", f.next.next.v)
	}

	// current time is passed on
	a.SetCurrent(now.Add(time.Duration(1000 * 1e6)))
	d := Mul(a, b)
	if *d.Val() != 8 {
		t.Errorf("expected 8 but got %+v", *d.Val())
	}
}

func TestValueSeriesMulConst(t *testing.T) {
	a := NewValueSeries()
	now := time.Now()
	a.Set(now, 1)
	a.Set(now.Add(time.Duration(1000*1e6)), 2)

	b := MulConst(a, 3)
	f := b.GetFirst()
	if f == nil {
		t.Fatalf("expected to be non nil but got nil")
	}
	if f.v != 3 {
		t.Errorf("expected %+v but got %+v", 3, f.v)
	}
	if f.next.v != 6 {
		t.Errorf("expected %+v but got %+v", 6, f.v)
	}

	// current time is passed on
	a.SetCurrent(now.Add(time.Duration(1000 * 1e6)))
	d := MulConst(a, 3)
	if *d.Val() != 6.0 {
		t.Errorf("expected 6 but got %+v", *d.Val())
	}
}

func TestValueSeriesSub(t *testing.T) {
	a := NewValueSeries()
	now := time.Now()
	nilTime := now.Add(time.Duration(3000 * 1e6))
	a.Set(now, 1)
	a.Set(now.Add(time.Duration(1000*1e6)), 2)
	a.Set(now.Add(time.Duration(2000*1e6)), 3)
	a.Set(nilTime, 4)

	b := NewValueSeries()
	b.Set(now, 4)
	b.Set(now.Add(time.Duration(1000*1e6)), 4)
	b.Set(now.Add(time.Duration(2000*1e6)), 1)

	c := Sub(a, b)
	c.SetCurrent(now)
	f := c.GetCurrent()
	if f == nil {
		t.Fatalf("expected to be non nil but got nil")
	}
	if f.v != -3 {
		t.Errorf("expected %+v but got %+v", -3, f.v)
	}
	if f.next.v != -2 {
		t.Errorf("expected %+v but got %+v", -2, f.next.v)
	}
	if f.next.next.v != 2 {
		t.Errorf("expected %+v but got %+v", 2, f.next.next.v)
	}
	n := c.Get(nilTime)
	if n != nil {
		t.Errorf("expected nil but got %+v", n.v)
	}

	// current time is passed on
	a.SetCurrent(now.Add(time.Duration(1000 * 1e6)))
	d := Sub(a, b)
	if *d.Val() != -2 {
		t.Errorf("expected -2 but got %+v", *d.Val())
	}
}

func TestValueSeriesSubConst(t *testing.T) {
	a := NewValueSeries()
	now := time.Now()
	a.Set(now, 1)
	a.Set(now.Add(time.Duration(1000*1e6)), 2)

	b := SubConst(a, 3)
	f := b.GetFirst()
	if f == nil {
		t.Fatalf("expected to be non nil but got nil")
	}
	if f.v != -2 {
		t.Errorf("expected %+v but got %+v", -2, f.v)
	}
	if f.next.v != -1 {
		t.Errorf("expected %+v but got %+v", -1, f.v)
	}

	// current time is passed on
	a.SetCurrent(now.Add(time.Duration(1000 * 1e6)))
	d := SubConst(a, 3)
	if *d.Val() != -1 {
		t.Errorf("expected -1 but got %+v", *d.Val())
	}
}

func TestValueSeriesGetFirst(t *testing.T) {

	s := NewValueSeries()
	now := time.Now()
	s.Set(now, 1)
	s.Set(now.Add(time.Duration(1000*1e6)), 2)
	s.SetCurrent(now)
	f := s.GetFirst()
	if f == nil {
		t.Errorf("expected to be non nil but got nil")
	}
	if f.next == nil {
		t.Errorf("expected next to be non nil but got nil")
	}
	if f.next.v != 2 {
		t.Errorf("expected next value to be 2 but got  %+v", f.next.v)
	}
}

func TestMemoryLeakArithmetic(t *testing.T) {
	v := 4.2351

	testMemoryLeak(t, func(o OHLCVSeries) error {
		c := OHLCVAttr(o, OHLCPropClose)
		op := OHLCVAttr(o, OHLCPropOpen)
		s1 := Add(c, op)
		s2 := AddConst(s1, v)
		s3 := Sub(s2, c)
		s4 := SubConst(s3, v)
		s5 := Mul(s4, s2)
		s6 := MulConst(s5, v)
		s7 := Div(s6, c)
		DivConst(s7, v)

		return nil
	})
}
