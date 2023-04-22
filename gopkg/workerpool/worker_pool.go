package workerpool

import "sync"

type none struct{}

type WorkerPool struct {
	name       string
	Dispatcher *Dispatcher
	closeOnce  sync.Once
}

func NewWorkerPool(name string, maxWorkers int, maxTaskCount int) *WorkerPool {
	p := &WorkerPool{
		name:       name,
		Dispatcher: NewDispatcher(name, maxWorkers, maxTaskCount),
	}

	return p
}

func (p *WorkerPool) Start() {
	p.Dispatcher.Run()
}

func (p *WorkerPool) AddTask(task Task) {
	p.Dispatcher.AddTask(task)
}

func (p *WorkerPool) Close() {
	p.closeOnce.Do(func() {
		p.Dispatcher.Close()
	})
}

func (p *WorkerPool) SetNewHandleFun(f NewHandleFun) {
	p.Dispatcher.SetNewHandleFun(f)
}
