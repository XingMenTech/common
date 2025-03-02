package logger

import (
	"context"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

var (
	logFileCleaner *LogFileCleaner
)

// LogFileCleaner /日志清理器
type LogFileCleaner struct {
	logDir         string
	maxLogFileSize int64
	maxLogDirSize  int64
	wg             *sync.WaitGroup
	ctx            context.Context
	cancel         context.CancelFunc
	log            *logrus.Entry
}

// FileInfoArray / 文件信息排序
type FileInfoArray []os.FileInfo

func (object FileInfoArray) Len() int { return len(object) }
func (object FileInfoArray) Less(i, j int) bool {
	return object[i].ModTime().Unix() < object[j].ModTime().Unix()
}
func (object FileInfoArray) Swap(i, j int) { object[i], object[j] = object[j], object[i] }

// NewLogFileCleaner /工厂方法
func NewLogFileCleaner(logDir string, maxLogFileSize, maxLogDirSize int64) *LogFileCleaner {
	object := &LogFileCleaner{
		logDir:         logDir,
		maxLogFileSize: maxLogFileSize,
		maxLogDirSize:  maxLogDirSize,
		wg:             &sync.WaitGroup{},
		log:            LOG.WithField("module", "LogFileCleaner"),
	}
	object.ctx, object.cancel = context.WithCancel(context.Background())
	return object
}

// Name /名字
func (object *LogFileCleaner) Name() string {
	return "LogFileCleaner"
}

// ShutdownPriority /关闭优先级
func (object *LogFileCleaner) ShutdownPriority() int {
	return 5
}

// BeforeShutdown /关闭钩子
func (object *LogFileCleaner) BeforeShutdown() {
	object.close()
}

func (object *LogFileCleaner) AfterShutdown() {
}

// / 关闭
func (object *LogFileCleaner) close() {
	object.cancel()
	object.wg.Wait()
	object.log.Infof("log file cleaner closed")
}

// CheckLogFileSizeAndRemove / 删除历史日志
func (object *LogFileCleaner) CheckLogFileSizeAndRemove() {
	object.wg.Add(1)
loop:
	for {
		select {
		case <-object.ctx.Done():
			break loop
		case <-time.After(5 * time.Minute):
			size, err := DirSize(object.logDir)
			if err != nil {
				object.log.Errorf("dir size: %s, errors: %s\n", object.logDir, err)
				continue
			}
			if size < object.maxLogDirSize {
				continue
			}

			arr, err := ioutil.ReadDir(object.logDir)
			if nil != err {
				object.log.Errorf("read dir: %s, errors: %s\n", object.logDir, err)
				continue
			}
			fiArr := FileInfoArray(arr)
			sort.Sort(fiArr)

			for _, fi := range fiArr {
				size, _ := DirSize(object.logDir)
				if size < object.maxLogDirSize {
					break
				}
				filePath, err := filepath.Abs(filepath.Join(object.logDir, fi.Name()))
				if nil != err {
					object.log.Errorf("filepath abs errors: %s\n", err)
					continue
				}
				fi, err = os.Stat(filePath)
				if err != nil {
					continue
				}
				if fi.IsDir() {
					err = os.RemoveAll(filePath)
				} else {
					err = os.Remove(filePath)
				}
				object.log.Warn("remove:%s errors:%s", filePath, err)
			}
		}
	}
	object.wg.Done()
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
