// Package power provides algorithms for simulating the power dissipation of
// concurrent applications running on multiprocessor platforms.
package power

import (
	"math"

	"github.com/go-eslab/persim/system"
	"github.com/go-eslab/persim/time"
)

// Self represents a power distributer configured for a particular system.
type Self struct {
	plat *system.Platform
	app  *system.Application
	dt   float64
}

// New returns a power distributor for the given platform, application, and
// sampling period.
func New(plat *system.Platform, app *system.Application, dt float64) *Self {
	return &Self{
		plat: plat,
		app:  app,
		dt:   dt,
	}
}

// Compute returns the power profile corresponding to the given schedule.
func (self *Self) Compute(sched *time.Schedule) []float64 {
	cores, tasks := self.plat.Cores, self.app.Tasks

	span, dt := sched.Span(), self.dt

	cc := uint32(len(cores))
	tc := len(tasks)
	sc := uint32(math.Floor(span / dt))

	P := make([]float64, cc*sc)

	for i := 0; i < tc; i++ {
		j := uint32(sched.Mapping[i])
		s := uint32(math.Floor(sched.Start[i] / dt))
		f := uint32(math.Floor(sched.Finish[i] / dt) - 1)

		for ; s <= f; s++ {
			P[s*cc+j] = cores[j].Power[tasks[i].Type]
		}
	}

	return P
}
