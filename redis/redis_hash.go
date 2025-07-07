package redis

import "reflect"

// 1	HDEL key field1 [field2] 删除一个或多个哈希表字段
func HDel(key string, fields ...string) error {
	return client.HDel(ctx, associate(key), fields...).Err()
}

// 2	HEXISTS key field 查看哈希表 key 中，指定的字段是否存在。
func HExists(key, field string) (bool, error) {
	exists := client.HExists(ctx, associate(key), field)
	return exists.Val(), exists.Err()
}

// 3	HGET key field 获取存储在哈希表中指定字段的值。
func HGet(key string, field string, to interface{}) error {
	cmd := client.HGet(ctx, associate(key), field)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	if err := cmd.Scan(to); err != nil {
		return decode([]byte(cmd.Val()), to)
	}
	return nil
}

// 4	HGETALL key 获取在哈希表中指定 key 的所有字段和值
func HGetAll(key string, to interface{}) (error, map[string]interface{}) {
	cmd := client.HGetAll(ctx, associate(key))
	if cmd.Err() != nil {
		return cmd.Err(), nil
	}
	resultType := reflect.TypeOf(to)

	result := make(map[string]interface{})
	for k, v := range cmd.Val() {
		i := reflect.New(resultType).Interface()
		if err := decode([]byte(v), i); err != nil {
			continue
		}
		result[k] = i
	}

	return nil, result
}

// 5	HINCRBY key field increment 为哈希表 key 中的指定字段的整数值加上增量 increment 。
func HIncrby(key string, field string, incr int64) (int64, error) {
	cmd := client.HIncrBy(ctx, associate(key), field, incr)
	return cmd.Val(), cmd.Err()
}

// 6	HINCRBYFLOAT key field increment 为哈希表 key 中的指定字段的浮点数值加上增量 increment 。
func HIncrbyfloat(key string, field string, incr float64) (float64, error) {
	cmd := client.HIncrByFloat(ctx, associate(key), field, incr)
	return cmd.Val(), cmd.Err()
}

// 7	HKEYS key 获取哈希表中的所有字段
func HKeys(key string) ([]string, error) {
	cmd := client.HKeys(ctx, associate(key))
	return cmd.Val(), cmd.Err()
}

// 8	HLEN key 获取哈希表中字段的数量
func HLen(key string) int64 {
	hLen := client.HLen(ctx, associate(key))
	if hLen.Err() != nil {
		return 0
	}
	return hLen.Val()
}

// 9	HMGET key field1 [field2] 获取所有给定字段的值
func HMGet(key string, fields ...string) []interface{} {
	cmd := client.HMGet(ctx, associate(key), fields...)
	return cmd.Val()
}

// 10	HMSET key field1 value1 [field2 value2 ] 同时将多个 field-value (域-值)对设置到哈希表 key 中。
func HMSet(key string, fields map[string]interface{}) error {
	args := make([]interface{}, 0)
	for k, v := range fields {
		encode, err := encode(v)
		if err != nil {
			return err
		}
		args = append(args, k, string(encode))
	}
	return client.HMSet(ctx, associate(key), args).Err()

}

// 11	HSET key field value 将哈希表 key 中的字段 field 的值设为 value 。
func HSet(key string, field string, val interface{}) error {
	valByte, err := encode(val)
	if err != nil {
		return err
	}
	return client.HSet(ctx, associate(key), field, string(valByte)).Err()
}

// 12	HSETNX key field value 只有在字段 field 不存在时，设置哈希表字段的值。
func HSetnx(key string, field string, val interface{}) (bool, error) {
	valByte, err := encode(val)
	if err != nil {
		return false, err
	}
	return client.HSetNX(ctx, associate(key), field, string(valByte)).Result()
}

// 13	HVALS key 获取哈希表中所有值。
func HVals(key string, to interface{}) ([]interface{}, error) {
	cmd := client.HVals(ctx, associate(key))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	result := make([]interface{}, len(cmd.Val()))
	resultType := reflect.TypeOf(to)
	for i, s := range cmd.Val() {
		val := reflect.New(resultType).Interface()
		if err := decode([]byte(s), val); err != nil {
			continue
		}
		result[i] = val
	}
	return result, nil
}
