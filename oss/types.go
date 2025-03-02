package oss

const (
	AliPlatformCode = "ali_oss"
	AwsPlatformCode = "aws_s3"
)

// 图片格式
const (
	GameImagesGif = ".gif"
	GameImageJpeg = ".jpeg"
	GameImagePng  = ".png"
	GameImageJpg  = ".jpg"
	GameImageBmp  = ".bmp"
	GameImageZip  = ".zip"
)

var GameImageType = map[string]int{
	GameImagesGif: 1,
	GameImageJpeg: 2,
	GameImagePng:  3,
	GameImageJpg:  4,
	GameImageBmp:  5,
	GameImageZip:  6,
}

var GameImageContentType = map[string]string{
	GameImagesGif: "image/gif",
	GameImageJpeg: "image/jpeg",
	GameImagePng:  "image/png",
	GameImageJpg:  "image/jpg",
	GameImageBmp:  "image/bmp",
	GameImageZip:  "application/zip",
}
