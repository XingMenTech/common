package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
)

// ZADD 向有序集合添加一个或多个成员，或者更新已存在成员的分数
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

// ZCARD 获取有序集合的成员数
func ZCard(key string) (int64, error) {
	cmd := client.ZCard(ctx, associate(key))
	return cmd.Val(), cmd.Err()
}

// ZCOUNT 计算在有序集合中指定区间分数的成员数
func ZCount(key, min, max string) (int64, error) {
	cmd := client.ZCount(ctx, associate(key), min, max)
	return cmd.Val(), cmd.Err()
}

// ZSCORE
func ZScore(key, member string) (string, error) {
	cmd := client.ZScore(ctx, associate(key), member)
	return fmt.Sprintf("%.2f", cmd.Val()), cmd.Err()
}

// ZINCRBY 有序集合中对指定成员的分数加上增量 increment
func ZIncrBy(key, member string, increment float64) (string, error) {
	cmd := client.ZIncrBy(ctx, associate(key), increment, member)
	return fmt.Sprintf("%.2f", cmd.Val()), cmd.Err()
}

// ZINTERSTORE 计算给定的一个或多个有序集的交集并将结果集存储在新的有序集合 destination 中
func ZInterStore(destination string, keys ...string) (int64, error) {
	cmd := client.ZInterStore(ctx, associate(destination), &redis.ZStore{
		Keys: keys,
	})
	return cmd.Val(), cmd.Err()
}

// ZLEXCOUNT 在有序集合中计算指定字典区间内成员数量
func ZLexCount(key, min, max string) (int64, error) {
	cmd := client.ZLexCount(ctx, associate(key), min, max)
	return cmd.Val(), cmd.Err()
}

// ZRANGE 通过索引区间返回有序集合指定区间内的成员
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

// ZRANGEBYLEX key min max [LIMIT offset count] 通过字典区间返回有序集合的成员
func ZRangeByLex(key, min, max string) ([]string, error) {
	cmd := client.ZRangeByLex(ctx, associate(key), &redis.ZRangeBy{
		Min: min,
		Max: max,
	})
	return cmd.Val(), cmd.Err()

}

// ZRANGEBYSCORE 通过分数返回有序集合指定区间内的成员
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

// ZRANK key member 返回有序集合中指定成员的索引
func ZRank(key, member string) (int64, error) {
	cmd := client.ZRank(ctx, associate(key), member)
	return cmd.Val(), cmd.Err()
}

// ZREVRANK key member 返回有序集合中指定成员的排名，有序集成员按分数值递减(从大到小)排序
func ZRevRank(key, member string) (int64, error) {
	cmd := client.ZRevRank(ctx, associate(key), member)
	return cmd.Val(), cmd.Err()
}

// ZREVRANGE 返回有序集中指定区间内的成员，通过索引，分数从高到低
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

// ZREVRANGEBYSCORE key max min [WITHSCORES] 返回有序集中指定分数区间内的成员，分数从高到低排序
func ZRevRangeByScore(key string, max, min int64, withscores bool) ([]string, error) {
	if withscores {
		cmd := client.ZRevRangeByScoreWithScores(ctx, associate(key), &redis.ZRangeBy{
			Min: strconv.FormatInt(max, 10),
			Max: strconv.FormatInt(min, 10),
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
		cmd := client.ZRevRangeByScore(ctx, associate(key), &redis.ZRangeBy{
			Min: strconv.FormatInt(max, 10),
			Max: strconv.FormatInt(min, 10),
		})
		return cmd.Val(), cmd.Err()
	}
}

// ZREM 移除有序集合中的一个或多个成员
func ZRem(key string, values ...interface{}) (int, error) {
	cmd := client.ZRem(ctx, associate(key), values...)
	return int(cmd.Val()), cmd.Err()
}

// ZREMRANGEBYLEX key min max 移除有序集合中给定的字典区间的所有成员
func ZRemRangeByLex(key, min, max string) (int, error) {
	cmd := client.ZRemRangeByLex(ctx, associate(key), min, max)
	return int(cmd.Val()), cmd.Err()
}

// ZREMRANGEBYRANK
// ZREMRANGEBYRANK key start stop 移除有序集合中给定的排名区间的所有成员
func ZRemRangeByRank(key string, start, stop int64) (int, error) {
	cmd := client.ZRemRangeByRank(ctx, associate(key), start, stop)
	return int(cmd.Val()), cmd.Err()
}

// ZREMRANGEBYSCORE key min max 移除有序集合中给定的分数区间的所有成员
func ZRemRangeByScore(key string, min, max int64) (int, error) {
	cmd := client.ZRemRangeByScore(ctx, associate(key), strconv.FormatInt(min, 10), strconv.FormatInt(max, 10))
	return int(cmd.Val()), cmd.Err()
}

// ZUNIONSTORE destination numkeys key [key ...] 计算给定的一个或多个有序集的并集，并存储在新的 key 中
func ZUnionStore(destination string, keys ...string) (int64, error) {
	cmd := client.ZUnionStore(ctx, associate(destination), &redis.ZStore{
		Keys: keys,
	})
	return cmd.Val(), cmd.Err()
}
