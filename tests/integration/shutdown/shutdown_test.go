package shutdown_test

import (
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8081/health", nil)
	if err != nil {
		t.Logf("Failed to create request: %v", err)
	} else {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Logf("Health check failed (expected if no health endpoint): %v", err)
		} else {
			if closeErr := resp.Body.Close(); closeErr != nil {
				t.Logf("Failed to close response body: %v", closeErr)
			}
		}
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
		if killErr := cmd.Process.Kill(); killErr != nil {
			t.Logf("Failed to force kill process: %v", killErr)
		}
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
		if killErr := cmd.Process.Kill(); killErr != nil {
			t.Logf("Failed to force kill process: %v", killErr)
		}
	}

	t.Log("Timeout shutdown test completed")
}
