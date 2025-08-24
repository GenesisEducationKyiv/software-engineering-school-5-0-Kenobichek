# Graceful Shutdown Implementation

This document describes the graceful shutdown pattern implemented in the Weather Forecast API application, based on the VictoriaMetrics blog post approach.

## Overview

The application implements a graceful shutdown pattern that ensures proper cleanup and termination when receiving shutdown signals (SIGINT, SIGTERM). This prevents data corruption, ensures in-flight requests are completed, and properly closes all resources.

## Implementation Details

### Signal Handling

The application uses Go's `signal.NotifyContext` to listen for OS interrupt signals:

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()
```

### Component Shutdown Order

When a shutdown signal is received, the application follows this shutdown sequence:

1. **HTTP Server Shutdown**: Gracefully stops accepting new connections and waits for existing requests to complete
2. **Scheduler Shutdown**: Stops the cron scheduler and waits for running jobs to finish
3. **Database Connection Cleanup**: Closes the database connection
4. **Application Termination**: Logs completion and exits

### Configuration

The graceful shutdown timeout is configurable via the `GRACEFUL_SHUTDOWN_TIMEOUT` environment variable:

```bash
GRACEFUL_SHUTDOWN_TIMEOUT=30s
```

Default value is 30 seconds.

### Key Components

#### HTTP Server
- Uses `http.Server.Shutdown()` with context timeout
- Waits for existing requests to complete
- Stops accepting new connections

#### Scheduler
- Implements `Stop()` method that calls `cron.Cron.Stop()`
- Waits for running cron jobs to complete
- Thread-safe with mutex protection

#### Database
- Properly closes database connections
- Logs any errors during cleanup

## Usage

The graceful shutdown is automatically handled when the application receives:
- `SIGINT` (Ctrl+C)
- `SIGTERM` (termination signal)

### Example Shutdown Sequence

```
2024/01/01 12:00:00 Server is running on :8080
2024/01/01 12:00:05 Shutdown signal received
2024/01/01 12:00:05 [Scheduler] Stopping scheduler...
2024/01/01 12:00:05 [Scheduler] Stopped successfully
2024/01/01 12:00:05 Application shutdown complete
```

## Benefits

1. **Data Integrity**: Prevents data corruption by allowing in-flight operations to complete
2. **Resource Cleanup**: Ensures all resources (database connections, file handles) are properly closed
3. **User Experience**: Existing requests are not abruptly terminated
4. **Monitoring**: Proper logging of shutdown process for debugging
5. **Configurable**: Shutdown timeout can be adjusted based on application needs

## Testing

The graceful shutdown functionality is tested in `internal/scheduler/scheduler_test.go`:

- `TestSchedulerStop`: Tests normal shutdown sequence
- `TestSchedulerStopWithoutStart`: Tests edge case of stopping without starting

Run tests with:
```bash
go test ./internal/scheduler/...
``` 