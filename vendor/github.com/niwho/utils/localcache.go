package utils

import (
	"time"

	"github.com/niwho/inmem"
)

// key=value 通用存储
type LocalCache struct {
	localCache inmem.Cache
	size       int
	expiretime time.Duration
	hit        int64
	count      int64
}

func NewLocalCache(size int, expireTime time.Duration) *LocalCache {
	localCache := &LocalCache{
		size:       size,
		expiretime: expireTime,
	}
	localCache.localCache = inmem.NewLocked(size)
	return localCache
}

// 这个函数可能有性能瓶颈
func (localCache *LocalCache) IsHit(key interface{}) bool {
	// true 表示重复了，有冲突
	return localCache.localCache.AddIfNotExist(key, struct{}{}, time.Now().Add(localCache.expiretime))

}

//
func (localCache *LocalCache) Get(key interface{}) (interface{}, bool) {
	return localCache.localCache.Get(key)
}

func (localCache *LocalCache) Remove(key interface{}) {
	localCache.localCache.Remove(key)
}

func (localCache *LocalCache) Set(key, val interface{}) {
	localCache.localCache.Add(key, val, time.Now().Add(localCache.expiretime))
}

func (localCache *LocalCache) SetWithTtl(key, val interface{}, ttl int) {
	localCache.localCache.Add(key, val, time.Now().Add(time.Duration(ttl)*time.Second))
}

func (localCache *LocalCache) Len() int {
	return localCache.localCache.Len()
}

// key=value 通用存储
type CommonCache struct {
	localCache inmem.Cache
	size       int
	expiretime time.Duration
	hit        int64
	count      int64
}

func NewCommonCache(size int, expireTime time.Duration) *CommonCache {
	localCache := &CommonCache{
		size:       size,
		expiretime: expireTime,
	}
	localCache.localCache = inmem.NewLocked(size)
	return localCache
}

// 这个函数可能有性能瓶颈
func (localCache *CommonCache) IsHit(key interface{}) bool {
	// true 表示重复了，有冲突
	return localCache.localCache.AddIfNotExist(key, struct{}{}, time.Now().Add(localCache.expiretime))

}

//
func (localCache *CommonCache) Get(key interface{}) (interface{}, bool) {
	return localCache.localCache.Get(key)
}

func (localCache *CommonCache) Remove(key interface{}) {
	localCache.localCache.Remove(key)
}

func (localCache *CommonCache) Set(key, val interface{}) {
	localCache.localCache.Add(key, val, time.Now().Add(localCache.expiretime))
}

func (localCache *CommonCache) SetWithTtl(key, val interface{}, ttl int) {
	localCache.localCache.Add(key, val, time.Now().Add(time.Duration(ttl)*time.Second))
}

func (localCache *CommonCache) Len() int {
	return localCache.localCache.Len()
}
