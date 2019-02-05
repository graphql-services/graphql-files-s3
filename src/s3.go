package src

import (
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

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

	s, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("S3_REGION"))})
	if err != nil {
		return err
	}

	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(c.Bucket),
		Key:                  aws.String(c.Key),
		ACL:                  aws.String("private"),
		Body:                 c.Body,
		ContentLength:        aws.Int64(c.Size),
		ContentType:          aws.String(c.ContentType),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})

	return err
}
