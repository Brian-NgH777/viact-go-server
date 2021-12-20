package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"mime/multipart"
)

const (
	AWS_S3_REGION = "ap-southeast-1"
	AWS_S3_BUCKET = "viact-bridge"
	AWS_S3_ACCESS_KEY_ID = "AKIAZRF7Z3PJ5HAOTROS"
	AWS_S3_SECRET_ACCESS_KEY = "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"
	//S3_ACL        = "public-read"
)

func UploadFile(file *multipart.FileHeader) error {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(AWS_S3_REGION),
		Credentials: credentials.NewStaticCredentials(
			AWS_S3_ACCESS_KEY_ID,
			AWS_S3_SECRET_ACCESS_KEY,
			""),
	})
	if err != nil {
		log.Fatal(err)
	}

	content, err := file.Open()
	// Upload Files
	_, err = s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(AWS_S3_BUCKET),
		Body:   content,
		Key:    aws.String(file.Filename),
		//ACL:    aws.String(S3_ACL),
	})
	return err
}
