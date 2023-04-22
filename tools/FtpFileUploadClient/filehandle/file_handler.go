package filehandle

import "go/gopkg/workerpool"

const DefaultMaxTaskCount = 10000

var GFileHandler *FileHandler

type FileHandler struct {
	Pool *workerpool.WorkerPool
}

func NewFileHandler(maxWorkers int, maxTaskCount int) *FileHandler {
	h := FileHandler{
		Pool: workerpool.NewWorkerPool("file", maxWorkers, maxTaskCount),
	}

	return &h
}

func GetFileHandler() *FileHandler {
	return GFileHandler
}

func (h *FileHandler) Start() {
	h.Pool.Start()
}

func (h *FileHandler) Dispatch(path string) {
	task := FileTask{
		path: path,
	}
	h.Pool.AddTask(&task)
}

func (h *FileHandler) SetNewHandlerFunc(f workerpool.NewHandleFun) {
	h.Pool.SetNewHandleFun(f)
}
