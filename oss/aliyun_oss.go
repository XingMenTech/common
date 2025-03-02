package oss

import (
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/zhangyuanCloud/common/logger"
	"io"
	"os"
	"path/filepath"
)

// ------------------[oss]------------------
type AliyunOssConfigStruct struct {
	Endpoint   string `yaml:"endpoint"`
	AccessId   string `yaml:"access_id"`
	AccessKey  string `yaml:"access_key"`
	BucketName string `yaml:"bucket"`
	OssUrl     string `yaml:"oss_url"`
}

type AliyunAdapter struct {
	config *AliyunOssConfigStruct
}

func NewAliyunAdapter() FileUploadAdapter {
	return &AliyunAdapter{}
}

func (a *AliyunAdapter) Upload(src io.Reader, name, uploadPath string) (path string, err error) {

	filePath, err := filepath.Abs(name)
	if err != nil {
		logger.LOG.Errorf("文件路径失败%v", err)
		return
	}

	imageFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		logger.LOG.Errorf("打开文件失败%v", err)
		return
	}

	_, err = io.Copy(imageFile, src)
	if err != nil {
		logger.LOG.Errorf("拷贝图片失败%v", err)
		return
	}

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
	err6 := bucket.PutObjectFromFile(fileDir, filePath)
	if err6 != nil {
		err = err6
		logger.LOG.Errorf("上传oss失败%v", err6)
		return
	}

	path = a.config.OssUrl + fileDir

	err7 := os.Remove(filePath)
	if err7 != nil {
		logger.LOG.Errorf("删除图片失败,原因%v", err7)
	}
	return
}

func (a *AliyunAdapter) StartAndGC(config interface{}) error {
	if config == nil {
		return errors.New("aliyun oss config invalid")
	}
	a.config = config.(*AliyunOssConfigStruct)
	return nil
}
