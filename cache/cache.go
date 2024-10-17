package cache

import (
	"encoding/gob"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
)

type Cache struct {
	ctx           *context.Context
	storageClient *storage.Client
	bucket        *storage.BucketHandle
}

func NewCache(ctx *context.Context, bucket string) (*Cache, error) {
	sc, err := storage.NewClient(*ctx)
	if err != nil {
		return nil, err
	}
	return &Cache{ctx: ctx, storageClient: sc, bucket: sc.Bucket(bucket)}, nil
}

func (r *Cache) Close() error {
	return r.storageClient.Close()
}

func (r *Cache) Get(key string, value any) error {
	rc, err := r.bucket.Object(key).NewReader(*r.ctx)
	if err != nil {
		return err
	}
	defer rc.Close()
	return gob.NewDecoder(rc).Decode(value)
}

func (r *Cache) Put(key string, value any) error {
	rc := r.bucket.Object(key).NewWriter(*r.ctx)
	defer rc.Close()
	return gob.NewEncoder(rc).Encode(value)
}

func (r *Cache) Remove(key string) error {
	return r.bucket.Object(key).Delete(*r.ctx)
}
