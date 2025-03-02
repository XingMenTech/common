package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

var (
	client *redis.Client
	ctx    context.Context
	Key    string
)

type Config struct {
	Prefix   string `yaml:"redis_prefix" json:"prefix"`
	Host     string `yaml:"redis_host" json:"host"`
	Password string `yaml:"redis_password" json:"password"`
	DbNum    int    `yaml:"redis_dbnum" json:"dbNum"`
}

func InitRedisCache(config *Config) error {
	ctx = context.Background()
	Key = config.Prefix

	cli, err := startAndGC(config.Host, config.Password, config.DbNum)
	if err != nil {
		return errors.New(fmt.Sprintf("can't connect redis service %v", err))
	}

	client = cli
	return nil
}

func associate(originKey interface{}) string {
	return fmt.Sprintf("%s:%s", Key, originKey)
}

// start gc routine based on config string settings.
func startAndGC(host, passWord string, dbNum int) (*redis.Client, error) {

	cli := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: passWord,
		DB:       dbNum,
	})
	cmd := cli.Ping(ctx)
	if cmd.Err() != nil {
		return nil, errors.New(fmt.Sprintf("redis connect errors: %v \n", cmd.Err()))
	}

	return cli, nil
}

func SetCache(key string, val interface{}, timeout int64) error {
	bytes, err := Encode(val)
	if err != nil {
		return err
	}
	cmd := client.Set(ctx, associate(key), string(bytes), time.Duration(timeout)*time.Second)
	return cmd.Err()
}

func GetCache(key string, to interface{}) error {
	cmd := client.Get(ctx, associate(key))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	if err := Decode([]byte(cmd.Val()), to); err != nil {
		return err
	}
	return nil
}

func Set(key string, val interface{}, timeout int64) error {
	cmd := client.Set(ctx, associate(key), val, time.Duration(timeout)*time.Second)
	return cmd.Err()
}

func Get(key string, to interface{}) error {
	cmd := client.Get(ctx, associate(key))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	if err := cmd.Scan(to); err != nil {
		if err != nil {
			err := Decode([]byte(cmd.Val()), to)
			return err
		}
	}
	return nil
}

// check if cached value exists or not.
func IsExist(key string) bool {
	val := client.Exists(ctx, associate(key)).Val()
	return val != 0
}

// delete cached value by key.
func Delete(key string) error {
	return client.Del(ctx, key).Err()
}

// HKEYS
func HKeys(key string) ([]string, error) {
	cmd := client.HKeys(ctx, associate(key))
	return cmd.Val(), cmd.Err()
}
func HExists(key, field string) (bool, error) {
	exists := client.HExists(ctx, associate(key), field)
	return exists.Val(), exists.Err()
}

// HSetAll
func HSetAll(key string, val interface{}) error {
	return client.HSet(ctx, associate(key), val).Err()
}

// HGetAll
func HGetAll(key string, to interface{}) error {
	cmd := client.HGetAll(ctx, associate(key))
	if err := cmd.Scan(to); err != nil {
		return err
	}
	return nil
}
func HVals(key string) ([]string, error) {
	cmd := client.HVals(ctx, associate(key))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return cmd.Val(), nil
}
func HLen(key string) int64 {
	hLen := client.HLen(ctx, associate(key))
	if hLen.Err() != nil {
		return 0
	}
	return hLen.Val()
}

// HSet
func HSet(key string, field string, val interface{}) error {
	valByte, err := Encode(val)
	if err != nil {
		return err
	}
	return client.HSet(ctx, associate(key), field, string(valByte)).Err()
}

// HGet
func HGet(key string, field string, to interface{}) error {
	cmd := client.HGet(ctx, associate(key), field)
	if err := cmd.Scan(to); err != nil {
		return err
	}
	return nil
}

// HIncrby
func HIncrby(key string, field string, incr int64) (int64, error) {
	cmd := client.HIncrBy(ctx, associate(key), field, incr)
	return cmd.Val(), cmd.Err()
}

// HDEL
func HDel(key string, fields ...string) error {
	return client.HDel(ctx, associate(key), fields...).Err()
}

// 订阅主题
func Subscribe(channel ...string) *redis.PubSub {
	return client.Subscribe(ctx, channel...)
}

