package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type StorageService interface {
	UploadFile(ctx context.Context, file io.Reader, fileName string, contentType string) (*UploadResult, error)
	DeleteFile(ctx context.Context, storageKey string) error
	GetFileURL(storageKey string) string
	GetSignedURL(storageKey string, expiration time.Duration) (string, error)
}

type UploadResult struct {
	URL        string
	StorageKey string
	FileName   string
	FileSize   int64
}

type S3Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	Region          string
	UseSSL          bool
	PublicURL       string
}

type s3Storage struct {
	client     *s3.S3
	bucketName string
	publicURL  string
}

func NewS3Storage(cfg S3Config) (StorageService, error) {
	awsConfig := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Region:           aws.String(cfg.Region),
		S3ForcePathStyle: aws.Bool(true),
	}

	if cfg.Endpoint != "" {
		awsConfig.Endpoint = aws.String(cfg.Endpoint)
	}

	if !cfg.UseSSL {
		awsConfig.DisableSSL = aws.Bool(true)
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 session: %w", err)
	}

	client := s3.New(sess)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = client.HeadBucketWithContext(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(cfg.BucketName),
	})
	if err != nil {
		_, err = client.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(cfg.BucketName),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &s3Storage{
		client:     client,
		bucketName: cfg.BucketName,
		publicURL:  cfg.PublicURL,
	}, nil
}

func (s *s3Storage) UploadFile(ctx context.Context, file io.Reader, fileName string, contentType string) (*UploadResult, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	fileSize := int64(len(data))

	if fileSize > 10*1024*1024 {
		return nil, fmt.Errorf("file size exceeds 10MB limit")
	}

	ext := filepath.Ext(fileName)
	storageKey := fmt.Sprintf("resources/%s/%s%s",
		time.Now().Format("2006/01"),
		uuid.New().String(),
		ext,
	)

	_, err = s.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(storageKey),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
		ACL:         aws.String("public-read"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to S3: %w", err)
	}

	url := s.GetFileURL(storageKey)

	return &UploadResult{
		URL:        url,
		StorageKey: storageKey,
		FileName:   fileName,
		FileSize:   fileSize,
	}, nil
}

func (s *s3Storage) DeleteFile(ctx context.Context, storageKey string) error {
	_, err := s.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(storageKey),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}
	return nil
}

func (s *s3Storage) GetFileURL(storageKey string) string {
	if s.publicURL != "" {
		return fmt.Sprintf("%s/%s/%s", s.publicURL, s.bucketName, storageKey)
	}
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucketName, storageKey)
}

func (s *s3Storage) GetSignedURL(storageKey string, expiration time.Duration) (string, error) {
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(storageKey),
	})

	url, err := req.Presign(expiration)
	if err != nil {
		return "", fmt.Errorf("failed to create signed URL: %w", err)
	}

	return url, nil
}

type localStorage struct {
	baseURL  string
	basePath string
}

func NewLocalStorage(basePath, baseURL string) StorageService {
	return &localStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}
}

func (l *localStorage) UploadFile(ctx context.Context, file io.Reader, fileName string, contentType string) (*UploadResult, error) {
	storageKey := fmt.Sprintf("resources/%s/%s", time.Now().Format("2006/01"), fileName)
	return &UploadResult{
		URL:        fmt.Sprintf("%s/%s", l.baseURL, storageKey),
		StorageKey: storageKey,
		FileName:   fileName,
		FileSize:   0,
	}, nil
}

func (l *localStorage) DeleteFile(ctx context.Context, storageKey string) error {
	return nil
}

func (l *localStorage) GetFileURL(storageKey string) string {
	return fmt.Sprintf("%s/%s", l.baseURL, storageKey)
}

func (l *localStorage) GetSignedURL(storageKey string, expiration time.Duration) (string, error) {
	return l.GetFileURL(storageKey), nil
}
