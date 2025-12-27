package awss3

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"gorm.io/gorm/logger"
)

type S3Storage struct {
	prefix string
	s3     *s3.S3
	mode   logger.LogLevel
	params StorageConfigureParams
}

func (s *S3Storage) Get() interface{} {
	return s.s3
}

func (s *S3Storage) Run() error {
	//key := "DO0027CKHVN42V7ZXU99"                           // Access key pair. You can create access key pairs using the control panel or API.
	//secret := "E3sA383lgwq6PUlbCPzUm45FOaK0vHBkzdVkmC8qmxg" // Secret access key defined through an environment variable.

	key := s.params.AccessKey
	secret := s.params.SecretKey
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String(s.params.Endpoint),
		Region:      aws.String(s.params.Region),
	}

	// Step 3: The new session validates your request and directs it to your Space's specified endpoint using the AWS SDK.
	newSession, _ := session.NewSession(s3Config)
	s3Client := s3.New(newSession)
	s.s3 = s3Client
	return nil
}

func (s *S3Storage) Configure(prefix string, params StorageConfigureParams) error {
	s.prefix = prefix
	s.params = params
	return nil
}

func (s *S3Storage) GetPrefix() string {
	return s.prefix
}

func (s *S3Storage) Stop() <-chan bool {
	stop := make(chan bool)
	go func() {
		stop <- true
	}()
	return stop
}

func (s *S3Storage) UploadFile(file *multipart.FileHeader, bucketName string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	buffer := make([]byte, file.Size)
	_, err = src.Read(buffer)
	if err != nil {
		return "", err
	}

	ext := filepath.Ext(file.Filename)
	uniqueID := uuid.New().String()
	filename := fmt.Sprintf("%s/%s%s", s.prefix, uniqueID, ext)

	_, err = s.s3.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(buffer),
		ACL:         aws.String("public-read"),
		ContentType: aws.String(file.Header.Get("Content-Type")),
	})

	if err != nil {
		return "", err
	}

	fileURL := fmt.Sprintf("https://%s.%s/%s", bucketName, s.params.Endpoint, filename)
	return fileURL, nil
}
