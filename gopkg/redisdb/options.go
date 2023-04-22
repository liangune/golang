package redisdb

import "time"

const (
	// 空闲连接数
	DefaultMinIdleConns = 5
	DefaultHost         = "127.0.0.1"
	DefaultPort         = "6379"
	// 空闲连接检查的周期
	DefaultIdleCheckFrequency = time.Minute
	// 空闲连接的超时时间
	DefaultIdleTimeout = 5 * time.Minute
)

const (
	IsNotCluster = iota
	IsCluster
)

type Options struct {
	Host               string
	Port               string
	ClusterHost        string
	IsCluster          int
	PassWord           string
	SSL                int
	MinIdleConns       int
	IdleCheckFrequency time.Duration
	IdleTimeout        time.Duration
}

func (c *Options) init() {
	if c.Host == "" {
		c.Host = DefaultHost
	}
	if c.Port == "" {
		c.Port = DefaultPort
	}
	if c.MinIdleConns <= 0 {
		c.MinIdleConns = DefaultMinIdleConns
	}
	if c.IdleCheckFrequency <= 0 {
		c.IdleCheckFrequency = DefaultIdleCheckFrequency
	}
	if c.IdleTimeout <= 0 {
		c.IdleTimeout = DefaultIdleTimeout
	}
}
