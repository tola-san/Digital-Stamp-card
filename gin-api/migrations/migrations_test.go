package migrations

import "testing"

func TestAllUpMigrationsAreEmbedded(t *testing.T) {
	for _, name := range Up {
		script, err := ReadUp(name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		if len(script) == 0 {
			t.Fatalf("migration %s is empty", name)
		}
	}
}
