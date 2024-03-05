package rrr

import "time"

// TimestampNow() function returns a the current timestamp in nanoseconds.
func TimestampNow() int64 {
	return time.Now().UnixNano()
}
