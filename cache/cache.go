package cache

import (
	"encoding/gob"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
)

type Cache interface {
	Get(key string, value any) error
	Put(key string, value any) error
	Remove(key string) error
	Close() error
}

type GoogleStorageCache struct {
	ctx           *context.Context
	storageClient *storage.Client
	bucket        *storage.BucketHandle
}

func NewGoogleStorageCache(ctx *context.Context, bucket string) (Cache, error) {
	sc, err := storage.NewClient(*ctx)
	if err != nil {
		return nil, err
	}
	return &GoogleStorageCache{ctx: ctx, storageClient: sc, bucket: sc.Bucket(bucket)}, nil
}

func (r *GoogleStorageCache) Close() error {
	return r.storageClient.Close()
}

func (r *GoogleStorageCache) Get(key string, value any) error {
	rc, err := r.bucket.Object(key).NewReader(*r.ctx)
	if err != nil {
		return err
	}
	defer rc.Close()
	return gob.NewDecoder(rc).Decode(value)
}

func (r *GoogleStorageCache) Put(key string, value any) error {
	rc := r.bucket.Object(key).NewWriter(*r.ctx)
	defer rc.Close()
	return gob.NewEncoder(rc).Encode(value)
}

func (r *GoogleStorageCache) Remove(key string) error {
	return r.bucket.Object(key).Delete(*r.ctx)
}
