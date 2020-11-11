package data

import (
	"fmt"
	"strconv"
)

// A Value represents some kind of number.
// A Value can be Added to another Value.
// A Value can be Less than another Value.
// A Value has a String representation.
//
// Generics would have been nice for this.
type Value interface {
	Add(Value) Value
	Less(Value) bool
	String() string
}

// Zero returns the zero value of the same
// underlying type as v.
func Zero(reference Value) Value {
	switch reference.(type) {
	case *Int:
		return NewInt(0)
	case *Float:
		return NewFloat(0)
	default:
		panic("reference value must be Int or Float")
	}
}

// Int wraps int.
type Int struct {
	i int
}

// NewInt creates a new Int of value i.
func NewInt(i int) *Int {
	return &Int{i: i}
}

// Add returns the sum of i and v.
func (i *Int) Add(v Value) Value {
	switch v := v.(type) {
	case *Int:
		a := i.i
		return NewInt(a + v.i)
	default:
		panic("cannot add Int and non-Int")
	}
}

// Less returns true if i is less than v.
func (i *Int) Less(v Value) bool {
	switch v := v.(type) {
	case *Int:
		a := i.i
		return a < v.i
	default:
		panic("cannot compare Int and non-Int")
	}
}

func (i *Int) String() string {
	return strconv.Itoa(i.i)
}

// Float wraps a float64.
type Float struct {
	f float64
}

// NewFloat creates a Float of value f.
func NewFloat(f float64) *Float {
	return &Float{f: f}
}

// Add returns the sum of f and v.
func (f *Float) Add(v Value) Value {
	switch v := v.(type) {
	case *Float:
		a := f.f
		return NewFloat(a + v.f)
	default:
		panic("cannot add Float and non-Float")
	}
}

// Less returns true if f is less than v.
func (f *Float) Less(v Value) bool {
	switch v := v.(type) {
	case *Float:
		a := f.f
		return a < v.f
	default:
		panic("cannot compare Float and non-Float")
	}
}

func (f *Float) String() string {
	return fmt.Sprintf("%.3f", f.f)
}
