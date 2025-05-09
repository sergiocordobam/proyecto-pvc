package gcp

import (
	"context"
	"document-service/internal/domain/models"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type StorageClient struct {
	Client     *storage.Client
	BucketName string
}

var (
	singleClient *StorageClient
	once         sync.Once
)

func NewStorageClient(ctx context.Context, BucketName string) (*StorageClient, error) {
	var err error
	var gcpClient *storage.Client
	if BucketName == "" {
		err = errors.New("bucket name is empty")
		return nil, err
	}
	if singleClient != nil {
		err = errors.New("gcp storage client already created")
		return nil, err
	}
	once.Do(func() {
		gcpClient, err = storage.NewClient(ctx)
		if err != nil {
			return
		}
		singleClient = &StorageClient{
			Client:     gcpClient,
			BucketName: BucketName,
		}
	},
	)
	if err != nil {
		return nil, err
	}
	return singleClient, nil
}

func (s *StorageClient) ReadObjectData(ctx context.Context, path string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StorageClient) ObjectExists(ctx context.Context, userID int, path string) bool {
	bkt := s.Client.Bucket(s.BucketName)
	userIDStr := strconv.Itoa(userID)
	folderName := fmt.Sprintf("%s/%s", userIDStr, path)
	it := bkt.Objects(ctx, &storage.Query{Prefix: folderName})
	_, errObject := it.Next()
	return errors.Is(errObject, iterator.Done)
}

func (s *StorageClient) DeleteObject(ctx context.Context, path string) error {
	obj := s.Client.Bucket(s.BucketName).Object(path)
	err := obj.Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageClient) SetObjectAttributes(ctx context.Context, objectHandler *storage.ObjectHandle, attrs storage.ObjectAttrsToUpdate) error {
	if _, err := objectHandler.Update(ctx, attrs); err != nil {
		return err
	}
	return nil
}

func (s *StorageClient) GenerateSignedURL(filename string, method string, metadata models.Metadata, expirationTime time.Time) (string, error) {
	methodMap := map[string]string{
		"up":   "PUT",
		"down": "GET",
	}
	if methodMap[method] == "" {
		return "", fmt.Errorf("invalid method: %s", method)
	}
	signedConfiguration := storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  methodMap[method],
		Expires: expirationTime,
		Headers: []string{
			"Content-Type",
			"x-goog-meta-status:" + metadata.Status,
			"x-goog-meta-document-type" + ":" + metadata.Type,
		},
	}
	url, err := s.Client.Bucket(s.BucketName).SignedURL(filename, &signedConfiguration)
	if err != nil {
		return "", fmt.Errorf("generateURL-Failed: Bucket(%q).SignedURL(%q, opts): %v", s.BucketName, filename, err)
	}
	return url, nil
}

func (s *StorageClient) ListObjectsWithPrefix(ctx context.Context, prefix string) ([]storage.ObjectAttrs, error) {
	query := &storage.Query{
		Prefix: prefix,
	}
	objectIterator := s.Client.Bucket(s.BucketName).Objects(ctx, query)
	objectInfoList := make([]storage.ObjectAttrs, 0)
	for {
		objAttrs, err := objectIterator.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating bucket objects: %w", err)
		}

		objectInfoList = append(objectInfoList, *objAttrs)
	}

	return objectInfoList, nil
}

func (s *StorageClient) Close() error {
	return s.Client.Close()
}
func (s *StorageClient) GetBucketPointer() *storage.BucketHandle {
	return s.Client.Bucket(s.BucketName)
}
