package numeric

import (
	"fmt"
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

	assert.Equal(temperature.Cores, uint(nc), t)
	assert.Equal(temperature.Nodes, uint(4*nc+12), t)

	assert.Equal(temperature.system.A, fixtureA, t)
	assert.Equal(temperature.system.B, fixtureB, t)
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
