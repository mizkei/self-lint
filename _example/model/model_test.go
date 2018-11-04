package model_test

import (
	"testing"
	"time"

	"github.com/mizkei/self-lint/_example/model"
	"github.com/mizkei/self-lint/_example/test"
	mtime "github.com/mizkei/self-lint/_example/time"
)

func TestGreet(t *testing.T) {
	defer mtime.ResetTime()
	for name, tc := range map[string]struct {
		time   time.Time
		expect string
	}{
		"morning": {
			time:   test.ParseDate(t, "2018-11-04 11:59:59"),
			expect: "good morning",
		},
		"evening": {
			time:   test.ParseDate(t, "2018-11-04 18:00:00"),
			expect: "good evening",
		},
		"afternoon": {
			time:   test.ParseDate(t, "2018-11-04 17:59:59"),
			expect: "good afternoon",
		},
	} {
		t.Run(name, func(t *testing.T) {
			mtime.SetTime(tc.time)
			res := model.Greet()
			if res != tc.expect {
				t.Fatalf("fail, got:%s, want:%s", res, tc.expect)
			}
		})
	}
}
