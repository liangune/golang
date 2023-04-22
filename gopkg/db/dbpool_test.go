package db

import "testing"

func TestDbpoolInit(t *testing.T) {
	DBPoolInit(&DBConfig{
		Username:     "gosun",
		Password:     "video",
		Host:         "192.168.33.183",
		Port:         5432,
		Dbname:       "gosun",
		MaxOpenConns: 100,
		MaxIdleConns: 20,
		DbType:       PostgreSQL,
		Timeout:      "10s",
	})
}
