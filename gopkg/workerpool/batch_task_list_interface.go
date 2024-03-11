package workerpool

type BatchTaskListInterface interface {
	ExecuteBatch(w WorkerInterface) error
	SetWorkerId(id int32)
	GetWorkerId() int32
	AddBatchTask(task BatchTaskInterface)
	Size() int
	Clear()
	ToBytes() ([]byte, error)
}

type NewBatchTaskListFun func() BatchTaskListInterface

type BatchTaskListBase struct {
	WorkerId       int32
	Status         TaskStatus
	BatchTaskSlice []BatchTaskInterface
}

func (t *BatchTaskListBase) ExecuteBatch(w WorkerInterface) error {
	return nil
}

func (t *BatchTaskListBase) SetWorkerId(id int32) {
	t.WorkerId = id
}

func (t *BatchTaskListBase) GetWorkerId() int32 {
	return t.WorkerId
}

func (t *BatchTaskListBase) AddBatchTask(task BatchTaskInterface) {
	t.BatchTaskSlice = append(t.BatchTaskSlice, task)
}

func (t *BatchTaskListBase) Size() int {
	return len(t.BatchTaskSlice)
}

func (t *BatchTaskListBase) Clear() {
	t.BatchTaskSlice = t.BatchTaskSlice[0:0]
}

func (t *BatchTaskListBase) ToBytes() ([]byte, error) {
	return nil, nil
}
