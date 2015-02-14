package temperature

import (
	"math"
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestLoad(t *testing.T) {
	temperature, _ := Load(findFixture("002.json"))

	assert.Equal(temperature.Cores, uint(2), t)
	assert.Equal(temperature.Nodes, uint(4*2+12), t)

	assert.AlmostEqual(temperature.system.D, fixtureD, t)

	assert.AlmostEqual(abs(temperature.system.U), abs(fixtureU), t)
	assert.AlmostEqual(temperature.system.Λ, fixtureΛ, t)

	assert.AlmostEqual(temperature.system.E, fixtureE, t)
	assert.AlmostEqual(temperature.system.F, fixtureF, t)
}

func abs(A []float64) []float64 {
	B := make([]float64, len(A))

	for i := range B {
		B[i] = math.Abs(A[i])
	}

	return B
}
