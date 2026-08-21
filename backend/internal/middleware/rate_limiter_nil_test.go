package middleware

import "testing"

func TestRateLimiterSafe(t *testing.T) {
	rl := NewRateLimiter(60)
	if !rl.allow("u1") {
		t.Fatal("first allow should succeed")
	}
	if !rl.allow("u1") {
		t.Fatal("second allow should succeed")
	}
}
