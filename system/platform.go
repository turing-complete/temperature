package system

// Platform represents a platform composed of a number of processing elements,
// referred to as cores, which is capable of running an application.
type Platform struct {
	Cores []Core
}

// Core represents a processing element of a platform. Each core is
// characterized by two vectors: execution time (Time) and power consumption
// (Power). Each entry in these vectors corresponds to a task Type, not to a
// task ID (see the ID and Type fields of the Task struct).
type Core struct {
	ID    uint32
	Time  []float64
	Power []float64
}
