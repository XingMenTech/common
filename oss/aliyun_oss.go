package oss

import (
	"errors"
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"gitlab.novgate.com/common/common/logger"
)

// ------------------[oss]------------------

type AliyunAdapter struct {
	config *Config
}

func NewAliyunAdapter() FileUploadAdapter {
	return &AliyunAdapter{}
}

func (a *AliyunAdapter) Upload(src io.Reader, name, uploadPath string) (path string, err error) {
	//获取oss服务器信息
	endpoint := a.config.Endpoint
	accessKeyId := a.config.AccessId
	accessKeySecret := a.config.AccessKey
	bucketName := a.config.BucketName

	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		logger.LOG.Errorf("创建oss失败%v", err)
		return
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		logger.LOG.Errorf("使用oss空间失败%v", err)
		return
	}

	fileDir := uploadPath + "/" + name
	err6 := bucket.PutObject(fileDir, src)
	if err6 != nil {
		err = err6
		logger.LOG.Errorf("上传oss失败%v", err6)
		return
	}

	ossUrl := a.config.OssUrl
	if ossUrl == "" {
		ossUrl = a.config.Endpoint
	}

	path = ossUrl + fileDir

	return
}

func (a *AliyunAdapter) startAndGC(config *Config) error {
	if config == nil {
		return errors.New("aliyun oss config invalid")
	}
	a.config = config
	return nil
}
