package clock

import (
    "time"
)

// A Clock implements the go-fed/activity/pub/Clock interface. It just wraps a time.Time object, but is explicitly defined in case we want to get fancy in the future.
type Clock struct {
    t *time.Time
}

func New() *Clock {
    return &Clock{
        t: new(time.Time),
    }
}

func (c *Clock) Now() *time.Time {
    return c.t.Now()
}
