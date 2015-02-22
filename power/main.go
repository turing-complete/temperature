// Package power provides algorithms for simulating the power dissipation of
// concurrent applications running on multiprocessor platforms.
package power

import (
	"github.com/ready-steady/simulation/system"
	"github.com/ready-steady/simulation/time"
)

// Power represents a power simulator configured for a particular system.
type Power struct {
	platform    *system.Platform
	application *system.Application
}

// New returns a power distributor for the given platform and application.
func New(platform *system.Platform, application *system.Application) *Power {
	return &Power{
		platform:    platform,
		application: application,
	}
}

// Compute calculates the power profile of a schedule. The sampling interval is
// specified by the Δt parameter, and the sc parameter controls the number of
// samples that the output matrix will contain; long schedules are truncated.
func (p *Power) Compute(schedule *time.Schedule, Δt float64, sc uint) []float64 {
	cores, tasks := p.platform.Cores, p.application.Tasks
	cc, tc := uint(len(cores)), uint(len(tasks))

	P := make([]float64, cc*sc)

	if count := uint(schedule.Span / Δt); count < sc {
		sc = count
	}

	for i := uint(0); i < tc; i++ {
		j := schedule.Mapping[i]
		p := cores[j].Power[tasks[i].Type]

		s := uint(schedule.Start[i] / Δt)
		f := uint(schedule.Finish[i] / Δt)
		if f > sc {
			f = sc
		}

		for ; s < f; s++ {
			P[s*cc+j] = p
		}
	}

	return P
}

// Process takes a schedule and returns a function func(time float64, power
// []float64) that computes the power dissipation at an arbitrary time moment
// according to the schedule.
func (p *Power) Process(schedule *time.Schedule) func(float64, []float64) {
	cores, tasks := p.platform.Cores, p.application.Tasks
	cc, tc := uint(len(cores)), uint(len(tasks))

	mapping := make([][]uint, cc)
	for i := uint(0); i < cc; i++ {
		mapping[i] = make([]uint, 0, tc)
		for j := uint(0); j < tc; j++ {
			if i == schedule.Mapping[j] {
				mapping[i] = append(mapping[i], j)
			}
		}
	}

	start, finish := schedule.Start, schedule.Finish

	return func(time float64, power []float64) {
		for i := uint(0); i < cc; i++ {
			power[i] = 0
			for _, j := range mapping[i] {
				if start[j] <= time && time <= finish[j] {
					power[i] = cores[i].Power[tasks[j].Type]
					break
				}
			}
		}
	}
}
