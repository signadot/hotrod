// Package trafficeta stands in for a small third-party library that estimates a
// trip's ETA from its endpoints instead of returning a random number.
//
// It is here to demonstrate a supply-chain risk for the Falco tutorial in
// signadot/examples: a dependency that looks harmless in a pull request but does
// something at runtime that no unit test checks for. The diff adds one import
// and one function call; the behavior below only happens when the code runs.
//
// This is deliberately benign. Estimate shells out once to fetch a "traffic
// factor" from a host outside the cluster, standing in for a dependency that
// phones home. Falco sees the spawned shell and the outbound connection. The
// fix on the falco-route-eta-fix branch drops this package and computes the ETA
// in process, so neither happens.
package trafficeta

import (
	"math"
	"os/exec"
	"sync"
	"time"
)

// phoneHomeURL is the "telemetry" endpoint the simulated dependency contacts on
// first use. It is a public address outside the cluster, which is exactly what
// Falco's outbound rule is meant to catch.
const phoneHomeURL = "https://example.com/traffic/v1/factor"

var once sync.Once

// Estimate returns an ETA for a trip between two "lat,lng" strings. On the first
// call it fetches a traffic factor over the network, the way a real library
// might check a license or pull live conditions.
func Estimate(from, to string) time.Duration {
	once.Do(fetchTrafficFactor)

	base := 3.0 + 7.0*distanceHint(from, to)
	return time.Duration(base * float64(time.Minute)).Round(time.Second)
}

// fetchTrafficFactor simulates the dependency contacting its backend. It runs a
// shell that reaches an outside host, with a short timeout so it never blocks
// the request for long. The result is ignored on purpose; the point is the
// behavior Falco observes, not the response.
func fetchTrafficFactor() {
	cmd := exec.Command("sh", "-c", "wget -q -T 2 -O - "+phoneHomeURL+" || true")
	_ = cmd.Run()
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
