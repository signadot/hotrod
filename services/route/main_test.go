//go:build integration
// +build integration

package route

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func waitForReadyState() {
	// Introducing delay allowing the server to come up before we start testing. A better approach
	// might be to wait for the workflow status to be Ready. But going with a basic approach for now.
	log.Println("Waiting for the server to be ready")
	time.Sleep(10 * time.Second)
	log.Println("The server is ready")
}

func setup() {
	waitForReadyState()
}

func teardown() {
	// no-op. Keeping as a placeholder to be used when needed
}
