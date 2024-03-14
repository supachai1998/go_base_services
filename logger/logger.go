package logger

import (
	"context"
	"fmt"
	"go_base/xerror"
	"log"
	"os"
	"path"
	"time"

	prettyconsole "github.com/thessem/zap-prettyconsole"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey struct{}

// MustInit initializes and replaces zap's global logger.
func MustInit() {
	logger := defaultConfig()
	zap.ReplaceGlobals(logger)
}

func StdLogger(logger *zap.Logger) *log.Logger {
	return zap.NewStdLog(logger)
}
func L() *zap.SugaredLogger {
	return zap.S()
}

func ContextWithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

func Ctx(ctx context.Context) *zap.SugaredLogger {
	l, ok := ctx.Value(ctxKey{}).(*zap.SugaredLogger)
	if ok {
		return l
	}

	return L()
}

func ReplaceGlobals(logger *zap.Logger) func() {
	return zap.ReplaceGlobals(logger)
}

// NewNoop returns no-op logger.
func NewNoop() *zap.Logger {
	return zap.NewNop()
}

func NewPretty(logFile *os.File) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Development = true
	cfg.Encoding = "pretty_console"
	cfg.EncoderConfig = prettyconsole.NewEncoderConfig()
	return cfg.Build()
}

func Warn(ctx context.Context, err error, message string) {
	xerr, _ := err.(*xerror.Xerror)
	Ctx(ctx).With(zap.Inline(xerr)).Warn(message)
}

func WithError(err error) zapcore.Field {
	xerr, ok := err.(*xerror.Xerror)
	if !ok {
		return zap.String("error", err.Error())
	}

	return zap.Inline(xerr)
}

func defaultConfig() *zap.Logger {
	now := time.Now()
	logfile := path.Join("./logger/log/app", fmt.Sprintf("%s.log", now.Format("2006-01-02")))

	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return zap.NewNop()
	}

	pe := zap.NewProductionEncoderConfig()

	fileEncoder := zapcore.NewJSONEncoder(pe)
	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	// cfg := ecszap.NewDefaultEncoderConfig()
	// core := ecszap.NewCore(cfg, os.Stdout, zap.ErrorLevel)
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(file), highPriority),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), highPriority),
	)

	logger := zap.New(core, zap.AddCaller())
	logger.Named("go-base")
	return logger
}

func GormLogger() *zap.Logger {
	return zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.Lock(os.Stdout),
		zap.InfoLevel,
	))
}
