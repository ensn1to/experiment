package log

import (
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Level = zapcore.Level

const (
	InfoLevel   Level = zap.InfoLevel   // 0, default level
	WarnLevel   Level = zap.WarnLevel   // 1
	ErrorLevel  Level = zap.ErrorLevel  // 2
	DPanicLevel Level = zap.DPanicLevel // 3, used in development log

	// PanicLevel logs a message, then panics
	PanicLevel Level = zap.PanicLevel // 4

	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel Level = zap.FatalLevel // 5
	DebugLevel Level = zap.DebugLevel // -1
)

type Field = zap.Field

type Option = zap.Option

type Logger struct {
	l     *zap.Logger
	level zap.AtomicLevel // log 动态日志level
}

// new logger. self-defined logger instead of std
func New(writer io.Writer, level Level, opts ...Option) *Logger {
	if writer == nil {
		panic("the writer is nil")
	}

	// zapcore.Level转成atomicLevel
	atomicLevel := zap.NewAtomicLevelAt(level)

	cfg := zap.NewProductionConfig()
	// 自定义encoder
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339))
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.AddSync(writer),
		atomicLevel,
	)

	return &Logger{
		l:     zap.New(core, opts...),
		level: atomicLevel,
	}
}

func NewProduction(outputPaths []string, opts ...Option) *Logger {
	cfg := zap.NewProductionConfig()
	// 自定义encoder
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02T15:04:05"))
	}
	if len(outputPaths) > 0 {
		cfg.OutputPaths = outputPaths
	}

	logger, err := cfg.Build(opts...)
	if err != nil {
		panic(err)
	}

	return &Logger{
		l: logger,
	}
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.l.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.l.Info(msg, fields...)
}

func (l *Logger) Infof(msg string, keysAndValues ...interface{}) {
	l.l.Sugar().Infof(msg, keysAndValues...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.l.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.l.Error(msg, fields...)
}

func (l *Logger) DPanic(msg string, fields ...Field) {
	l.l.DPanic(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...Field) {
	l.l.Panic(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	l.l.Fatal(msg, fields...)
}

func (l *Logger) Sync() error {
	return l.l.Sync()
}

// SetLevel alters the logging level on runtime
// it is concurrent-safe
func (l *Logger) SetLevel(level Level) error {
	l.level.SetLevel(level)
	return nil
}

func Sync() error {
	if std != nil {
		return std.Sync()
	}

	return nil
}

var std = New(os.Stderr, InfoLevel, WithCaller(true))

// default logger
func Default() *Logger {
	return std
}

// not safe for concurrent use
func ResetDefault(l *Logger) {
	std = l
	Info = std.Info
	Warn = std.Warn
	Error = std.Error
	DPanic = std.DPanic
	Panic = std.Panic
	Fatal = std.Fatal
	Debug = std.Debug
}

// 将zap的一些配合log输出的sugar函数暴露给用户
var (
	Skip        = zap.Skip
	Binary      = zap.Binary
	Bool        = zap.Bool
	Boolp       = zap.Boolp
	ByteString  = zap.ByteString
	Complex128  = zap.Complex128
	Complex128p = zap.Complex128p
	Complex64   = zap.Complex64
	Complex64p  = zap.Complex64p
	Float64     = zap.Float64
	Float64p    = zap.Float64p
	Float32     = zap.Float32
	Float32p    = zap.Float32p
	Int         = zap.Int
	Intp        = zap.Intp
	Int64       = zap.Int64
	Int64p      = zap.Int64p
	Int32       = zap.Int32
	Int32p      = zap.Int32p
	Int16       = zap.Int16
	Int16p      = zap.Int16p
	Int8        = zap.Int8
	Int8p       = zap.Int8p
	String      = zap.String
	Stringp     = zap.Stringp
	Uint        = zap.Uint
	Uintp       = zap.Uintp
	Uint64      = zap.Uint64
	Uint64p     = zap.Uint64p
	Uint32      = zap.Uint32
	Uint32p     = zap.Uint32p
	Uint16      = zap.Uint16
	Uint16p     = zap.Uint16p
	Uint8       = zap.Uint8
	Uint8p      = zap.Uint8p
	Uintptr     = zap.Uintptr
	Uintptrp    = zap.Uintptrp
	Reflect     = zap.Reflect
	Namespace   = zap.Namespace
	Stringer    = zap.Stringer
	Time        = zap.Time
	Timep       = zap.Timep
	Stack       = zap.Stack
	StackSkip   = zap.StackSkip
	Duration    = zap.Duration
	Durationp   = zap.Durationp
	Any         = zap.Any
)

// 将std实例方法以包级函数形式暴露
var (
	Info   = std.Info
	Warn   = std.Warn
	Error  = std.Error
	DPanic = std.DPanic
	Panic  = std.Panic
	Fatal  = std.Fatal
	Debug  = std.Debug
)

var (
	AddStacktrace = zap.AddStacktrace
	WithCaller    = zap.WithCaller
)

type LevelEnablerFunc func(level Level) bool

type TeeOption struct {
	W   io.Writer
	Lef LevelEnablerFunc
}

func NewTee(tops []TeeOption, opts ...Option) *Logger {
	var cores []zapcore.Core
	cfg := zap.NewProductionConfig()
	// 自定义encoder
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02T15:04:05"))
	}

	for _, top := range tops {
		top := top
		if top.W == nil {
			panic("the writer is nil")
		}

		lv := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return top.Lef(Level(level))
		})

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.AddSync(top.W),
			lv,
		)

		cores = append(cores, core)

	}

	return &Logger{
		l: zap.New(zapcore.NewTee(cores...), opts...),
	}
}

type RotateOptions struct {
	MaxSize    int  // 日志文件最大值
	MaxAge     int  // 日志文件存活最长时间
	MaxBackups int  // 日志文件最大备份数量
	Compress   bool // 是否压缩
}

type TeeOptionWithRotate struct {
	Filename string
	Ropt     RotateOptions
	Lef      LevelEnablerFunc
}

func NewTeeWithRotate(tops []TeeOptionWithRotate, opts ...Option) *Logger {
	var cores []zapcore.Core
	cfg := zap.NewProductionConfig()
	// 自定义encoder
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02T15:04:05"))
	}

	for _, top := range tops {
		top := top

		lv := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return top.Lef(Level(level))
		})

		// 通过WriteSyncer接口接入
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   top.Filename,
			MaxSize:    top.Ropt.MaxSize,
			MaxBackups: top.Ropt.MaxBackups,
			MaxAge:     top.Ropt.MaxAge,
			Compress:   top.Ropt.Compress,
		})

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.AddSync(w),
			lv,
		)

		cores = append(cores, core)

	}

	return &Logger{
		l: zap.New(zapcore.NewTee(cores...), opts...),
	}
}
