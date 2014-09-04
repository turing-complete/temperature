package system

// Application represents an application as a collection of dependent tasks.
// The tasks are assumed to form a directed acyclic graph.
type Application struct {
	Tasks []Task
}

// Task represents a task of an application. A task can have a number of
// children, which are tasks that depend on the current one (they can only
// proceed when this task is done). Each task is also given a type (Type),
// which is used for looking up the execution time and power consumption of the
// task when it is being executed of a core (see the definition of Core).
type Task struct {
	ID       uint16
	Type     uint16
	Parents  []uint16
	Children []uint16
}

// Roots returns the IDs of the tasks without parents.
func (a *Application) Roots() []uint16 {
	size := len(a.Tasks)
	roots := make([]uint16, 0, 1)

	for i := 0; i < size; i++ {
		if len(a.Tasks[i].Parents) == 0 {
			roots = append(roots, uint16(i))
		}
	}

	return roots
}

// Leafs returns the IDs of the tasks without children.
func (a *Application) Leafs() []uint16 {
	size := len(a.Tasks)
	leafs := make([]uint16, 0, size/2+1)

	for i := 0; i < size; i++ {
		if len(a.Tasks[i].Children) == 0 {
			leafs = append(leafs, uint16(i))
		}
	}

	return leafs
}
