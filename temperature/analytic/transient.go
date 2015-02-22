package analytic

import (
	"github.com/ready-steady/linear/matrix"
)

// ComputeTransient calculates the temperature profile corresponding to the
// given power profile. The sc parameter controls the number of steps/samples
// that the output matrix will contain; power profiles longer than this value
// get truncated.
func (t *Temperature) ComputeTransient(P []float64, sc uint) []float64 {
	cc := t.Cores
	nc := t.Nodes

	Q := make([]float64, cc*sc)
	S := make([]float64, nc*sc)

	matrix.Multiply(t.system.F, P, S, nc, cc, sc)

	var i, j, k uint

	for i, j, k = 1, 0, nc; i < sc; i++ {
		matrix.MultiplyAdd(t.system.E, S[j:k], S[k:k+nc], S[k:k+nc], nc, nc, 1)
		j += nc
		k += nc
	}

	for i = 0; i < cc; i++ {
		for j = 0; j < sc; j++ {
			Q[cc*j+i] = t.system.D[i]*S[nc*j+i] + t.Config.AmbientTemp
		}
	}

	return Q
}
