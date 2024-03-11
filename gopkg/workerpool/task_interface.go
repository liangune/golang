package workerpool

type TaskStatus int8

const (
	TaskStatusOK     TaskStatus = 0
	TaskStatusFailed TaskStatus = 1
)

type TaskInterface interface {
	Execute(w WorkerInterface) error
	SetWorkerId(id int32)
	GetWorkerId() int32
	GetStatus() TaskStatus
	SetStatus(nStatus TaskStatus)
}

type TaskBase struct {
	WorkerId int32      `json:"omitempty"`
	status   TaskStatus `json:"omitempty"`
}

func (t *TaskBase) Execute(w WorkerInterface) error {
	return nil
}

func (t *TaskBase) SetWorkerId(id int32) {
	t.WorkerId = id
}
func (t *TaskBase) GetWorkerId() int32 {
	return t.WorkerId
}

func (t *TaskBase) GetStatus() TaskStatus {
	return t.status
}

func (t *TaskBase) SetStatus(status TaskStatus) {
	t.status = status
}
