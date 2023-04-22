package vglog

import (
	"testing"
	"time"
)

func TestVglog(t *testing.T) {
	VglogInit("./logs", InfoLog, DebugMode)
	Error("this is error log")

	NewGlogCleaner(InitOption{
		Path:           "./logs/",
		Interval:       2,
		Reserve:        1,
		Compress:       true,
		CompressMethod: CompressMethodGzip,
	})

	time.Sleep(time.Second * 10)
	FlushLog()
}
