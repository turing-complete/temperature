package analytic

import (
	"math"

	"github.com/ready-steady/hotspot"
	"github.com/ready-steady/linear/matrix"
)

// Fluid represents an integrator of a thermal system with a fluid time step.
type Fluid struct {
	nc uint
	nn uint

	D []float64
	U []float64
	Λ []float64

	qamb float64
}

// NewFluid returns a new integrator with a fluid time step.
func NewFluid(config *Config) (*Fluid, error) {
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

	temperature := &Fluid{
		nc: nc,
		nn: nn,

		D: D,

		Λ: Λ,
		U: U,

		qamb: config.Ambience,
	}

	return temperature, nil
}

// Compute calculates the temperature profile corresponding to a power profile.
// The power profile is specified by a matrix P capturing the power dissipation
// and a vector ΔT assigning durations to the corresponding power samples.
func (self *Fluid) Compute(P []float64, ΔT []float64) []float64 {
	nc, nn, ns := self.nc, self.nn, uint(len(ΔT))

	D, U, Λ, qamb := self.D, self.U, self.Λ, self.qamb

	diag := make([]float64, nn)
	temp := make([]float64, nn*nn)

	E := make([]float64, nn*nn)
	F := make([]float64, nn*nc)

	S1 := make([]float64, nn)
	S2 := make([]float64, nn)

	Q := make([]float64, nc*ns)

	for i := uint(0); i < ns; i++ {
		Δt := ΔT[i]

		for j := uint(0); j < nn; j++ {
			diag[j] = math.Exp(Δt * Λ[j])
			for k := uint(0); k < nn; k++ {
				temp[k*nn+j] = diag[j] * U[j*nn+k]
			}
		}
		matrix.Multiply(U, temp, E, nn, nn, nn)

		for j := uint(0); j < nn; j++ {
			diag[j] = (diag[j] - 1) / Λ[j]
			for k := uint(0); k < nc; k++ {
				temp[k*nn+j] = diag[j] * U[j*nn+k] * D[k]
			}
		}
		matrix.Multiply(U, temp, F, nn, nn, nc)

		matrix.Multiply(F, P[i*nc:(i+1)*nc], S1, nn, nc, 1)
		matrix.MultiplyAdd(E, S2, S1, S1, nn, nn, 1)

		for j := uint(0); j < nc; j++ {
			Q[i*nc+j] = D[j]*S1[j] + qamb
		}

		S1, S2 = S2, S1
	}

	return Q
}
