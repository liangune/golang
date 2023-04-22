package zaplog

import (
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestLogger(t *testing.T) {
	ZaplogInit("./logs/server.log", zapcore.InfoLevel, 1, 4, 7, true)
	Infof("my log test : %v", "zap")
}

/*
func TestLoggerBackups(t *testing.T) {
	LoggerInit("./logs/server.log", zapcore.InfoLevel, 1, 30, 7, true)
	i := 0
	for {
		i++
		Infof("TestLoggerBackups TestLoggerBackups number: %d", i)
	}
}
*/
