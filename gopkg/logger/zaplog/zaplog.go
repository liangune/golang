package zaplog

import (
	"fmt"
	"os"

	"github.com/unknwon/goconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type sugarLogger = zap.SugaredLogger

var (
	sugar *sugarLogger
	//Logger *zap.Logger
)

type Config struct {
	Format        int    //日志格式, 0:普通一行文件，1:json格式文件
	Level         int    //0:debug,1:info,2:warn,3:error, 默认:1
	Compress      bool   // 是否压缩
	MaxSize       int    // 每个日志文件保存的最大尺寸 单位：M
	MaxBackups    int    // 日志文件最多保存多少个文件备份
	MaxAge        int    // 文件最多保存多少天
	Filename      string // 日志文件路径
	TimeKey       string
	LevelKey      string
	NameKey       string
	CallerKey     string
	MessageKey    string
	StacktraceKey string
}

func InitFromCfg(cfg *Config) {
	_init(cfg)
}

func GetDefaultConfig() *Config {
	return &Config{
		Filename:      "logs/default.log",
		MaxSize:       10,
		MaxBackups:    1000,
		MaxAge:        30,
		Compress:      true,
		TimeKey:       "Time",
		LevelKey:      "Level",
		NameKey:       "logger",
		CallerKey:     "LN",
		MessageKey:    "Msg",
		StacktraceKey: "Stacktrace",
		Level:         1,
	}
}

func init() {
	_init(nil)
}

func initFromCfgFile(cfg *Config) {
	var configFilePath string
	if len(os.Args) >= 2 {
		configFilePath = os.Args[1]
	} else {
		configFilePath = "logconfig.ini"
	}

	if c, err := goconfig.LoadConfigFile(configFilePath); err != nil {
		fmt.Println("Default Config")
	} else {
		// 读取配置文件, 获得参数
		if val, err := c.GetValue("", "LogPathFileName"); err == nil {
			cfg.Filename = val
		}

		if val, err := c.Int("", "LogMaxSize"); err == nil {
			cfg.MaxSize = val
		}

		if val, err := c.Int("", "LogMaxBackups"); err == nil {
			cfg.MaxBackups = val
		}

		if val, err := c.Int("", "LogMaxAge"); err == nil {
			cfg.MaxAge = val
		}

		if val, err := c.Int("", "LogCompress"); err == nil {
			if val <= 0 {
				cfg.Compress = false
			}
		}

		if val, err := c.Int("", "LogFormat"); err == nil {
			cfg.Format = val
		}

		if val, err := c.Int("", "LogLevel"); err == nil {
			cfg.Level = val
		}
	}
}

func _init(cfg *Config) {
	if cfg == nil {
		cfg = GetDefaultConfig()
	}

	// 跟据文件进行初始化
	initFromCfgFile(cfg)

	hook := lumberjack.Logger{
		Filename:   cfg.Filename,   // 日志文件路径
		MaxSize:    cfg.MaxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: cfg.MaxBackups, // 日志文件最多保存多少个文件备份
		MaxAge:     cfg.MaxAge,     // 文件最多保存多少天
		Compress:   cfg.Compress,   // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       cfg.TimeKey,
		LevelKey:      cfg.LevelKey,
		NameKey:       cfg.NameKey,
		CallerKey:     cfg.CallerKey,
		MessageKey:    cfg.MessageKey,
		StacktraceKey: cfg.StacktraceKey,
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:    zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		//EncodeTime:     zapcore.RFC3339TimeEncoder,
		//EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		//EncodeCaller:   zapcore.FullCallerEncoder,        // 全路径编码器
		EncodeCaller: zapcore.ShortCallerEncoder, // 当前路径编码器
		EncodeName:   zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	switch cfg.Level {
	case 0:
		atomicLevel.SetLevel(zap.DebugLevel)
	case 1:
		atomicLevel.SetLevel(zap.InfoLevel)
	case 2:
		atomicLevel.SetLevel(zap.WarnLevel)
	case 3:
		atomicLevel.SetLevel(zap.ErrorLevel)
	default:
		atomicLevel.SetLevel(zap.InfoLevel)
	}

	var encoder zapcore.Encoder
	if cfg.Format == 1 {
		encoder = zapcore.NewJSONEncoder(encoderConfig) // 编码器配置,Json格式
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	core := zapcore.NewCore(
		encoder, // 普通的一行字符格式
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	// 开启开发模式，堆栈跟踪,把文件和行号记录下来
	// 注意下面这个写法，实现了Field的中间件写法,skip往上面跳一帧，函数
	caller := zap.AddCaller()
	callerSkip := zap.AddCallerSkip(1)
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	//filed := zap.Fields(zap.String("ServerName", "Let is go"))
	// 构造日志
	Logger := zap.New(core, caller, callerSkip, development)
	sugar = Logger.Sugar()
	//Logger.Info("msg", zap.String())
}
