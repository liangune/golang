package redisdb

import (
	"github.com/go-redis/redis"
	"time"
)

//------------------------------------------------------------------------------
/****** key ******/
func (c *RedisClient) Del(keys ...string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.Del(keys...).Result()
	}
	return c.client.Del(keys...).Result()
}

func (c *RedisClient) Unlink(keys ...string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.Unlink(keys...).Result()
	}
	return c.client.Unlink(keys...).Result()
}

func (c *RedisClient) Exists(keys ...string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.Exists(keys...).Result()
	}
	return c.client.Exists(keys...).Result()
}

func (c *RedisClient) Expire(key string, expiration time.Duration) (bool, error) {
	if c.isCluster {
		return c.clusterClient.Expire(key, expiration).Result()
	}
	return c.client.Expire(key, expiration).Result()
}

func (c *RedisClient) ExpireAt(key string, tm time.Time) (bool, error) {
	if c.isCluster {
		return c.clusterClient.ExpireAt(key, tm).Result()
	}
	return c.client.ExpireAt(key, tm).Result()
}

func (c *RedisClient) TTL(key string) (time.Duration, error) {
	if c.isCluster {
		return c.clusterClient.TTL(key).Result()
	}
	return c.client.TTL(key).Result()
}

func (c *RedisClient) PTTL(key string) (time.Duration, error) {
	if c.isCluster {
		return c.clusterClient.PTTL(key).Result()
	}
	return c.client.PTTL(key).Result()
}

func (c *RedisClient) Type(key string) (string, error) {
	if c.isCluster {
		return c.clusterClient.Type(key).Result()
	}
	return c.client.Type(key).Result()
}

func (c *RedisClient) Rename(key, newkey string) (string, error) {
	if c.isCluster {
		return c.clusterClient.Rename(key, newkey).Result()
	}
	return c.client.Rename(key, newkey).Result()
}

func (c *RedisClient) Keys(pattern string) ([]string, error) {
	if c.isCluster {
		return c.clusterClient.Keys(pattern).Result()
	}
	return c.client.Keys(pattern).Result()
}

func (c *RedisClient) Scan(cursor uint64, match string, count int64) ([]string, uint64, error) {
	if c.isCluster {
		return c.clusterClient.Scan(cursor, match, count).Result()
	}
	return c.client.Scan(cursor, match, count).Result()
}

func (c *RedisClient) Dump(key string) (string, error) {
	if c.isCluster {
		return c.clusterClient.Dump(key).Result()
	}
	return c.client.Dump(key).Result()
}

//------------------------------------------------------------------------------
/***** string ******/
func (c *RedisClient) Append(key, value string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.Append(key, value).Result()
	}
	return c.client.Append(key, value).Result()
}

func (c *RedisClient) Incr(key string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.Incr(key).Result()
	}
	return c.client.Incr(key).Result()
}

func (c *RedisClient) IncrBy(key string, value int64) (int64, error) {
	if c.isCluster {
		return c.clusterClient.IncrBy(key, value).Result()
	}
	return c.client.IncrBy(key, value).Result()
}

func (c *RedisClient) IncrByFloat(key string, value float64) (float64, error) {
	if c.isCluster {
		return c.clusterClient.IncrByFloat(key, value).Result()
	}
	return c.client.IncrByFloat(key, value).Result()
}

func (c *RedisClient) Decr(key string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.Decr(key).Result()
	}
	return c.client.Decr(key).Result()
}

func (c *RedisClient) DecrBy(key string, decrement int64) (int64, error) {
	if c.isCluster {
		return c.clusterClient.DecrBy(key, decrement).Result()
	}
	return c.client.DecrBy(key, decrement).Result()
}

func (c *RedisClient) GetSet(key string, value interface{}) (string, error) {
	if c.isCluster {
		return c.clusterClient.GetSet(key, value).Result()
	}
	return c.client.GetSet(key, value).Result()
}

