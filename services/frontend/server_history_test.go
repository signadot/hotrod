package frontend

import (
	"testing"
)

func TestExtractDispatchedDriver(t *testing.T) {
	t.Parallel()

	requestID, driverPlate, ok := extractDispatchedDriver(
		"req-42-dispatched-driver",
		"Driver A-123 arriving in 2m3s",
	)
	if !ok {
		t.Fatalf("expected dispatched-driver notification to match")
	}
	if requestID != 42 {
		t.Fatalf("unexpected request id: got %d want 42", requestID)
	}
	if driverPlate != "A-123" {
		t.Fatalf("unexpected driver plate: got %q want %q", driverPlate, "A-123")
	}
}

func TestExtractDispatchedDriverInvalidID(t *testing.T) {
	t.Parallel()

	_, _, ok := extractDispatchedDriver(
		"req-abc-dispatched-driver",
		"Driver A-123 arriving in 2m3s",
	)
	if ok {
		t.Fatalf("expected invalid notification ID to fail parsing")
	}
}

func TestExtractDispatchedDriverInvalidBody(t *testing.T) {
	t.Parallel()

	_, _, ok := extractDispatchedDriver(
		"req-42-dispatched-driver",
		"Finding an available driver",
	)
	if ok {
		t.Fatalf("expected non-dispatch body to fail parsing")
	}
}
