// Package system provides an abstraction of an electronic system for the
// purpose of simulation. Such a system is composed of a platform and an
// application. A platform is a collection of processing elements, referred to
// as cores, and an application is a collection of data-dependent tasks,
// forming a directed acyclic graph.
package system

// Load constructs a platform and an application based on the specification
// given in a file. The only supported format is TGFF.
func Load(path string) (*Platform, *Application, error) {
	return loadTGFF(path)
}
