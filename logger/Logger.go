// 日志单实例
// Neo
// 提供文件日志分割功能，提供日志目录大小监视功能(当日志文件目录大小超过最大值，优先删除最早的历史日志)

package logger

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

var (
	LOG *logrus.Logger
)

const (
	defaultLogDirName = "logs"   //日志目录
	maxLogFileSize    = 10 << 20 //最大日志文件大小10MB
	maxLogDirSize     = 1 << 30  //最大日志文件目录大小1GB
)

// --------------------[log]-------------------
type LogConfig struct {
	Path  string `yaml:"path" json:"path" comment:"日志文件路径"`
	Level int    `yaml:"level" json:"level" comment:"日志级别"`
}

// LogrusFileLoggerHook /文件日志Hook
type LogrusFileLoggerHook struct {
	maxLogSize int64
	logDir     string
	logPath    string
	level      logrus.Level
	logFile    *os.File
}

// NewLogrusFileLoggerHook /工厂方法
func NewLogrusFileLoggerHook(logDir string, maxLogSize int64, level logrus.Level) (hook *LogrusFileLoggerHook, err error) {
	object := &LogrusFileLoggerHook{
		maxLogSize: maxLogSize,
		logDir:     logDir,
		level:      level,
	}
	return object, object.makeLogFile()
}

// /创建日志文件
func (object *LogrusFileLoggerHook) makeLogFile() error {
	var err error
	if filepath.IsAbs(object.logDir) {
		object.logPath = object.logDir
	} else {
		dir, err := os.Getwd()
		if nil != err {
			panic(err)
		}
		object.logPath = filepath.Join(dir, object.logDir)
	}
	err = os.MkdirAll(object.logDir, 0777)
	if nil != err {
		panic(err)
	}
	logTAG := "UNKNOWN"
	switch object.level {
	case logrus.PanicLevel:
		logTAG = "PANIC"
	case logrus.FatalLevel:
		logTAG = "FATAL"
	case logrus.ErrorLevel:
		logTAG = "ERROR"
	case logrus.WarnLevel:
		logTAG = "WARN"
	case logrus.InfoLevel:
		logTAG = "INFO"
	case logrus.DebugLevel:
		logTAG = "DEBUG"
	case logrus.TraceLevel:
		logTAG = "TRACE"
	}
	now := time.Now()
	object.logPath += fmt.Sprintf("%s%s-%s-%04d-%02d-%02dT%02d-%02d-%02d.%d.log",
		string(filepath.Separator),
		logTAG,
		"BACKEND",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		os.Getpid())
	if nil != object.logFile {
		err = object.logFile.Close()
		if nil != err {
			return err
		}
	}
	object.logFile, err = os.OpenFile(object.logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if nil != err {
		return err
	}
	return nil
}

// Levels /日志等级回调
func (object *LogrusFileLoggerHook) Levels() []logrus.Level {
	return []logrus.Level{object.level}
}

// Fire /激发
func (object *LogrusFileLoggerHook) Fire(entry *logrus.Entry) error {
	if size, err := FileSize(object.logPath); nil != err {
		return err
	} else if size > object.maxLogSize {
		if err = object.makeLogFile(); nil != err {
			return err
		}
	}
	content, err := entry.String()
	if nil != err {
		return err
	}
	_, err = object.logFile.Write([]byte(content))
	if nil != err {
		return err
	}
	return nil
}

// /清理空日志文件
func cleanEmptyLogFile(logDir string) {
	arr, err := os.ReadDir(logDir)
	if nil != err {
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "ioutil.ReadDir errors: %s\n", err)
		}
		return
	}
	for _, fi := range arr {
		info, err := fi.Info()
		if nil != err {
			fmt.Fprintf(os.Stderr, "filepath.Abs errors: %s\n", err)
			continue
		}

		if 0 == info.Size() {
			absPath, err := filepath.Abs(filepath.Join(logDir, fi.Name()))
			if nil != err {
				fmt.Fprintf(os.Stderr, "filepath.Abs errors: %s\n", err)
				continue
			}
			if err = os.Remove(absPath); nil != err {
				fmt.Fprintf(os.Stderr, "os.Remove %s errors: %s\n", absPath, err)
				continue
			}
		}
	}
}

// InitializeLogger /初始化日志
func InitializeLogger(logConfig *LogConfig) {
	if nil == LOG {
		//清理空文件
		logDir := logConfig.Path

		if "" == logDir {
			logDir = defaultLogDirName
		}
		cleanEmptyLogFile(logDir)

		//日志等级
		logLevelConfig := int(logrus.DebugLevel)
		if logConfig.Level > 0 {
			logLevelConfig = logConfig.Level
		}
		logLevel := logrus.Level(logLevelConfig)

		//初始化日志对象
		LOG = logrus.New()
		LOG.SetLevel(logLevel)
		LOG.SetFormatter(&Formatter{TimestampFormat: time.RFC3339})
		//if config.GrayLogConfig != nil {
		//	graylogHook := graylog.NewGraylogHook(config.GrayLogConfig.Api, map[string]interface{}{
		//		"stream": config.GrayLogConfig.IdaBackend,
		//	})
		//	LOG.AddHook(graylogHook)
		//}

		//创建Hook
		for level := logrus.PanicLevel; level <= logLevel; level++ {
			hook, err := NewLogrusFileLoggerHook(logDir, maxLogFileSize, level)
			if nil != err {
				panic(err)
			}
			LOG.AddHook(hook)
		}
	}
}

type errorHook struct{}

func (hook *errorHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel}
}

func (hook *errorHook) Fire(entry *logrus.Entry) error {
	return nil
}

// /文件大小
func FileSize(path string) (size int64, err error) {
	if 0 >= len(path) {
		err = errors.New("invalid path")
		return
	}

	var fi os.FileInfo
	fi, err = os.Stat(path)
	if nil == err {
		size = fi.Size()
	}

	return
}
