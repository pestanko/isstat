package core

import (
	"fmt"
	"io/ioutil"
	"path"
	"time"
)

type Cache struct {
	CacheDir string
	Disabled bool
}

func NewCache(cacheDir string) Cache {
	return Cache{CacheDir: cacheDir, Disabled: false}
}

func (cache *Cache) StoreFile(filename string, data []byte) error {
	fullpath := path.Join(cache.CacheDir, filename)
	return ioutil.WriteFile(fullpath, data, 0644)
}

func (cache *Cache) Store(name string, data []byte) error {
	fullname := fmt.Sprintf("%s-%s", name, GetCurrentTimestamp())
	return cache.StoreFile(fullname, data)
}

func GetCurrentTimestamp() string {
	return time.Now().Format("2006-01-02T15-04-05")
}