package db

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBPool struct {
	db *gorm.DB
}

func DBPoolInit(dbConfig *DBConfig) (*DBPool, error) {
	username := dbConfig.Username
	password := dbConfig.Password
	host := dbConfig.Host
	port := dbConfig.Port
	dbname := dbConfig.Dbname
	timeout := dbConfig.Timeout
	if timeout <= 0 {
		timeout = defaultConnectTimeout
	}

	var dsn string
	switch dbConfig.DbType {
	case MySQL:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%ds", username, password, host, port, dbname, timeout)
	case PostgreSQL:
		//sslmode是安全验证模式
		if dbConfig.SSL == EnableSSL {
			dsn = fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=require password=%s connect_timeout=%d", host, port, username, dbname, password, timeout)
		} else {
			dsn = fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s connect_timeout=%d", host, port, username, dbname, password, timeout)
		}
	default:
		return nil, errors.New("DBType is not support, please check db package")
	}

	p := &DBPool{}

	var err error
	switch dbConfig.DbType {
	case MySQL:
		p.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	case PostgreSQL:
		p.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	default:
		return nil, errors.New("DBType is not support, please check db package")
	}

	if err != nil {
		return nil, fmt.Errorf("连接数据库失败, error=%v", err)
	}

	sqlDB, err := p.db.DB()
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	//连接池最大连接数
	if dbConfig.MaxOpenConns <= 0 {
		sqlDB.SetMaxOpenConns(defaultMaxOpenConns)
	} else {
		sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	}
	//连接池最大允许的空闲连接数
	if dbConfig.MaxIdleConns <= 0 {
		sqlDB.SetMaxOpenConns(defaultMaxIdleConns)
	} else {
		sqlDB.SetMaxOpenConns(dbConfig.MaxIdleConns)
	}
	//连接最大可重用的时间
	if dbConfig.MaxLifetime <= 0 {
		sqlDB.SetConnMaxLifetime(defaultMaxLifetime)
	} else {
		sqlDB.SetConnMaxLifetime(dbConfig.MaxLifetime)
	}

	return p, nil
}

func (p *DBPool) GetDB() *gorm.DB {
	return p.db
}

func (p *DBPool) Close() {
	sqlDB, _ := p.db.DB()
	sqlDB.Close()
}
