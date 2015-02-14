// Package power provides algorithms for simulating the power dissipation of
// concurrent applications running on multiprocessor platforms.
package power

// #include <string.h>
import "C"

import (
	"errors"
	"unsafe"

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
// sampling period.
func New(platform *system.Platform, application *system.Application, Δt float64) (*Power, error) {
	if Δt <= 0 {
		return nil, errors.New("the time step is invalid")
	}

	power := &Power{
		platform:    platform,
		application: application,
		Δt:          Δt,
	}

	return power, nil
}

// Compute constructs the power profile of the given schedule and stores it in a
// cc-by-sc matrix P where cc is the number of cores and sc is the maximal
// number of steps (samples) that the matrix can accommodate.
func (p *Power) Compute(schedule *time.Schedule, P []float64, sc uint) {
	cores, tasks := p.platform.Cores, p.application.Tasks
	Δt := p.Δt

	cc := uint(len(cores))
	tc := uint(len(tasks))
	if count := uint(schedule.Span / Δt); count < sc {
		sc = count
	}

	// FIXME: Bad, bad, bad!
	C.memset(unsafe.Pointer(&P[0]), 0, C.size_t(8*cc*sc))

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
}
