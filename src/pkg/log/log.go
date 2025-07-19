package log

import (
	"context"
	"os"
	"path"
	"sync"
	"sync/atomic"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/* Global Default Logger */

var (
	initMark      sync.Once
	defaultLogger atomic.Value

	cfg = zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "func",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
)

// 后续可以把logger作到配置文件中，不需要lazy init
func logger() *zap.SugaredLogger {
	if l, ok := defaultLogger.Load().(*zap.SugaredLogger); ok {
		return l
	} else {
		initMark.Do(InitDebugLogger)
		return defaultLogger.Load().(*zap.SugaredLogger)
	}
}

func Child(name string) *zap.Logger {
	return logger().Desugar().Named(name).WithOptions(zap.AddCallerSkip(-1))
}

func Close() {
	_ = logger().Sync()
}

/* Logger Initialization */

func InitLogger() {
	_ = os.MkdirAll("../log", 0777)
	w := zapcore.AddSync(
		&lumberjack.Logger{
			Filename:   "../log/" + path.Base(os.Args[0]) + ".log",
			MaxSize:    500, // MB
			MaxBackups: 100,
			LocalTime:  true,
		},
	)
	core := zapcore.NewCore(zapcore.NewJSONEncoder(cfg), w, zap.DebugLevel)
	l := zap.New(core, zap.ErrorOutput(w), zap.AddStacktrace(zap.DPanicLevel), zap.AddCaller(), zap.AddCallerSkip(1))
	defaultLogger.Store(l.Sugar())
}

func InitDebugLogger() {
	core := zapcore.NewCore(zapcore.NewJSONEncoder(cfg), zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
	defaultLogger.Store(zap.New(core, zap.AddCallerSkip(1), zap.AddCaller()).Sugar())
}

func Debugw(msg string, strKeysAndValues ...interface{}) {
	logger().Debugw(msg, strKeysAndValues...)
}

func Infow(msg string, strKeysAndValues ...interface{}) {
	logger().Infow(msg, strKeysAndValues...)
}

func Errorw(msg string, strKeysAndValues ...interface{}) {
	logger().Errorw(msg, strKeysAndValues...)
}

func DebugEx(ctx context.Context, msg string, strKeysAndValues ...interface{}) {
	logger().Debugw(msg, append(LogFields2Interfaces(ctx), strKeysAndValues...)...)
}

func InfoEx(ctx context.Context, msg string, strKeysAndValues ...interface{}) {
	logger().Infow(msg, append(LogFields2Interfaces(ctx), strKeysAndValues...)...)
}

func ErrorEx(ctx context.Context, msg string, strKeysAndValues ...interface{}) {
	logger().Errorw(msg, append(LogFields2Interfaces(ctx), strKeysAndValues...)...)
}
