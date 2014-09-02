package system

type App struct {
	Tasks []*Task
}

type Task struct {
	Type uint32
}

func NewApp(graph Graph) (*App, error) {
	app := &App{
	}

	return app, nil
}
