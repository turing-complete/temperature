// Package provides an Runge–Kutta integrator for systems of ordinary
// differential equations modeling temperature of electronic systems.
package numeric

import (
	"github.com/ready-steady/hotspot"
	"github.com/ready-steady/numeric/integration/ode"
	"github.com/ready-steady/simulation/temperature"
)

// Temperature represents an integrator.
type Temperature struct {
	Cores uint
	Nodes uint

	system system

	integrator *ode.DormandPrince
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

	integrator, err := ode.NewDormandPrince(&ode.Config{
		MaximalStep:       0,
		InitialStep:       0,
		AbsoluteTolerance: 1e-3,
		RelativeTolerance: 1e-3,
	})
	if err != nil {
		return nil, err
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

		integrator: integrator,
	}

	return temperature, nil
}
