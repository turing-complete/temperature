package system

type Application struct {
	Tasks []Task
}

type Task struct {
	ID       uint32
	Type     uint32
	Children []*Task
}
