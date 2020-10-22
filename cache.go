package main

import (
	"fmt"
	"sync"
)

type Cache interface {
	Insert(id int, url string) error
	Get(id int) (string, error)
}
type localCache struct {
	urls map[int]string
	sync.Mutex
}

func NewLocalCache() *localCache {
	return &localCache{urls: map[int]string{}}
}
func (lc *localCache) Insert(id int, url string) error {
	lc.Lock()
	defer lc.Unlock()

	if _, ok := lc.urls[id]; ok != false {
		return fmt.Errorf("Item with ID &v already existsin cache")
	}

	lc.urls[id] = url

	return nil
}

func (lc *localCache) Get(id int) (string, error) {
	lc.Lock()
	defer lc.Unlock()

	if url, ok := lc.urls[id]; ok {
		return url, nil
	}

	return "", fmt.Errorf("Item &v is not present in cache")
}
