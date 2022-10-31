// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package cache provide a Cache interface and some implement engine
// Usage:
//
// import(
//   "github.com/astaxie/beego/cache"
// )
//
// bm, err := cache.NewCache("memory", `{"interval":60}`)
//
// Use it like this:
//
//	bm.Put("astaxie", 1, 10 * time.Second)
//	bm.Get("astaxie")
//	bm.IsExist("astaxie")
//	bm.Delete("astaxie")
//
//  more docs http://beego.me/docs/module/cache.md
package cache

import (
	"fmt"
	"time"
)

// Cache interface contains all behaviors for cache adapter.
// usage:
//	cache.Register("file",cache.NewFileCache) // this operation is run in init method of file.go.
//	c,err := cache.NewCache("file","{....}")
//	c.Put("key",value, 3600 * time.Second)
//	v := c.Get("key")
//
//	c.Incr("counter")  // now is 1
//	c.Incr("counter")  // now is 2
//	count := c.Get("counter").(int)
type Cache interface {
	// get cached value by key.
	Get(key string) interface{}
	// GetMulti is a batch version of Get.
	GetMulti(keys []string) []interface{}
	// set cached value with key.
	Set(key string, value interface{}, seconds, milliseconds int, mustExists, mustNotExists bool) error
	// set cached value with key and expire time.
	Put(key string, val interface{}, timeout time.Duration) error

	Exists(key string) (bool, error)

	Expire(key string, time int64) (bool, error)
	// delete cached value by key.
	Delete(key string) error
	// increase cached int value by key, as a counter.
	Incr(key string) (int, error)
	// increase cached int value by key, as a counter.
	IncrBy(key string, increment int) (int, error)
	// decrease cached int value by key, as a counter.
	Decr(key string) (int, error)
	// decrease cached int value by key, as a counter.
	DecrBy(key string, increment int) (int, error)
	// check if cached value exists or not.
	IsExist(key string) bool
	// SETNX
	Setnx(key string, value interface{}) (bool, error)
	// SADD
	SAdd(key string, members ...interface{}) (int, error)
	//SPop
	SPop(key string) (string, error)
	// SIsMember
	SIsMember(key, member string) (bool, error)
	//SMembers
	SMembers(key string) ([]string, error)
	// SDIFF
	SDiff(keys ...interface{}) ([]string, error)
	// SMOVE
	SMove(source, destination, member string) (bool, error)
	// SREM
	SRem(key string, members ...interface{}) (int, error)
	// SUNION
	SUnion(keys ...interface{}) ([]string, error)
	// ZADD
	ZAdd(key string, pairs map[string]float64) error
	// ZSCORE
	ZScore(key, member string) (string, error)
	// ZRANGE
	ZRange(key string, start, stop int, withscores bool) ([]string, error)
	// ZRANGEBYSCORE
	ZRangeByScore(key string, min, max int64, withscores bool) ([]string, error)
	// ZREVRANGE
	ZRevRange(key string, start, stop int, withscores bool) ([]string, error)
	// ZREM
	ZRem(key string, values ...interface{}) (int, error)

	ZIncrby(key, member string, increment int64) (result int64, err error)
	// ZREMRANGEBYRANK
	ZRemRangeByRank(key string, start, stop int) (int, error)

	ZRemRangeByScore(key string, min, max int64) (result int, err error)
	// HGet
	HGet(key, field string) (string, error)
	// HSet
	HSet(key, field, value string) (bool, error)

	HInCrBy(key, field string, val int64) (int64, error)

	HExists(key, field string) (bool, error)

	HMGet(key string, fields ...string) ([]string, error)

	HMSet(key string, params ...string) (string, error)

	HVals(key string) ([]string, error)

	HGetAll(key string) ([]string, error)
	// HDel
	HDel(key string, fields ...string) (int, error)
	// clear all cache.
	ClearAll() error
	// start gc routine based on config string settings.
	StartAndGC(config string) error
	// push to list
	RPush(key string, values ...string) error

	LPush(key string, values ...string) error
	// pop from list
	LPop(key string) (string, error)
	// pop from list with block
	BLPop(key string, timeout int) (string, error)
	// remove from list
	LRem(key string, count int, value string) error
	// range list
	Lrange(key string, start, stop int) ([]string, error)
}

// Instance is a function create a new Cache Instance
type Instance func() Cache

var adapters = make(map[string]Instance)

// Register makes a cache adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("cache: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("cache: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// NewCache Create a new cache driver by adapter name and config string.
// config need to be correct JSON as string: {"interval":360}.
// it will start gc automatically.
func NewCache(adapterName, config string) (adapter Cache, err error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("cache: unknown adapter name %q (forgot to import?)", adapterName)
		return
	}
	adapter = instanceFunc()
	err = adapter.StartAndGC(config)
	if err != nil {
		adapter = nil
	}
	return
}
