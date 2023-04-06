package log

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.uber.org/fx/fxevent"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/log_mock.go . Logger

type ctxKey int

const (
	RequestIDKey ctxKey = iota + 1
)

type Field zapcore.Field

func String(key, value string) Field {
	return Field(zap.String(key, value))
}

func Any(key string, value interface{}) Field {
	return Field(zap.Any(key, value))
}

func Int64(key string, value int64) Field {
	return Field(zap.Int64(key, value))
}
func Uint64(key string, value uint64) Field {
	return Field(zap.Uint64(key, value))
}

type Logger interface {
	fxevent.Logger
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Print(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Warning(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
	Logger() *zap.Logger
}

type Log struct {
	logger *zap.Logger
}

func (l Log) Logger() *zap.Logger {
	return l.logger
}

func (l Log) Debug(msg string, fields ...Field) {
	var zf []zap.Field
	for _, f := range fields {
		zf = append(zf, zap.Field(f))
	}
	l.logger.Debug(msg, zf...)
}

func (l Log) Info(msg string, fields ...Field) {
	var zf []zap.Field
	for _, f := range fields {
		zf = append(zf, zap.Field(f))
	}
	l.logger.Info(msg, zf...)
}

func (l Log) Print(msg string, fields ...Field) {
	var zf []zap.Field
	for _, f := range fields {
		zf = append(zf, zap.Field(f))
	}
	l.logger.Info(msg, zf...)
}

func (l Log) Warn(msg string, fields ...Field) {
	var zf []zap.Field
	for _, f := range fields {
		zf = append(zf, zap.Field(f))
	}
	l.logger.Warn(msg, zf...)
}

func (l Log) Warning(msg string, fields ...Field) {
	var zf []zap.Field
	for _, f := range fields {
		zf = append(zf, zap.Field(f))
	}
	l.logger.Warn(msg, zf...)
}

func (l Log) Error(msg string, fields ...Field) {
	var zf []zap.Field
	for _, f := range fields {
		zf = append(zf, zap.Field(f))
	}
	l.logger.Error(msg, zf...)
}

func (l Log) Fatal(msg string, fields ...Field) {
	var zf []zap.Field
	for _, f := range fields {
		zf = append(zf, zap.Field(f))
	}
	l.logger.Fatal(msg, zf...)
}

func (l Log) Panic(msg string, fields ...Field) {
	var zf []zap.Field
	for _, f := range fields {
		zf = append(zf, zap.Field(f))
	}
	l.logger.Panic(msg, zf...)
}

func NewLog(level string) (Logger, error) {
	config := zap.NewProductionConfig()
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}
	config.Level = lvl
	config.Development = config.Level.Level() == zapcore.DebugLevel
	config.DisableStacktrace = config.Level.Level() != zapcore.DebugLevel
	config.DisableCaller = config.Level.Level() != zapcore.DebugLevel
	config.EncoderConfig.MessageKey = "message"
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	return &Log{logger: logger}, nil
}

func (l Log) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.logger.Info("OnStart hook executing",
			zap.String("callee", e.FunctionName),
			zap.String("caller", e.CallerName),
		)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.logger.Info("OnStart hook failed",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.Error(e.Err),
			)
		} else {
			l.logger.Info("OnStart hook executed",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.String("runtime", e.Runtime.String()),
			)
		}
	case *fxevent.OnStopExecuting:
		l.logger.Info("OnStop hook executing",
			zap.String("callee", e.FunctionName),
			zap.String("caller", e.CallerName),
		)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.logger.Info("OnStop hook failed",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.Error(e.Err),
			)
		} else {
			l.logger.Info("OnStop hook executed",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.String("runtime", e.Runtime.String()),
			)
		}
	case *fxevent.Supplied:
		l.logger.Info("supplied", zap.String("type", e.TypeName), zap.Error(e.Err))
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.logger.Info("provided",
				zap.String("constructor", e.ConstructorName),
				zap.String("type", rtype),
			)
		}
		if e.Err != nil {
			l.logger.Error("error encountered while applying options",
				zap.Error(e.Err))
		}
	case *fxevent.Invoking:
		// Do not log stack as it will make logs hard to read.
		l.logger.Info("invoking",
			zap.String("function", e.FunctionName))
	case *fxevent.Invoked:
		if e.Err != nil {
			l.logger.Error("invoke failed",
				zap.Error(e.Err),
				zap.String("stack", e.Trace),
				zap.String("function", e.FunctionName))
		}
	case *fxevent.Stopping:
		l.logger.Info("received signal",
			zap.String("signal", strings.ToUpper(e.Signal.String())))
	case *fxevent.Stopped:
		if e.Err != nil {
			l.logger.Error("stop failed", zap.Error(e.Err))
		}
	case *fxevent.RollingBack:
		l.logger.Error("start failed, rolling back", zap.Error(e.StartErr))
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.logger.Error("rollback failed", zap.Error(e.Err))
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.logger.Error("start failed", zap.Error(e.Err))
		} else {
			l.logger.Info("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.logger.Error("custom logger initialization failed", zap.Error(e.Err))
		} else {
			l.logger.Info("initialized custom fxevent.Logger", zap.String("function", e.ConstructorName))
		}
	}
}
