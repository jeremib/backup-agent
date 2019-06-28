package aws

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func Upload(accessKey string, secret string, region string, bucket string, from string, to string) string {

	file, err := os.Open(from)
	defer file.Close()

	//select Region to use.
	conf := aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secret, ""),
	}
	sess := session.New(&conf)
	svc := s3manager.NewUploader(sess)
	t := time.Now()

	fmt.Println("Uploading file to S3...")
	result, err := svc.Upload(&s3manager.UploadInput{
		Bucket:  aws.String(bucket),
		Key:     aws.String(to),
		Body:    file,
		Tagging: aws.String("uploaded-at=" + t.Format(time.UnixDate)),
	})
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

	return result.Location
}
