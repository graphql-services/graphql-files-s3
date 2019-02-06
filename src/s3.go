package src

import (
	"io"
	"os"

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

// UploadToS3Config ...
type UploadToS3Config struct {
	Bucket      string
	Key         string
	Body        io.ReadSeeker
	Size        int64
	ContentType string //http.DetectContentType(buffer)
}

// UploadToS3 ...
func UploadToS3(c UploadToS3Config) error {

	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(c.Bucket),
		Key:                aws.String(c.Key),
		ACL:                aws.String("private"),
		Body:               c.Body,
		ContentLength:      aws.Int64(c.Size),
		ContentType:        aws.String(c.ContentType),
		ContentDisposition: aws.String("attachment"),
		// ServerSideEncryption: aws.String("AES256"),
	})

	return err
}

// GetS3ObjectConfig ...
type GetS3ObjectConfig struct {
	Bucket string
	Key    string
}

// GetS3Object ...
func GetS3Object(c GetS3ObjectConfig) (*s3.GetObjectOutput, error) {
	res, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(c.Key),
	})

	return res, err
}
