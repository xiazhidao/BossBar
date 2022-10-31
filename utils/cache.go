package utils

import (
	"BossBar/cache"
	_ "BossBar/cache/redis"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/astaxie/beego"
)

var cc cache.Cache

func InitCache() {
	prefix := beego.AppConfig.String("cache::redis_prefix")
	host := beego.AppConfig.String("cache::redis_host")
	password := beego.AppConfig.String("cache::redis_password")
	dbNum := beego.AppConfig.String("cache::redis_dbnum")
	maxIdle := beego.AppConfig.String("cache::redis_maxidle")
	var err error
	defer func() {
		if r := recover(); r != nil {
			cc = nil
		}
	}()
	cc, err = cache.NewCache("redis", `{"key":"`+prefix+`","conn":"`+host+`","password":"`+password+`","dbNum":"`+dbNum+`","maxIdle":"`+maxIdle+`"}`)
	if err != nil {
		log.Errorf("Connect to the redis host %s failed, err:%s", host, err.Error())
	}
}

// SetCache
func SetCache(key string, value interface{}, timeout int) error {
	data, err := Encode(value)
	if err != nil {
		return err
	}
	if cc == nil {
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			log.Errorf("[cache] SetCache panic, err:%v", r)
			cc = nil
		}
	}()
	if timeout > 0 {
		timeouts := time.Duration(timeout) * time.Second
		err = cc.Put(key, data, timeouts)
	} else {
		err = cc.Set(key, data, 0, 0, false, false)
	}

	if err != nil {
		log.Errorf("SetCache failed, key:%s, err:%s", key, err.Error())
		return err
	} else {
		return nil
	}
}

func GetCache(key string, to interface{}) error {
	if cc == nil {
		return errors.New("cc is nil")
	}

	data := cc.Get(key)
	if data == nil {
		return errors.New(fmt.Sprintf("GetCache not found, key:%s", key))
	}

	err := Decode(data.([]byte), to)
	if err != nil {
		log.Errorf("GetCache failed, key:%s, err:%s", key, err.Error())
	}

	return err
}

// ExpireCache
func ExpireCache(key string, time int64) bool {
	if cc == nil {
		return false
	}
	result, _ := cc.Expire(key, time)
	return result
}

// ExpireCache
func ExistsCache(key string) (bool, error) {
	if cc == nil {
		return false, errors.New("cc is nil")
	}

	result, err := cc.Exists(key)
	if err != nil {
		log.Errorf("ExistsCache failed, key:%s, err:%s", key, err.Error())
	}
	return result, err
}

// DelCache
func DelCache(key string) error {
	if cc == nil {
		return errors.New("cc is nil")
	}
	err := cc.Delete(key)
	if err != nil {
		return errors.New("Cache删除失败")
	} else {
		return nil
	}
}

func GetPureCache(key string) (result string, err error) {
	if cc == nil {
		return "", errors.New("cc is nil")
	}
	if data := cc.Get(key); data == nil {
		return "", errors.New("Cache不存在")
	} else {
		result = string(data.([]byte))
	}
	return
}

// SetNxExCache
func SetNxExCache(key string, value interface{}, timeout int) error {
	if cc == nil {
		return errors.New("cc is nil")
	}
	return cc.Set(key, value, timeout, 0, false, true)
}

// IncrByCache
func IncrByCache(key string, increment int) (result int, err error) {
	if cc == nil {
		return 0, errors.New("cc is nil")
	}
	result, err = cc.IncrBy(key, increment)
	if err != nil {
		log.Errorf("IncrByCache failed, key:%s, err:%s", key, err.Error())
	}
	return
}

// DecrByCache
func DecrByCache(key string, decrement int) (result int, err error) {
	if cc == nil {
		return 0, errors.New("cc is nil")
	}
	result, err = cc.DecrBy(key, decrement)
	if err != nil {
		log.Errorf("DecrByCache failed, key:%s, err:%s", key, err.Error())
	}
	return
}

