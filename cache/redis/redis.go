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

// Package redis for cache provider
//
// depend on github.com/gomodule/redigo/redis
//
// go install github.com/gomodule/redigo/redis
//
// Usage:
// import(
//
//	_ "github.com/astaxie/beego/cache/redis"
//	"github.com/astaxie/beego/cache"
//
// )
//
//	bm, err := cache.NewCache("redis", `{"conn":"127.0.0.1:11211"}`)
//
//	more docs http://beego.me/docs/module/cache.md
package redis

import (
	"BossBar/cache"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	// DefaultKey the collection name of redis for cache adapter.
	DefaultKey = "beecacheRedis"
)

// Cache is Redis cache adapter.
type Cache struct {
	p        *redis.Pool // redis connection pool
	conninfo string
	dbNum    int
	key      string
	password string
	maxIdle  int
}

// NewRedisCache create new redis cache with default collection name.
func NewRedisCache() cache.Cache {
	return &Cache{key: DefaultKey}
}

// actually do the redis cmds, args[0] must be the key name.
func (rc *Cache) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	args[0] = rc.associate(args[0])
	c := rc.p.Get()
	defer c.Close()
	return c.Do(commandName, args...)
}

// associate with config key.
func (rc *Cache) associate(originKey interface{}) string {
	return fmt.Sprintf("%s:%s", rc.key, originKey)
}

// Get cache from redis.
func (rc *Cache) Get(key string) interface{} {
	if v, err := rc.do("GET", key); err == nil {
		return v
	}
	return nil
}

// Get cache from redis.
func (rc *Cache) Expire(key string, time int64) (bool, error) {
	return redis.Bool(rc.do("EXPIRE", key, time))
}

func (rc *Cache) Exists(key string) (bool, error) {
	return redis.Bool(rc.do("EXISTS", key))
}

// GetMulti get cache from redis.
func (rc *Cache) GetMulti(keys []string) []interface{} {
	c := rc.p.Get()
	defer c.Close()
	var args []interface{}
	for _, key := range keys {
		args = append(args, rc.associate(key))
	}
	values, err := redis.Values(c.Do("MGET", args...))
	if err != nil {
		return nil
	}
	return values
}

// Set set cache to redis.
func (rc *Cache) Set(key string, value interface{}, seconds, milliseconds int, mustExists, mustNotExists bool) error {
	var args []interface{}
	args = append(args, key, value)
	if seconds > 0 {
		args = append(args, "EX", seconds)
	}
	if milliseconds > 0 {
		args = append(args, "PX", milliseconds)
	}
	if mustExists {
		args = append(args, "XX")
	} else if mustNotExists {
		args = append(args, "NX")
	}
	_, err := redis.String(rc.do("SET", args...))
	return err
}

// HSet
func (rc *Cache) HSet(key, field, value string) (bool, error) {
	return redis.Bool(rc.do("HSET", key, field, value))
}

// HGet
func (rc *Cache) HGet(key, field string) (string, error) {
	return redis.String(rc.do("HGET", key, field))
}

func (rc *Cache) HInCrBy(key, field string, val int64) (int64, error) {
	return redis.Int64(rc.do("HINCRBY", key, field, val))
}

func (rc *Cache) HExists(key, field string) (bool, error) {
	return redis.Bool(rc.do("HEXISTS", key, field))
}

func (rc *Cache) HMGet(key string, fields ...string) ([]string, error) {
	var args []interface{}
	args = append(args, key)
	for _, v := range fields {
		args = append(args, v)
	}

	return redis.Strings(rc.do("HMGET", args...))
}

func (rc *Cache) HMSet(key string, params ...string) (string, error) {
	var args []interface{}
	args = append(args, key)
	for _, v := range params {
		args = append(args, v)
	}

	return redis.String(rc.do("HMSET", args...))
}

// HGetVals
func (rc *Cache) HVals(key string) ([]string, error) {
	return redis.Strings(rc.do("HVALS", key))
}

// HGetVals
func (rc *Cache) HGetAll(key string) ([]string, error) {
	return redis.Strings(rc.do("HGETALL", key))
}

// HDel
func (rc *Cache) HDel(key string, fields ...string) (int, error) {
	var args []interface{}
	args = append(args, key)
	for _, v := range fields {
		args = append(args, v)
	}
	return redis.Int(rc.do("HDEL", args...))
}

// Put put cache to redis.
func (rc *Cache) Put(key string, val interface{}, timeout time.Duration) error {
	_, err := rc.do("SETEX", key, int64(timeout/time.Second), val)
	return err
}

