package oss

import (
	"errors"
	"io"

	"github.com/sirupsen/logrus"
	"gitlab.novgate.com/common/common/logger"
)

type Config struct {
	Endpoint   string `yaml:"endpoint" json:"endpoint" comment:"接口地址"`
	AccessId   string `yaml:"access_id" json:"accessId" comment:"accessId"`
	AccessKey  string `yaml:"access_key" json:"accessKey" comment:"accessKey"`
	BucketName string `yaml:"bucket" json:"bucketName" comment:"存储桶"`
	OssUrl     string `yaml:"oss_url" json:"ossUrl" comment:"CDN域名"`
	Region     string `yaml:"region" json:"region" comment:"区域"`
}

func init() {
	Register(AliPlatformCode, NewAliyunAdapter)
	Register(AwsPlatformCode, NewAwsAdapter)
}

type FileUploadAdapter interface {
	Upload(src io.Reader, name, uploadPath string) (path string, err error)
	startAndGC(config *Config) error
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

func NewOssAdapter(platformCode string, config *Config) (FileUploadAdapter, error) {

	instanceFunc, ok := adapters[platformCode]
	if !ok {
		logger.LOG.WithFields(logrus.Fields{"platformCode": platformCode}).Error("unexpected platform code")
		return nil, errors.New("invalid platform code")
	}

	adapter := instanceFunc()
	err := adapter.startAndGC(config)
	if err != nil {
		adapter = nil
	}
	return adapter, nil
}
