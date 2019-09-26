package utils

import (
	"context"
	"fmt"
	"math"
	"sync"
)

type CacheManager struct {
	sync.Pool
	cacheChains GetDataChain
}

func NewCacheManager() *CacheManager {
	cm := &CacheManager{}
	cm.New = func() interface{} {
		return &CacheRequest{}
	}
	return cm
}

//context.WithValue(ctx, userKey, u)
//u, ok := ctx.Value(userKey).(*User)
func (cm *CacheManager) GetData(ctx context.Context, key interface{}) interface{} {
	c := cm.Get().(*CacheRequest)
	c.reset()
	c.ctx = ctx
	c.key = key
	c.cacheChains = cm.cacheChains

	c.Next()

	dat, _ := c.GetData()
	cm.Put(c)

	return dat
}

// 注意按添加的顺序依次执行
func (cm *CacheManager) AddChains(chains ...GetData) {
	cm.cacheChains = append(cm.cacheChains, chains...)
}

type GetData func(*CacheRequest)
type GetDataChain []GetData

const (
	// 有符号，所以是一半
	abortIndex int8 = math.MaxInt8 / 2
)

type CacheRequest struct {
	index       int8
	cacheChains GetDataChain

	key   interface{}
	data  interface{}
	isSet bool

	ctx context.Context
}

func (c *CacheRequest) reset() {
	c.index = -1
	c.cacheChains = nil
	c.key = ""
	c.data = ""
	c.isSet = false
}

func (c *CacheRequest) GetKey() interface{} {
	return c.key
}

func (c *CacheRequest) GetData() (interface{}, bool) {
	return c.data, c.isSet
}

func (c *CacheRequest) SetData(data interface{}) {
	c.data = data
	c.isSet = true
}

func (c *CacheRequest) Next() {
	c.index++
	s := int8(len(c.cacheChains))
	for ; c.index < s; c.index++ {
		fmt.Println("index~~~~~~", c.index)
		c.cacheChains[c.index](c)
	}
}

func (c *CacheRequest) Abort() {
	c.index = abortIndex
}
