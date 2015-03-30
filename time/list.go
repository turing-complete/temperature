package time

import (
	"github.com/ready-steady/simulation/system"
)

// List represents a list scheduler.
type List struct {
	platform    *system.Platform
	application *system.Application
	roots       []uint
}

// NewList creates a new list scheduler for the given platform and application.
func NewList(platform *system.Platform, application *system.Application) *List {
	return &List{
		platform:    platform,
		application: application,
		roots:       application.Roots(),
	}
}

// Compute constructs a schedule according to the given priority vector. The
// length of this vector equals to the number of tasks in the system, and
// smaller values correspond to higher priorities.
func (l *List) Compute(priority []float64) *Schedule {
	cores := l.platform.Cores
	tasks := l.application.Tasks

	nc := uint(len(cores))
	nt := uint(len(tasks))

	mapping := make([]uint, nt)
	order := make([]uint, nt)
	start := make([]float64, nt)
	finish := make([]float64, nt)

	scheduled := make([]bool, nt)

	ctime := make([]float64, nc)
	ttime := make([]float64, nt)

	var i, j, k, cid, tid, kid, pid uint
	var span float64
	var ready bool

	size := uint(len(l.roots))

	// According to the benchmarks, keeping it sorted is not worth it.
	pool := make([]uint, size, nt)
	copy(pool, l.roots)

	for size > 0 {
		// Find the earliest available core.
		cid = 0
		for i = 1; i < nc; i++ {
			if ctime[i] < ctime[cid] {
				cid = i
			}
		}

		// Find the highest priority task.
		j, tid = 0, pool[0]
		for i = 1; i < size; i++ {
			if priority[pool[i]] < priority[tid] {
				j, tid = i, pool[i]
			}
		}

		// Remove the task from the pool.
		copy(pool[j:], pool[j+1:])
		pool = pool[:size-1]

		mapping[tid] = cid
		order[k] = tid
		k++
		if ctime[cid] > ttime[tid] {
			start[tid] = ctime[cid]
		} else {
			start[tid] = ttime[tid]
		}
		finish[tid] = start[tid] + cores[cid].Time[tasks[tid].Type]

		scheduled[tid] = true

		// Update the time when the core is again available.
		ctime[cid] = finish[tid]

		if span < finish[tid] {
			span = finish[tid]
		}

		for _, kid = range tasks[tid].Children {
			// Update the time when the child can potentially start executing.
			if ttime[kid] < finish[tid] {
				ttime[kid] = finish[tid]
			}

			// Push the child into the pool if it has become ready for
			// scheduling, that is, if all its parents have been scheduled.
			ready = true

			for _, pid = range tasks[kid].Parents {
				if !scheduled[pid] {
					ready = false
					break
				}
			}

			if !ready {
				continue
			}

			pool = append(pool, kid)
		}

		size = uint(len(pool))
	}

	return &Schedule{
		Mapping: mapping,
		Order:   order,
		Start:   start,
		Finish:  finish,
		Span:    span,
	}
}

// Delay constructs a new schedule based on an old one by adding delays to the
// execution times of the tasks.
func (l *List) Delay(schedule *Schedule, delay []float64) *Schedule {
	cores := l.platform.Cores
	tasks := l.application.Tasks

	nt := uint(len(tasks))

	duration := make([]float64, nt)

	for i := uint(0); i < nt; i++ {
		tid := schedule.Order[i]
		cid := schedule.Mapping[tid]
		duration[tid] = cores[cid].Time[tasks[tid].Type] + delay[tid]
	}

	return l.Update(schedule, duration)
}

// Update constructs a new schedule based on an old one by setting the execution
// times of the tasks to new values.
func (l *List) Update(schedule *Schedule, duration []float64) *Schedule {
	tasks := l.application.Tasks

	nc := uint(len(l.platform.Cores))
	nt := uint(len(tasks))

	start := make([]float64, nt)
	finish := make([]float64, nt)

	ctime := make([]float64, nc)
	ttime := make([]float64, nt)

	span := 0.0

	for i := uint(0); i < nt; i++ {
		tid := schedule.Order[i]
		cid := schedule.Mapping[tid]

		if ctime[cid] > ttime[tid] {
			start[tid] = ctime[cid]
		} else {
			start[tid] = ttime[tid]
		}
		finish[tid] = start[tid] + duration[tid]

		ctime[cid] = finish[tid]
		if span < finish[tid] {
			span = finish[tid]
		}

		for _, kid := range tasks[tid].Children {
			if ttime[kid] < finish[tid] {
				ttime[kid] = finish[tid]
			}
		}
	}

	return &Schedule{
		// FIXME: Do not be greedy! Make a copy!
		Mapping: schedule.Mapping,
		Order:   schedule.Order,
		Start:   start,
		Finish:  finish,
		Span:    span,
	}
}
