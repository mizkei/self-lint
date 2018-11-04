package model

import (
	"github.com/mizkei/self-lint/_example/test"
	"github.com/mizkei/self-lint/_example/time"
)

func Greet() string {
	now := time.Now()

	switch {
	case now.Hour() < 12:
		return "good morning"
	case now.Hour() < 18:
		return "good afternoon"
	default:
		return "good evening"
	}
}

func GreetIllegal() string {
	time.SetTime(test.ParseDate(nil, "2000-01-01 00:00:00"))
	return Greet()
}
