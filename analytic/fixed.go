package analytic

import (
	"errors"
	"math"

	"github.com/ready-steady/linear/decomposition"
	"github.com/ready-steady/linear/matrix"
	"github.com/turing-complete/hotspot"
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
	if config.TimeStep <= 0.0 {
		return nil, errors.New("the time step should be positive")
	}

	model := hotspot.New((*hotspot.Config)(&config.Config))
	nc, nn := model.Cores, model.Nodes

	// Reusing model.G to store A and model.C to store D.
	A := model.G
	D := model.C
	for i := uint(0); i < nn; i++ {
		D[i] = math.Sqrt(1.0 / model.C[i])
	}
	for i := uint(0); i < nn; i++ {
		for j := uint(0); j < nn; j++ {
			A[j*nn+i] = -D[i] * D[j] * A[j*nn+i]
		}
	}

	// Reusing A (which is model.G) to store U.
	U := A
	Λ := make([]float64, nn)
	if err := decomposition.SymmetricEigen(A, U, Λ, nn); err != nil {
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
		diag[i] = (diag[i] - 1.0) / Λ[i]
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
// The power profile is specified by a matrix P capturing the power dissipation
// at a number of equidistant time moments according to TimeStep in Config.
func (self *Fixed) Compute(P []float64) []float64 {
	nc, nn := self.nc, self.nn
	ns := uint(len(P)) / nc

	D, E, F, qamb := self.D, self.E, self.F, self.qamb

	S := make([]float64, nn*ns)
	matrix.Multiply(F, P, S, nn, nc, ns)

	for i, j, k := uint(1), uint(0), nn; i < ns; i++ {
		matrix.MultiplyAdd(E, S[j:k], S[k:k+nn], S[k:k+nn], nn, nn, 1)
		j += nn
		k += nn
	}

	Q := make([]float64, nc*ns)
	for i := uint(0); i < nc; i++ {
		for j := uint(0); j < ns; j++ {
			Q[j*nc+i] = D[i]*S[j*nn+i] + qamb
		}
	}

	return Q
}
