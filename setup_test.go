package httpmock

import (
	"os"
	"testing"
	"time"
)

func setup() {
	// Sleep for 60 seconds for spring to boot up.
	time.Sleep(0 * time.Second)

	// flag.Parse()
}

func shutdown() {}

// TestMain has custom setup and shutdown
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()

	os.Exit(code)
}
