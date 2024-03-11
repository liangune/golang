package workerpool

import (
	"sync"
	"time"
)

const (
	// 默认批量处理等待时间, 单位毫秒
	DefaultBatchWaitTimeout = 1000 * time.Millisecond
	// 默认批量值
	DefaultMaxBatch = 20
)

type BatchWorkerMgr struct {
	name            string
	batchDispatcher *BatchDispatcher
	closeOnce       sync.Once
}

func NewBatchWorkerMgr(name string, maxWorkers int, maxTaskCount int, f1 NewHandleFun, f2 NewBatchTaskListFun) *BatchWorkerMgr {
	p := &BatchWorkerMgr{
		name:            name,
		batchDispatcher: NewBatchDispatcher(name, maxWorkers, maxTaskCount, f1, f2),
	}

	return p
}

func (p *BatchWorkerMgr) Start() {
	p.batchDispatcher.Run()
}

func (p *BatchWorkerMgr) AddTask(task BatchTaskInterface) {
	p.batchDispatcher.AddTask(task)
}

func (p *BatchWorkerMgr) Submit(task BatchTaskInterface) {
	p.batchDispatcher.Submit(task)
}

func (p *BatchWorkerMgr) Close() {
	p.closeOnce.Do(func() {
		p.batchDispatcher.Close()
	})
}

func (p *BatchWorkerMgr) SetNewHandleFun(f NewHandleFun) {
	p.batchDispatcher.SetNewHandleFun(f)
}

func (p *BatchWorkerMgr) SetMaxBatch(maxBatch uint32) {
	p.batchDispatcher.SetMaxBatch(maxBatch)
}

func (p *BatchWorkerMgr) GetMaxBatch() uint32 {
	return p.batchDispatcher.GetMaxBatch()
}

func (p *BatchWorkerMgr) SetWaitTimeout(waitTimeout time.Duration) {
	p.batchDispatcher.SetWaitTimeout(waitTimeout)
}

func (p *BatchWorkerMgr) GetWaitTimeout() time.Duration {
	return p.batchDispatcher.GetWaitTimeout()
}

func (p *BatchWorkerMgr) SetBatchTaskListFunc(f NewBatchTaskListFun) {
	p.batchDispatcher.SetBatchTaskListFunc(f)
}

func (p *BatchWorkerMgr) GetTaskCount() int {
	return len(p.batchDispatcher.taskChan)
}
