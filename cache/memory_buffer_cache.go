package cache

import (
	"bytes"
	"encoding/gob"
)

type MemoryBufferCache struct {
	buf bytes.Buffer
}

func (c *MemoryBufferCache) Get(key string, value any) error {
	return gob.NewDecoder(&c.buf).Decode(value)
}

func (c *MemoryBufferCache) Put(key string, value any) error {
	c.buf.Reset()
	gob.NewEncoder(&c.buf).Encode(value)
	return nil
}

func (c *MemoryBufferCache) Remove(key string) error {
	c.buf.Reset()
	return nil
}

func (c *MemoryBufferCache) Close() error {
	return nil
}

func NewMemoryCache() Cache {
	return &MemoryBufferCache{}
}
