package numeric

import (
	"github.com/ready-steady/ode"
	"github.com/turing-complete/hotspot"
)

// Temperature is an integrator of a thermal system.
type Temperature struct {
	nc uint
	nn uint

	system     system
	integrator ode.Integrator
}

// New returns a new integrator.
func New(config *Config, integrator ode.Integrator) *Temperature {
	model := hotspot.New((*hotspot.Config)(&config.Config))
	nc, nn := model.Cores, model.Nodes

	// Reusing model.G to store A and model.C to store B.
	A := model.G
	B := model.C
	for i := uint(0); i < nn; i++ {
		B[i] = 1 / model.C[i]
	}
	for i := uint(0); i < nn; i++ {
		for j := uint(0); j < nn; j++ {
			A[j*nn+i] = -B[i] * A[j*nn+i]
		}
	}

	return &Temperature{
		nc: nc,
		nn: nn,

		system: system{
			A: A,
			B: B,

			Qamb: config.Ambience,
		},

		integrator: integrator,
	}
}
