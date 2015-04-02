// Package power provides algorithms for simulating the power dissipation of
// concurrent applications running on multiprocessor platforms.
package power

import (
	"github.com/ready-steady/simulation/system"
	"github.com/ready-steady/simulation/time"
	"github.com/ready-steady/sort"
)

// Power represents a power simulator configured for a particular system.
type Power struct {
	platform    *system.Platform
	application *system.Application
}

// New returns a power distributor for a platform and an application.
func New(platform *system.Platform, application *system.Application) *Power {
	return &Power{
		platform:    platform,
		application: application,
	}
}

// Partition computes a power profile with a variable time step dictated by the
// time moments of power switches (the start and finish times of the tasks) and
// a number of additional time moments gathered in points.
func (p *Power) Partition(schedule *time.Schedule, points []float64,
	ε float64) ([]float64, []float64, []uint) {

	cores, tasks := p.platform.Cores, p.application.Tasks
	nc, nt, np := uint(len(cores)), uint(len(tasks)), uint(len(points))

	time := make([]float64, 2*nt+np)
	copy(time[:nt], schedule.Start)
	copy(time[nt:], schedule.Finish)
	copy(time[2*nt:], points)

	ΔT, index := traverse(time, ε)
	sindex, findex, pindex := index[:nt], index[nt:2*nt], index[2*nt:]

	ns := uint(len(ΔT))

	P := make([]float64, nc*ns)

	for i := uint(0); i < nt; i++ {
		j := schedule.Mapping[i]
		p := cores[j].Power[tasks[i].Type]

		s, f := sindex[i], findex[i]

		for ; s < f; s++ {
			P[s*nc+j] = p
		}
	}

	return P, ΔT, pindex
}

// Sample computes a power profile with respect to a sampling interval Δt. The
// required number of samples is specified by ns; short schedules are extended
// while long ones are truncated.
func (p *Power) Sample(schedule *time.Schedule, Δt float64, ns uint) []float64 {
	cores, tasks := p.platform.Cores, p.application.Tasks
	nc, nt := uint(len(cores)), uint(len(tasks))

	P := make([]float64, nc*ns)

	if count := uint(schedule.Span / Δt); count < ns {
		ns = count
	}

	for i := uint(0); i < nt; i++ {
		j := schedule.Mapping[i]
		p := cores[j].Power[tasks[i].Type]

		s := uint(schedule.Start[i] / Δt)
		f := uint(schedule.Finish[i] / Δt)
		if f > ns {
			f = ns
		}

		for ; s < f; s++ {
			P[s*nc+j] = p
		}
	}

	return P
}

// Progress returns a function func(time float64, power []float64) that computes
// the power dissipation at an arbitrary time moment according to a schedule.
func (p *Power) Progress(schedule *time.Schedule) func(float64, []float64) {
	cores, tasks := p.platform.Cores, p.application.Tasks
	nc, nt := uint(len(cores)), uint(len(tasks))

	mapping := make([][]uint, nc)
	for i := uint(0); i < nc; i++ {
		mapping[i] = make([]uint, 0, nt)
		for j := uint(0); j < nt; j++ {
			if i == schedule.Mapping[j] {
				mapping[i] = append(mapping[i], j)
			}
		}
	}

	start, finish := schedule.Start, schedule.Finish

	return func(time float64, power []float64) {
		for i := uint(0); i < nc; i++ {
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

func traverse(points []float64, ε float64) ([]float64, []uint) {
	np := uint(len(points))
	order, _ := sort.Quick(points)

	Δ := make([]float64, np-1)
	index := make([]uint, np)

	j := uint(0)

	for i, x := uint(1), points[0]; i < np; i++ {
		if δ := points[i] - x; δ < ε {
			index[order[i]] = index[order[i-1]]
		} else {
			Δ[j] = δ
			j++
			x = points[i]
			index[order[i]] = j
		}
	}

	return Δ[:j], index
}
