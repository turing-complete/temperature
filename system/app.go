package system

type App struct {
	Tasks []Task
}

type Task struct {
	ID       uint32
	Type     uint32
	Children []*Task
}
