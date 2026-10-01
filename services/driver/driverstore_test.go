package driver

import (
	"testing"
)

func TestFormatDriverIDAddsDriverPlatePrefix(t *testing.T) {
	t.Setenv("DRIVER_ID_PREFIX", "test-")
	t.Setenv("DRIVER_ID_SUFFIX", "-tail")

	got := formatDriverID(42)
	want := "test-sd-driver-T700042C-tail"
	if got != want {
		t.Fatalf("formatDriverID() = %q, want %q", got, want)
	}
}
