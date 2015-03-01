package analytic

import (
	"github.com/ready-steady/linear/matrix"
)

// Compute calculates the temperature profile corresponding to a power profile.
// The power profile is specified by a matrix capturing the power dissipation at
// a number of equidistant time moments (see TimeStep in Config). The ns
// parameter controls the number of samples that the temperature profile will
// contain; if the power profile contains more samples than needed, it will be
// accordingly truncated.
func (t *Temperature) Compute(P []float64, ns uint) []float64 {
	nc, nn := t.nc, t.nn

	S := make([]float64, nn*ns)
	matrix.Multiply(t.system.F, P, S, nn, nc, ns)

	for i, j, k := uint(1), uint(0), nn; i < ns; i++ {
		matrix.MultiplyAdd(t.system.E, S[j:k], S[k:k+nn], S[k:k+nn], nn, nn, 1)
		j += nn
		k += nn
	}

	Q := make([]float64, nc*ns)
	for i := uint(0); i < nc; i++ {
		for j := uint(0); j < ns; j++ {
			Q[j*nc+i] = t.system.D[i]*S[j*nn+i] + t.system.Qamb
		}
	}

	return Q
}
