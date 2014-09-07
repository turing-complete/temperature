package time

// Schedule represents a schedule of an application on a platform.
type Schedule struct {
	Mapping []uint16
	Start   []float64
	Finish  []float64
}

// Span returns the completion time of the last task in the schedule.
func (s *Schedule) Span() float64 {
	span := float64(0)

	for _, finish := range s.Finish {
		if span < finish {
			span = finish
		}
	}

	return span
}
