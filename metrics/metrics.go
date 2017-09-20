package metrics

import (
	"time"
)

/*
Package metrics wants to help us to know what's going on.

We'll try to keep it simple.
*/

type Data struct {
	// Basically, we have a slice of time signatures per key,
	// with a ClickTimer("my-timer") function to add them
	// and a set of APIs to give you time differences and whatnot.
	timers map[string][]int64

	counters map[string]int
}

// The global variable here
var data Data

func init() {
	data = Data{}
	data.timers = make(map[string][]int64)
	data.counters = make(map[string]int)
}

/*
  TIMERS
*/

func NewTimer(key string) {
	if _, ok := data.timers[key]; !ok {
		data.timers[key] = []int64{time.Now().UnixNano()}
	}
}

func ClickTimer(key string) {
	if _, ok := data.timers[key]; ok {
		data.timers[key] = append(data.timers[key], time.Now().UnixNano())
	}
}

func GetTotalDiffTimer(key string) int64 {
	if _, ok := data.timers[key]; ok {
		var start int
		l := len(data.timers[key])
		switch {
		case l == 1:
			return 0
		case l == 2:
			start = 0
		default:
			start = 1
		}

		return (data.timers[key][l-1] -
			data.timers[key][start]) / (1000 * 1000)
	}
	return 0
}

func GetAverageDiffTimer(key string) (int, int64, float64) {
	if _, ok := data.timers[key]; ok {
		n := len(data.timers[key])
		sum := int64(0)

		for i := 1; i < n; i++ {
			sum += data.timers[key][i] - data.timers[key][i-1]
		}

		return n, sum, float64(sum) / float64(n)
	}
	return 0, 0, 0
}

func GetCountTimer(key string) int {
	if _, ok := data.timers[key]; ok {
		return len(data.timers[key]) - 1
	}
	return 0
}

/*
 COUNTERS
*/

func NewCounter(key string) {
	if _, ok := data.counters[key]; !ok {
		data.counters[key] = 0
	}
}

func IncCounter(key string) {
	if _, ok := data.counters[key]; ok {
		data.counters[key] += 1
	}
}

func GetCounter(key string) int {
	if val, ok := data.counters[key]; ok {
		return val
	}
	return 0
}
