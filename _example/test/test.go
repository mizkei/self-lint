package test

import (
	"testing"
	"time"
)

func IntPt(i int) *int          { return &i }
func StringPt(s string) *string { return &s }
func ParseDate(t *testing.T, v string) time.Time {
	t.Helper()
	tm, err := time.Parse("2006-01-02 15:04:05", v)
	if err != nil {
		t.Fatal(err)
	}
	return tm
}
