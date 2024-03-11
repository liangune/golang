package redisdb

import (
	"runtime"
	"time"
)

const (
	// 空闲连接数
	DefaultMinIdleConnections = 5
	DefaultHost               = "127.0.0.1"
	DefaultPort               = "6379"
	// 空闲连接检查的周期
	DefaultIdleCheckFrequency = time.Minute
	// 空闲连接的超时时间
	DefaultIdleTimeout = 5 * time.Minute
	// 连接池数量
	DefaultPoolSize = 100
)

const (
	IsNotCluster = iota
	IsCluster
)

const (
	DisableSSL int = 0
	EnableSSL  int = 1
)

type Options struct {
	Host               string
	Port               string
	ClusterHost        string
	IsCluster          int
	Password           string
	SSL                int
	MinIdleConns       int
	IdleCheckFrequency time.Duration
	IdleTimeout        time.Duration
	PoolSize           int
}

func NewOptions() *Options {
	op := &Options{}
	op.Host = DefaultHost
	op.Port = DefaultPort
	op.IsCluster = IsNotCluster
	op.MinIdleConns = DefaultMinIdleConnections
	op.IdleCheckFrequency = DefaultIdleCheckFrequency
	op.IdleTimeout = DefaultIdleTimeout
	op.SSL = DisableSSL
	if DefaultPoolSize > 10*runtime.NumCPU() {
		op.PoolSize = DefaultPoolSize
	}

	return op
}
