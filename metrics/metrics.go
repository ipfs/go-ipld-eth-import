package metrics

import (
	"time"
)

/*
Package metrics wants to help us to know what's going on.

We'll try to keep it simple.
*/

// Data has all the metrics data in memory. It has counters and loggers.
// The formers can only be incremented or decreased, while the latter
// can used to get time differences.
type Data struct {
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
	data.counters = make(map[string]int)
	data.loggers = make(map[string][]int64)
}

/*
 COUNTERS
*/

// NewCounter returns a counter with the given key.
func NewCounter(key string) {
	if _, ok := data.counters[key]; !ok {
		data.counters[key] = 0
	}
}

// IncCounter increments the given counter by 1.
func IncCounter(key string) {
	if _, ok := data.counters[key]; ok {
		data.counters[key]++
	}
}

// GetCounter returns the current value of the given counter.
func GetCounter(key string) int {
	if val, ok := data.counters[key]; ok {
		return val
	}
	return 0
}

/*
  LOGGERS
*/

// NewLogger returns a logger.
func NewLogger(key string) {
	if _, ok := data.loggers[key]; !ok {
		var slice []int64
		data.loggers[key] = slice
	}
}

// AddLog adds an int64 value to the logger. Useful for
// aggregations, such as the total number of bytes stored.
func AddLog(key string, val int64) {
	if _, ok := data.loggers[key]; ok {
		data.loggers[key] = append(data.loggers[key], val)
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

// StopLogDiff completed the functionality documented by StartLogDiff.
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
				n++
			}
		}

		return n, sum, float64(sum) / float64(n)
	}
	return 0, 0, 0
}
