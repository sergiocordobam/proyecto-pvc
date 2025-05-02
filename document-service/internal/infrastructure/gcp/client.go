package gcp

import (
	"context"
	"errors"
	"sync"
	"time"

	"cloud.google.com/go/storage"
)

type StorageClient struct {
	Client *storage.Client
}

var (
	singleClient *StorageClient
	once         sync.Once
)

func NewStorageClient(ctx context.Context) (*StorageClient, error) {
	var err error
	var gcpClient *storage.Client
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
			Client: gcpClient,
		}
	},
	)
	if err != nil {
		return nil, err
	}
	return singleClient, nil
}

func (s StorageClient) WriteObject(ctx context.Context, path string, content []byte) error {
	//TODO implement me
	panic("implement me")
}

func (s StorageClient) ReadObject(ctx context.Context, path string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (s StorageClient) ObjectExists(ctx context.Context, prefix string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s StorageClient) DeleteObject(ctx context.Context, path string) error {
	//TODO implement me
	panic("implement me")
}

func (s StorageClient) SetObjectAttributes(ctx context.Context, path string, attrs storage.ObjectAttrsToUpdate) error {
	//TODO implement me
	panic("implement me")
}

func (s StorageClient) GenerateSignedURL(path string, method string, expiry time.Duration) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s StorageClient) ListObjectsWithPrefix(ctx context.Context, prefix string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (s StorageClient) UserExits(ctx context.Context, userID uint64) bool {
	//TODO implement me
	panic("implement me")
}
