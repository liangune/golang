package vglog

import (
	"bytes"
	"fmt"
	"go/gopkg/logger/vglog/glog"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

const (
	SevDebug                    = "debug"
	SevInfo                     = "info"
	SevNotice                   = "notice"
	SevWarn                     = "warning"
	SevError                    = "error"
	SevAccess                   = "access"
	SevInterfaceAverageDuration = "average"
	SevFatal                    = "fatal"
)

type Severity = int32

// 日记等级, 一般设置InfoLog
const (
	DebugLog                    Severity = 1 // debug
	InfoLog                              = 2 // info
	NoticeLog                            = 3 // notice
	WarnLog                              = 4 // warning
	ErrorLog                             = 5 // error
	AccessLog                            = 6 // access
	InterfaceAverageDurationLog          = 7 // average
	FatalLog                             = 8 // fatal
)

var (
	//日志等级
	allowLogLevel int32 = DebugLog
	//日志文件权限
	logFileMode os.FileMode
	// 是否调试模式
	vglogModeName = ReleaseMode
	vglogModeCode = releaseCode
)

const (
	// DebugMode indicates vglog mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates vglog mode is release.
	ReleaseMode = "release"

	// Unix log file default permission bits
	defaultLogFileMode os.FileMode = 0640
)
const (
	debugCode = iota
	releaseCode
)

var logger *glog.Logger

/*
@brief  日记的初始化函数
@param  logdir  日记存放目录
@param  logLevel 日记等级, 一般设置为InfoLog, 比当前设置等级低的日记不会输出
@param  mode 日记的模式, 分为两种模式: debug和release, 设置为debug模式时, 控制台也是输出日记
@return 无
*/
func VglogInit(logdir string, logLevel Severity, mode string) {
	allowLogLevel = logLevel

	SetLogMode(mode)

	logger = glog.NewLogger().
		AlsoToStderr(IsDebugging()).
		LogDir(logdir).
		EnableLogHeader(true).
		EnableLogLink(false).
		FlushInterval(time.Second).
		HeaderFormat(headerFormatFunc).
		FileNameFormat(fileNameFormatFunc).
		LogFileMode(defaultLogFileMode).
		Init()

	// 退出前flush
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		os.Stdout.WriteString("\b\bflushing log ...\n")
		logger.Flush()
		os.Stdout.WriteString("flush log done\n")
		//signal.Reset(os.Interrupt)

		//proc, err := os.FindProcess(syscall.Getpid())
		//if err != nil {
		//	panic(err)
		//}
		os.Exit(0)
		//err = proc.Signal(os.Interrupt)
		//if err != nil {
		//panic(err)
		//	fmt.Println(err)
		//}

	}()

}

func headerFormatFunc(buf *bytes.Buffer, l glog.Severity, ts time.Time, pid int, file string, line int) {
	switch l {
	case glog.InfoLog:
		fmt.Fprintf(buf, "[%s][%s:%d][INFO] ", ts.Format("2006-01-02 15:04:05"), file, line)
	case glog.DebugLog:
		fmt.Fprintf(buf, "[%s][%s:%d][DEBUG] ", ts.Format("2006-01-02 15:04:05"), file, line)
	case glog.WarnLog:
		fmt.Fprintf(buf, "[%s][%s:%d][WARN] ", ts.Format("2006-01-02 15:04:05"), file, line)
	case glog.ErrorLog:
		fmt.Fprintf(buf, "[%s][%s:%d][ERROR] ", ts.Format("2006-01-02 15:04:05"), file, line)
	case glog.NoticeLog:
		fmt.Fprintf(buf, "[%s][%s:%d][notice] ", ts.Format("2006-01-02 15:04:05"), file, line)
	case glog.AccessLog:
		fmt.Fprintf(buf, "[%s] ", ts.Format("2006-01-02 15:04:05"))
	case glog.InterfaceAverageDurationLog:
		fmt.Fprintf(buf, "[%s] ", ts.Format("2006-01-02 15:04:05"))
	}
}

func SetLogLevel(level int) {
	atomic.StoreInt32(&allowLogLevel, int32(level))
}

func GetLogLevel() int {
	return int(atomic.LoadInt32(&allowLogLevel))
}

// Unix log file permission bits
func SetLogFileMode(mode os.FileMode) {
	logFileMode = mode
	logger.LogFileMode(logFileMode)
}

// log tag
func logTag(severityLevel string) string {
	var tag string = SevInfo
	switch severityLevel {
	case glog.SevDebug:
		tag = SevDebug
	case glog.SevInfo:
		tag = SevInfo
	case glog.SevWarn:
		tag = SevWarn
	case glog.SevError:
		tag = SevError
	case glog.SevAccess:
		tag = SevAccess
	case glog.SevNotice:
		tag = SevNotice
	case glog.SevFatal:
		tag = SevFatal
	case glog.SevInterfaceAverageDuration:
		tag = SevInterfaceAverageDuration
	}

	return tag
}

// log file name
func fileNameFormatFunc(severityLevel string, ts time.Time) string {
	var filename string
	tag := logTag(severityLevel)
	filename = fmt.Sprintf("%s.%04d-%02d-%02d.log",
		tag,
		ts.Year(),
		ts.Month(),
		ts.Day())
	return filename
}

// SetLogMode sets vglog mode according to input string.
func SetLogMode(value string) {
	switch value {
	case DebugMode, "":
		vglogModeCode = debugCode
	case ReleaseMode:
		vglogModeCode = releaseCode
	default:
		panic(fmt.Sprintf("vglog mode unknown: %s", value))
	}
	if value == "" {
		value = DebugMode
	}
	vglogModeName = value
}

// IsDebugging returns true if vglog is running in debug mode.
// Use vglog.ReleaseMode to disable debug mode.
func IsDebugging() bool {
	return vglogModeCode == debugCode
}

func Debug(format string, args ...interface{}) {
	if atomic.LoadInt32(&allowLogLevel) <= DebugLog {
		logger.DebugDepth(1, fmt.Sprintf(format, args...))
	}
}

func Info(format string, args ...interface{}) {
	if atomic.LoadInt32(&allowLogLevel) <= InfoLog {
		logger.InfoDepth(1, fmt.Sprintf(format, args...))
	}
}

func Notice(format string, args ...interface{}) {
	if atomic.LoadInt32(&allowLogLevel) <= NoticeLog {
		logger.NoticeDepth(1, fmt.Sprintf(format, args...))
	}
}

func Warn(format string, args ...interface{}) {
	if atomic.LoadInt32(&allowLogLevel) <= WarnLog {
		logger.WarningDepth(1, fmt.Sprintf(format, args...))
	}
}

func Error(format string, args ...interface{}) {
	if atomic.LoadInt32(&allowLogLevel) <= ErrorLog {
		logger.ErrorDepth(1, fmt.Sprintf(format, args...))
	}
}

func Access(format string, args ...interface{}) {
	if atomic.LoadInt32(&allowLogLevel) <= AccessLog {
		logger.AccessDepth(1, fmt.Sprintf(format, args...))
	}
}

func InterfaceAverageDuration(format string, args ...interface{}) {
	if atomic.LoadInt32(&allowLogLevel) <= InterfaceAverageDurationLog {
		logger.InterfaceAverageDurationDepth(1, fmt.Sprintf(format, args...))
	}
}

func FlushLog() {
	logger.Flush()
}
