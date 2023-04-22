package workerpool

type Task interface {
	Execute(w WorkerInterface) error
}
