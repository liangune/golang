package redisdb

import (
	"github.com/go-redis/redis"
	"time"
)

type Cmd interface {
	/****** key ******/
	// 删除指定一个key或者多个key
	Del(keys ...string) (int64, error)

	// 异步删除指定一个key或者多个key, 不阻塞
	Unlink(keys ...string) (int64, error)

	Exists(keys ...string) (int64, error)
	Expire(key string, expiration time.Duration) (bool, error)
	ExpireAt(key string, tm time.Time) (bool, error)
	TTL(key string) (time.Duration, error)
	PTTL(key string) (time.Duration, error)
	Type(key string) (string, error)
	Rename(key, newkey string) (string, error)
	Keys(pattern string) ([]string, error)

	/*
		功能: 迭代当前数据库中的数据库键
		参数: cursor 游标, 游标参数被设置为0时, 开始新的一次迭代
		参数: match 匹配字符串, 支持通配符*
		参数: count 默认值是10
		返回值: key数组, 游标, 错误
	*/
	Scan(cursor uint64, match string, count int64) ([]string, uint64, error)
	//
	Dump(key string) (string, error)

	/****** String ******/
	Append(key, value string) (int, error)

	Incr(key string) (int, error)
	IncrBy(key string, value int64) (int, error)
	IncrByFloat(key string, value float64) (float64, error)
	Decr(key string) (int, error)
	DecrBy(key string, decrement int64) (int, error)

	GetSet(key string, value interface{}) (string, error)
	Get(key string) (string, error)
	GetRange(key string, start, end int64) (string, error)
	MGet(keys ...string) ([]interface{}, error)
	MSet(pairs ...interface{}) (string, error)
	MSetNX(pairs ...interface{}) (bool, error)
	Set(key string, value interface{}, expiration time.Duration) (string, error)
	SetNX(key string, value interface{}, expiration time.Duration) (bool, error)
	SetXX(key string, value interface{}, expiration time.Duration) (bool, error)
	SetRange(key string, offset int64, value string) (int64, error)
	StrLen(key string) (int, error)

	/* Bit */
	GetBit(key string, offset int64) (int64, error)
	SetBit(key string, offset int64, value int) (int64, error)
	BitCount(key string, bitCount *redis.BitCount) (int64, error)
	BitOpAnd(destKey string, keys ...string) (int64, error)
	BitOpOr(destKey string, keys ...string) (int64, error)
	BitOpXor(destKey string, keys ...string) (int64, error)
	BitOpNot(destKey string, key string) (int64, error)
	BitPos(key string, bit int64, pos ...int64) (int64, error)

	/****** Hash ******/
	HDel(key string, fields ...string) (int64, error)
	HExists(key, field string) (bool, error)
	HGet(key, field string) (string, error)
	HGetAll(key string) (map[string]string, error)
	HIncrBy(key, field string, incr int64) (int64, error)
	HIncrByFloat(key, field string, incr float64) (float64, error)
	HKeys(key string) ([]string, error)
	HLen(key string) (int64, error)
	HMGet(key string, fields ...string) ([]interface{}, error)
	HMSet(key string, fields map[string]interface{}) (string, error)
	HSet(key, field string, value interface{}) (bool, error)
	HSetNX(key, field string, value interface{}) (bool, error)
	HVals(key string) ([]string, error)
	HScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error)

	/****** List ******/
	BLPop(timeout time.Duration, keys ...string) ([]string, error)
	BRPop(timeout time.Duration, keys ...string) ([]string, error)
	BRPopLPush(source, destination string, timeout time.Duration) (string, error)

	LIndex(key string, index int64) (string, error)
	LInsert(key, op string, pivot, value interface{}) (int64, error)
	LInsertBefore(key string, pivot, value interface{}) (int64, error)
	LInsertAfter(key string, pivot, value interface{}) (int64, error)
	LLen(key string) (int64, error)
	LPop(key string) (string, error)

	LPush(key string, values ...interface{}) (int64, error)
	LPushX(key string, value interface{}) (int64, error)
	LRange(key string, start, stop int64) ([]string, error)
	LRem(key string, count int64, value interface{}) (int64, error)
	LSet(key string, index int64, value interface{}) (string, error)
	LTrim(key string, start, stop int64) (string, error)
	RPop(key string) (string, error)
	RPopLPush(source, destination string) (string, error)
	RPush(key string, values ...interface{}) (int64, error)
	RPushX(key string, value interface{}) (int64, error)

	/****** Set ******/
	SAdd(key string, members ...interface{}) (int64, error)
	SCard(key string) (int64, error)
	SDiff(keys ...string) ([]string, error)
	SDiffStore(destination string, keys ...string) (int64, error)
	SInter(keys ...string) ([]string, error)
	SInterStore(destination string, keys ...string) (int64, error)
	SIsMember(key string, member interface{}) (bool, error)
	SMembers(key string) ([]string, error)
	SMembersMap(key string) (map[string]struct{}, error)
	SMove(source, destination string, member interface{}) (bool, error)
	SPop(key string) (string, error)
	SPopN(key string, count int64) ([]string, error)
	SRandMember(key string) (string, error)
	SRandMemberN(key string, count int64) ([]string, error)
	SRem(key string, members ...interface{}) (int64, error)
	SUnion(keys ...string) ([]string, error)
	SUnionStore(destination string, keys ...string) (int64, error)
	SScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error)

	/****** ZSet ******/
	ZAdd(key string, members ...redis.Z) (int64, error)
	ZAddNX(key string, members ...redis.Z) (int64, error)
	ZAddXX(key string, members ...redis.Z) (int64, error)
	ZAddCh(key string, members ...redis.Z) (int64, error)
	ZAddNXCh(key string, members ...redis.Z) (int64, error)
	ZAddXXCh(key string, members ...redis.Z) (int64, error)

	ZIncr(key string, member redis.Z) (float64, error)
	ZIncrNX(key string, member redis.Z) (float64, error)
	ZIncrXX(key string, member redis.Z) (float64, error)
	ZCard(key string) (int64, error)
	ZCount(key, min, max string) (int64, error)
	ZLexCount(key, min, max string) (int64, error)
	ZIncrBy(key string, increment float64, member string) (float64, error)
	ZInterStore(destination string, store redis.ZStore, keys ...string) (int64, error)
	ZPopMax(key string, count ...int64) ([]redis.Z, error)
	ZPopMin(key string, count ...int64) ([]redis.Z, error)
	ZRange(key string, start, stop int64) ([]string, error)
	ZRangeWithScores(key string, start, stop int64) ([]redis.Z, error)
	ZRangeByScore(key string, opt redis.ZRangeBy) ([]string, error)
	ZRangeByLex(key string, opt redis.ZRangeBy) ([]string, error)
	ZRangeByScoreWithScores(key string, opt redis.ZRangeBy) ([]redis.Z, error)
	ZRank(key, member string) (int64, error)
	ZRem(key string, members ...interface{}) (int64, error)
	ZRemRangeByRank(key string, start, stop int64) (int64, error)
	ZRemRangeByScore(key, min, max string) (int64, error)
	ZRemRangeByLex(key, min, max string) (int64, error)
	ZRevRange(key string, start, stop int64) ([]string, error)
	ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error)
	ZRevRangeByScore(key string, opt redis.ZRangeBy) ([]string, error)
	ZRevRangeByLex(key string, opt redis.ZRangeBy) ([]string, error)
	ZRevRangeByScoreWithScores(key string, opt redis.ZRangeBy) ([]redis.Z, error)
	ZRevRank(key, member string) (int64, error)
	ZScore(key, member string) (float64, error)
	ZUnionStore(dest string, store redis.ZStore, keys ...string) (int64, error)

	/****** HyperLogLog ******/
	PFAdd(key string, els ...interface{}) (int64, error)
	PFCount(keys ...string) (int64, error)
	PFMerge(dest string, keys ...string) (string, error)

	/****** Stream ******/

	/****** Transaction ******/
}
