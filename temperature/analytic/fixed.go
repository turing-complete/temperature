package analytic

import (
	"errors"
	"math"

	"github.com/ready-steady/hotspot"
	"github.com/ready-steady/linear/matrix"
)

// Fixed represents an integrator of a thermal system with a fixed time step.
type Fixed struct {
	nc uint
	nn uint

	D []float64
	E []float64
	F []float64

	qamb float64
}

// NewFixed returns a new integrator with a fixed time step.
func NewFixed(config *Config) (*Fixed, error) {
	if config.TimeStep <= 0 {
		return nil, errors.New("the time step should be positive")
	}

	model := hotspot.New((*hotspot.Config)(&config.Config))
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
	if err := matrix.SymmetricEigen(A, U, Λ, nn); err != nil {
		return nil, err
	}

	Δt := config.TimeStep

	diag := make([]float64, nn)
	temp := make([]float64, nn*nn)

	E := make([]float64, nn*nn)
	for i := uint(0); i < nn; i++ {
		diag[i] = math.Exp(Δt * Λ[i])
		for j := uint(0); j < nn; j++ {
			temp[j*nn+i] = diag[i] * U[i*nn+j]
		}
	}
	matrix.Multiply(U, temp, E, nn, nn, nn)

	F := make([]float64, nn*nc)
	for i := uint(0); i < nn; i++ {
		diag[i] = (diag[i] - 1) / Λ[i]
		for j := uint(0); j < nc; j++ {
			temp[j*nn+i] = diag[i] * U[i*nn+j] * D[j]
		}
	}
	matrix.Multiply(U, temp, F, nn, nn, nc)

	temperature := &Fixed{
		nc: nc,
		nn: nn,

		D: D,

		E: E,
		F: F,

		qamb: config.Ambience,
	}

	return temperature, nil
}

// Compute calculates the temperature profile corresponding to a power profile.
// The power profile is specified by a matrix capturing the power dissipation at
// a number of equidistant time moments (see TimeStep in Config). The ns
// parameter controls the number of samples that the temperature profile will
// contain; if the power profile contains more samples than needed, it will be
// accordingly truncated.
func (self *Fixed) Compute(P []float64, ns uint) []float64 {
	nc, nn := self.nc, self.nn

	S := make([]float64, nn*ns)
	matrix.Multiply(self.F, P, S, nn, nc, ns)

	for i, j, k := uint(1), uint(0), nn; i < ns; i++ {
		matrix.MultiplyAdd(self.E, S[j:k], S[k:k+nn], S[k:k+nn], nn, nn, 1)
		j += nn
		k += nn
	}

	Q := make([]float64, nc*ns)
	for i := uint(0); i < nc; i++ {
		for j := uint(0); j < ns; j++ {
			Q[j*nc+i] = self.D[i]*S[j*nn+i] + self.qamb
		}
	}

	return Q
}
