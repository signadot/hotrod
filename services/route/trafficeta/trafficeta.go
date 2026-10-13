// Package trafficeta estimates a trip's ETA from its endpoints instead of
// returning a random number.
//
// This is the fixed version for the Falco tutorial in signadot/examples. The
// earlier version simulated a compromised dependency that shelled out to a host
// outside the cluster on first use. This version computes the ETA in process:
// no shell, no network, nothing for Falco's sandbox rules to flag.
package trafficeta

import (
	"math"
	"time"
)

// Estimate returns an ETA for a trip between two "lat,lng" strings, computed
// locally from the endpoints.
func Estimate(from, to string) time.Duration {
	base := 3.0 + 7.0*distanceHint(from, to)
	return time.Duration(base * float64(time.Minute)).Round(time.Second)
}

// distanceHint turns two coordinate-ish strings into a small, stable factor in
// [0,1] so the ETA is deterministic for a given trip.
func distanceHint(from, to string) float64 {
	h := 0
	for _, c := range from + "|" + to {
		h = (h*31 + int(c)) & 0xffff
	}
	return math.Abs(math.Sin(float64(h)))
}
