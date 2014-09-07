package time

// Schedule represents a schedule of an application on a platform.
type Schedule struct {
	Mapping []uint16
	Start   []float64
	Finish  []float64
}
