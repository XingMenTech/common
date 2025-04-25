package redis

import "time"

// 1	BLPOP key1 [key2 ] timeout 移出并获取列表的第一个元素， 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止。
func Blpop(key string, timeout int) (string, error) {
	cmd := client.BLPop(ctx, time.Duration(timeout)*time.Second, associate(key))
	return cmd.Val()[0], cmd.Err()
}

// 2	BRPOP key1 [key2 ] timeout 移出并获取列表的最后一个元素， 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止。
func Brpop(key string, timeout int) (string, error) {
	cmd := client.BRPop(ctx, time.Duration(timeout)*time.Second, associate(key))
	return cmd.Val()[0], cmd.Err()
}

// 3	BRPOPLPUSH source destination timeout 从列表中弹出一个值，将弹出的元素插入到另外一个列表中并返回它； 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止。
func Brpoplpush(source, destination string, timeout int) (string, error) {
	cmd := client.BRPopLPush(ctx, associate(source), associate(destination), time.Duration(timeout)*time.Second)
	return cmd.Val(), cmd.Err()
}

// 4	LINDEX key index 通过索引获取列表中的元素
func Lindex(key string, index int) (string, error) {
	cmd := client.LIndex(ctx, associate(key), int64(index))
	return cmd.Val(), cmd.Err()
}

// 5	LINSERT key BEFORE|AFTER pivot value 在列表的元素前或者后插入元素
func Linsert(key string, before string, pivot, val interface{}) error {
	return client.LInsert(ctx, associate(key), before, pivot, val).Err()
}

// 6	LLEN key 获取列表长度
func Llen(key string) (int, error) {
	cmd := client.LLen(ctx, associate(key))
	return int(cmd.Val()), cmd.Err()
}

// 7	LPOP key 移出并获取列表的第一个元素
func Lpop(key string) (string, error) {
	cmd := client.LPop(ctx, associate(key))
	return cmd.Val(), cmd.Err()
}

// 8	LPUSH key value1 [value2] 将一个或多个值插入到列表头部
func LPush(key string, val interface{}) error {
	return client.LPush(ctx, associate(key), val).Err()
}

// 9	LPUSHX key value 将一个值插入到已存在的列表头部
func LPushX(key string, val interface{}) error {
	return client.LPushX(ctx, associate(key), val).Err()
}

// 10	LRANGE key start stop 获取列表指定范围内的元素
func Lrange(key string, start, stop int) ([]string, error) {
	cmd := client.LRange(ctx, associate(key), int64(start), int64(stop))
	return cmd.Val(), cmd.Err()
}

// 11	LREM key count value 移除列表元素
func LRem(key string, index int, val interface{}) (int, error) {
	cmd := client.LRem(ctx, associate(key), int64(index), val)
	return int(cmd.Val()), cmd.Err()
}

// 12	LSET key index value 通过索引设置列表元素的值
func Lset(key string, index int, val interface{}) error {
	return client.LSet(ctx, associate(key), int64(index), val).Err()
}

// 13	LTRIM key start stop 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
func Ltrim(key string, start, stop int) error {
	return client.LTrim(ctx, associate(key), int64(start), int64(stop)).Err()
}

// 14	RPOP key 移除列表的最后一个元素，返回值为移除的元素。
func RPop(key string) (string, error) {
	cmd := client.RPop(ctx, associate(key))
	return cmd.Val(), cmd.Err()
}

// 15	RPOPLPUSH source destination 移除列表的最后一个元素，并将该元素添加到另一个列表并返回
func RPopLPush(source, destination string) (string, error) {
	cmd := client.RPopLPush(ctx, associate(source), associate(destination))
	return cmd.Val(), cmd.Err()
}

// 16	RPUSH key value1 [value2] 在列表中添加一个或多个值到列表尾部
func Rpush(key string, val interface{}) error {
	return client.RPush(ctx, associate(key), val).Err()
}

// 17	RPUSHX key value 为已存在的列表添加值
func RPushX(key string, val interface{}) error {
	return client.RPushX(ctx, associate(key), val).Err()
}
