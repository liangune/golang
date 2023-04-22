package workerpool

import (
	"sync"
)

type WorkerInterface interface {
	Start()
	Close()
	GetHandle() HandleInterface
}

type Worker struct {
	workerId       int32 //协程ID
	taskChan       chan Task
	workerChanPool chan chan Task
	closer, closed chan none
	closeOnce      sync.Once
	handle         HandleInterface
}

func NewWorker(workerChanPool chan chan Task, Id int32, f NewHandleFun) *Worker {
	worker := Worker{
		workerId:       Id,
		taskChan:       make(chan Task),
		workerChanPool: workerChanPool,
		closer:         make(chan none),
		closed:         make(chan none),
	}
	if f != nil {
		worker.handle = f()
		worker.handle.Init()
	}
	return &worker
}

func (w *Worker) Start() {
	go func() {
		for {
			w.workerChanPool <- w.taskChan
			select {
			case task := <-w.taskChan:
				task.Execute(w)
			case <-w.closer:
				close(w.closed)
				return
			}
		}
	}()
}

func (w *Worker) Close() {
	w.closeOnce.Do(func() {
		close(w.closer)
		<-w.closed
		close(w.taskChan)
	})
}

func (w *Worker) GetHandle() HandleInterface {
	return w.handle
}
