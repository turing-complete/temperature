package numeric

import (
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestCompute(t *testing.T) {
	temperature := load("002")

	cc := uint(2)
	sc := uint(len(fixtureP)) / cc
	Δt := temperature.system.Δt

	power := func(time float64, power []float64) {
		k := uint(time / Δt)
		for i := uint(0); i < cc; i++ {
			power[i] = fixtureP[k*cc+i]
		}
	}

	Q, err := temperature.Compute(power, sc)

	assert.Success(err, t)
	assert.EqualWithin(Q, fixtureQ, 2e-10, t)
}
