package observability

import (
    "io"
    "math/rand"
    "os"
    "strconv"
    "time"

    "github.com/sirupsen/logrus"
)

// Logger is the singleton instance used across the service.
var Logger = logrus.New()

var sampleRate int

func init() {
    Logger.SetOutput(os.Stdout)
    Logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})

    // configure level via LOG_LEVEL (default info)
    lvl, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
    if err != nil {
        lvl = logrus.InfoLevel
    }
    Logger.SetLevel(lvl)

    // sampling rate via LOG_SAMPLE_RATE, default 10 (=> 1/10 msgs)
    sampleRate = 10
    if srEnv := os.Getenv("LOG_SAMPLE_RATE"); srEnv != "" {
        if v, err := strconv.Atoi(srEnv); err == nil && v > 0 {
            sampleRate = v
        }
    }

    rand.Seed(time.Now().UnixNano())

    // Replace the standard library logger output with logrus at info level.
    // The std logger is used in legacy code; redirect it to our logger.
    logrus.SetOutput(Logger.WriterLevel(logrus.InfoLevel))
}

func sampled() bool {
    if sampleRate <= 1 {
        return true // log everything
    }
    return rand.Intn(sampleRate) == 0
}

// Helper wrappers providing sampling for Info and Debug levels.
// Warn and Error are always logged.

func Infof(format string, args ...interface{}) {
    if sampled() {
        Logger.Infof(format, args...)
    }
}

func Infow(msg string, fields logrus.Fields) {
    if sampled() {
        Logger.WithFields(fields).Info(msg)
    }
}

func Debugf(format string, args ...interface{}) {
    if sampled() {
        Logger.Debugf(format, args...)
    }
}

func Warnf(format string, args ...interface{}) { Logger.Warnf(format, args...) }
func Errorf(format string, args ...interface{}) { Logger.Errorf(format, args...) }
func Fatalf(format string, args ...interface{}) { Logger.Fatalf(format, args...) }

// Writer returns an io.Writer that writes to the underlying logger at info level.
func Writer() io.Writer { return Logger.WriterLevel(logrus.InfoLevel) }
