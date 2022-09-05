package log

import (
	"io"

	"go.uber.org/zap/zapcore"
)

func NewFileSyncer(writer io.Writer) zapcore.WriteSyncer {
	if ws, ok := writer.(zapcore.WriteSyncer); ok {
		return ws
	}

	// 支持并发写
	return zapcore.Lock(zapcore.AddSync(writer))
}
