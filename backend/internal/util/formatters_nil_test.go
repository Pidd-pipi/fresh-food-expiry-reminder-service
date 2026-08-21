package util

import "testing"

func TestCountStringsSafe(t *testing.T) {
	counts := CountStrings([]string{"expiring", "expired", "expiring"})
	if counts["expiring"] != 2 {
		t.Fatalf("expiring = %d, want 2", counts["expiring"])
	}
	if counts["expired"] != 1 {
		t.Fatalf("expired = %d, want 1", counts["expired"])
	}
}

func TestGroupStringsSafe(t *testing.T) {
	groups := GroupStrings([]string{"a", "b", "a"})
	if len(groups) != 2 {
		t.Fatalf("groups = %d, want 2", len(groups))
	}
	if groups["a"]["a"] != 2 {
		t.Fatalf("groups[a][a] = %d, want 2", groups["a"]["a"])
	}
}
