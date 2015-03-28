package numeric

import (
	"fmt"
	"path"
	"testing"

	"github.com/ready-steady/assert"
	"github.com/ready-steady/fixture"
	"github.com/ready-steady/ode/dopri"
)

func TestNew(t *testing.T) {
	const (
		nc = 2
	)

	temperature := load(nc)

	assert.Equal(temperature.nc, uint(nc), t)
	assert.Equal(temperature.nn, uint(4*nc+12), t)

	assert.Equal(temperature.system.A, fixtureA, t)
	assert.Equal(temperature.system.B, fixtureB, t)
}

func load(nc uint) *Temperature {
	config := &Config{}
	fixture.Load(findFixture(fmt.Sprintf("%03d.json", nc)), config)

	integrator, _ := dopri.New(&dopri.Config{
		MaxStep:  0,
		TryStep:  0,
		AbsError: 1e-3,
		RelError: 1e-3,
	})

	return New(config, integrator)
}

func findFixture(name string) string {
	return path.Join("fixtures", name)
}
