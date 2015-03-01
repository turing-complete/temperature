// Package temperature provides an exponential-integrator-based solver for
// systems of differential-algebraic equations modeling temperature of
// electronic systems.
package analytic

import (
	"errors"
	"math"

	"github.com/ready-steady/hotspot"
	"github.com/ready-steady/linear/decomposition"
	"github.com/ready-steady/linear/matrix"
)

// Temperature represents an integrator.
type Temperature struct {
	Cores uint
	Nodes uint

	system system
}

// New returns an integrator set up according to the given configuration.
func New(c *Config) (*Temperature, error) {
	if c.TimeStep <= 0 {
		return nil, errors.New("the time step should be positive")
	}

	model := hotspot.New((*hotspot.Config)(&c.Config))
	nc, nn := model.Cores, model.Nodes

	// Reusing model.G to store A and model.C to store D.
	A := model.G
	D := model.C
	for i := uint(0); i < nn; i++ {
		D[i] = math.Sqrt(1 / model.C[i])
	}
	for i := uint(0); i < nn; i++ {
		for j := uint(0); j < nn; j++ {
			A[j*nn+i] = -D[i] * D[j] * A[j*nn+i]
		}
	}

	// Reusing A (which is model.G) to store U.
	U := A
	Λ := make([]float64, nn)
	if err := decomposition.SymEig(A, U, Λ, nn); err != nil {
		return nil, err
	}

	Δt := c.TimeStep

	coef := make([]float64, nn)
	temp := make([]float64, nn*nn)

	for i := uint(0); i < nn; i++ {
		coef[i] = math.Exp(Δt * Λ[i])
	}
	for i := uint(0); i < nn; i++ {
		for j := uint(0); j < nn; j++ {
			temp[j*nn+i] = coef[i] * U[i*nn+j]
		}
	}

	E := make([]float64, nn*nn)
	matrix.Multiply(U, temp, E, nn, nn, nn)

	// Technically, temp = temp[0 : nn*nc].
	for i := uint(0); i < nn; i++ {
		coef[i] = (coef[i] - 1) / Λ[i]
	}
	for i := uint(0); i < nn; i++ {
		for j := uint(0); j < nc; j++ {
			temp[j*nn+i] = coef[i] * U[i*nn+j] * D[j]
		}
	}

	F := make([]float64, nn*nc)
	matrix.Multiply(U, temp, F, nn, nn, nc)

	temperature := &Temperature{
		Cores: nc,
		Nodes: nn,

		system: system{
			D: D,

			Λ: Λ,
			U: U,

			E: E,
			F: F,

			Qamb: c.Ambience,
		},
	}

	return temperature, nil
}
