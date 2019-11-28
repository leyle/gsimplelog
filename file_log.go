package gsimplelog

import (
	"fmt"
	"os"
	"sync"
	"path/filepath"
	"time"
	"io"
	"strconv"
)

func NewFileLogger(filePrefix, filePath string, level, multiSize int) (ILogger, error) {
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		fmt.Println("创建日志目录失败 ", err.Error())
		return nil, err
	}

	lf := &LogFile{
		path:      filePath,
		prefix:    filePrefix,
		multiSize: multiSize,
	}
	err = lf.Create()
	if err != nil {
		return nil, err
	}
	l := &FileLogger{
		Level: level,
		Out:   lf,
	}
	return l, nil
}

type FileLogger struct {
	Out   io.WriteCloser
	Level int
}

func (l *FileLogger) SetLevel(level int) {
	l.Level = level
}

func (l *FileLogger) Trace(params ...interface{}) {
	if l.Level <= LogTrace {
		l.Out.Write([]byte((fmt.Sprintf("%s [TRACE] %s\n", Now(), fmt.Sprint(params ...)))))
	}
}
func (l *FileLogger) Debug(params ...interface{}) {
	if l.Level <= LogDebug {
		l.Out.Write([]byte((fmt.Sprintf("%s [DEBUG] %s\n", Now(), fmt.Sprint(params ...)))))
	}
}
func (l *FileLogger) Info(params ...interface{}) {
	if l.Level <= LogInfo {
		l.Out.Write([]byte((fmt.Sprintf("%s [INFO] %s\n", Now(), fmt.Sprint(params ...)))))
	}
}
func (l *FileLogger) Error(params ...interface{}) {
	if l.Level <= LogError {
		l.Out.Write([]byte((fmt.Sprintf("%s [ERROR] %s\n", Now(), fmt.Sprint(params ...)))))
	}
}
func (l *FileLogger) Tracef(format string, params ...interface{}) {
	if l.Level <= LogTrace {
		l.Out.Write([]byte((fmt.Sprintf("%s [TRACE] %s\n", Now(), fmt.Sprintf(format, params...)))))
	}
}
func (l *FileLogger) Debugf(format string, params ...interface{}) {
	if l.Level <= LogDebug {
		l.Out.Write([]byte((fmt.Sprintf("%s [DEBUG] %s\n", Now(), fmt.Sprintf(format, params...)))))
	}
}
func (l *FileLogger) Infof(format string, params ...interface{}) {
	if l.Level <= LogInfo {
		l.Out.Write([]byte((fmt.Sprintf("%s [INFO] %s\n", Now(), fmt.Sprintf(format, params...)))))
	}
}
func (l *FileLogger) Errorf(format string, params ...interface{}) {
	if l.Level <= LogError {
		l.Out.Write([]byte((fmt.Sprintf("%s [ERROR] %s\n", Now(), fmt.Sprintf(format, params...)))))
	}
}

func (l *FileLogger) Close() error {
	return l.Out.Close()
}

type LogFile struct {
	path      string
	prefix 	  string
	file      *os.File
	multiSize int
	size      int
	index     int
	sync.RWMutex
}

func (l *LogFile) Write(b []byte) (int, error) {
	l.Lock()
	defer l.Unlock()
	n, err := l.file.Write(b)
	l.size += n
	if l.size > l.multiSize {
		l.file.Close()
		l.Create()
		l.size = 0
	}
	return n, err
}

func (l *LogFile) Create() error {
	l.index ++

	var filename string
	if l.index <= 1 {
		filename = l.prefix + "_" + time.Now().Format("20060102") + ".log"
	} else {
		filename = l.prefix + "_" + time.Now().Format("20060102") + "_" + strconv.Itoa(l.index) + ".log"
	}

	fullpath := filepath.Join(l.path, filename)
	file, err := os.OpenFile(fullpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	l.file = file
	return nil
}

func (l *LogFile) Close() error {
	l.Lock()
	defer l.Unlock()
	return l.file.Close()
}
