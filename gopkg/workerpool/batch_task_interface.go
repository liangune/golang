package workerpool

type BatchTaskInterface interface {
	SetWorkerId(id int32)
	GetWorkerId() int32
	GetStatus() TaskStatus
	SetStatus(nStatus TaskStatus)
	ToJson() map[string]interface{}
}

type BatchTaskBase struct {
	WorkerId int32      `json:"omitempty"`
	Status   TaskStatus `json:"omitempty"`
}

func (t *BatchTaskBase) SetWorkerId(id int32) {
	t.WorkerId = id
}
func (t *BatchTaskBase) GetWorkerId() int32 {
	return t.WorkerId
}

func (t *BatchTaskBase) GetStatus() TaskStatus {
	return t.Status
}

func (t *BatchTaskBase) SetStatus(status TaskStatus) {
	t.Status = status
}

func (t *BatchTaskBase) ToJson() map[string]interface{} {
	return nil
}
