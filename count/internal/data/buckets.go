// Package data provides time series buckets of accumulated values.
//
// All of the complex logic in this package was copied from recounter
// which was copied from RecentEventsCounter which was derived from
// some notes etched on a mysterious stone tablet found the Himalayas.
package data

import (
	"bytes"
	"time"
)

const (
	separator = ", "
)

// Buckets contains a history of accumulating values over
// periodic intervals of time.
//
// Not threadsafe, bring your own lock.
type Buckets struct {
	buckets    []Value
	ticker     Ticker
	oldestIdx  int
	oldestTick int
}

// NewBuckets creates a new Buckets of size length, where each bucket contains
// values for spanning the duration specified by period.
func NewBuckets(period Ticker, length int, zero Value) *Buckets {
	return &Buckets{
		buckets:    fill(zero, length),
		ticker:     period,
		oldestIdx:  1,
		oldestTick: period() - length + 1,
	}
}

func fill(v Value, length int) []Value {
	values := make([]Value, length)
	for i := 0; i < length; i++ {
		values[i] = v
	}
	return values
}

// Increment the current bucket by delta.
func (bs *Buckets) Increment(delta Value) {
	tick := bs.ticker()

	newestTick := bs.oldestTick + len(bs.buckets) - 1
	if newestTick == tick {
		newestIdx := 0
		if bs.oldestIdx == 0 {
			newestIdx = len(bs.buckets) - 1
		} else {
			newestIdx = bs.oldestIdx - 1
		}
		bs.buckets[newestIdx] = bs.buckets[newestIdx].Add(delta)
		return
	}

	numToExpire := tick - newestTick
	if numToExpire >= len(bs.buckets) {
		// zero everything and start over
		bs.oldestTick = tick - len(bs.buckets) + 1
		for i := range bs.buckets {
			if i == 0 {
				bs.buckets[i] = delta
			} else {
				bs.buckets[i] = Zero(delta)
			}
		}
		bs.oldestIdx = 1
		return
	}

	// overwrite the gap with zeros
	for i := 0; i < numToExpire-1; i++ {
		bs.buckets[bs.oldestIdx] = Zero(delta)
		bs.oldestIdx++
		if bs.oldestIdx == len(bs.buckets) {
			bs.oldestIdx = 0
		}
	}
	bs.oldestTick += numToExpire

	bs.buckets[bs.oldestIdx] = delta
	bs.oldestIdx++
	if bs.oldestIdx == len(bs.buckets) {
		bs.oldestIdx = 0
	}
}

// Compare will apply f over the most recently completed
// bucket (not the currently active bucket) and v.
//
// In this example the most recently completed bucket
// contains the value of 3.
//
//	"[1, 3, 5, 2]"
func (bs *Buckets) Compare(f Cmp, v Value) bool {
	// force refresh
	bs.Increment(Zero(bs.buckets[0]))

	lastCompletedIdx := bs.oldestIdx - 2
	if lastCompletedIdx == -2 {
		lastCompletedIdx = len(bs.buckets) - 2
	} else if lastCompletedIdx == -1 {
		lastCompletedIdx = len(bs.buckets) - 1
	}

	return f(bs.buckets[lastCompletedIdx], v)
}

// Cmp is a predicate of a and b.
type Cmp func(a, b Value) bool

// LessEq returns true iff a <= b.
func LessEq(a, b Value) bool {
	return a.Less(b) || (!a.Less(b) && !b.Less(a))
}

// GreaterEq returns true iff a >= b
func GreaterEq(a, b Value) bool {
	return b.Less(a) || (!b.Less(a) && !a.Less(b))
}

// Sum all the accumulated values in all buckets of bs.
func (bs *Buckets) Sum() Value {
	sum := bs.buckets[0]
	for i := 1; i < len(bs.buckets); i++ {
		sum = sum.Add(bs.buckets[i])
	}
	return sum
}

func (bs *Buckets) String() string {
	var buf bytes.Buffer
	buf.WriteString("[")

	// force refresh
	bs.Increment(Zero(bs.buckets[0]))

	if bs.oldestIdx == 0 {
		for i := len(bs.buckets) - 1; i >= 0; i-- {
			buf.WriteString(bs.buckets[i].String())
			if i != 0 {
				buf.WriteString(separator)
			}
		}
	} else {
		for i := bs.oldestIdx - 1; i >= 0; i-- {
			buf.WriteString(bs.buckets[i].String())
			buf.WriteString(separator)
		}
		for i := len(bs.buckets) - 1; i >= bs.oldestIdx; i-- {
			buf.WriteString(bs.buckets[i].String())
			if i != bs.oldestIdx {
				buf.WriteString(separator)
			}
		}
	}

	buf.WriteString("]")
	return buf.String()
}

// A Ticker is used to determine which bucket is currently
// accumulating values. As time passes, the values in each
// bucket "slide" over, as the Ticker rolls over.
type Ticker func() int

// Convenience values of Ticker which are used most often.
var (
	OneSecond Ticker = func() int {
		return int(time.Now().UnixNano() / int64(time.Second))
	}

	OneMinute Ticker = func() int {
		return int(time.Now().UnixNano() / int64(time.Minute))
	}

	FiveMinutes Ticker = func() int {
		return int(time.Now().UnixNano() / int64(5*time.Minute))
	}

	FifteenMinutes Ticker = func() int {
		return int(time.Now().UnixNano() / int64(15*time.Minute))
	}

	OneHour Ticker = func() int {
		return int(time.Now().UnixNano() / int64(1*time.Hour))
	}
)