func (c *RedisClient) Get(key string) (string, error) {
	if c.isCluster {
		return c.clusterClient.Get(key).Result()
	}
	return c.client.Get(key).Result()
}

func (c *RedisClient) GetRange(key string, start, end int64) (string, error) {
	if c.isCluster {
		return c.clusterClient.GetRange(key, start, end).Result()
	}
	return c.client.GetRange(key, start, end).Result()
}

func (c *RedisClient) MGet(keys ...string) ([]interface{}, error) {
	if c.isCluster {
		return c.clusterClient.MGet(keys...).Result()
	}
	return c.client.MGet(keys...).Result()
}

func (c *RedisClient) MSet(pairs ...interface{}) (string, error) {
	if c.isCluster {
		return c.clusterClient.MSet(pairs...).Result()
	}
	return c.client.MSet(pairs...).Result()
}

func (c *RedisClient) MSetNX(pairs ...interface{}) (bool, error) {
	if c.isCluster {
		return c.clusterClient.MSetNX(pairs...).Result()
	}
	return c.client.MSetNX(pairs...).Result()
}

func (c *RedisClient) Set(key string, value interface{}, expiration time.Duration) (string, error) {
	if c.isCluster {
		return c.clusterClient.Set(key, value, expiration).Result()
	}
	return c.client.Set(key, value, expiration).Result()
}

func (c *RedisClient) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	if c.isCluster {
		return c.clusterClient.SetNX(key, value, expiration).Result()
	}
	return c.client.SetNX(key, value, expiration).Result()
}

func (c *RedisClient) SetXX(key string, value interface{}, expiration time.Duration) (bool, error) {
	if c.isCluster {
		return c.clusterClient.SetXX(key, value, expiration).Result()
	}
	return c.client.MSetNX(key, value, expiration).Result()
}

func (c *RedisClient) SetRange(key string, offset int64, value string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.SetRange(key, offset, value).Result()
	}
	return c.client.SetRange(key, offset, value).Result()
}

func (c *RedisClient) StrLen(key string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.StrLen(key).Result()
	}
	return c.client.StrLen(key).Result()
}

//------------------------------------------------------------------------------
/* Bit */
func (c *RedisClient) GetBit(key string, offset int64) (int64, error) {
	if c.isCluster {
		return c.clusterClient.GetBit(key, offset).Result()
	}
	return c.client.GetBit(key, offset).Result()
}

func (c *RedisClient) SetBit(key string, offset int64, value int) (int64, error) {
	if c.isCluster {
		return c.clusterClient.SetBit(key, offset, value).Result()
	}
	return c.client.SetBit(key, offset, value).Result()
}

func (c *RedisClient) BitCount(key string, bitCount *redis.BitCount) (int64, error) {
	if c.isCluster {
		return c.clusterClient.BitCount(key, bitCount).Result()
	}
	return c.client.BitCount(key, bitCount).Result()
}

func (c *RedisClient) BitOpAnd(destKey string, keys ...string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.BitOpAnd(destKey, keys...).Result()
	}
	return c.client.BitOpAnd(destKey, keys...).Result()
}

func (c *RedisClient) BitOpOr(destKey string, keys ...string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.BitOpOr(destKey, keys...).Result()
	}
	return c.client.BitOpOr(destKey, keys...).Result()
}

func (c *RedisClient) BitOpXor(destKey string, keys ...string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.BitOpXor(destKey, keys...).Result()
	}
	return c.client.BitOpXor(destKey, keys...).Result()
}

func (c *RedisClient) BitOpNot(destKey string, key string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.BitOpNot(destKey, key).Result()
	}
	return c.client.BitOpNot(destKey, key).Result()
}

func (c *RedisClient) BitPos(key string, bit int64, pos ...int64) (int64, error) {
	if c.isCluster {
		return c.clusterClient.BitPos(key, bit, pos...).Result()
	}
	return c.client.HDel(key).Result()
}

