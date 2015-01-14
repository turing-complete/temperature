package temperature

import (
	"math"
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestLoad(t *testing.T) {
	solver, _ := Load(findFixture("002.json"))

	assert.Equal(solver.Cores, uint32(2), t)
	assert.Equal(solver.Nodes, uint32(4*2+12), t)

	assert.AlmostEqual(solver.system.D, fixtureD, t)

	assert.AlmostEqual(abs(solver.system.U), abs(fixtureU), t)
	assert.AlmostEqual(solver.system.Λ, fixtureΛ, t)

	assert.AlmostEqual(solver.system.E, fixtureE, t)
	assert.AlmostEqual(solver.system.F, fixtureF, t)
}

func abs(A []float64) []float64 {
	B := make([]float64, len(A))

	for i := range B {
		B[i] = math.Abs(A[i])
	}

	return B
}
