package src

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	s3Client *s3.S3
)

func init() {
	region := os.Getenv("S3_REGION")
	s, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		panic(err)
	}
	s3Client = s3.New(s)
}

func PutObjectPresignedURL(bucket, key string) (url string, err error) {
	req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return req.Presign(15 * time.Minute)
}
func GetObjectPresignedURL(bucket, key, contentDisposition string) (url string, err error) {
	req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket:                     aws.String(bucket),
		Key:                        aws.String(key),
		ResponseContentDisposition: aws.String(contentDisposition),
	})
	return req.Presign(15 * time.Minute)
}
