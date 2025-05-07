package redis

// 1	SADD key member1 [member2] 向集合添加一个或多个成员
func SAdd(key string, members ...interface{}) (int, error) {
	cmd := client.SAdd(ctx, associate(key), members...)
	return int(cmd.Val()), cmd.Err()
}

// 2	SCARD key 获取集合的成员数
func SCard(key string) int64 {
	cmd := client.SCard(ctx, associate(key))
	return cmd.Val()
}

// 3	SDIFF key1 [key2] 返回第一个集合与其他集合之间的差异。
func SDiff(keys ...string) ([]string, error) {
	cmd := client.SDiff(ctx, keys...)
	return cmd.Val(), cmd.Err()
}

// 4	SDIFFSTORE destination key1 [key2] 返回给定所有集合的差集并存储在 destination 中
func SDiffStore(destination string, keys ...string) (int, error) {
	cmd := client.SDiffStore(ctx, associate(destination), keys...)
	return int(cmd.Val()), cmd.Err()
}

// 5	SINTER key1 [key2] 返回给定所有集合的交集
func SInter(keys ...string) ([]string, error) {
	cmd := client.SInter(ctx, keys...)
	return cmd.Val(), cmd.Err()
}

// 6	SINTERSTORE destination key1 [key2] 返回给定所有集合的交集并存储在 destination 中
func SInterStore(destination string, keys ...string) (int, error) {
	cmd := client.SInterStore(ctx, associate(destination), keys...)
	return int(cmd.Val()), cmd.Err()
}

// 7	SISMEMBER key member 判断 member 元素是否是集合 key 的成员
func SIsMember(key string, member interface{}) (bool, error) {
	cmd := client.SIsMember(ctx, associate(key), member)
	return cmd.Val(), cmd.Err()
}

// 8	SMEMBERS key 返回集合中的所有成员
func SMembers(key string) ([]string, error) {
	cmd := client.SMembers(ctx, associate(key))
	return cmd.Val(), cmd.Err()
}

// 9	SMOVE source destination member 将 member 元素从 source 集合移动到 destination 集合
func SMove(source, destination string, member interface{}) (bool, error) {
	cmd := client.SMove(ctx, associate(source), associate(destination), member)
	return cmd.Val(), cmd.Err()
}

// 10	SPOP key 移除并返回集合中的一个随机元素
func SPop(key string) (string, error) {
	cmd := client.SPop(ctx, associate(key))
	return cmd.Val(), cmd.Err()
}

// 11	SRANDMEMBER key [count] 返回集合中一个或多个随机数
func SRandMember(key string, count int) ([]string, error) {
	cmd := client.SRandMemberN(ctx, associate(key), int64(count))
	return cmd.Val(), cmd.Err()
}

// 12	SREM key member1 [member2] 移除集合中一个或多个成员
func SRem(key string, members ...interface{}) (int, error) {
	cmd := client.SRem(ctx, associate(key), members...)
	return int(cmd.Val()), cmd.Err()
}

// 13	SUNION key1 [key2] 返回所有给定集合的并集
func SUnion(keys ...string) ([]string, error) {
	cmd := client.SUnion(ctx, keys...)
	return cmd.Val(), cmd.Err()
}

// 14	SUNIONSTORE destination key1 [key2] 所有给定集合的并集存储在 destination 集合中
func SUnionStore(destination string, keys ...string) (int, error) {
	cmd := client.SUnionStore(ctx, associate(destination), keys...)
	return int(cmd.Val()), cmd.Err()
}
