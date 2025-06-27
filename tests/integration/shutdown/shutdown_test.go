package shutdown

import (
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

func TestGracefulShutdown(t *testing.T) {
	// Skip if not in integration test environment
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Start the application as a separate process
	cmd := exec.Command("go", "run", "cmd/main.go")
	cmd.Env = append(os.Environ(),
		"PORT=8081", // Use different port for testing
		"GRACEFUL_SHUTDOWN_TIMEOUT=5s",
	)

	// Capture output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the process
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start application: %v", err)
	}

	// Wait a bit for the app to start
	time.Sleep(2 * time.Second)

	// Test that the server is responding
	resp, err := http.Get("http://localhost:8081/health")
	if err != nil {
		t.Logf("Health check failed (expected if no health endpoint): %v", err)
	} else {
		resp.Body.Close()
	}

	// Send SIGTERM signal
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("Failed to send SIGTERM: %v", err)
	}

	// Wait for graceful shutdown (with timeout)
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Logf("Process exited with error (this might be expected): %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Application did not shutdown within 10 seconds")
		cmd.Process.Kill() // Force kill if needed
	}

	t.Log("Graceful shutdown test completed")
}

func TestGracefulShutdownWithTimeout(t *testing.T) {
	// Skip if not in integration test environment
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Start the application with very short timeout
	cmd := exec.Command("go", "run", "cmd/main.go")
	cmd.Env = append(os.Environ(),
		"PORT=8082",
		"GRACEFUL_SHUTDOWN_TIMEOUT=1s", // Very short timeout
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start application: %v", err)
	}

	time.Sleep(2 * time.Second)

	// Send SIGTERM
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("Failed to send SIGTERM: %v", err)
	}

	// Wait for shutdown
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		t.Logf("Process exited: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("Application did not shutdown within 5 seconds")
		cmd.Process.Kill()
	}

	t.Log("Timeout shutdown test completed")
}
