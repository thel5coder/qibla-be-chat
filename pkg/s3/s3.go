package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	s3pkg "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	uuid "github.com/satori/go.uuid"
	"mime/multipart"
	"strings"
	"time"
)

// Credential ...
type Credential struct {
	URL       string
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
}

// Upload ...
func (cred *Credential) Upload(path string, file *multipart.FileHeader) (res string, err error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(cred.URL),
		Region:      aws.String(cred.Region),
		Credentials: credentials.NewStaticCredentials(cred.AccessKey, cred.SecretKey, ""),
	})
	if err != nil {
		return res, err
	}

	uploader := s3manager.NewUploader(sess)

	f, err := file.Open()
	if err != nil {
		return res, err
	}
	defer f.Close()

	res = path + "/" + time.Now().Format("20060102T150405") + "_" + uuid.NewV4().String() + strings.Replace(file.Filename, " ", "", -1)

	// Upload the file to S3.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(cred.Bucket),
		Key:    aws.String(res),
		Body:   f,
	})
	if err != nil {
		return res, err
	}

	// res = aws.StringValue(&result.Location)

	return res, err
}

// GetURL ...
func (cred *Credential) GetURL(key string) (res string, err error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(cred.URL),
		Region:      aws.String(cred.Region),
		Credentials: credentials.NewStaticCredentials(cred.AccessKey, cred.SecretKey, ""),
	})
	if err != nil {
		return res, err
	}

	svc := s3pkg.New(sess)
	req, _ := svc.GetObjectRequest(&s3pkg.GetObjectInput{
		Bucket: aws.String(cred.Bucket),
		Key:    aws.String(key),
	})

	res, err = req.Presign(15 * time.Minute)

	return res, err
}

// Delete ...
func (cred *Credential) Delete(key string) (res string, err error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(cred.URL),
		Region:      aws.String(cred.Region),
		Credentials: credentials.NewStaticCredentials(cred.AccessKey, cred.SecretKey, ""),
	})
	if err != nil {
		return res, err
	}

	svc := s3pkg.New(sess)
	_, err = svc.DeleteObject(&s3pkg.DeleteObjectInput{
		Bucket: aws.String(cred.Bucket),
		Key:    aws.String(key),
	})

	return res, err
}
