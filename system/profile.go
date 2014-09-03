package system

import (
	"math"
)

type Profile struct {
	ASAP     []float64 // As Soon As Possible
	ALAP     []float64 // As Late As Possible
	Mobility []float64

	time []float64
}

func CollectProfile(plat *Platform, app *Application) *Profile {
	cc := len(plat.Cores)
	tc := len(app.Tasks)

	p := &Profile{
		ASAP:     make([]float64, tc),
		ALAP:     make([]float64, tc),
		Mobility: make([]float64, tc),

		time: make([]float64, tc),
	}

	// Initialize ASAP and ALAP and compute the average execution time.
	for i := 0; i < tc; i++ {
		if i == 0 {
			p.ASAP[i] = math.Inf(-1)
			p.ALAP[i] = math.Inf(1)
		} else {
			p.ASAP[i] = p.ASAP[0]
			p.ALAP[i] = p.ALAP[0]
		}

		for j := 0; j < cc; j++ {
			p.time[i] += plat.Cores[j].Time[app.Tasks[i].Type]
		}
		p.time[i] /= float64(cc)
	}

	// Compute ASAP starting from the roots.
	for _, id := range app.Roots() {
		p.propagateASAP(app, id, 0)
	}

	leafs := app.Leafs()
	totalASAP := float64(0)
	for _, id := range leafs {
		if end := p.ASAP[id] + p.time[id]; end > totalASAP {
			totalASAP = end
		}
	}

	// Compute ASAP starting from the leafs.
	for _, id := range leafs {
		p.propagateALAP(app, id, totalASAP)
	}

	return p
}

func (p *Profile) propagateASAP(app *Application, id uint32, time float64) {
	if p.ASAP[id] >= time {
		return
	}

	p.ASAP[id] = time
	time += p.time[id]

	for _, id := range app.Tasks[id].Children {
		p.propagateASAP(app, id, time)
	}
}

func (p *Profile) propagateALAP(app *Application, id uint32, time float64) {
	time = time - p.time[id]
	if time < 0 {
		time = 0
	}

	if time >= p.ALAP[id] {
		return
	}

	p.ALAP[id] = time

	p.Mobility[id] = time - p.ASAP[id]
	if p.Mobility[id] < 0 {
		p.Mobility[id] = 0
	}

	for _, id := range app.Tasks[id].Parents {
		p.propagateALAP(app, id, time)
	}
}
