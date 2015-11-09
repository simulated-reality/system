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
	ID       uint
	Type     uint
	Parents  []uint
	Children []uint
}

// Len returns the number of tasks.
func (self *Application) Len() int {
	return len(self.Tasks)
}

// Roots returns the IDs of the tasks without parents.
func (self *Application) Roots() []uint {
	size := uint(len(self.Tasks))
	roots := make([]uint, 0, 1)

	for i := uint(0); i < size; i++ {
		if len(self.Tasks[i].Parents) == 0 {
			roots = append(roots, i)
		}
	}

	return roots
}

// Leafs returns the IDs of the tasks without children.
func (self *Application) Leafs() []uint {
	size := uint(len(self.Tasks))
	leafs := make([]uint, 0, size/2+1)

	for i := uint(0); i < size; i++ {
		if len(self.Tasks[i].Children) == 0 {
			leafs = append(leafs, i)
		}
	}

	return leafs
}
