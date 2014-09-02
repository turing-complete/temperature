package system

type App struct {
	Tasks []Task
}

type Task struct {
	Type     uint32
	Children []*Task
}
