// Package power provides algorithms for simulating the power dissipation of
// concurrent applications running on multiprocessor platforms.
package power

import (
	"github.com/ready-steady/persim/system"
	"github.com/ready-steady/persim/time"
)

// Power represents a power simulator configured for a particular system.
type Power struct {
	platform    *system.Platform
	application *system.Application
	Δt          float64
}

// New returns a power distributor for the given platform, application, and
// sampling period.
func New(platform *system.Platform, application *system.Application, Δt float64) *Power {
	return &Power{
		platform:    platform,
		application: application,
		Δt:          Δt,
	}
}

// Compute constructs the power profile of the given schedule and stores it in a
// cc-by-sc matrix P where cc is the number of cores and sc is the maximal
// number of steps (samples) that the matrix can accommodate. P is assumed to be
// zeroed.
func (p *Power) Compute(schedule *time.Schedule, P []float64, sc uint32) {
	cores, tasks := p.platform.Cores, p.application.Tasks
	Δt := p.Δt

	cc := uint32(len(cores))
	tc := uint16(len(tasks))
	if count := uint32(schedule.Span / Δt); count < sc {
		sc = count
	}

	for i := uint16(0); i < tc; i++ {
		j := uint32(schedule.Mapping[i])
		p := cores[j].Power[tasks[i].Type]

		s := uint32(schedule.Start[i] / Δt)
		f := uint32(schedule.Finish[i] / Δt)
		if f > sc {
			f = sc
		}

		for ; s < f; s++ {
			P[s*cc+j] = p
		}
	}
}
