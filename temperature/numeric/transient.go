package numeric

import (
	"errors"

	"github.com/ready-steady/linear/matrix"
)

// Compute calculates the temperature profile corresponding to a power profile
// and a timeline. The power profile is specified by a function func(time
// float64, power []float64) evaluating the power dissipation at an arbitrary
// time moment. The timeline should be an increasing sequence that contains at
// least two elements with the first one being the initial time.
func (t *Temperature) Compute(power func(float64, []float64),
	time []float64) ([]float64, error) {

	nc, nn := t.Cores, t.Nodes

	ns := uint(len(time))
	if ns < 2 {
		return nil, errors.New("the timeline should have at least two points")
	}

	A, B := t.system.A, t.system.B

	P := make([]float64, nc)

	derivative := func(time float64, S, dS []float64) {
		matrix.Multiply(A, S, dS, nn, nn, 1)
		power(time, P)
		for i := uint(0); i < nc; i++ {
			dS[i] += B[i] * P[i]
		}
	}

	S, _, _, err := t.integrator.Compute(derivative, time, make([]float64, nn))
	if err != nil {
		return nil, err
	}

	Q := make([]float64, ns*nc)
	for i := uint(0); i < nc; i++ {
		for j := uint(0); j < ns; j++ {
			Q[j*nc+i] = S[j*nn+i] + t.system.Qamb
		}
	}

	return Q, nil
}
