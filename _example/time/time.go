package time

import (
	stime "time"
)

var (
	defaultNow = stime.Now
	now        = defaultNow
)

func Now() stime.Time {
	return now()
}

func SetTime(t stime.Time) {
	now = func() stime.Time { return t }
}

func ResetTime() {
	now = defaultNow
}
