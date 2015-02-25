package numeric

import (
	"fmt"
	"path"
	"testing"

	"github.com/ready-steady/simulation/temperature"
	"github.com/ready-steady/support/assert"
)

func TestNew(t *testing.T) {
	temperature := load("002")

	assert.Equal(temperature.Cores, uint(2), t)
	assert.Equal(temperature.Nodes, uint(4*2+12), t)

	assert.Equal(temperature.system.A, fixtureA, t)
	assert.Equal(temperature.system.B, fixtureB, t)
}

func load(path string) *Temperature {
	config, _ := temperature.LoadConfig(findFixture(fmt.Sprintf("%s.json", path)))
	temperature, _ := New(config)
	return temperature
}

func findFixture(name string) string {
	return path.Join("fixtures", name)
}
