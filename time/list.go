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

// Compute constructs a schedule according to the given priority vector.
// The length of this vector equals to the number of tasks in the system, and
// smaller values correspond to higher priorities.
func (l *List) Compute(priority []float64) *Schedule {
	cores := l.platform.Cores
	tasks := l.application.Tasks

	cc := uint(len(cores))
	tc := uint(len(tasks))

	mapping := make([]uint, tc)
	order := make([]uint, tc)
	start := make([]float64, tc)
	finish := make([]float64, tc)

	scheduled := make([]bool, tc)

	ctime := make([]float64, cc)
	ttime := make([]float64, tc)

	var i, j, k, cid, tid, kid, pid uint
	var span float64
	var ready bool

	size := uint(len(l.roots))

	// According to the benchmarks, keeping it sorted is not worth it.
	pool := make([]uint, size, tc)
	copy(pool, l.roots)

	for size > 0 {
		// Find the earliest available core.
		cid = 0
		for i = 1; i < cc; i++ {
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

// Recompute constructs a new schedule based on an old one by adding a delay to
// the execution time of the tasks.
func (l *List) Recompute(schedule *Schedule, delay []float64) *Schedule {
	cores := l.platform.Cores
	tasks := l.application.Tasks

	cc := uint(len(l.platform.Cores))
	tc := uint(len(tasks))

	start := make([]float64, tc)
	finish := make([]float64, tc)

	ctime := make([]float64, cc)
	ttime := make([]float64, tc)

	var i, cid, tid, kid uint
	var span float64

	for ; i < tc; i++ {
		tid = schedule.Order[i]
		cid = schedule.Mapping[tid]

		if ctime[cid] > ttime[tid] {
			start[tid] = ctime[cid]
		} else {
			start[tid] = ttime[tid]
		}
		finish[tid] = start[tid] + cores[cid].Time[tasks[tid].Type] + delay[tid]

		ctime[cid] = finish[tid]
		if span < finish[tid] {
			span = finish[tid]
		}

		for _, kid = range tasks[tid].Children {
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
