package utils

import "testing"

func TestUUID(t *testing.T) {
	uuid := UUID()
	if len(uuid) <= 0 {
		t.Fatalf("Generate failed")
	}
	t.Logf("uuid=%s\n", uuid)
}
