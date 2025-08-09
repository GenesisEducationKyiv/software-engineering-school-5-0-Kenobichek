# 006: Logger Implementation with Zap

## Context
A high-performance, structured, and configurable logger is required for the microservices. The logging mechanism must be fast to minimize its impact on application performance, especially under high load. The system also needs the flexibility to set different log levels (e.g., debug for development, info for production) and to produce structured data that can be easily analyzed by centralized logging systems.

## Decision

The go.uber.org/zap library will be used.

A wrapper struct, zapLogger, has been created to provide a simple, consistent interface for common logging levels (Info, Error, Debug). This design choice isolates the rest of the codebase from the specifics of zap, making it easier to swap out the logging library in the future if needed.

## How it works

The NewZapLogger function initializes the logging core components:

1. An encoder is created using zapcore.NewJSONEncoder with a production-ready configuration. This ensures logs are formatted as structured JSON with ISO 8601 timestamps.

2. A writer is configured to output logs to standard output (os.Stdout), which is a common pattern for applications running in containers.

3. A baseCore is established, combining the encoder, writer, and a base logging level of zap.DebugLevel.

4. A sampledCore is wrapped around the baseCore using zapcore.NewSamplerWithOptions. This component implements the log sampling logic, controlling the rate of messages logged after the initial burst.

5. Finally, a new zap.New logger is instantiated with the sampledCore and the zap.AddCaller() option. The .Sugar() method is then called to provide a flexible and easy-to-use logging interface.

The zapLogger struct holds the sugar logger, and its methods (Info, Error, Debug, Sync) simply delegate calls to the underlying sugar instance.



## Justification

* High Performance: zap is recognized as one of the fastest Go loggers available due to its use of reflection-free encoding and caching. This is critical for performance-sensitive applications.
* Structured Logging: zap supports structured logging in JSON format, which simplifies the aggregation and analysis of logs in a centralized system. The implementation uses zapcore.JSONEncoder for this purpose.
* Flexible Configuration: Log output (os.Stdout), log level (zap.DebugLevel), and time format (zapcore.ISO8601TimeEncoder) can be easily configured.
* Sampling: zap includes built-in log sampling to reduce noise from repetitive log messages. The implementation is configured to log the first 100 messages within a second, and then only one out of every 100 messages of the same type. This prevents the logging system from being overwhelmed by a flood of identical logs.
* Sugared Logger: The "sugared" version of the logger (zap.Sugar()) is used to provide a more convenient and familiar fmt.Printf-like interface. This makes it easier to log messages with formatted strings and arguments instead of explicit key-value pairs.
* Caller Information: zap.AddCaller() is included to automatically embed the file and line number where the logger was called. This is invaluable for debugging and tracing the source of a log message.