// Delete delete cache in redis.
func (rc *Cache) Delete(key string) error {
	_, err := rc.do("DEL", key)
	return err
}

// IsExist check cache's existence in redis.
func (rc *Cache) IsExist(key string) bool {
	v, err := redis.Bool(rc.do("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

// Incr increase counter in redis.
func (rc *Cache) Incr(key string) (int, error) {
	return redis.Int(rc.do("INCR", key))
}

// IncrBy increase counter in redis.
func (rc *Cache) IncrBy(key string, increment int) (int, error) {
	return redis.Int(rc.do("INCRBY", key, increment))
}

// Decr decrease counter in redis.
func (rc *Cache) Decr(key string) (int, error) {
	return redis.Int(rc.do("DECR", key))
}

// DecrBy increase counter in redis.
func (rc *Cache) DecrBy(key string, decrement int) (int, error) {
	return redis.Int(rc.do("DECRBY", key, decrement))
}

// Setnx.
func (rc *Cache) Setnx(key string, value interface{}) (result bool, err error) {
	result, err = redis.Bool(rc.do("SETNX", key, value))
	return
}

// SAdd.
func (rc *Cache) SAdd(key string, members ...interface{}) (result int, err error) {
	var args []interface{}
	args = append(args, key)
	for _, m := range members {
		args = append(args, m)
	}
	result, err = redis.Int(rc.do("SADD", args...))
	return
}

// SPop.
func (rc *Cache) SPop(key string) (result string, err error) {
	var args []interface{}
	args = append(args, key)
	result, err = redis.String(rc.do("SPop", args...))
	return
}

// SIsMember.
func (rc *Cache) SIsMember(key, member string) (result bool, err error) {
	result, err = redis.Bool(rc.do("SISMEMBER", key, member))
	return
}

// SMembers.
func (rc *Cache) SMembers(key string) (result []string, err error) {
	result, err = redis.Strings(rc.do("SMEMBERS", key))
	return
}

// SDiff.
func (rc *Cache) SDiff(keys ...interface{}) (result []string, err error) {
	var args []interface{}
	for index, key := range keys {
		if index > 0 {
			args = append(args, rc.associate(key))
		} else {
			args = append(args, key)
		}
	}
	result, err = redis.Strings(rc.do("SDIFF", args...))
	return
}

// SMove.
func (rc *Cache) SMove(source, destination, member string) (result bool, err error) {
	destination = rc.associate(destination)
	result, err = redis.Bool(rc.do("SMOVE", source, destination, member))
	return
}

// SRem.
func (rc *Cache) SRem(key string, members ...interface{}) (result int, err error) {
	var args []interface{}
	args = append(args, key)
	for _, v := range members {
		args = append(args, v)
	}
	result, err = redis.Int(rc.do("SREM", args...))
	return
}

// SUnion.
func (rc *Cache) SUnion(keys ...interface{}) (result []string, err error) {
	var args []interface{}
	for index, key := range keys {
		if index > 0 {
			args = append(args, rc.associate(key))
		} else {
			args = append(args, key)
		}
	}
	result, err = redis.Strings(rc.do("SUNION", args...))
	return
}

// ZAdd.
func (rc *Cache) ZAdd(key string, pairs map[string]float64) error {
	var args []interface{}
	args = append(args, key)
	for k, v := range pairs {
		args = append(args, v)
		args = append(args, k)
	}
	_, err := redis.Bool(rc.do("ZADD", args...))
	return err
}

// ZScore.
func (rc *Cache) ZScore(key, member string) (result string, err error) {
	result, err = redis.String(rc.do("ZSCORE", key, member))
	return
}

// ZRange.
func (rc *Cache) ZRange(key string, start, stop int, withscores bool) (result []string, err error) {
	if withscores {
		result, err = redis.Strings(rc.do("ZRANGE", key, start, stop, "withscores"))
	} else {
		result, err = redis.Strings(rc.do("ZRANGE", key, start, stop))
	}
	return
}

// ZRangeByScore.
func (rc *Cache) ZRangeByScore(key string, min, max int64, withscores bool) (result []string, err error) {
	if withscores {
		result, err = redis.Strings(rc.do("ZRANGEBYSCORE", key, min, max, "withscores"))
	} else {
		result, err = redis.Strings(rc.do("ZRANGEBYSCORE", key, min, max))
	}
	return
}

// ZRevRange.
func (rc *Cache) ZRevRange(key string, start, stop int, withscores bool) (result []string, err error) {
	if withscores {
		result, err = redis.Strings(rc.do("ZREVRANGE", key, start, stop, "withscores"))
	} else {
		result, err = redis.Strings(rc.do("ZREVRANGE", key, start, stop))
	}
	return
}

// ZREM.
func (rc *Cache) ZRem(key string, values ...interface{}) (result int, err error) {
	var args []interface{}
	args = append(args, key)
	for _, v := range values {
		args = append(args, v)
	}
	result, err = redis.Int(rc.do("ZREM", args...))
	return
}

func (rc *Cache) ZIncrby(key, member string, increment int64) (result int64, err error) {
	result, err = redis.Int64(rc.do("ZINCRBY", key, increment, member))
	return
}

// ZREMRANGEBYRANK.
func (rc *Cache) ZRemRangeByRank(key string, start, stop int) (result int, err error) {
	result, err = redis.Int(rc.do("ZREMRANGEBYRANK", key, start, stop))
	return
}

// ZRemRangeByScore.
func (rc *Cache) ZRemRangeByScore(key string, min, max int64) (result int, err error) {
	result, err = redis.Int(rc.do("ZREMRANGEBYSCORE", key, min, max))
	return
}

// ClearAll clean all cache in redis. delete this redis collection.
func (rc *Cache) ClearAll() error {
	c := rc.p.Get()
	defer c.Close()
	cachedKeys, err := redis.Strings(c.Do("KEYS", rc.key+":*"))
	if err != nil {
		return err
	}
	for _, str := range cachedKeys {
		if _, err = c.Do("DEL", str); err != nil {
			return err
		}
	}
	return err
}

// StartAndGC start redis cache adapter.
// config is like {"key":"collection key","conn":"connection info","dbNum":"0"}
// the cache item in redis are stored forever,
// so no gc operation.
func (rc *Cache) StartAndGC(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)

	if _, ok := cf["key"]; !ok {
		cf["key"] = DefaultKey
	}
	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}

	// Format redis://<password>@<host>:<port>
	cf["conn"] = strings.Replace(cf["conn"], "redis://", "", 1)
	if i := strings.Index(cf["conn"], "@"); i > -1 {
		cf["password"] = cf["conn"][0:i]
		cf["conn"] = cf["conn"][i+1:]
	}

	if _, ok := cf["dbNum"]; !ok {
		cf["dbNum"] = "0"
	}
	if _, ok := cf["password"]; !ok {
		cf["password"] = ""
	}
	if _, ok := cf["maxIdle"]; !ok {
		cf["maxIdle"] = "3"
	}
	rc.key = cf["key"]
	rc.conninfo = cf["conn"]
	rc.dbNum, _ = strconv.Atoi(cf["dbNum"])
	rc.password = cf["password"]
	rc.maxIdle, _ = strconv.Atoi(cf["maxIdle"])

	rc.connectInit()

	c := rc.p.Get()
	defer c.Close()

	return c.Err()
}

// push to list
func (rc *Cache) RPush(key string, values ...string) error {
	var args []interface{}
	args = append(args, key)
	for _, v := range values {
		args = append(args, v)
	}
	_, err := redis.Int(rc.do("RPUSH", args...))
	return err
}

// push to list
func (rc *Cache) LPush(key string, values ...string) error {
	var args []interface{}
	args = append(args, key)
	for _, v := range values {
		args = append(args, v)
	}
	_, err := redis.Int(rc.do("LPUSH", args...))
	return err
}

// pop from list
func (rc *Cache) LPop(key string) (string, error) {
	return redis.String(rc.do("LPOP", key))
}

// pop from list with block
func (rc *Cache) BLPop(key string, timeout int) (string, error) {
	reply, err := redis.Strings(rc.do("BLPOP", key, timeout))
	if err != nil {
		return "", err
	}
	return reply[1], nil
}

// remove from list
func (rc *Cache) LRem(key string, count int, value string) error {
	_, err := redis.Int(rc.do("LREM", key, count, value))
	return err
}

// range list
func (rc *Cache) Lrange(key string, start, stop int) ([]string, error) {
	c := rc.p.Get()
	defer c.Close()
	return redis.Strings(c.Do("LRANGE", key, start, stop))
}

// connect to redis.
func (rc *Cache) connectInit() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", rc.conninfo)
		if err != nil {
			return nil, err
		}

		if rc.password != "" {
			if _, err := c.Do("AUTH", rc.password); err != nil {
				c.Close()
				return nil, err
			}
		}

		_, selecterr := c.Do("SELECT", rc.dbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	// initialize a new pool
	rc.p = &redis.Pool{
		MaxIdle:     rc.maxIdle,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
}

func init() {
	cache.Register("redis", NewRedisCache)
}
