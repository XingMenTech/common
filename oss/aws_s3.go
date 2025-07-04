package oss

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"path/filepath"
	"strings"
)

// ------------------[aws s3]------------------
type AwsS3ConfigStruct struct {
	Endpoint   string `yaml:"endpoint" json:"endpoint" comment:"接口地址"`
	AccessId   string `yaml:"access_id" json:"accessId" comment:"accessId"`
	AccessKey  string `yaml:"access_key" json:"accessKey" comment:"accessKey"`
	BucketName string `yaml:"bucket" json:"bucketName" comment:"s3库名"`
	OssUrl     string `yaml:"oss_url" json:"ossUrl" comment:"CDN域名"`
	Region     string `yaml:"region" json:"region" comment:"区域"`
}

type AwsAdapter struct {
	config *AwsS3ConfigStruct
}

func NewAwsAdapter() FileUploadAdapter {
	return &AwsAdapter{}
}

func (a *AwsAdapter) Upload(src io.Reader, name, uploadPath string) (path string, err error) {
	conf := &aws.Config{
		Credentials: credentials.NewStaticCredentials(a.config.AccessId, a.config.AccessKey, ""),
		Region:      aws.String(a.config.Region),
	}
	sess, err := session.NewSession(conf)
	if err != nil {
		return
	}
	uploader := s3manager.NewUploader(sess)
	fileDir := uploadPath + "/" + name

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(a.config.BucketName),
		Key:         aws.String(fileDir),
		Body:        src,
		ContentType: aws.String(GameImageContentType[filepath.Ext(name)]),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	path = strings.Replace(result.Location, a.config.Endpoint, a.config.OssUrl, 1)

	return
}

func (a *AwsAdapter) StartAndGC(config interface{}) error {
	if config == nil {
		return errors.New("aws s3 config invalid")
	}
	a.config = config.(*AwsS3ConfigStruct)
	return nil
}
