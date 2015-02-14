package time

// Schedule represents a schedule of an application on a platform.
type Schedule struct {
	Mapping []uint
	Order   []uint
	Start   []float64
	Finish  []float64
	Span    float64
}
