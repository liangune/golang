package workerpool

import (
	"sync"
)

type Dispatcher struct {
	dispatcherId string
	maxWorkers   int //协程最大数量
	workers      []*Worker
	closer       chan none
	closed       chan none
	//endSignal      chan os.Signal
	taskChan       chan TaskInterface
	workerChanPool chan chan TaskInterface
	closeOnce      sync.Once
	newHandleFun   NewHandleFun
}

func NewDispatcher(dispatcherId string, maxWorkers, maxTaskCount int, f NewHandleFun) *Dispatcher {
	//endSignal := make(chan os.Signal)
	//signal.Notify(endSignal, syscall.SIGINT, syscall.SIGTERM)
	return &Dispatcher{
		dispatcherId: dispatcherId,
		maxWorkers:   maxWorkers,
		closer:       make(chan none),
		closed:       make(chan none),
		//endSignal:      endSignal,
		taskChan:       make(chan TaskInterface, maxTaskCount),
		workerChanPool: make(chan chan TaskInterface, maxWorkers),
		newHandleFun:   f,
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.workerChanPool, int32(i), d.newHandleFun)
		d.workers = append(d.workers, worker)
		worker.Start()
	}
	go d.Dispatch()
}

func (d *Dispatcher) Dispatch() {
	defer close(d.closed)

	for {
		select {
		//case endSignal := <-d.endSignal:
		//	fmt.Printf("收到[%v]信号, %s协程池关闭... \n", endSignal, d.dispatcherId)
		//	return
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

func (d *Dispatcher) AddTask(task TaskInterface) {
	d.taskChan <- task
}

func (d *Dispatcher) Submit(task TaskInterface) {
	d.taskChan <- task
}

func (d *Dispatcher) Close() {
	d.closeOnce.Do(func() {
		close(d.closer)
		<-d.closed
		close(d.taskChan)
		close(d.workerChanPool)
	})
}

func (d *Dispatcher) SetNewHandleFun(f NewHandleFun) {
	d.newHandleFun = f
}
