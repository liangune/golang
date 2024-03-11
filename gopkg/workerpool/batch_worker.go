package workerpool

import (
	"sync"
	"time"
)

type BatchWorker struct {
	workerId       int32 //协程ID
	taskChan       chan BatchTaskInterface
	workerChanPool chan chan BatchTaskInterface
	closer, closed chan none
	closeOnce      sync.Once
	handle         HandleInterface
	maxBatch       uint32
	waitTimeout    time.Duration //超时时间,单位毫秒
	batchTaskList  BatchTaskListInterface
}

func NewBatchWorker(workerChanPool chan chan BatchTaskInterface, Id int32, f1 NewHandleFun, f2 NewBatchTaskListFun) *BatchWorker {
	worker := BatchWorker{
		workerId:       Id,
		taskChan:       make(chan BatchTaskInterface),
		workerChanPool: workerChanPool,
		closer:         make(chan none),
		closed:         make(chan none),
	}
	if f1 != nil {
		worker.handle = f1()
		worker.handle.Init()
	}
	if f2 != nil {
		worker.batchTaskList = f2()
	}
	return &worker
}

func (w *BatchWorker) Start() {
	timeout := time.NewTicker(time.Millisecond * w.waitTimeout)
	go func() {
		w.workerChanPool <- w.taskChan
		for {
			select {
			case task := <-w.taskChan:
				task.SetWorkerId(w.workerId)
				w.batchTaskList.AddBatchTask(task)
				if uint32(w.batchTaskList.Size()) >= w.maxBatch {
					w.batchTaskList.ExecuteBatch(w)
					w.batchTaskList.Clear()
				}
				w.workerChanPool <- w.taskChan
			case <-timeout.C:
				if w.batchTaskList.Size() > 0 {
					w.batchTaskList.ExecuteBatch(w)
					w.batchTaskList.Clear()
				}
			case <-w.closer:
				close(w.closed)
				return
			}
		}
	}()
}

func (w *BatchWorker) Close() {
	w.closeOnce.Do(func() {
		close(w.closer)
		<-w.closed
		close(w.taskChan)
	})
}

func (w *BatchWorker) GetHandle() HandleInterface {
	return w.handle
}

func (w *BatchWorker) SetMaxBatch(maxBatch uint32) {
	w.maxBatch = maxBatch
}

func (w *BatchWorker) GetMaxBatch() uint32 {
	return w.maxBatch
}

func (w *BatchWorker) SetWaitTimeout(waitTimeout time.Duration) {
	w.waitTimeout = waitTimeout
}

func (w *BatchWorker) GetWaitTimeout() time.Duration {
	return w.waitTimeout
}

func (w *BatchWorker) SetBatchTaskList(batchTaskList BatchTaskListInterface) {
	w.batchTaskList = batchTaskList
}
