// Package power provides algorithms for simulating the power dissipation of
// concurrent applications running on multiprocessor platforms.
package power

import (
	"math"

	"github.com/go-eslab/persim/system"
	"github.com/go-eslab/persim/time"
	"github.com/go-math/linal/matrix"
)

// Power represents a power distributer configured for a particular system.
type Power struct {
	plat *system.Platform
	app  *system.Application
	dt   float64
}

// New returns a power distributor for the given platform, application, and
// sampling period.
func New(plat *system.Platform, app *system.Application, dt float64) *Power {
	return &Power{
		plat: plat,
		app:  app,
		dt:   dt,
	}
}

// Compute returns the power profile corresponding to the given schedule.
func (p *Power) Compute(sched *time.Schedule) *matrix.Matrix {
	cores, tasks := p.plat.Cores, p.app.Tasks

	span, dt := sched.Span(), p.dt

	cc := uint32(len(cores))
	tc := len(tasks)
	sc := uint32(math.Floor(span / dt))

	data := make([]float64, cc*sc)

	for i := 0; i < tc; i++ {
		j := uint32(sched.Mapping[i])
		s := uint32(math.Floor(sched.Start[i] / dt))
		f := uint32(math.Floor(sched.Finish[i] / dt) - 1)

		for ; s <= f; s++ {
			data[s*cc+j] = cores[j].Power[tasks[i].Type]
		}
	}

	return matrix.New(cc, sc, data)
}
