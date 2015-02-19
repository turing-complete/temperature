// Package power provides algorithms for simulating the power dissipation of
// concurrent applications running on multiprocessor platforms.
package power

import (
	"errors"

	"github.com/ready-steady/simulation/system"
	"github.com/ready-steady/simulation/time"
)

// Power represents a power simulator configured for a particular system.
type Power struct {
	platform    *system.Platform
	application *system.Application
	Δt          float64
}

// New returns a power distributor for the given platform, application, and
// sampling interval.
func New(platform *system.Platform, application *system.Application, Δt float64) (*Power, error) {
	if Δt <= 0 {
		return nil, errors.New("the time step should be positive")
	}

	power := &Power{
		platform:    platform,
		application: application,
		Δt:          Δt,
	}

	return power, nil
}

// Compute calculates the power profile corresponding to the given schedule. The
// sc parameter controls the number of steps/samples that the output matrix will
// contain; schedules longer than this value (multiplied by the sampling
// interval Δt passed to New) get truncated.
func (p *Power) Compute(schedule *time.Schedule, sc uint) []float64 {
	cores, tasks := p.platform.Cores, p.application.Tasks
	cc, tc := uint(len(cores)), uint(len(tasks))
	Δt := p.Δt

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
