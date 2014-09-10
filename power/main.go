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

// Compute constructs the power profile of the given schedule and stores it in
// a cc-by-sc matrix P where cc is the number of cores and sc is the maximal
// number of steps (samples) that the matrix can accommodate.
func (self *Self) Compute(sched *time.Schedule, P []float64, sc uint32) {
	cores, tasks := self.plat.Cores, self.app.Tasks
	dt := self.dt

	cc := uint32(len(cores))
	tc := uint16(len(tasks))
	if count := uint32(math.Floor(sched.Span() / dt)); count < sc {
		sc = count
	}

	var j, s, f uint32

	for i := uint16(0); i < tc; i++ {
		j = uint32(sched.Mapping[i])
		s = uint32(math.Floor(sched.Start[i] / dt))
		f = uint32(math.Floor(sched.Finish[i] / dt) - 1)
		if f >= sc {
			f = sc - 1
		}

		for ; s <= f; s++ {
			P[s*cc+j] = cores[j].Power[tasks[i].Type]
		}
	}
}