// 订阅主题
func PSubscribe(channel ...string) *redis.PubSub {
	return client.PSubscribe(ctx, channel...)
}

// 发布主题消息
func Publish(channel string, msg interface{}) error {
	msgByte, err := Encode(msg)
	if err != nil {
		return err
	}
	return client.Publish(ctx, channel, string(msgByte)).Err()
}

func ReceiveMessage(pubSub *redis.PubSub) (*redis.Message, error) {
	return pubSub.ReceiveMessage(ctx)
}
func GetForNoPrefix(key string) interface{} {
	cmd := client.Get(ctx, key)
	return cmd.Val()
}

func SetForNoPrefix(key string, val interface{}, dur time.Duration) error {
	return client.Set(ctx, key, val, dur).Err()
}

// GetMulti is a batch version of Get.
func GetMulti(keys []string) []interface{} {
	var args []string
	for _, key := range keys {
		args = append(args, associate(key))
	}
	return client.MGet(ctx, args...).Val()
}

// set cached value with key.
func Put(key string, val interface{}, timeout time.Duration) error {

	return client.SetEX(ctx, associate(key), val, timeout).Err()
}

// increase cached int value by key, as a counter.
func Incr(key string) (int64, error) {
	cmd := client.Incr(ctx, associate(key))
	return cmd.Val(), cmd.Err()
}

// increase cached int value by key, as a counter.
func IncrBy(key string, val int64) (int64, error) {
	cmd := client.IncrBy(ctx, associate(key), val)
	return cmd.Val(), cmd.Err()
}

// decrease cached int value by key, as a counter.
func Decr(key string) error {
	return client.Decr(ctx, associate(key)).Err()
}

// SETNX
func Setnx(key string, value interface{}) (bool, error) {
	cmd := client.SetNX(ctx, associate(key), value, 0)
	return cmd.Val(), cmd.Err()
}

// SETNX WITH EXPIRE
func SetnxExpire(key string, value interface{}, expire int64) (bool, error) {
	cmd := client.SetNX(ctx, associate(key), value, time.Duration(expire))
	return cmd.Val(), cmd.Err()
}

// SADD
func SAdd(key string, members ...interface{}) (int, error) {
	cmd := client.SAdd(ctx, associate(key), members...)
	return int(cmd.Val()), cmd.Err()
}

// SIsMember
func SIsMember(key, member string) (bool, error) {
	cmd := client.SIsMember(ctx, associate(key), member)
	return cmd.Val(), cmd.Err()
}

// SDIFF
func SDiff(keys ...string) ([]string, error) {
	cmd := client.SDiff(ctx, keys...)
	return cmd.Val(), cmd.Err()
}

// SMOVE
func SMove(source, destination, member string) (bool, error) {
	cmd := client.SMove(ctx, associate(source), associate(destination), member)
	return cmd.Val(), cmd.Err()
}

// SREM
func SRem(key string, members ...interface{}) (int, error) {
	cmd := client.SRem(ctx, associate(key), members...)
	return int(cmd.Val()), cmd.Err()
}

// SUNION
func SUnion(keys ...string) ([]string, error) {
	cmd := client.SUnion(ctx, keys...)
	return cmd.Val(), cmd.Err()
}

// ZADD
func ZAdd(key string, pairs map[string]float64) error {
	var args []*redis.Z
	for k, v := range pairs {
		args = append(args, &redis.Z{
			Score:  v,
			Member: k,
		})
	}

	return client.ZAdd(ctx, associate(key), args...).Err()
}

// ZSCORE
func ZScore(key, member string) (string, error) {

	cmd := client.ZScore(ctx, associate(key), member)
	return fmt.Sprintf("%.2f", cmd.Val()), cmd.Err()
}

// ZRANGE
func ZRange(key string, start, stop int, withscores bool) ([]string, error) {

	if withscores {
		cmd := client.ZRangeWithScores(ctx, associate(key), int64(start), int64(stop))
		if cmd.Err() != nil {
			return []string{}, cmd.Err()
		}
		var res = make([]string, len(cmd.Val()))
		for i, z := range cmd.Val() {
			res[i] = z.Member.(string)
		}
		return res, nil
	} else {
		cmd := client.ZRange(ctx, associate(key), int64(start), int64(stop))
		return cmd.Val(), cmd.Err()
	}
}

