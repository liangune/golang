package workerpool

import (
	"sync"
	"time"
)

type BatchDispatcher struct {
	dispatcherId        string
	maxWorkers          int //协程最大数量
	workers             []*BatchWorker
	closer              chan none
	closed              chan none
	taskChan            chan BatchTaskInterface
	workerChanPool      chan chan BatchTaskInterface
	closeOnce           sync.Once
	newHandleFun        NewHandleFun
	maxBatch            uint32
	waitTimeout         time.Duration //超时时间,单位毫秒
	newBatchTaskListFun NewBatchTaskListFun
}

func NewBatchDispatcher(dispatcherId string, maxWorkers, maxTaskCount int, f1 NewHandleFun, f2 NewBatchTaskListFun) *BatchDispatcher {
	return &BatchDispatcher{
		dispatcherId:        dispatcherId,
		maxWorkers:          maxWorkers,
		closer:              make(chan none),
		closed:              make(chan none),
		taskChan:            make(chan BatchTaskInterface, maxTaskCount),
		workerChanPool:      make(chan chan BatchTaskInterface, maxWorkers),
		newHandleFun:        f1,
		maxBatch:            DefaultMaxBatch,
		waitTimeout:         DefaultBatchWaitTimeout,
		newBatchTaskListFun: f2,
	}
}

func (d *BatchDispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewBatchWorker(d.workerChanPool, int32(i), d.newHandleFun, d.newBatchTaskListFun)
		worker.SetMaxBatch(d.maxBatch)
		worker.SetWaitTimeout(d.waitTimeout)

		d.workers = append(d.workers, worker)
		worker.Start()
	}
	go d.Dispatch()
}

func (d *BatchDispatcher) Dispatch() {
	defer close(d.closed)

	for {
		select {
		case task, ok := <-d.taskChan:
			if true == ok {
				WorkerTaskChan := <-d.workerChanPool
				WorkerTaskChan <- task
			}
		case <-d.closer:
			for _, w := range d.workers {
				w.Close()
			}
			return
		}
	}
}

func (d *BatchDispatcher) AddTask(task BatchTaskInterface) {
	d.taskChan <- task
}

func (d *BatchDispatcher) Submit(task BatchTaskInterface) {
	d.taskChan <- task
}

func (d *BatchDispatcher) Close() {
	d.closeOnce.Do(func() {
		close(d.closer)
		<-d.closed
		close(d.taskChan)
		close(d.workerChanPool)
	})
}

func (d *BatchDispatcher) SetNewHandleFun(f NewHandleFun) {
	d.newHandleFun = f
}

func (d *BatchDispatcher) SetMaxBatch(maxBatch uint32) {
	d.maxBatch = maxBatch
}

func (d *BatchDispatcher) GetMaxBatch() uint32 {
	return d.maxBatch
}

func (d *BatchDispatcher) SetWaitTimeout(waitTimeout time.Duration) {
	d.waitTimeout = waitTimeout
}

func (d *BatchDispatcher) GetWaitTimeout() time.Duration {
	return d.waitTimeout
}

func (d *BatchDispatcher) SetBatchTaskListFunc(f NewBatchTaskListFun) {
	d.newBatchTaskListFun = f
}
