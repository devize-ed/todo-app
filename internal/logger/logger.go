package logger

import (
	"errors"
	"os"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// global is the singleton logger instance shared across the application.
// Initialize() sets it; before that, it is a safe no-op nop logger so the
// package-level helpers never panic on a nil pointer.
var global = &Logger{SugaredLogger: zap.NewNop().Sugar()}

type Logger struct {
	*zap.SugaredLogger
}

// Initialize builds the singleton logger at the given level and stores it
// in the package-level global. It also returns the instance for callers that
// prefer dependency injection.
func Initialize(level string) (*Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05")
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.MessageKey = "msg"
	cfg.EncoderConfig.LevelKey = "level"
	cfg.DisableStacktrace = true

	zl, err := cfg.Build(
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCaller(),
		// AddCallerSkip(1) so callers see logger.go:NN replaced by the
		// real call site when using the package-level wrappers below.
		zap.AddCallerSkip(1),
	)
	if err != nil {
		return nil, err
	}

	global = &Logger{SugaredLogger: zl.Sugar()}
	return global, nil
}

// Get returns the singleton logger instance.
func Get() *Logger {
	return global
}

// SafeSync flushes the singleton logger, ignoring the harmless EINVAL/ENOTTY
// errors that zap returns when stdout/stderr is not a real file.
func SafeSync() {
	if global == nil || global.SugaredLogger == nil {
		return
	}
	if err := global.Sync(); err != nil {
		var pe *os.PathError
		if errors.As(err, &pe) && (errors.Is(pe.Err, syscall.EINVAL) || errors.Is(pe.Err, syscall.ENOTTY)) {
			return
		}
		if errors.Is(err, syscall.EINVAL) || errors.Is(err, syscall.ENOTTY) {
			return
		}
		global.Errorf("failed to sync logger: %v", err)
	}
}

// Package-level wrappers around the singleton SugaredLogger.
// They are thin shims so callers can write `logger.Warn(...)` directly.

func Debug(args ...any)                    { global.Debug(args...) }
func Info(args ...any)                     { global.Info(args...) }
func Warn(args ...any)                     { global.Warn(args...) }
func Error(args ...any)                    { global.Error(args...) }
func Fatal(args ...any)                    { global.Fatal(args...) }
func Debugf(template string, args ...any)  { global.Debugf(template, args...) }
func Infof(template string, args ...any)   { global.Infof(template, args...) }
func Warnf(template string, args ...any)   { global.Warnf(template, args...) }
func Errorf(template string, args ...any)  { global.Errorf(template, args...) }
func Fatalf(template string, args ...any)  { global.Fatalf(template, args...) }
func Debugw(msg string, kv ...any)         { global.Debugw(msg, kv...) }
func Infow(msg string, kv ...any)          { global.Infow(msg, kv...) }
func Warnw(msg string, kv ...any)          { global.Warnw(msg, kv...) }
func Errorw(msg string, kv ...any)         { global.Errorw(msg, kv...) }