//------------------------------------------------------------------------------
/****** Hash ******/
func (c *RedisClient) HDel(key string, fields ...string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.HDel(key, fields...).Result()
	}
	return c.client.HDel(key, fields...).Result()
}

func (c *RedisClient) HExists(key, field string) (bool, error) {
	if c.isCluster {
		return c.clusterClient.HExists(key, field).Result()
	}
	return c.client.HExists(key, field).Result()
}

func (c *RedisClient) HGet(key, field string) (string, error) {
	if c.isCluster {
		return c.clusterClient.HGet(key, field).Result()
	}
	return c.client.HGet(key, field).Result()
}

func (c *RedisClient) HGetAll(key string) (map[string]string, error) {
	if c.isCluster {
		return c.clusterClient.HGetAll(key).Result()
	}
	return c.client.HGetAll(key).Result()
}

func (c *RedisClient) HIncrBy(key, field string, incr int64) (int64, error) {
	if c.isCluster {
		return c.clusterClient.HIncrBy(key, field, incr).Result()
	}
	return c.client.HIncrBy(key, field, incr).Result()
}

func (c *RedisClient) HIncrByFloat(key, field string, incr float64) (float64, error) {
	if c.isCluster {
		return c.clusterClient.HIncrByFloat(key, field, incr).Result()
	}
	return c.client.HIncrByFloat(key, field, incr).Result()
}

func (c *RedisClient) HKeys(key string) ([]string, error) {
	if c.isCluster {
		return c.clusterClient.HKeys(key).Result()
	}
	return c.client.HKeys(key).Result()
}

func (c *RedisClient) HLen(key string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.HLen(key).Result()
	}
	return c.client.HLen(key).Result()
}

func (c *RedisClient) HMGet(key string, fields ...string) ([]interface{}, error) {
	if c.isCluster {
		return c.clusterClient.HMGet(key, fields...).Result()
	}
	return c.client.HMGet(key, fields...).Result()
}

func (c *RedisClient) HMSet(key string, fields map[string]interface{}) (string, error) {
	if c.isCluster {
		return c.clusterClient.HMSet(key, fields).Result()
	}
	return c.client.HMSet(key, fields).Result()
}

func (c *RedisClient) HSet(key, field string, value interface{}) (bool, error) {
	if c.isCluster {
		return c.clusterClient.HSet(key, field, value).Result()
	}
	return c.client.HSet(key, field, value).Result()
}

func (c *RedisClient) HSetNX(key, field string, value interface{}) (bool, error) {
	if c.isCluster {
		return c.clusterClient.HSetNX(key, field, value).Result()
	}
	return c.client.HSetNX(key, field, value).Result()
}

func (c *RedisClient) HVals(key string) ([]string, error) {
	if c.isCluster {
		return c.clusterClient.HVals(key).Result()
	}
	return c.client.HVals(key).Result()
}

func (c *RedisClient) HScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	if c.isCluster {
		return c.clusterClient.HScan(key, cursor, match, count).Result()
	}
	return c.client.HScan(key, cursor, match, count).Result()
}

//------------------------------------------------------------------------------
/****** List ******/
func (c *RedisClient) BLPop(timeout time.Duration, keys ...string) ([]string, error) {
	if c.isCluster {
		return c.clusterClient.BLPop(timeout, keys...).Result()
	}
	return c.client.BLPop(timeout, keys...).Result()
}

func (c *RedisClient) BRPop(timeout time.Duration, keys ...string) ([]string, error) {
	if c.isCluster {
		return c.clusterClient.BRPop(timeout, keys...).Result()
	}
	return c.client.BRPop(timeout, keys...).Result()
}

func (c *RedisClient) BRPopLPush(source, destination string, timeout time.Duration) (string, error) {
	if c.isCluster {
		return c.clusterClient.BRPopLPush(source, destination, timeout).Result()
	}
	return c.client.BRPopLPush(source, destination, timeout).Result()
}

