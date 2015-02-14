package system

import (
	"math"
)

// Profile captures various statistics about the tasks of an application
// running on a platform.
type Profile struct {
	ASAP     []float64 // As Soon As Possible, the earliest start time
	ALAP     []float64 // As Late As Possible, the latest start time
	Mobility []float64 // max(0, ALAP - ASAP)

	time []float64
}

// NewProfile collects a profile of the given system. Since the mapping of
// the tasks onto the cores is assumed to be unknown at this stage, the profile
// is based on the average execution time of the tasks across all the cores.
func NewProfile(platform *Platform, application *Application) *Profile {
	cc := len(platform.Cores)
	tc := len(application.Tasks)

	p := &Profile{
		ASAP:     make([]float64, tc),
		ALAP:     make([]float64, tc),
		Mobility: make([]float64, tc),

		time: make([]float64, tc),
	}

	for i := 0; i < tc; i++ {
		if i == 0 {
			p.ASAP[i] = math.Inf(-1)
			p.ALAP[i] = math.Inf(1)
		} else {
			p.ASAP[i] = p.ASAP[0]
			p.ALAP[i] = p.ALAP[0]
		}

		for j := 0; j < cc; j++ {
			p.time[i] += platform.Cores[j].Time[application.Tasks[i].Type]
		}
		p.time[i] /= float64(cc)
	}

	// Compute ASAP starting from the roots.
	for _, i := range application.Roots() {
		p.propagateASAP(application, i, 0)
	}

	leafs := application.Leafs()

	totalASAP := float64(0)
	for _, i := range leafs {
		if end := p.ASAP[i] + p.time[i]; end > totalASAP {
			totalASAP = end
		}
	}

	// Compute ASAP starting from the leafs.
	for _, i := range leafs {
		p.propagateALAP(application, i, totalASAP)
	}

	return p
}

func (p *Profile) propagateASAP(application *Application, i uint, time float64) {
	if p.ASAP[i] >= time {
		return
	}

	p.ASAP[i] = time
	time += p.time[i]

	for _, i = range application.Tasks[i].Children {
		p.propagateASAP(application, i, time)
	}
}

func (p *Profile) propagateALAP(application *Application, i uint, time float64) {
	if time > p.time[i] {
		time = time - p.time[i]
	} else {
		time = 0
	}

	if time >= p.ALAP[i] {
		return
	}

	p.ALAP[i] = time

	if time > p.ASAP[i] {
		p.Mobility[i] = time - p.ASAP[i]
	} else {
		p.Mobility[i] = 0
	}

	for _, i = range application.Tasks[i].Parents {
		p.propagateALAP(application, i, time)
	}
}