// SAddCache
func SAddCache(key string, members ...interface{}) (result int, err error) {
	if cc == nil {
		return 0, errors.New("cc is nil")
	}
	result, err = cc.SAdd(key, members...)
	if err != nil {
		log.Errorf("SAddCache failed, key:%s, err:%s", key, err.Error())
	}
	return
}

// SAddCache
func SPopCache(key string) (result string, err error) {
	if cc == nil {
		return "", errors.New("cc is nil")
	}
	result, err = cc.SPop(key)
	if err != nil {
		log.Errorf("SPopCache failed, key:%s, err:%s", key, err.Error())
	}
	return
}

// SIsMemberCache
func SIsMemberCache(key, member string) (result bool, err error) {
	if cc == nil {
		return false, errors.New("cc is nil")
	}
	result, err = cc.SIsMember(key, member)
	if err != nil {
		log.Errorf("SIsMemberCache failed, key:%s, err:%s", key, err.Error())
	}
	return
}

// SDiffCache
func SDiffCache(keys ...interface{}) (result []string, err error) {
	if cc == nil {
		return nil, errors.New("cc is nil")
	}
	result, err = cc.SDiff(keys...)
	return
}

// SMembersCache
func SMembersCache(key string) (result []string, err error) {
	if cc == nil {
		return nil, errors.New("cc is nil")
	}
	result, err = cc.SMembers(key)
	return
}

// SMoveCache
func SMoveCache(source, destination, member string) (result bool, err error) {
	if cc == nil {
		return false, errors.New("cc is nil")
	}
	result, err = cc.SMove(source, destination, member)
	return
}

// SRemCache
func SRemCache(key string, members ...interface{}) (result int, err error) {
	if cc == nil {
		return 0, errors.New("cc is nil")
	}
	result, err = cc.SRem(key, members...)
	return
}

// SUnionCache
func SUnionCache(keys ...interface{}) (result []string, err error) {
	if cc == nil {
		return nil, errors.New("cc is nil")
	}
	result, err = cc.SUnion(keys...)
	return
}

// ZAddCache
func ZAddCache(key string, pairs map[string]float64) error {
	if cc == nil {
		return errors.New("cc is nil")
	}
	err := cc.ZAdd(key, pairs)
	if err != nil {
		log.Errorf("ZAddCache err, key:%s, err:%s", key, err.Error())
	}
	return err
}

// ZScoreCache
func ZScoreCache(key, member string) (result string, err error) {
	if cc == nil {
		return "", errors.New("cc is nil")
	}
	result, err = cc.ZScore(key, member)
	return
}

// ZRangeCache
func ZRangeCache(key string, start, stop int, withscores bool) (result []string, err error) {
	if cc == nil {
		return nil, errors.New("cc is nil")
	}
	result, err = cc.ZRange(key, start, stop, withscores)
	return
}

// ZRangeByScoreCache
func ZRangeByScoreCache(key string, min, max int64, withscores bool) (result []string, err error) {
	if cc == nil {
		return nil, errors.New("cc is nil")
	}
	result, err = cc.ZRangeByScore(key, min, max, withscores)
	return
}

// ZRevRangeCache
func ZRevRangeCache(key string, start, stop int, withscores bool) (result []string, err error) {
	if cc == nil {
		return nil, errors.New("cc is nil")
	}
	result, err = cc.ZRevRange(key, start, stop, withscores)
	return
}

// ZRemCache
func ZRemCache(key string, values ...interface{}) (result int, err error) {
	if cc == nil {
		return 0, errors.New("cc is nil")
	}
	result, err = cc.ZRem(key, values...)
	return
}

// IncrByCache
func ZIncrByCache(key, member string, increment int64) (result int64, err error) {
	if cc == nil {
		return 0, errors.New("cc is nil")
	}
	result, err = cc.ZIncrby(key, member, increment)
	if err != nil {
		log.Errorf("ZIncrByCache failed, key:%s, err:%s", key, err.Error())
	}
	return
}

