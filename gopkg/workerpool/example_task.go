package workerpool

import (
	"fmt"
)

type ExampleTask struct {
	Message string
}

func (t *ExampleTask) SetMsg(message string) {
	t.Message = message
}

func (t *ExampleTask) Execute(w WorkerInterface) error {
	fmt.Println(t.Message)
	//time.Sleep(time.Second * 1)
	return nil
}
