package numeric

import (
	"github.com/ready-steady/linear/matrix"
)

// Compute calculates the temperature profile corresponding to the given power
// profile. The power profile is specified by a function func(time float64,
// power []float64) that computes the power dissipation at an arbitrary time
// moment. The sc parameter controls the number of samples that the temperature
// profile will contain (taking into account the time step given in Config).
func (t *Temperature) Compute(power func(float64, []float64), sc uint) ([]float64, error) {
	cc, nc := t.Cores, t.Nodes

	A, B := t.system.A, t.system.B

	P := make([]float64, cc)

	derivative := func(time float64, S, dS []float64) {
		matrix.Multiply(A, S, dS, nc, nc, 1)
		power(time, P)
		for i := uint(0); i < cc; i++ {
			dS[i] += B[i] * P[i]
		}
	}

	points := make([]float64, sc)
	for i := uint(0); i < sc; i++ {
		points[i] = float64(i) * t.system.Δt
	}

	S, _, err := t.integrator.Compute(derivative, points, make([]float64, nc))
	if err != nil {
		return nil, err
	}

	Q := make([]float64, sc*cc)
	for i := uint(0); i < cc; i++ {
		for j := uint(0); j < sc; j++ {
			Q[j*cc+i] = S[j*nc+i] + t.system.Qamb
		}
	}

	return Q, nil
}
