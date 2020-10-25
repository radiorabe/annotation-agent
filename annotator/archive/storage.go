package archive

import (
	"context"
	"log"

	"github.com/sirupsen/logrus"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	uuid "github.com/satori/go.uuid"
)

// StorageInterface ...
type StorageInterface interface {
	Init() StorageInterface
	Store(url string, filePath string, contentType string) (string, error)
}

// Storage ...
type Storage struct {
	bucketName string

	client *minio.Client
}

// NewStorage ...
func NewStorage(endpoint string, accessKey string, accessSecret string, useSSL bool, bucketName string) StorageInterface {
	return &Storage{
		bucketName: bucketName,
		client: getStorageClient(
			endpoint,
			accessKey,
			accessSecret,
			useSSL,
		),
	}
}

// Init ...
func (s *Storage) Init() StorageInterface {
	log := logrus.
		WithField("bucketName", s.bucketName)

	exists, err := s.client.BucketExists(context.TODO(), s.bucketName)
	if err != nil {
		log.WithError(err).Error("Failed to check if bucket exists.")
	}
	if exists {
		log.Info("Bucket exists.")
		return s
	}
	err = s.client.MakeBucket(context.TODO(), s.bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		log.Fatal("Failed to create bucket.")
	}
	return s
}

// Store ...
func (s *Storage) Store(url string, filePath string, contentType string) (string, error) {
	objectName := uuid.NewV5(uuid.NamespaceURL, url).String()
	_, err := s.client.FPutObject(
		context.TODO(),
		s.bucketName,
		objectName,
		filePath,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		logrus.WithError(err).Fatalln(err)
		return "", err
	}

	return objectName, nil
}

func getStorageClient(endpoint string, accessKey string, accessSecret string, useSSL bool) *minio.Client {
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, accessSecret, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return minioClient
}
