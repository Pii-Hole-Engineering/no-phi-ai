package rrr

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimestampNow(t *testing.T) {
	t.Parallel()

	start := time.Now().UnixNano()
	actual := TimestampNow()

	assert.IsType(t, int64(0), actual, "TimestampNow() should return a timestamp as an int64")
	assert.GreaterOrEqual(t, actual, start, "TimestampNow() should return a timestamp in nanoseconds greater than the time of the start of the test")
	assert.LessOrEqual(t, actual, time.Now().UnixNano(), "TimestampNow() should return a timestamp in nanoseconds less than the current time in nanoseconds")
}
