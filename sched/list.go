package sched

import (
	"github.com/goesd/persim/system"
)

// List represents a list scheduler.
type List struct {
	plat  *system.Platform
	app   *system.Application
	roots []uint16
}

// NewList creates a new list scheduler for the given platform and application.
func NewList(plat *system.Platform, app *system.Application) *List {
	return &List{
		plat:  plat,
		app:   app,
		roots: app.Roots(),
	}
}

// Schedule constructs a schedule according to the given priority vector.
// The length of this vector equals to the number of tasks in the system, and
// smaller values correspond to higher priorities.
func (l *List) Schedule(priority []float64) *Schedule {
	cores := l.plat.Cores
	tasks := l.app.Tasks

	cc := uint16(len(cores))
	tc := uint16(len(tasks))

	mapping := make([]uint16, tc)
	start := make([]float64, tc)
	duration := make([]float64, tc)

	pushed := make([]bool, tc)
	scheduled := make([]bool, tc)

	ctime := make([]float64, cc)
	ttime := make([]float64, tc)

	var i, tid, cid uint16
	var finish float64
	var ready bool

	// Always kept sorted according to the priority.
	pool := newPool(uint16(tc))
	for _, tid = range l.roots {
		pool.push(tid, priority[tid])
		pushed[tid] = true
	}

	for len(pool) > 0 {
		// Pull the task with the highest priority.
		tid = pool.pop()

		// Find the earliest available core for the task.
		for cid, i = 0, 1; i < cc; i++ {
			if ctime[i] < ctime[cid] {
				cid = i
			}
		}

		mapping[tid] = cid

		if ctime[cid] > ttime[tid] {
			start[tid] = ctime[cid]
		} else {
			start[tid] = ttime[tid]
		}

		duration[tid] = cores[cid].Time[tasks[tid].Type]

		scheduled[tid] = true

		// Update the time when the core is again available.
		finish = start[tid] + duration[tid]
		ctime[cid] = finish

		for _, cid = range tasks[tid].Children {
			// Update the time when the child can potentially start executing.
			if ttime[cid] < finish {
				ttime[cid] = finish
			}

			if pushed[cid] {
				continue
			}

			// Push the child into the pool if it has become ready for
			// scheduling, that is, if all its parents have been scheduled.
			ready = true

			for _, tid = range tasks[cid].Parents {
				if !scheduled[tid] {
					ready = false
					break
				}
			}

			if !ready {
				continue
			}

			pool.push(cid, priority[cid])
			pushed[cid] = true
		}
	}

	return &Schedule{
		Mapping:  mapping,
		Start:    start,
		Duration: duration,
	}
}

// Reschedule constructs a schedule based on another schedule.
func (l *List) Reschedule(s *Schedule) *Schedule {
	return s
}