func (c *RedisClient) LIndex(key string, index int64) (string, error) {
	if c.isCluster {
		return c.clusterClient.LIndex(key, index).Result()
	}
	return c.client.LIndex(key, index).Result()
}

func (c *RedisClient) LInsert(key, op string, pivot, value interface{}) (int64, error) {
	if c.isCluster {
		return c.clusterClient.LInsert(key, op, pivot, value).Result()
	}
	return c.client.LInsert(key, op, pivot, value).Result()
}

func (c *RedisClient) LInsertBefore(key string, pivot, value interface{}) (int64, error) {
	if c.isCluster {
		return c.clusterClient.LInsertBefore(key, pivot, value).Result()
	}
	return c.client.LInsertBefore(key, pivot, value).Result()
}

func (c *RedisClient) LInsertAfter(key string, pivot, value interface{}) (int64, error) {
	if c.isCluster {
		return c.clusterClient.LInsertAfter(key, pivot, value).Result()
	}
	return c.client.LInsertAfter(key, pivot, value).Result()
}

func (c *RedisClient) LLen(key string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.LLen(key).Result()
	}
	return c.client.LLen(key).Result()
}

func (c *RedisClient) LPop(key string) (string, error) {
	if c.isCluster {
		return c.clusterClient.LPop(key).Result()
	}
	return c.client.LPop(key).Result()
}

func (c *RedisClient) LPush(key string, values ...interface{}) (int64, error) {
	if c.isCluster {
		return c.clusterClient.LPush(key, values...).Result()
	}
	return c.client.LPush(key, values...).Result()
}

func (c *RedisClient) LPushX(key string, value interface{}) (int64, error) {
	if c.isCluster {
		return c.clusterClient.LPushX(key, value).Result()
	}
	return c.client.LPushX(key, value).Result()
}

func (c *RedisClient) LRange(key string, start, stop int64) ([]string, error) {
	if c.isCluster {
		return c.clusterClient.LRange(key, start, stop).Result()
	}
	return c.client.LRange(key, start, stop).Result()
}

func (c *RedisClient) LRem(key string, count int64, value interface{}) (int64, error) {
	if c.isCluster {
		return c.clusterClient.LRem(key, count, value).Result()
	}
	return c.client.LRem(key, count, value).Result()
}

func (c *RedisClient) LSet(key string, index int64, value interface{}) (string, error) {
	if c.isCluster {
		return c.clusterClient.LSet(key, index, value).Result()
	}
	return c.client.LSet(key, index, value).Result()
}

func (c *RedisClient) LTrim(key string, start, stop int64) (string, error) {
	if c.isCluster {
		return c.clusterClient.LTrim(key, start, stop).Result()
	}
	return c.client.LTrim(key, start, stop).Result()
}

func (c *RedisClient) RPop(key string) (string, error) {
	if c.isCluster {
		return c.clusterClient.RPop(key).Result()
	}
	return c.client.RPop(key).Result()
}

func (c *RedisClient) RPopLPush(source, destination string) (string, error) {
	if c.isCluster {
		return c.clusterClient.RPopLPush(source, destination).Result()
	}
	return c.client.RPopLPush(source, destination).Result()
}

func (c *RedisClient) RPush(key string, values ...interface{}) (int64, error) {
	if c.isCluster {
		return c.clusterClient.RPush(key, values...).Result()
	}
	return c.client.RPush(key, values...).Result()
}

func (c *RedisClient) RPushX(key string, value interface{}) (int64, error) {
	if c.isCluster {
		return c.clusterClient.RPushX(key, value).Result()
	}
	return c.client.RPushX(key, value).Result()
}

//------------------------------------------------------------------------------
/****** Set ******/
func (c *RedisClient) SAdd(key string, members ...interface{}) (int64, error) {
	if c.isCluster {
		return c.clusterClient.SAdd(key, members...).Result()
	}
	return c.client.SAdd(key, members...).Result()
}

func (c *RedisClient) SCard(key string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.SCard(key).Result()
	}
	return c.client.SCard(key).Result()
}

