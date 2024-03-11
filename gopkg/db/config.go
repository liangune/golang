package db

import "time"

const (
	DisableSSL int = 0
	EnableSSL  int = 1
)

const (
	defaultMaxLifetime    = 5 * time.Minute
	defaultMaxOpenConns   = 100
	defaultMaxIdleConns   = 5
	defaultConnectTimeout = 15
)

// support databases MySQL, PostgreSQL, SQLite, SQL Server
const (
	MySQL      = "MySQL"
	PostgreSQL = "PostgreSQL"
	SQLite     = "SQLite"
	SQLServer  = "SQLServer"
	ClickHouse = "ClickHouse"
)

type DBConfig struct {
	Host         string        // 数据库地址
	Port         int           // 数据库的端口号
	Username     string        // 用户名
	Password     string        // 密码
	Dbname       string        // 数据库名称
	DbType       string        // 数据库类型
	Timeout      int           // 超时时间, 单位秒
	SSL          int           // 是否启用SSL连接
	MaxOpenConns int           // 数据库连接池最大连接数
	MaxIdleConns int           // 连接池最大允许的空闲连接数
	MaxLifetime  time.Duration // 连接可复用的最大时间
}

/*
DBConfig {
	Username: "root",
	Password: "123456",
	Host: "127.0.0.1",
	Port: 3306,
	Dbname: "dbtest",
	DbType: "postgres"
	Timeout: 10,
	SSL: 0,
	MaxOpenConns: 100,
	MaxIdleConns: 20,
	MaxLifetime: 5 * time.Minute,
}
*/
