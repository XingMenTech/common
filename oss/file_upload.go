package oss

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gitlab.novgate.com/common/common/logger"
	"io"
)

func init() {
	Register(AliPlatformCode, NewAliyunAdapter)
	Register(AwsPlatformCode, NewAwsAdapter)
}

type FileUploadAdapter interface {
	Upload(src io.Reader, name, uploadPath string) (path string, err error)
	StartAndGC(config interface{}) error
}

type newAdapterFunc func() FileUploadAdapter

var adapters = make(map[string]newAdapterFunc)

func Register(name string, adapter newAdapterFunc) {
	if adapter == nil {
		panic("upload: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("upload: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

func NewFileUploadAdapter(platformCode string, config interface{}) (FileUploadAdapter, error) {

	instanceFunc, ok := adapters[platformCode]
	if !ok {
		logger.LOG.WithFields(logrus.Fields{"platformCode": platformCode}).Error("unexpected platform code")
		return nil, errors.New("invalid platform code")
	}

	adapter := instanceFunc()
	err := adapter.StartAndGC(config)
	if err != nil {
		adapter = nil
	}
	return adapter, nil
}