func (c *RedisClient) SDiff(keys ...string) ([]string, error) {
	if c.isCluster {
		return c.clusterClient.SDiff(keys...).Result()
	}
	return c.client.SDiff(keys...).Result()
}

func (c *RedisClient) SDiffStore(destination string, keys ...string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.SDiffStore(destination, keys...).Result()
	}
	return c.client.SDiffStore(destination, keys...).Result()
}

func (c *RedisClient) SInter(keys ...string) ([]string, error) {
	if c.isCluster {
		return c.clusterClient.SInter(keys...).Result()
	}
	return c.client.SInter(keys...).Result()
}

func (c *RedisClient) SInterStore(destination string, keys ...string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.SInterStore(destination, keys...).Result()
	}
	return c.client.SInterStore(destination, keys...).Result()
}

func (c *RedisClient) SIsMember(key string, member interface{}) (bool, error) {
	if c.isCluster {
		return c.clusterClient.SIsMember(key, member).Result()
	}
	return c.client.SIsMember(key, member).Result()
}

func (c *RedisClient) SMembers(key string) ([]string, error) {
	if c.isCluster {
		return c.clusterClient.SMembers(key).Result()
	}
	return c.client.SMembers(key).Result()
}

func (c *RedisClient) SMembersMap(key string) (map[string]struct{}, error) {
	if c.isCluster {
		return c.clusterClient.SMembersMap(key).Result()
	}
	return c.client.SMembersMap(key).Result()
}

func (c *RedisClient) SMove(source, destination string, member interface{}) (bool, error) {
	if c.isCluster {
		return c.clusterClient.SMove(source, destination, member).Result()
	}
	return c.client.SMove(source, destination, member).Result()
}

func (c *RedisClient) SPop(key string) (string, error) {
	if c.isCluster {
		return c.clusterClient.SPop(key).Result()
	}
	return c.client.SPop(key).Result()
}

func (c *RedisClient) SPopN(key string, count int64) ([]string, error) {
	if c.isCluster {
		return c.clusterClient.SPopN(key, count).Result()
	}
	return c.client.SPopN(key, count).Result()
}

func (c *RedisClient) SRandMember(key string) (string, error) {
	if c.isCluster {
		return c.clusterClient.SRandMember(key).Result()
	}
	return c.client.SRandMember(key).Result()
}

func (c *RedisClient) SRandMemberN(key string, count int64) ([]string, error) {
	if c.isCluster {
		return c.clusterClient.SRandMemberN(key, count).Result()
	}
	return c.client.SRandMemberN(key, count).Result()
}

func (c *RedisClient) SRem(key string, members ...interface{}) (int64, error) {
	if c.isCluster {
		return c.clusterClient.SRem(key, members...).Result()
	}
	return c.client.SRem(key, members...).Result()
}

func (c *RedisClient) SUnion(keys ...string) ([]string, error) {
	if c.isCluster {
		return c.clusterClient.SUnion(keys...).Result()
	}
	return c.client.SUnion(keys...).Result()
}

func (c *RedisClient) SUnionStore(destination string, keys ...string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.SUnionStore(destination, keys...).Result()
	}
	return c.client.SUnionStore(destination, keys...).Result()
}

func (c *RedisClient) SScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	if c.isCluster {
		return c.clusterClient.SScan(key, cursor, match, count).Result()
	}
	return c.client.SScan(key, cursor, match, count).Result()
}

//------------------------------------------------------------------------------
/****** ZSet ******/
func (c *RedisClient) ZAdd(key string, members ...redis.Z) (int64, error) {
	if c.isCluster {
		return c.clusterClient.ZAdd(key, members...).Result()
	}
	return c.client.ZAdd(key, members...).Result()
}

func (c *RedisClient) ZAddNX(key string, members ...redis.Z) (int64, error) {
	if c.isCluster {
		return c.clusterClient.ZAddNX(key, members...).Result()
	}
	return c.client.ZAddNX(key, members...).Result()
}

