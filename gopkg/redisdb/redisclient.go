package redisdb

import (
	"fmt"
	"github.com/go-redis/redis"
	"strings"
)

type RedisClient struct {
	client        *redis.Client
	clusterClient *redis.ClusterClient
	isCluster     bool
}

func NewRedisClient(opt *Options) (*RedisClient, error) {
	opt.init()
	c := RedisClient{
		client:        nil,
		clusterClient: nil,
		isCluster:     false,
	}
	switch opt.IsCluster {
	case IsNotCluster:
		opt := redis.Options{
			Addr:      fmt.Sprintf("%s:%s", opt.Host, opt.Port),
			Password:  opt.PassWord,
			TLSConfig: nil,
		}
		c.client = redis.NewClient(&opt)

	case IsCluster:
		addrs := strings.Split(opt.ClusterHost, ",")
		if len(addrs) < 1 {
			return nil, fmt.Errorf("redis cluster host is empty")
		}
		opt := redis.ClusterOptions{
			Addrs:     addrs,
			Password:  opt.PassWord,
			TLSConfig: nil,
		}
		c.clusterClient = redis.NewClusterClient(&opt)
		c.isCluster = true
	default:
		return nil, fmt.Errorf("redis config IsCluster is %d, must be 0 or 1", opt.IsCluster)
	}

	_, err := c.Ping()
	if err != nil {
		return nil, fmt.Errorf("redis client ping fail: %v", err)
	}
	return &c, nil
}

func (c *RedisClient) IsCluster() bool {
	return c.isCluster
}

func (c *RedisClient) GetClient() *redis.Client {
	return c.client
}

func (c *RedisClient) GetClusterClient() *redis.ClusterClient {
	return c.clusterClient
}

func (c *RedisClient) Ping() (string, error) {
	if c.isCluster {
		return c.clusterClient.Ping().Result()
	}
	return c.client.Ping().Result()
}

func (c *RedisClient) Do(args ...interface{}) (reply interface{}, err error) {
	if c.isCluster {
		return c.clusterClient.Do(args...).Result()
	}
	return c.client.Do(args...).Result()
}
