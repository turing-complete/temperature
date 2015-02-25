package analytic

import (
	"fmt"
	"math"
	"path"
	"testing"

	"github.com/ready-steady/support/assert"
	"github.com/ready-steady/support/fixture"
)

func TestNew(t *testing.T) {
	temperature := load("002")

	assert.Equal(temperature.Cores, uint(2), t)
	assert.Equal(temperature.Nodes, uint(4*2+12), t)

	assert.EqualWithin(temperature.system.D, fixtureD, 1e-14, t)

	assert.EqualWithin(abs(temperature.system.U), abs(fixtureU), 1e-9, t)
	assert.EqualWithin(temperature.system.Λ, fixtureΛ, 1e-9, t)

	assert.EqualWithin(temperature.system.E, fixtureE, 1e-9, t)
	assert.EqualWithin(temperature.system.F, fixtureF, 1e-9, t)
}

func abs(A []float64) []float64 {
	B := make([]float64, len(A))

	for i := range B {
		B[i] = math.Abs(A[i])
	}

	return B
}

func load(name string) *Temperature {
	config := &Config{}
	fixture.Load(findFixture(fmt.Sprintf("%s.json", name)), config)
	temperature, _ := New(config)
	return temperature
}

func findFixture(name string) string {
	return path.Join("fixtures", name)
}
