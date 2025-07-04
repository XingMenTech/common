package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	client *redis.Client
	ctx    context.Context
	Key    string
)

type Config struct {
	Prefix   string `yaml:"prefix" json:"prefix" comment:"KEY前缀"`
	Host     string `yaml:"host" json:"host" comment:"主机名"`
	Password string `yaml:"password" json:"password" comment:"密码"`
	DbNum    int    `yaml:"dbNum" json:"dbNum" comment:"数据库"`
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

// check if cached value exists or not.
func IsExist(key string) bool {
	val := client.Exists(ctx, associate(key)).Val()
	return val != 0
}

// delete cached value by key.
func Delete(key string) error {
	return client.Del(ctx, associate(key)).Err()
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

// clear all cache.
func ClearAll() error {
	return client.FlushAll(ctx).Err()
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
