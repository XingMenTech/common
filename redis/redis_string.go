package redis

import (
	"reflect"
	"time"
)

// 1	SET key value 设置指定 key 的值。
func Set(key string, val interface{}, timeout int64) error {
	cmd := client.Set(ctx, associate(key), val, time.Duration(timeout)*time.Second)
	return cmd.Err()
}
func SetCache(key string, val interface{}, timeout int64) error {
	bytes, err := encode(val)
	if err != nil {
		return err
	}
	cmd := client.Set(ctx, associate(key), string(bytes), time.Duration(timeout)*time.Second)
	return cmd.Err()
}
func SetForNoPrefix(key string, val interface{}, dur int64) error {
	bytes, err := encode(val)
	if err != nil {
		return err
	}
	cmd := client.Set(ctx, key, string(bytes), time.Duration(dur)*time.Second)
	return cmd.Err()
}

// 2	GET key 获取指定 key 的值。
func Get(key string, to interface{}) error {
	cmd := client.Get(ctx, associate(key))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	if err := cmd.Scan(to); err != nil {
		return decode([]byte(cmd.Val()), to)
	}
	return nil
}
func GetCache(key string, to interface{}) error {
	cmd := client.Get(ctx, associate(key))
	if cmd.Err() != nil {
		return cmd.Err()
	}

	if err := cmd.Scan(to); err != nil {
		return decode([]byte(cmd.Val()), to)
	}
	return nil
}
func GetForNoPrefix(key string, to interface{}) error {
	cmd := client.Get(ctx, key)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	if err := cmd.Scan(to); err != nil {
		return decode([]byte(cmd.Val()), to)
	}
	return nil
}

// GetRange 3	GETRANGE key start end 返回 key 中字符串值的子字符
func GetRange(key string, start, end int64) (string, error) {
	cmd := client.GetRange(ctx, associate(key), start, end)
	return cmd.Val(), cmd.Err()
}

// GetSet 4	GETSET key value 将给定 key 的值设为 value ，并返回 key 的旧值(old value)。
func GetSet(key string, val interface{}) (string, error) {
	cmd := client.GetSet(ctx, associate(key), val)
	return cmd.Val(), cmd.Err()
}

// GetBit 5	GETBIT key offset 对 key 所储存的字符串值，获取指定偏移量上的位(bit)。
func GetBit(key string, offset int64) (int64, error) {
	cmd := client.GetBit(ctx, associate(key), offset)
	return cmd.Val(), cmd.Err()
}

// MGet 6	MGET key1 [key2..] 获取所有(一个或多个)给定 key 的值。
func MGet(to interface{}, keys ...string) []interface{} {
	var args []string
	for _, key := range keys {
		args = append(args, associate(key))
	}
	cmd := client.MGet(ctx, args...)
	result := make([]interface{}, len(cmd.Val()))
	resultType := reflect.TypeOf(to)
	for i, s := range cmd.Val() {
		val := reflect.New(resultType).Interface()
		if err := decode([]byte(s.(string)), val); err != nil {
			continue
		}
		result[i] = val
	}
	return result
}

// SetBit 7	SETBIT key offset value 对 key 所储存的字符串值，设置或清除指定偏移量上的位(bit)。
func SetBit(key string, offset int64, value int) (int64, error) {
	cmd := client.SetBit(ctx, associate(key), offset, value)
	return cmd.Val(), cmd.Err()
}

// SetEX 8	SETEX key seconds value 将值 value 关联到 key ，并将 key 的过期时间设为 seconds (以秒为单位)。
func SetEX(key string, val interface{}, timeout int64) error {
	return client.SetEX(ctx, associate(key), val, time.Duration(timeout)*time.Second).Err()
}

// Setnx 9	SETNX key value 只有在 key 不存在时设置 key 的值。
func Setnx(key string, value interface{}) (bool, error) {
	cmd := client.SetNX(ctx, associate(key), value, 0)
	return cmd.Val(), cmd.Err()
}

