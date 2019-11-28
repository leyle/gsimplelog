package gsimplelog

import (
	"errors"
	"time"
)

const (
	LogModeOff     = "off"
	LogModeConsole = "console"
	LogModeFile    = "file"
)

type ILogConfig interface {
	GetLogLevel() string
}

func InitLogger(logLevel int, logMode, logPath, logPrefix string, rotateSize int) (err error) {
	var l ILogger
	switch logMode {
	case LogModeOff:
		l, err = NewSkipLogger()
		if err != nil {
			return errors.New("init logger failed:" + err.Error())
		}
	case LogModeConsole:
		l, err = NewStdLogger(logLevel)
		if err != nil {
			return errors.New("init logger failed:" + err.Error())
		}
	case LogModeFile:
		//multiSize: 50MB
		l, err = NewFileLogger(logPrefix, logPath, logLevel, rotateSize)
		if err != nil {
			return errors.New("init logger failed:" + err.Error())
		}
	default:
		return errors.New("not support LogMode:" + logMode)
	}
	Logger = l
	return
}

func ApplyConfig(logConfig ILogConfig) error {
	levelFlag, ok := LevelMap[logConfig.GetLogLevel()]
	if !ok {
		return errors.New("not support LogLevel:" + logConfig.GetLogLevel())
	}
	Logger.SetLevel(levelFlag)
	return nil
}

var Logger ILogger = &StdLogger{Level: LogDebug}

func SetLogger(logger ILogger) {
	Logger = logger
}

const (
	LogTrace = 0
	LogDebug = 1
	LogInfo  = 2
	LogError = 3
	LogOff   = 4
)

var LevelMap = map[string]int{
	"trace": LogTrace,
	"debug": LogDebug,
	"info":  LogInfo,
	"error": LogError,
	"off":   LogOff,
}

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

type ILogger interface {
	SetLevel(int)
	Trace(...interface{})
	Debug(...interface{})
	Info(...interface{})
	Error(...interface{})
	Tracef(string, ...interface{})
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
	Close() error
}
