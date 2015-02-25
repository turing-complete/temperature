// Package provides an Runge–Kutta integrator for systems of ordinary
// differential equations modeling temperature of electronic systems.
package numeric

import (
	"github.com/ready-steady/hotspot"
	"github.com/ready-steady/simulation/temperature"
)

// Temperature represents an integrator.
type Temperature struct {
	Cores uint
	Nodes uint

	system system
}

// New returns an integrator set up according to the given configuration.
func New(c *temperature.Config) (*Temperature, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	model := hotspot.New(&c.HotSpot)
	cc, nc := model.Cores, model.Nodes

	// Reusing model.G to store A and model.C to store B.
	A := model.G
	B := model.C
	for i := uint(0); i < nc; i++ {
		B[i] = 1 / model.C[i]
	}
	for i := uint(0); i < nc; i++ {
		for j := uint(0); j < nc; j++ {
			A[j*nc+i] = -B[i] * A[j*nc+i]
		}
	}

	temperature := &Temperature{
		Cores: cc,
		Nodes: nc,

		system: system{
			A: A,
			B: B,

			Δt:   c.TimeStep,
			Qamb: c.Ambience,
		},
	}

	return temperature, nil
}
