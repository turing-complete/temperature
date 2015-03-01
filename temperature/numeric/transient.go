package numeric

import (
	"github.com/ready-steady/linear/matrix"
)

// Compute calculates the temperature profile corresponding to a power profile.
// The power profile is specified by a function func(time float64, power
// []float64) evaluating the power dissipation at an arbitrary time moment. The
// time moments for which the temperature profile is computed are specified by
// the time array; see the corresponding ODE solver for further details.
//
// http://godoc.org/github.com/ready-steady/numeric/integration/ode#DormandPrince.Compute
func (t *Temperature) Compute(power func(float64, []float64),
	time []float64) ([]float64, []float64, error) {

	nc, nn := t.Cores, t.Nodes

	A, B := t.system.A, t.system.B

	P := make([]float64, nc)

	derivative := func(time float64, S, dS []float64) {
		matrix.Multiply(A, S, dS, nn, nn, 1)
		power(time, P)
		for i := uint(0); i < nc; i++ {
			dS[i] += B[i] * P[i]
		}
	}

	S, time, _, err := t.integrator.Compute(derivative, time, make([]float64, nn))
	if err != nil {
		return nil, nil, err
	}

	ns := uint(len(time))

	Q := make([]float64, ns*nc)
	for i := uint(0); i < nc; i++ {
		for j := uint(0); j < ns; j++ {
			Q[j*nc+i] = S[j*nn+i] + t.system.Qamb
		}
	}

	return Q, time, nil
}