// ZRANGEBYSCORE
func ZRangeByScore(key, min, max string, withscores bool) ([]string, error) {

	if withscores {
		cmd := client.ZRangeByScoreWithScores(ctx, associate(key), &redis.ZRangeBy{
			Min: min,
			Max: max,
		})
		if cmd.Err() != nil {
			return []string{}, cmd.Err()
		}
		var res = make([]string, len(cmd.Val()))
		for i, z := range cmd.Val() {
			res[i] = z.Member.(string)
		}
		return res, nil
	} else {
		cmd := client.ZRangeByScore(ctx, associate(key), &redis.ZRangeBy{
			Min: min,
			Max: max,
		})

		return cmd.Val(), cmd.Err()
	}

}

// ZREVRANGE
func ZRevRange(key string, start, stop int, withscores bool) ([]string, error) {
	if withscores {
		cmd := client.ZRevRangeByScoreWithScores(ctx, associate(key), &redis.ZRangeBy{
			Min:    strconv.Itoa(stop),
			Max:    strconv.Itoa(start),
			Offset: 0,
			Count:  0,
		})
		if cmd.Err() != nil {
			return []string{}, cmd.Err()
		}
		var res = make([]string, len(cmd.Val()))
		for i, z := range cmd.Val() {
			res[i] = z.Member.(string)
		}
		return res, nil
	} else {
		cmd := client.ZRevRange(ctx, associate(key), int64(start), int64(stop))
		return cmd.Val(), cmd.Err()
	}
}

// ZREM
func ZRem(key string, values ...interface{}) (int, error) {
	cmd := client.ZRem(ctx, associate(key), values...)
	return int(cmd.Val()), cmd.Err()
}

// ZREMRANGEBYRANK
func ZRemRangeByRank(key string, start, stop int) (int, error) {
	cmd := client.ZRemRangeByRank(ctx, associate(key), int64(start), int64(stop))
	return int(cmd.Val()), cmd.Err()
}

// clear all cache.
func ClearAll() error {
	return client.FlushAll(ctx).Err()
}

// push to list
func Rpush(key string, val interface{}) error {
	return client.RPush(ctx, associate(key), val).Err()
}

// pop from list
func Lpop(key string) (string, error) {
	cmd := client.LPop(ctx, associate(key))
	return cmd.Val(), cmd.Err()
}

// pop from list with block
func Blpop(key string, timeout int) (string, error) {
	cmd := client.BLPop(ctx, time.Duration(timeout)*time.Second, associate(key))
	return cmd.Val()[0], cmd.Err()
}

// range list
func Lrange(key string, start, stop int) ([]string, error) {
	cmd := client.LRange(ctx, associate(key), int64(start), int64(stop))
	return cmd.Val(), cmd.Err()
}

// delete all form list
func Ldel(key string, index int, val interface{}) (int, error) {
	cmd := client.LRem(ctx, associate(key), int64(index), val)
	return int(cmd.Val()), cmd.Err()
}
func LPush(key string, val interface{}) error {
	return client.LPush(ctx, associate(key), val).Err()
}
func RPop(key string) (string, error) {
	cmd := client.RPop(ctx, associate(key))
	return cmd.Val(), cmd.Err()
}
func ExpireAt(key string, t time.Time) error {
	return client.ExpireAt(ctx, associate(key), t).Err()
}
func ExpireIn(key string, d time.Duration) error {
	return client.Expire(ctx, associate(key), d).Err()
}

func Encode(data interface{}) ([]byte, error) {
	if data == nil {
		return nil, errors.New("data is nil")
	}

	switch data.(type) {
	case string:
		return []byte(data.(string)), nil
	case []byte:
		return data.([]byte), nil
	default:
		return json.Marshal(data)
	}
}

func Decode(data []byte, to interface{}) error {
	if data == nil {
		return errors.New("data is nil")
	}
	switch to.(type) {
	case string:
		to = string(data)
		return nil
	case []byte:
		to = data
		return nil
	default:
		return json.Unmarshal(data, to)
	}
}