func (c *RedisClient) ZAddXX(key string, members ...redis.Z) (int64, error) {
	if c.isCluster {
		return c.clusterClient.ZAddXX(key, members...).Result()
	}
	return c.client.ZAddXX(key, members...).Result()
}

func (c *RedisClient) ZAddCh(key string, members ...redis.Z) (int64, error) {
	if c.isCluster {
		return c.clusterClient.ZAddCh(key, members...).Result()
	}
	return c.client.ZAddCh(key, members...).Result()
}

func (c *RedisClient) ZAddNXCh(key string, members ...redis.Z) (int64, error) {
	if c.isCluster {
		return c.clusterClient.ZAddNXCh(key, members...).Result()
	}
	return c.client.ZAddNXCh(key, members...).Result()
}

func (c *RedisClient) ZAddXXCh(key string, members ...redis.Z) (int64, error) {
	if c.isCluster {
		return c.clusterClient.ZAddXXCh(key, members...).Result()
	}
	return c.client.ZAddXXCh(key, members...).Result()
}

func (c *RedisClient) ZIncr(key string, member redis.Z) (float64, error) {
	if c.isCluster {
		return c.clusterClient.ZIncr(key, member).Result()
	}
	return c.client.ZIncr(key, member).Result()
}

func (c *RedisClient) ZIncrNX(key string, member redis.Z) (float64, error) {
	if c.isCluster {
		return c.clusterClient.ZIncrNX(key, member).Result()
	}
	return c.client.ZIncrNX(key, member).Result()
}

func (c *RedisClient) ZIncrXX(key string, member redis.Z) (float64, error) {
	if c.isCluster {
		return c.clusterClient.ZIncrXX(key, member).Result()
	}
	return c.client.ZIncrXX(key, member).Result()
}

func (c *RedisClient) ZCard(key string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.ZCard(key).Result()
	}
	return c.client.ZCard(key).Result()
}

func (c *RedisClient) ZCount(key, min, max string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.ZCount(key, min, max).Result()
	}
	return c.client.ZCount(key, min, max).Result()
}

func (c *RedisClient) ZLexCount(key, min, max string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.ZLexCount(key, min, max).Result()
	}
	return c.client.ZLexCount(key, min, max).Result()
}

func (c *RedisClient) ZIncrBy(key string, increment float64, member string) (float64, error) {
	if c.isCluster {
		return c.clusterClient.ZIncrBy(key, increment, member).Result()
	}
	return c.client.ZIncrBy(key, increment, member).Result()
}

func (c *RedisClient) ZInterStore(destination string, store redis.ZStore, keys ...string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.ZInterStore(destination, store, keys...).Result()
	}
	return c.client.ZInterStore(destination, store, keys...).Result()
}

func (c *RedisClient) ZPopMax(key string, count ...int64) ([]redis.Z, error) {
	if c.isCluster {
		return c.clusterClient.ZPopMax(key, count...).Result()
	}
	return c.client.ZPopMax(key, count...).Result()
}

func (c *RedisClient) ZPopMin(key string, count ...int64) ([]redis.Z, error) {
	if c.isCluster {
		return c.clusterClient.ZPopMax(key, count...).Result()
	}
	return c.client.ZPopMax(key, count...).Result()
}

func (c *RedisClient) ZRange(key string, start, stop int64) ([]string, error) {
	if c.isCluster {
		return c.clusterClient.ZRange(key, start, stop).Result()
	}
	return c.client.ZRange(key, start, stop).Result()
}

func (c *RedisClient) ZRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	if c.isCluster {
		return c.clusterClient.ZRangeWithScores(key, start, stop).Result()
	}
	return c.client.ZRangeWithScores(key, start, stop).Result()
}

func (c *RedisClient) ZRangeByScore(key string, opt redis.ZRangeBy) ([]string, error) {
	if c.isCluster {
		return c.clusterClient.ZRangeByScore(key, opt).Result()
	}
	return c.client.ZRangeByScore(key, opt).Result()
}

