package workerpool

type HandleInterface interface {
	Init() error
}

type NewHandleFun func() HandleInterface