// SetnxExpire SETNX WITH EXPIRE (Second)
func SetnxExpire(key string, value interface{}, expire int64) (bool, error) {
	cmd := client.SetNX(ctx, associate(key), value, time.Duration(expire)*time.Second)
	return cmd.Val(), cmd.Err()
}

// SetRange 10	SETRANGE key offset value 用 value 参数覆写给定 key 所储存的字符串值，从偏移量 offset 开始。
func SetRange(key string, offset int64, value interface{}) error {
	encode, err := encode(value)
	if err != nil {
		return err
	}
	cmd := client.SetRange(ctx, associate(key), offset, string(encode))
	return cmd.Err()
}

// Strlen 11	STRLEN key 返回 key 所储存的字符串值的长度。
func Strlen(key string) (int64, error) {
	cmd := client.StrLen(ctx, associate(key))
	return cmd.Val(), cmd.Err()
}

// MSet 12	MSET key value [key value ...] 同时设置一个或多个 key-value 对。
func MSet(keysAndValues map[string]interface{}) error {
	var args []interface{}
	for key, value := range keysAndValues {
		encode, _ := encode(value)
		args = append(args, associate(key), string(encode))
	}

	cmd := client.MSet(ctx, args...)
	return cmd.Err()
}

// MSetnx 13	MSETNX key value [key value ...] 同时设置一个或多个 key-value 对，当且仅当所有给定 key 都不存在。
func MSetnx(keysAndValues map[string]interface{}) (bool, error) {
	var args []interface{}
	for key, value := range keysAndValues {
		encode, _ := encode(value)
		args = append(args, associate(key), string(encode))
	}
	cmd := client.MSetNX(ctx, args...)
	return cmd.Val(), cmd.Err()
}

// PSetEX 14	PSETEX key milliseconds value 这个命令和 SETEX 命令相似，但它以毫秒为单位设置 key 的生存时间，而不是像 SETEX 命令那样，以秒为单位。
func PSetEX(key string, val interface{}, timeout int64) error {
	encode, _ := encode(val)
	cmd := client.Set(ctx, associate(key), string(encode), time.Duration(timeout)*time.Millisecond)
	return cmd.Err()
}

// Incr 15	INCR key 将 key 中储存的数字值增一。
func Incr(key string) (int64, error) {
	cmd := client.Incr(ctx, associate(key))
	return cmd.Val(), cmd.Err()
}

// IncrBy 16	INCRBY key increment 将 key 所储存的值加上给定的增量值（increment） 。
func IncrBy(key string, val int64) (int64, error) {
	cmd := client.IncrBy(ctx, associate(key), val)
	return cmd.Val(), cmd.Err()
}

// IncrByFloat 17	INCRBYFLOAT key increment 将 key 所储存的值加上给定的浮点增量值（increment） 。
func IncrByFloat(key string, val float64) (float64, error) {
	cmd := client.IncrByFloat(ctx, associate(key), val)
	return cmd.Val(), cmd.Err()
}

// Decr 18	DECR key 将 key 中储存的数字值减一。
func Decr(key string) (int64, error) {
	cmd := client.Decr(ctx, associate(key))
	return cmd.Val(), cmd.Err()
}

// DecrBy 19	DECRBY key decrement key 所储存的值减去给定的减量值（decrement） 。
func DecrBy(key string, val int64) (int64, error) {
	cmd := client.DecrBy(ctx, associate(key), val)
	return cmd.Val(), cmd.Err()
}

// Append 20	APPEND key value 如果 key 已经存在并且是一个字符串， APPEND 命令将指定的 value 追加到该 key 原来值（value）的末尾。
func Append(key string, val interface{}) (int64, error) {
	encode, _ := encode(val)
	cmd := client.Append(ctx, associate(key), string(encode))
	return cmd.Val(), cmd.Err()
}
