package analytic

import (
	"errors"
	"math"

	"github.com/ready-steady/linear/decomposition"
	"github.com/ready-steady/linear/matrix"
	"github.com/turing-complete/hotspot"
)

// Fixed is an integrator of a thermal system with a fixed time step.
type Fixed struct {
	nc uint
	nn uint

	D []float64
	E []float64
	F []float64

	qamb float64
}

// NewFixed returns a new integrator.
func NewFixed(config *Config) (*Fixed, error) {
	if config.TimeStep <= 0.0 {
		return nil, errors.New("the time step should be positive")
	}

	model := hotspot.New((*hotspot.Config)(&config.Config))
	nc, nn := model.Cores, model.Nodes

	A := model.G // Reuse model.G to store A.
	D := model.C // Reuse model.C to store D.
	for i := uint(0); i < nn; i++ {
		D[i] = math.Sqrt(1.0 / model.C[i])
	}
	for i := uint(0); i < nn; i++ {
		for j := uint(0); j < nn; j++ {
			A[j*nn+i] = -D[i] * D[j] * A[j*nn+i]
		}
	}

	U := A // Reuse A (which is model.G) to store U.
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
//
// The power profile is specified by a matrix P containing power samples at a
// number of equidistant time moments (see TimeStep in Config).
func (self *Fixed) Compute(P []float64) []float64 {
	nc, nn := self.nc, self.nn
	ns := uint(len(P)) / nc

	S := make([]float64, nn*ns)
	Q := make([]float64, nc*ns)
	for i, n, q := uint(0), nc*ns, self.qamb; i < n; i++ {
		Q[i] = q
	}

	D, E, F := self.D, self.E, self.F
	matrix.Multiply(F, P, S, nn, nc, ns)
	{
		Si := S[:nn]
		Qi := Q[:nc]
		for k := uint(0); k < nc; k++ {
			Qi[k] += D[k] * Si[k]
		}
	}
	for i := uint(1); i < ns; i++ {
		Sj := S[(i-1)*nn : i*nn]
		Si := S[i*nn : (i+1)*nn]
		Qi := Q[i*nc : (i+1)*nc]
		matrix.MultiplyAdd(E, Sj, Si, Si, nn, nn, 1)
		for k := uint(0); k < nc; k++ {
			Qi[k] += D[k] * Si[k]
		}
	}

	return Q
}

// ComputeWithLeakage calculates the temperature profile and the total power
// profile corresponding to a dynamic power profile taking into account the
// leakage power.
//
// The dynamic power profile is specified by a matrix P containing power samples
// at a number of equidistant time moments (see TimeStep in Config). The dynamic
// power profile is overwritten with the total power profile.
func (self *Fixed) ComputeWithLeakage(P []float64, leak func([]float64, []float64)) []float64 {
	nc, nn := self.nc, self.nn
	ns := uint(len(P)) / nc

	S := make([]float64, nn*ns)
	Q := make([]float64, nc*ns)
	for i, n, q := uint(0), nc*ns, self.qamb; i < n; i++ {
		Q[i] = q
	}

	D, E, F := self.D, self.E, self.F
	{
		Si := S[:nn]
		Qi := Q[:nc]
		Pi := P[:nc]
		leak(Qi, Pi) // Use index 0 as if -1.
		matrix.Multiply(F, Pi, Si, nn, nc, 1)
		for k := uint(0); k < nc; k++ {
			Qi[k] += D[k] * Si[k]
		}
	}
	for i := uint(1); i < ns; i++ {
		Sj := S[(i-1)*nn : i*nn]
		Qj := Q[(i-1)*nc : i*nc]
		Si := S[i*nn : (i+1)*nn]
		Qi := Q[i*nc : (i+1)*nc]
		Pi := P[i*nc : (i+1)*nc]
		leak(Qj, Pi)
		matrix.Multiply(F, Pi, Si, nn, nc, 1)
		matrix.MultiplyAdd(E, Sj, Si, Si, nn, nn, 1)
		for k := uint(0); k < nc; k++ {
			Qi[k] += D[k] * Si[k]
		}
	}

	return Q
}