// ZRemRangeByRankCache
func ZRemRangeByRankCache(key string, start, stop int) (result int, err error) {
	if cc == nil {
		return 0, errors.New("cc is nil")
	}
	result, err = cc.ZRemRangeByRank(key, start, stop)
	return
}

// ZRemRangeByRankCache
func ZRemRangeByScoreCache(key string, min, max int64) (result int, err error) {
	if cc == nil {
		return 0, errors.New("cc is nil")
	}
	result, err = cc.ZRemRangeByScore(key, min, max)
	return
}

func HSetCache(key, field, value string) (result bool, err error) {
	if cc == nil {
		return false, errors.New("cc is null")
	}
	if result, err = cc.HSet(key, field, value); err != nil {
		log.Errorf("HSetCache failed, key:%s, field:%s, err:%s", key, field, err.Error())
	}
	return
}

func HMSetCache(key string, params ...string) (result string, err error) {
	if cc == nil {
		return "", errors.New("cc is null")
	}
	if result, err = cc.HMSet(key, params...); err != nil {
		log.Errorf("HMSetCache failed, key:%s, params:%s, err:%s", key, params, err.Error())
	}
	return
}

func HInCrBy(key, field string, val int64) (result int64, err error) {
	if cc == nil {
		return 0, errors.New("cc is null")
	}
	if result, err = cc.HInCrBy(key, field, val); err != nil {
		log.Errorf("HInCrBy failed, key:%s, field:%s, err:%s", key, field, err.Error())
	}
	return
}

func HExists(key, field string) (result bool, err error) {
	if cc == nil {
		return false, errors.New("cc is null")
	}
	if result, err = cc.HExists(key, field); err != nil {
		log.Errorf("HExists failed, key:%s, field:%s, err:%s", key, field, err.Error())
	}
	return
}

func HGetCache(key, field string) (result string, err error) {
	if cc == nil {
		return "", errors.New("cc is nil")
	}
	result, err = cc.HGet(key, field)
	if err != nil {
		log.Errorf("HGetCache failed, key:%s, field:%s, err:%s", key, field, err.Error())
	}
	return
}

func HMGetCache(key string, fields ...string) (result []string, err error) {
	if cc == nil {
		return nil, errors.New("cc is nil")
	}
	result, err = cc.HMGet(key, fields...)
	if err != nil {
		log.Errorf("HMGetCache failed, key:%s, fields:%v, err:%s", key, fields, err.Error())
	}
	return
}

func HValsCache(key string) (result []string, err error) {
	if cc == nil {
		return nil, errors.New("cc is nil")
	}
	result, err = cc.HVals(key)
	return
}

func HGetAllCache(key string) (result []string, err error) {
	if cc == nil {
		return nil, errors.New("cc is nil")
	}
	result, err = cc.HGetAll(key)
	return
}

// HDelCache
func HDelCache(key string, fields ...string) (result int, err error) {
	if cc == nil {
		return 0, errors.New("cc is nil")
	}
	if result, err = cc.HDel(key, fields...); err != nil {
		log.Errorf("HDelCache failed, key:%s, fields:%v, err:%s", key, fields, err.Error())
	}
	return
}

func RPushCache(key string, values ...string) error {
	if cc == nil {
		return errors.New("cc is nil")
	}
	return cc.RPush(key, values...)
}

func LPushCache(key string, values ...string) error {
	if cc == nil {
		return errors.New("cc is nil")
	}
	return cc.RPush(key, values...)
}

func LPopCache(key string) (result string, err error) {
	if cc == nil {
		return "", errors.New("cc is nil")
	}
	return cc.LPop(key)
}

func BLPopCache(key string, timeout int) (result string, err error) {
	if cc == nil {
		return "", errors.New("cc is nil")
	}

	return cc.BLPop(key, timeout)
}

func LRemCache(key string, count int, value string) error {
	if cc == nil {
		return errors.New("cc is nil")
	}
	return cc.LRem(key, count, value)
}

// Encode
// 用gob进行数据编码
func Encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode
// 用gob进行数据解码
func Decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}
