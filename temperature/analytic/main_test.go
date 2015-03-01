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
	const (
		nc = 2
	)

	temperature := load(nc)

	assert.Equal(temperature.nc, uint(nc), t)
	assert.Equal(temperature.nn, uint(4*nc+12), t)

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

func load(nc uint) *Temperature {
	config := &Config{}
	fixture.Load(findFixture(fmt.Sprintf("%03d.json", nc)), config)
	temperature, _ := New(config)
	return temperature
}

func findFixture(name string) string {
	return path.Join("fixtures", name)
}
