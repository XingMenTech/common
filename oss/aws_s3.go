package oss

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ------------------[aws s3]------------------

type AwsAdapter struct {
	config   *Config
	s3Client *s3.Client
}

func NewAwsAdapter() FileUploadAdapter {
	return &AwsAdapter{}
}

func (a *AwsAdapter) Upload(src io.Reader, name, uploadPath string) (path string, err error) {
	if a.s3Client == nil {
		return "", errors.New("s3 client not initialized")
	}

	fileDir := uploadPath + "/" + name
	contentType := GameImageContentType[filepath.Ext(name)]
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 创建上传输入
	input := &s3.PutObjectInput{
		Bucket:      aws.String(a.config.BucketName),
		Key:         aws.String(fileDir),
		Body:        src,
		ContentType: aws.String(contentType),
	}

	// 直接上传文件（S3 客户端会自动处理大文件）
	_, err = a.s3Client.PutObject(context.TODO(), input)
	if err != nil {
		fmt.Println("Upload error:", err)
		return "", err
	}

	// 构建返回路径
	path = fmt.Sprintf("%s/%s", a.config.OssUrl, fileDir)
	return path, nil
}

func (a *AwsAdapter) startAndGC(cfg *Config) error {
	if cfg == nil {
		return errors.New("aws s3 config invalid")
	}
	a.config = cfg

	// 创建 AWS 配置
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRetryMode(aws.RetryModeStandard),
		awsconfig.WithRetryMaxAttempts(3),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessId,
			cfg.AccessKey,
			"",
		)),
		awsconfig.WithRegion(cfg.Region),
	)
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %v", err)
	}

	// 如果提供了自定义 Endpoint，则配置它
	if cfg.Endpoint != "" {
		awsCfg.BaseEndpoint = aws.String(cfg.Endpoint)
	}

	// 创建 S3 客户端
	a.s3Client = s3.NewFromConfig(awsCfg)

	return nil
}