func (c *RedisClient) ZRangeByLex(key string, opt redis.ZRangeBy) ([]string, error) {
	if c.isCluster {
		return c.clusterClient.ZRangeByLex(key, opt).Result()
	}
	return c.client.ZRangeByLex(key, opt).Result()
}

func (c *RedisClient) ZRangeByScoreWithScores(key string, opt redis.ZRangeBy) ([]redis.Z, error) {
	if c.isCluster {
		return c.clusterClient.ZRangeByScoreWithScores(key, opt).Result()
	}
	return c.client.ZRangeByScoreWithScores(key, opt).Result()
}

func (c *RedisClient) ZRank(key, member string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.ZRank(key, member).Result()
	}
	return c.client.ZRank(key, member).Result()
}

func (c *RedisClient) ZRem(key string, members ...interface{}) (int64, error) {
	if c.isCluster {
		return c.clusterClient.ZRem(key, members).Result()
	}
	return c.client.ZRem(key, members).Result()
}

func (c *RedisClient) ZRemRangeByRank(key string, start, stop int64) (int64, error) {
	if c.isCluster {
		return c.clusterClient.ZRemRangeByRank(key, start, stop).Result()
	}
	return c.client.ZRemRangeByRank(key, start, stop).Result()
}

func (c *RedisClient) ZRemRangeByScore(key, min, max string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.ZRemRangeByScore(key, min, max).Result()
	}
	return c.client.ZRemRangeByScore(key, min, max).Result()
}

func (c *RedisClient) ZRemRangeByLex(key, min, max string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.ZRemRangeByLex(key, min, max).Result()
	}
	return c.client.ZRemRangeByLex(key, min, max).Result()
}

func (c *RedisClient) ZRevRange(key string, start, stop int64) ([]string, error) {
	if c.isCluster {
		return c.clusterClient.ZRevRange(key, start, stop).Result()
	}
	return c.client.ZRevRange(key, start, stop).Result()
}

func (c *RedisClient) ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	if c.isCluster {
		return c.clusterClient.ZRevRangeWithScores(key, start, stop).Result()
	}
	return c.client.ZRevRangeWithScores(key, start, stop).Result()
}

func (c *RedisClient) ZRevRangeByScore(key string, opt redis.ZRangeBy) ([]string, error) {
	if c.isCluster {
		return c.clusterClient.ZRevRangeByScore(key, opt).Result()
	}
	return c.client.ZRevRangeByScore(key, opt).Result()
}

func (c *RedisClient) ZRevRangeByLex(key string, opt redis.ZRangeBy) ([]string, error) {
	if c.isCluster {
		return c.clusterClient.ZRevRangeByLex(key, opt).Result()
	}
	return c.client.ZRevRangeByLex(key, opt).Result()
}

func (c *RedisClient) ZRevRangeByScoreWithScores(key string, opt redis.ZRangeBy) ([]redis.Z, error) {
	if c.isCluster {
		return c.clusterClient.ZRevRangeByScoreWithScores(key, opt).Result()
	}
	return c.client.ZRevRangeByScoreWithScores(key, opt).Result()
}

func (c *RedisClient) ZRevRank(key, member string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.ZRevRank(key, member).Result()
	}
	return c.client.ZRevRank(key, member).Result()
}

func (c *RedisClient) ZScore(key, member string) (float64, error) {
	if c.isCluster {
		return c.clusterClient.ZScore(key, member).Result()
	}
	return c.client.ZScore(key, member).Result()
}

func (c *RedisClient) ZUnionStore(dest string, store redis.ZStore, keys ...string) (int64, error) {
	if c.isCluster {
		return c.clusterClient.ZUnionStore(dest, store, keys...).Result()
	}
	return c.client.ZUnionStore(dest, store, keys...).Result()
}

//------------------------------------------------------------------------------
