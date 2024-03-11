package asyncdb

import (
	"go/gopkg/logger/vglog"
	"go/gopkg/utils/queue"
	"gorm.io/gorm"
	"time"
)

type AsyncWriteDB struct {
	syncQueue *queue.SyncQueue
	dbpool    *gorm.DB
}

const (
	DefaultInterval = time.Millisecond * 1000
)

func NewAsyncWriteDB(dbPool *gorm.DB) *AsyncWriteDB {
	mgr := &AsyncWriteDB{
		syncQueue: queue.NewSyncQueue(),
		dbpool:    dbPool,
	}

	go mgr.Run()

	return mgr
}

func (self *AsyncWriteDB) Run() {
	self.OnWriteSqlEvent()
}

func (self *AsyncWriteDB) OnWriteSqlEvent() {
	ticker := time.NewTicker(DefaultInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			self.WirteSqlToDb()
		}
	}
}

func (self *AsyncWriteDB) AddSqlToDb(sql string) {
	var in interface{}
	in = sql
	self.syncQueue.Push(in)
}

func (self *AsyncWriteDB) WirteSqlToDb() {
	if self.syncQueue.Len() > 0 {
		for {
			if self.syncQueue.Len() <= 0 {
				break
			}
			sql, ok := self.syncQueue.TryPop()
			if ok && sql != nil {
				strsql := sql.(string)

				err := self.dbpool.Exec(strsql).Error
				if err != nil {
					vglog.Error("exec sql fail :%v, sql: %s", err, strsql)
				}
			}
		}
	}
}
