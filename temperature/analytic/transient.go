package analytic

import (
	"github.com/ready-steady/linear/matrix"
)

// Compute calculates the temperature profile corresponding to a power profile.
// The power profile is specified by a matrix capturing the power dissipation at
// a number of equidistant time moments (see TimeStep in Config). The sc
// parameter controls the number of samples that the temperature profile will
// contain; if the power profile contains more samples than needed, it will be
// accordingly truncated.
func (t *Temperature) Compute(P []float64, sc uint) []float64 {
	cc, nc := t.Cores, t.Nodes

	S := make([]float64, nc*sc)
	matrix.Multiply(t.system.F, P, S, nc, cc, sc)

	for i, j, k := uint(1), uint(0), nc; i < sc; i++ {
		matrix.MultiplyAdd(t.system.E, S[j:k], S[k:k+nc], S[k:k+nc], nc, nc, 1)
		j += nc
		k += nc
	}

	Q := make([]float64, cc*sc)
	for i := uint(0); i < cc; i++ {
		for j := uint(0); j < sc; j++ {
			Q[j*cc+i] = t.system.D[i]*S[j*nc+i] + t.system.Qamb
		}
	}

	return Q
}
