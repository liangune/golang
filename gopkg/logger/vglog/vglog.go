package vglog

import (
	"bytes"
	"fmt"
	"go/gopkg/logger/vglog/glog"
	"io"
	"os"
	"os/signal"
	"sync"
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
	SevTrace                    = "trace"
	SevInterfaceAverageDuration = "average"
	SevFatal                    = "fatal"
)

type Severity = int32

// 日记等级, 一般设置InfoLog
const (
	DebugLog                    Severity = 0 // debug
	InfoLog                              = 1 // info
	NoticeLog                            = 2 // notice
	WarnLog                              = 3 // warning
	ErrorLog                             = 4 // error
	AccessLog                            = 5 // access
	TraceLog                             = 6 // trace
	InterfaceAverageDurationLog          = 7 // average
	FatalLog                             = 8 // fatal
	NumSeverity                          = 9
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

type Logger struct {
	glogLogger *glog.Logger

	//日志等级
	allowLogLevel int32
	//日志文件权限
	logFileMode os.FileMode
	// 是否调试模式
	logModeName  string
	logModeCode  int
	fileIoWriter [NumSeverity]io.Writer

	mutex MutexWrap
}

type MutexWrap struct {
	lock     sync.Mutex
	disabled bool
}

func (mw *MutexWrap) Lock() {
	if !mw.disabled {
		mw.lock.Lock()
	}
}

func (mw *MutexWrap) Unlock() {
	if !mw.disabled {
		mw.lock.Unlock()
	}
}

func (mw *MutexWrap) Disable() {
	mw.disabled = true
}

func NewLogger() *Logger {
	l := &Logger{
		glogLogger:    nil,
		allowLogLevel: DebugLog,
		logFileMode:   defaultLogFileMode,
		logModeName:   ReleaseMode,
		logModeCode:   releaseCode,
	}
	return l
}

var DefaultLogger *Logger = NewLogger()

func (l *Logger) SetLogFileMode(mode os.FileMode) {
	l.logFileMode = mode
	l.glogLogger.LogFileMode(l.logFileMode)
}

// SetLogMode sets vglog mode according to input string.
func (l *Logger) SetLogMode(value string) {
	switch value {
	case DebugMode, "":
		l.logModeCode = debugCode
	case ReleaseMode:
		l.logModeCode = releaseCode
	default:
		panic(any(fmt.Sprintf("vglog mode unknown: %s", value)))
	}
	if value == "" {
		value = DebugMode
	}
	l.logModeName = value
}

// IsDebugging returns true if vglog is running in debug mode.
// Use vglog.ReleaseMode to disable debug mode.
func (l *Logger) IsDebugging() bool {
	return l.logModeCode == debugCode
}

func (l *Logger) SetLogLevel(level int) {
	atomic.StoreInt32(&l.allowLogLevel, int32(level))
}

func (l *Logger) GetLogLevel() int {
	return int(atomic.LoadInt32(&l.allowLogLevel))
}

// SetOutput sets the logger output.
func (l *Logger) SetOutput(logLevel Severity, output io.Writer) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.fileIoWriter[logLevel] = output
}

/*
@brief  日记的初始化函数
@param  logdir  日记存放目录
@param  logLevel 日记等级, 一般设置为InfoLog, 比当前设置等级低的日记不会输出
@param  mode 日记的模式, 分为两种模式: debug和release, 设置为debug模式时, 控制台也是输出日记
@return 无
*/
func VglogInit(logDir string, logLevel Severity, mode string) {
	DefaultLogger.allowLogLevel = logLevel

	DefaultLogger.SetLogMode(mode)

	DefaultLogger.glogLogger = glog.NewLogger().
		AlsoToStderr(DefaultLogger.IsDebugging()).
		LogDir(logDir).
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
		DefaultLogger.glogLogger.Flush()
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
	case glog.TraceLog:
		fmt.Fprintf(buf, "[%s] ", ts.Format("2006-01-02 15:04:05"))
	case glog.InterfaceAverageDurationLog:
		fmt.Fprintf(buf, "[%s] ", ts.Format("2006-01-02 15:04:05"))
	}
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
	case glog.SevTrace:
		tag = SevTrace
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
	if severityLevel == glog.SevAccess || severityLevel == glog.SevNotice || severityLevel == glog.SevTrace {
		filename = fmt.Sprintf("%s.%04d-%02d-%02d.%02d.log",
			tag,
			ts.Year(),
			ts.Month(),
			ts.Day(),
			ts.Hour())
	} else {
		filename = fmt.Sprintf("%s.%04d-%02d-%02d.log",
			tag,
			ts.Year(),
			ts.Month(),
			ts.Day())
	}
	return filename
}

func Debug(format string, args ...interface{}) {
	if atomic.LoadInt32(&DefaultLogger.allowLogLevel) <= DebugLog {
		DefaultLogger.glogLogger.DebugDepth(1, fmt.Sprintf(format, args...))
	}
}

func Info(format string, args ...interface{}) {
	if atomic.LoadInt32(&DefaultLogger.allowLogLevel) <= InfoLog {
		DefaultLogger.glogLogger.InfoDepth(1, fmt.Sprintf(format, args...))
	}
}

func Notice(format string, args ...interface{}) {
	if atomic.LoadInt32(&DefaultLogger.allowLogLevel) <= NoticeLog {
		DefaultLogger.glogLogger.NoticeDepth(1, fmt.Sprintf(format, args...))
	}
}

func Warn(format string, args ...interface{}) {
	if atomic.LoadInt32(&DefaultLogger.allowLogLevel) <= WarnLog {
		DefaultLogger.glogLogger.WarningDepth(1, fmt.Sprintf(format, args...))
	}
}

func Error(format string, args ...interface{}) {
	if atomic.LoadInt32(&DefaultLogger.allowLogLevel) <= ErrorLog {
		DefaultLogger.glogLogger.ErrorDepth(1, fmt.Sprintf(format, args...))
	}
}

func Access(format string, args ...interface{}) {
	if atomic.LoadInt32(&DefaultLogger.allowLogLevel) <= AccessLog {
		DefaultLogger.glogLogger.AccessDepth(1, fmt.Sprintf(format, args...))
	}
}

func Trace(format string, args ...interface{}) {
	if atomic.LoadInt32(&DefaultLogger.allowLogLevel) <= TraceLog {
		DefaultLogger.glogLogger.TraceDepth(1, fmt.Sprintf(format, args...))
	}
}

func InterfaceAverageDuration(format string, args ...interface{}) {
	if atomic.LoadInt32(&DefaultLogger.allowLogLevel) <= InterfaceAverageDurationLog {
		DefaultLogger.glogLogger.InterfaceAverageDurationDepth(1, fmt.Sprintf(format, args...))
	}
}

func FlushLog() {
	DefaultLogger.glogLogger.Flush()
}
