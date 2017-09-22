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

	// This is just a map of `+1` counters. You know.
	// How many iterations? How many cases of A? etc.
	counters map[string]int

	// Loggers have two use cases:
	// * Store time differences (with StartLogDiff / StopLogDiff)
	//   + Useful for RPC Calls and DB Queries
	// * Store series of values (mem / CPU / active goroutines / etc)
	loggers map[string][]int64
}

// The global variable here
var data Data

func init() {
	data = Data{}
	data.timers = make(map[string][]int64)
	data.counters = make(map[string]int)
	data.loggers = make(map[string][]int64)
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

/*
  LOGGERS
*/

func NewLogger(key string) {
	if _, ok := data.loggers[key]; !ok {
		var slice []int64
		data.loggers[key] = slice
	}
}

// StartLogDiff returns the index of the value logged,
// So you can get the time difference with it uwing StopLogDiff().
// It will store your time value as a negative number.
// Successive functions to get averages will ignore the negative values,
// deeming them as "incomplete logs".
func StartLogDiff(key string) int {
	if _, ok := data.loggers[key]; ok {
		data.loggers[key] = append(data.loggers[key], -1*time.Now().UnixNano())
		return len(data.loggers[key]) - 1
	}
	// No key found
	return -1
}

func StopLogDiff(key string, idx int) {
	if _, ok := data.loggers[key]; ok {
		if len(data.loggers[key]) > idx {
			// The value created at StartLogDiff is a negative one
			data.loggers[key][idx] = data.loggers[key][idx] + time.Now().UnixNano()
		}
	}
}

// GetAverageLogDiff will calculate the average of the log differences,
// discarding the negative ones, as those will be deemed as incomplete ops.
func GetAverageLogDiff(key string) (int, int64, float64) {
	if _, ok := data.loggers[key]; ok {
		n := 0
		sum := int64(0)

		for _, v := range data.loggers[key] {
			if v >= 0 {
				sum += v
				n += 1
			}
		}

		return n, sum, float64(sum) / float64(n)
	}
	return 0, 0, 0
}
